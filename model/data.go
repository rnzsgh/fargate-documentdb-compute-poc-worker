package model

import (
	"context"
	"fmt"
	"time"

	"github.com/mongodb/mongo-go-driver/bson"
	"github.com/mongodb/mongo-go-driver/bson/primitive"
	"github.com/mongodb/mongo-go-driver/mongo"

	docdb "github.com/rnzsgh/fargate-documentdb-compute-poc-worker/db"
)

type Data struct {
	Id      *primitive.ObjectID `json:"id" bson:"_id"`
	P       []float64           `json:"p" bson:"p"`
	Entropy float64             `json:"entropy" bson:"entropy"`
}

func DataFindById(id *primitive.ObjectID) (*Data, error) {
	ctx, _ := context.WithTimeout(context.Background(), 5*time.Second)
	data := &Data{}
	err := dataCollection().FindOne(ctx, bson.D{{"_id", id}}).Decode(data)
	if err == mongo.ErrNoDocuments {
		return nil, nil
	}
	return data, err
}

func dataCollection() *mongo.Collection {
	return docdb.Client.Database("work").Collection("data")
}

func DataUpdateEntropy(id *primitive.ObjectID, e float64) error {
	ctx, _ := context.WithTimeout(context.Background(), 5*time.Second)
	res, err := dataCollection().UpdateOne(
		ctx,
		bson.D{{"_id", id}},
		bson.D{{"$set", bson.D{{"entopy", e}}}})

	if err == nil {
		if res.MatchedCount != 1 && res.ModifiedCount != 1 {
			return fmt.Errorf("Data results not updated - data : %s", id.Hex())
		}
		return nil
	}

	return err
}
