package dao

import (
	"context"
	"coolcar/shared/id"
	mgutil "coolcar/shared/mongo"
	"coolcar/shared/mongo/objid"
	mongotesting "coolcar/shared/mongo/testing"
	"os"
	"testing"

	"go.mongodb.org/mongo-driver/bson"
)

// var mongoURI string

func TestResolveAccountID(t *testing.T) {
	c := context.Background()
	// mc, err := mongo.Connect(c, options.Client().ApplyURI(mongoURI))
	mc, err := mongotesting.NewClient(c)
	if err != nil {
		t.Fatalf("cannot connect mongodb: %v", err)
	}

	m := NewMongo(mc.Database("coolcar"))
	_, err = m.col.InsertMany(c, []interface{}{
		bson.M{
			// mgutil.IDFieldName: mustObjID("645218f71c4e0349d1012c6f"),
			mgutil.IDFieldName: objid.MustFromID(id.AccountID("645218f71c4e0349d1012c6f")),
			openIDField:        "openid_1",
		},
		bson.M{
			mgutil.IDFieldName: objid.MustFromID(id.AccountID("645218f71c4e0349d1012c70")),
			openIDField:        "openid_2",
		},
	})
	if err != nil {
		t.Fatalf("cannot insert initial values: %v", err)
	}

	// 替换函数测试
	// mgutil.NewObjID = func() primitive.ObjectID {
	// 	return objid.MustFromID(id.AccountID("645218f71c4e0349d1012c71")) // 固定插入的ID，使每次测试的条件相同，稳定重现
	// }
	mgutil.NewObjectIDWithValue(id.AccountID("645218f71c4e0349d1012c71"))

	cases := []struct {
		name   string
		openID string
		want   string
	}{
		{name: "existing_user", openID: "openid_1", want: "645218f71c4e0349d1012c6f"},
		{name: "another_existing_user", openID: "openid_2", want: "645218f71c4e0349d1012c70"},
		{name: "new_user", openID: "openid_3", want: "645218f71c4e0349d1012c71"},
	}

	for _, cc := range cases {
		// t.Run运行子测试 _为index
		t.Run(cc.name, func(t *testing.T) {
			id, err := m.ResolveAccountID(context.Background(), cc.openID)
			if err != nil {
				t.Errorf("faile resolve account id for %q: %v", id, err)
			}
			if id.String() != cc.want {
				t.Errorf("resolve account id want: %q, get: %q", cc.want, id)
			}
		})
	}
}

// func mustObjID(hex string) primitive.ObjectID {
// 	objID, err := primitive.ObjectIDFromHex(hex)
// 	if err != nil {
// 		panic(err)
// 	}
// 	return objID
// }

func TestMain(m *testing.M) {
	os.Exit(mongotesting.RunWithMongoInDocker(m))
}
