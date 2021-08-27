package database

import (
	"context"
	"craft/models"
	"fmt"
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type MongoUtilsInterfce interface {
	Insert(models.CraftConfiguration) (models.CraftConfiguration, error)
	FetchMyData(params models.Params) (models.IDS, error)
}

type MongoUtilsClient struct {
	Ctx  context.Context
	Coll *mongo.Collection
}

func (mu *MongoUtilsClient) Insert(payload models.CraftConfiguration) (craftConfig models.CraftConfiguration, err error) {
	timestamp := makeTimestamp()
	fmt.Println("inside Insert timestamp: ", timestamp)
	payload.Timestamp = timestamp
	fmt.Println("inside Insert metgod payload:", payload)
	res, err := mu.Coll.InsertOne(mu.Ctx, payload)
	if err != nil {
		return
	}
	fmt.Println("resposne: ", res)
	return
}

func (mu *MongoUtilsClient) FetchMyData(params models.Params) (IDs models.IDS, err error) {
	fmt.Println("inside FetchMyData ")
	// IDs = append(IDs, "dasdasda")
	// IDs = append(IDs, "dsd23")
	var pipeline []primitive.M
	filter := primitive.M{"$and": []primitive.M{{"attributes.key": params.Key, "attributes.value": params.Value, "timestamp": params.Timestamp}}}
	//filter := primitive.M{"$and": []primitive.M{{"timestamp": params.Timestamp}}}

	matchStage := primitive.M{"$match": filter}
	unwindStage := primitive.M{"$unwind": "$attributes"}
	projectStage := primitive.M{"$project": primitive.M{"id": 1}}
	groupStage := primitive.M{"$group": primitive.M{"_id": nil, "ids": primitive.M{"$addToSet": "$id"}}}
	//([{"$unwind":"$attributes"},{"$match":{"attributes.key":"mykey"}},{"$project":{"id":1}},{"$group":{"_id":"nil","ids":{"$push":"$id"}}}])
	pipeline = append(pipeline, unwindStage)
	pipeline = append(pipeline, matchStage)

	pipeline = append(pipeline, projectStage)
	pipeline = append(pipeline, groupStage)

	opt := options.Aggregate()
	opt.SetAllowDiskUse(true)
	fmt.Println("inside pipeline:  ", pipeline)

	cur, err := mu.Coll.Aggregate(mu.Ctx, pipeline, opt)
	// defer func() {
	// 	fmt.Println(cur.Close(mu.Ctx), "could not close cursor")
	// }()
	if err != nil {
		return IDs, err
	}
	fmt.Println("cursoesasdad: ", cur)

	for cur.Next(mu.Ctx) {
		fmt.Println("cursoe: ", cur)
		err := cur.Decode(&IDs)
		fmt.Println("cursoesasdad: ", IDs)
		if err != nil {
			fmt.Println("error decoding: ", err)
			return IDs, err
		}
	}
	return
}

func makeTimestamp() int64 {
	return time.Now().UnixNano() / 1000000
}
