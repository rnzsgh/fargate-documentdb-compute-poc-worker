package main

import (
	"context"
	"encoding/json"
	"flag"
	"time"

	log "github.com/golang/glog"
	"github.com/mongodb/mongo-go-driver/bson/primitive"
	docdb "github.com/rnzsgh/fargate-documentdb-compute-poc-worker/db"
	"github.com/rnzsgh/fargate-documentdb-compute-poc-worker/model"
	"gonum.org/v1/gonum/stat"
)

func init() {
	flag.Parse()
	flag.Lookup("logtostderr").Value.Set("true")
}

func main() {

	msgs := docdb.TaskReceiveQueue.Listen(2)

	log.Info("Listening for tasks")

	for msg := range msgs {

		var vals map[string]interface{}
		if err := json.Unmarshal([]byte(msg.Payload), &vals); err != nil {
			log.Errorf("Unable to unmarshal msg payload - reason: %v", err)
			continue
		}

		log.Infof(
			"Task received - job: %s - task: %s - data: %s",
			vals["jobId"].(string),
			vals["taskId"].(string),
			vals["dataId"].(string),
		)

		dataId, err := primitive.ObjectIDFromHex(vals["dataId"].(string))
		if err != nil {
			log.Errorf("Unable to data - reason: %v", err)
			continue
		}

		data, err := model.DataFindById(&dataId)
		if err != nil {
			log.Errorf("Unable to load data by id: %s - reason: %v", dataId, err)
			continue
		}

		log.Infof("Starting data processing - id: %s", dataId.Hex())
		start := time.Now()

		data.Entropy = stat.Entropy(data.P)
		t := time.Now()
		elapsed := t.Sub(start)

		if err := model.DataUpdateEntropy(data.Id, data.Entropy); err != nil {
			log.Errorf("Problem updating entropy - id: %s - reason %v", dataId.Hex(), err)
		}

		ctx, _ := context.WithTimeout(context.Background(), 5*time.Second)
		if err := msg.Done(ctx); err != nil {
			log.Errorf("Unable to mark message as done - id: %s - reason %v", dataId.Hex(), err)
		}

		if err := docdb.TaskCompletedQueue.Enqueue(ctx, msg.Payload, 30); err != nil {
			log.Errorf("Problem sending task completed msg - id: %s - reason %v", dataId.Hex(), err)
		}

		log.Infof("Data id processed in: %v - id: %s", elapsed, dataId.Hex())
		log.Infof("Data results updated in database and task is completed - data id: %s", dataId.Hex())
	}

	log.Flush()
}
