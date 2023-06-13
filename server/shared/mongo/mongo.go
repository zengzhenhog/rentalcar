package mgutil

import (
	"coolcar/shared/mongo/objid"
	"fmt"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

// IDField defines the field name for mongo document id
const (
	IDFieldName        = "_id"
	UpdatedAtFieldName = "updatedat"
)

// objID defines the object id field
type IDField struct {
	ID primitive.ObjectID `bson:"_id"`
}

type UpdatedAtField struct {
	UpdatedAt int64 `bson:"updatedat"`
}

var NewObjID = primitive.NewObjectID

func NewObjectIDWithValue(id fmt.Stringer) {
	NewObjID = func() primitive.ObjectID {
		return objid.MustFromID(id)
	}
}

var UpdatedAt = func() int64 {
	return time.Now().UnixNano()
}

// set returns a $set update document
func Set(v interface{}) bson.M {
	return bson.M{"$set": v}
}

// setOnInsert查到了直接返回，查不到才插入
func SetOnInsert(v interface{}) bson.M {
	return bson.M{"$setOnInsert": v}
}
