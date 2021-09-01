package database

import (
	"context"
	"craft/models"
	"errors"
	"fmt"
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

const (
	limit = 500
)

type MongoUtilsInterfce interface {
	Insert(payload models.CraftConfiguration) (models.CraftConfiguration, error)
	FetchMyData(params models.Params) (models.IDS, error)
	UpdateMydataConfig(id string, keys []string, UpdateContent models.CraftConfiguration) (models.IDS, string, error)
}

type MongoUtilsClient struct {
	Ctx  context.Context
	Coll *mongo.Collection
}

func (mu *MongoUtilsClient) Insert(payload models.CraftConfiguration) (craftConfig models.CraftConfiguration, err error) {
	timestamp := makeTimestamp()
	payload.Timestamp = timestamp
	res, err := mu.Coll.InsertOne(mu.Ctx, payload)
	if err != nil {
		return
	}
	fmt.Println("resposne: ", res)
	return
}

func (mu *MongoUtilsClient) FetchMyData(params models.Params) (IDs models.IDS, err error) {
	fmt.Println("inside FetchMyData ")

	var pipeline []primitive.M
	filter := primitive.M{"$and": []primitive.M{{"attributes.key": params.Key, "attributes.value": params.Value, "timestamp": params.Timestamp}}}
	limitStage := primitive.M{"$limit": params.Limit}
	matchStage := primitive.M{"$match": filter}
	unwindStage := primitive.M{"$unwind": "$attributes"}
	projectStage := primitive.M{"$project": primitive.M{"id": 1}}
	groupStage := primitive.M{"$group": primitive.M{"_id": nil, "ids": primitive.M{"$addToSet": "$id"}}}
	pipeline = append(pipeline, unwindStage)
	pipeline = append(pipeline, matchStage)

	pipeline = append(pipeline, projectStage)
	pipeline = append(pipeline, limitStage)
	pipeline = append(pipeline, groupStage)

	opt := options.Aggregate()
	opt.SetAllowDiskUse(true)
	fmt.Println("inside pipeline:  ", pipeline)

	cur, err := mu.Coll.Aggregate(mu.Ctx, pipeline, opt)

	if err != nil {
		return IDs, err
	}
	defer cur.Close(mu.Ctx)
	for cur.Next(mu.Ctx) {
		err := cur.Decode(&IDs)
		if err != nil {
			fmt.Println("error decoding: ", err)
			return IDs, err
		}
	}
	return
}

func (mu *MongoUtilsClient) UpdateMydataConfig(id string, keys []string, UpdateContent models.CraftConfiguration) (IDs models.IDS, msg string, err error) {
	fmt.Println("inside UpdateMydataConfig ")
	var pipeline []primitive.M
	filter := primitive.M{"$and": []primitive.M{{"attributes.key": primitive.M{"$in": keys}}, {"id": id}}}

	matchStage := primitive.M{"$match": filter}
	unwindStage := primitive.M{"$unwind": "$attributes"}
	groupStage := primitive.M{"$group": primitive.M{"_id": nil, "ids": primitive.M{"$addToSet": "$attributes.key"}}}

	pipeline = append(pipeline, unwindStage)
	pipeline = append(pipeline, matchStage)

	pipeline = append(pipeline, groupStage)

	opt := options.Aggregate()
	opt.SetAllowDiskUse(true)
	cur, err := mu.Coll.Aggregate(mu.Ctx, pipeline, opt)

	if err != nil {
		return
	}

	defer cur.Close(mu.Ctx)
	for cur.Next(mu.Ctx) {
		err := cur.Decode(&IDs)
		custErr := errors.New("Bad Request")
		if err != nil {
			fmt.Println("error decoding: ", err)
			return IDs, msg, err
		}
		msg = "Bad Request"
		return IDs, msg, custErr
	}

	timestamp := makeTimestamp()
	UpdateContent.Timestamp = timestamp
	var craftConfig models.CraftConfiguration
	opts := options.FindOneAndUpdate().SetUpsert(true)

	mu.Coll.FindOneAndUpdate(mu.Ctx, primitive.M{"id": id}, primitive.M{"$addToSet": primitive.M{"attributes": primitive.M{"$each": UpdateContent.Attributes}}, "$set": primitive.M{"timestamp": timestamp}}, opts).Decode(&craftConfig)
	if err != nil {
		return
	}
	fmt.Println("resposne: ", craftConfig)
	return
}

func makeTimestamp() int64 {
	return time.Now().UnixNano() / 1000000
}
