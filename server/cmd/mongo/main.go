package main

import (
	"context"
	"fmt"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func main() {
	c := context.Background()
	mc, err := mongo.Connect(c, options.Client().ApplyURI("mongodb://admin:123456@localhost:27017/?authSource=admin&readPreference=primary&ssl=false&directConnection=true"))
	if err != nil {
		panic(err)
	}
	col := mc.Database("coolcar").Collection("account")

	// insertRows(c, col)
	// findOneRows(c, col)
	findRows(c, col)
}

// 查找全部
func findRows(c context.Context, col *mongo.Collection) {
	cur, err := col.Find(c, bson.M{})
	if err != nil {
		panic(err)
	}

	for cur.Next(c) {
		var row struct {
			ID     primitive.ObjectID `bson:"_id"`
			OpenID string             `bson:"open_id"`
		}
		err = cur.Decode(&row)
		if err != nil {
			panic(err)
		}
		fmt.Printf("%+v\n", row)
	}
}

// 查找一条
func findOneRows(c context.Context, col *mongo.Collection) {
	res := col.FindOne(c, bson.M{"open_id": "123"})
	fmt.Printf("%+v\n", res)
	var row struct {
		ID     primitive.ObjectID `bson:"_id"`
		OpenID string             `bson:"open_id"`
	}
	err := res.Decode(&row)
	if err != nil {
		panic(err)
	}
	fmt.Printf("%+v\n", row)
}

func insertRows(c context.Context, col *mongo.Collection) {
	res, err := col.InsertMany(c, []interface{}{
		bson.M{"open_id": "123"},
		bson.M{"open_id": "456"},
	})

	if err != nil {
		panic(err)
	}

	fmt.Printf("%+v", res)
}
