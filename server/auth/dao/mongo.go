package dao

import (
	"context"
	"coolcar/shared/id"
	mgutil "coolcar/shared/mongo"
	"coolcar/shared/mongo/objid"
	"fmt"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

const openIDField = "open_id"

type Mongo struct {
	col *mongo.Collection
}

func NewMongo(db *mongo.Database) *Mongo {
	return &Mongo{
		col: db.Collection("account"),
	}
}

func (m *Mongo) ResolveAccountID(c context.Context, openID string) (id.AccountID, error) {
	insertedID := mgutil.NewObjID()
	res := m.col.FindOneAndUpdate(
		c,
		bson.M{openIDField: openID},
		mgutil.SetOnInsert(bson.M{mgutil.IDFieldName: insertedID, openIDField: openID}), // setOnInsert查到了直接返回，查不到才插入
		options.FindOneAndUpdate().SetUpsert(true).SetReturnDocument(options.After),
	)

	if err := res.Err(); err != nil {
		return "", fmt.Errorf("cannot FindOneAndUpdate: %v", err)
	}

	var row mgutil.IDField
	err := res.Decode(&row)
	if err != nil {
		return "", fmt.Errorf("cannot decode result: %v", err)
	}

	// return row.ID.Hex(), nil
	return objid.ToAccountID(row.ID), nil
}
