package db

import (
	"context"
	"fmt"
	"os"
	"time"

	"github.com/rnzsgh/fargate-documentdb-compute-poc-worker/cloud"

	"github.com/mongodb/mongo-go-driver/mongo"
	"github.com/mongodb/mongo-go-driver/mongo/options"
	"github.com/mongodb/mongo-go-driver/mongo/readpref"

	log "github.com/golang/glog"
	queue "github.com/rnzsgh/documentdb-queue"
)

var Client *mongo.Client

var TaskReceiveQueue *queue.Queue
var TaskCompletedQueue *queue.Queue

func init() {

	endpoint := os.Getenv("DOCUMENT_DB_ENDPOINT")
	port := os.Getenv("DOCUMENT_DB_PORT")
	user := os.Getenv("DOCUMENT_DB_USER")
	caFile := os.Getenv("DOCUMENT_DB_PEM")

	password := cloud.Secrets.DatabasePassword

	connectionUri := fmt.Sprintf("mongodb://%s:%s@%s:%s/work?ssl=true", user, password, endpoint, port)

	if len(os.Getenv("LOCAL")) == 0 {
		connectionUri = connectionUri + "&replicaSet=rs0"
	}

	var err error
	if TaskReceiveQueue, err = queue.NewQueue(
		"work",
		"dispatchQueue",
		connectionUri,
		caFile, 5*time.Second,
	); err != nil {
		log.Errorf("Unable to create work dispatch queue - endpoint: %s - reason: %v", endpoint, err)
	}

	if TaskCompletedQueue, err = queue.NewQueue(
		"work",
		"responseQueue",
		connectionUri,
		caFile,
		5*time.Second,
	); err != nil {
		log.Errorf("Unable to create work response queue - endpoint: %s - reason: %v", endpoint, err)
	}

	Client, err = mongo.NewClientWithOptions(
		connectionUri,
		options.Client().SetSSL(
			&options.SSLOpt{
				Enabled:  true,
				Insecure: true,
				CaFile:   caFile,
			},
		),
	)

	if err != nil {
		log.Errorf("Unable to create new db client: %v", err)
	}

	ctx, _ := context.WithTimeout(context.Background(), 10*time.Second)
	err = Client.Connect(ctx)

	if err != nil {
		log.Errorf("Unable to connect to db: %v", err)
	}

	if err = ping(); err != nil {
		log.Errorf("Unable to ping db: %v", err)
	}
}

func ping() error {
	ctx, _ := context.WithTimeout(context.Background(), 10*time.Second)
	return Client.Ping(ctx, readpref.Primary())

}
