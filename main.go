package main

import (
	"flag"
	"os"
	"time"

	log "github.com/golang/glog"
	"github.com/mongodb/mongo-go-driver/bson/primitive"
	appmath "github.com/rnzsgh/fargate-documentdb-compute-poc-worker/math"
	"github.com/rnzsgh/fargate-documentdb-compute-poc-worker/model"
)

func init() {
	flag.Parse()
	flag.Lookup("logtostderr").Value.Set("true")
}

type Data struct {
	Id      *primitive.ObjectID `json:"id" bson:"_id"`
	X       [][]float32         `json:"x" bson:"x"`
	W       [][]float32         `json:"w" bson:"w"`
	Results [][]float32         `json:"results" bson:"results"`
}

func main() {
	dataId, err := primitive.ObjectIDFromHex(os.Getenv("DATA_ID"))
	if err != nil {
		log.Errorf("Unable to load DATA_ID from env variable - value: %s", os.Getenv("DATA_ID"))
		return
	}

	data, err := model.DataFindById(&dataId)
	if err != nil {
		log.Errorf("Unable to load data by id: %s - reason: %v", dataId.Hex(), err)
		return
	}

	log.Infof("Starting data processing - id: %s", dataId.Hex())
	start := time.Now()
	results, err := appmath.Multiply(data.X, appmath.Transpose(data.W))
	t := time.Now()
	elapsed := t.Sub(start)
	if err != nil {
		log.Errorf("Unable to multiply - data id %s - reason: %v", dataId.Hex(), err)
		return
	}
	log.Infof("Data id processed in: %v - id: %s", elapsed, dataId.Hex())

	if err := model.DataUpdateResults(&dataId, results); err != nil {
		log.Errorf("Problem updating data results - id: %s - reason %v", dataId.Hex(), err)
		return
	}

	log.Infof("Data results updated in database and task is completed - data id: %s", dataId.Hex())

	log.Flush()
}
