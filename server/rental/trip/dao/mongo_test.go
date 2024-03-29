package dao

import (
	"context"
	rentalpb "coolcar/rental/api/gen/v1"
	"coolcar/shared/id"
	mgutil "coolcar/shared/mongo"
	"coolcar/shared/mongo/objid"
	mongotesting "coolcar/shared/mongo/testing"
	"os"
	"testing"

	"github.com/google/go-cmp/cmp"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"google.golang.org/protobuf/testing/protocmp"
)

func TestCreateTrip(t *testing.T) {
	c := context.Background()
	mc, err := mongotesting.NewClient(c)
	if err != nil {
		t.Fatalf("cannot connect mongodb: %v", err)
	}

	db := mc.Database("coolcar")
	err = mongotesting.SetupIndexes(c, db)
	if err != nil {
		t.Fatalf("cannot setup indexes: %v", err)
	}
	m := NewMongo(db)

	cases := []struct {
		name       string
		tripID     string
		accountID  string
		tripStatus rentalpb.TripStatus
		wantErr    bool
	}{
		{
			name:       "finished",
			tripID:     "6485d7452cf384b4e217ae91",
			accountID:  "account1",
			tripStatus: rentalpb.TripStatus_FINISHED,
		},
		{
			name:       "another_finished",
			tripID:     "6485d7452cf384b4e217ae92",
			accountID:  "account1",
			tripStatus: rentalpb.TripStatus_FINISHED,
		},
		{
			name:       "in_progress",
			tripID:     "6485d7452cf384b4e217ae99",
			accountID:  "account1",
			tripStatus: rentalpb.TripStatus_IN_PROGRESS,
		},
		{
			name:       "another_in_progress",
			tripID:     "6485d7452cf384b4e217ae95",
			accountID:  "account1",
			tripStatus: rentalpb.TripStatus_IN_PROGRESS,
			wantErr:    true,
		},
		{
			name:       "in_progress_by_another_account",
			tripID:     "6485d7452cf384b4e217ae98",
			accountID:  "account2",
			tripStatus: rentalpb.TripStatus_IN_PROGRESS,
		},
	}

	for _, cc := range cases {
		// 需要按顺序执行, t.Run不能保证子测试的顺序, 所以不使用t.Run
		// mgutil.NewObjID = func() primitive.ObjectID {
		// 	return objid.MustFromID(id.TripID(cc.tripID))
		// }
		mgutil.NewObjectIDWithValue(id.TripID(cc.tripID))
		tr, err := m.CreateTrip(c, &rentalpb.Trip{
			AccountId: cc.accountID,
			Status:    cc.tripStatus,
		})
		if cc.wantErr {
			if err == nil {
				t.Errorf("%s: error expected; got none", cc.name)
			}
			continue
		}
		if err != nil {
			t.Errorf("%s: error creating trip: %v", cc.name, err)
			continue
		}
		if tr.ID.Hex() != cc.tripID {
			t.Errorf("%s: incorrect trip id; want: %q; got: %q", cc.name, cc.tripID, tr.ID.Hex())
		}
	}
}

// var mongoURI string
func TestGetTrip(t *testing.T) {
	// mongoURI = "mongodb://admin:123456@localhost:27017"
	c := context.Background()
	// mc, err := mongo.Connect(c, options.Client().ApplyURI(mongoURI))
	mc, err := mongotesting.NewClient(c)
	if err != nil {
		t.Fatalf("cannot connect mongodb: %v", err)
	}

	m := NewMongo(mc.Database("coolcar"))
	acct := id.AccountID("account1")
	// 一起测试时，TestCreateTrip会固定NewObjID(最后一个子测试)，需要重新设置回去
	mgutil.NewObjID = primitive.NewObjectID

	tr, err := m.CreateTrip(c, &rentalpb.Trip{
		AccountId: acct.String(),
		CarId:     "car1",
		Start: &rentalpb.LocationStatus{
			PoiName: "startpoint",
			Location: &rentalpb.Location{
				Latitude:  30,
				Longitude: 120,
			},
		},
		End: &rentalpb.LocationStatus{
			PoiName:  "endpoint",
			FeeCent:  10000,
			KmDriven: 35,
			Location: &rentalpb.Location{
				Latitude:  35,
				Longitude: 120,
			},
		},
		Status: rentalpb.TripStatus_FINISHED,
	})

	if err != nil {
		t.Fatalf("cannot create trip: %v", err)
	}
	// t.Errorf("%+v", tr)
	// t.Errorf("inserted row %s with updatedat %v", tr.ID, tr.UpdatedAt)

	// got, err := m.GetTrip(c, tr.ID.Hex(), "account1")
	got, err := m.GetTrip(c, objid.ToTripID(tr.ID), acct)
	if err != nil {
		t.Errorf("cannot get trip: %v", err)
	}
	// t.Errorf("got trip: %+v", got)

	// got.Trip.Start.PoiName = "badstart"
	// 安装 go get -u github.com/google/go-cmp/cmp，tr内容里含有proto buffers，需要加参数protocmp.Transform()
	if diff := cmp.Diff(tr, got, protocmp.Transform()); diff != "" {
		t.Errorf("result differs; -want +got: %s", diff)
	}
}

func TestGetTrips(t *testing.T) {
	rows := []struct {
		id        string
		accountID string
		status    rentalpb.TripStatus
	}{
		{
			id:        "6485d7452cf384b4e217ae91",
			accountID: "account_id_for_get_trips",
			status:    rentalpb.TripStatus_FINISHED,
		},
		{
			id:        "6485d7452cf384b4e217ae92",
			accountID: "account_id_for_get_trips",
			status:    rentalpb.TripStatus_FINISHED,
		},
		{
			id:        "6485d7452cf384b4e217ae93",
			accountID: "account_id_for_get_trips",
			status:    rentalpb.TripStatus_FINISHED,
		},
		{
			id:        "6485d7452cf384b4e217ae94",
			accountID: "account_id_for_get_trips",
			status:    rentalpb.TripStatus_IN_PROGRESS,
		},
		{
			id:        "6485d7452cf384b4e217ae95",
			accountID: "account_id_for_get_trips_1",
			status:    rentalpb.TripStatus_FINISHED,
		},
	}

	c := context.Background()
	mc, err := mongotesting.NewClient(c)
	if err != nil {
		t.Fatalf("cannot connect mongodb: %v", err)
	}

	m := NewMongo(mc.Database("coolcar"))

	for _, r := range rows {
		mgutil.NewObjectIDWithValue(id.TripID(r.id))
		_, err := m.CreateTrip(c, &rentalpb.Trip{
			AccountId: r.accountID,
			Status:    r.status,
		})
		if err != nil {
			t.Fatalf("cannot create rows: %v", err)
		}
	}

	cases := []struct {
		name       string
		accountID  string
		status     rentalpb.TripStatus
		wantCount  int
		wantOnlyID string
	}{
		{
			name:      "get_all",
			accountID: "account_id_for_get_trips",
			status:    rentalpb.TripStatus_TS_NOT_SPECIFIED,
			wantCount: 4,
		},
		{
			name:       "get_in_progress",
			accountID:  "account_id_for_get_trips",
			status:     rentalpb.TripStatus_IN_PROGRESS,
			wantCount:  1,
			wantOnlyID: "6485d7452cf384b4e217ae94",
		},
	}

	for _, cc := range cases {
		t.Run(cc.name, func(t *testing.T) {
			res, err := m.GetTipps(context.Background(), id.AccountID(cc.accountID), cc.status)
			if err != nil {
				t.Errorf("cannot get trips: %v", err)
			}

			if cc.wantCount != len(res) {
				t.Errorf("incorrect result count; want: %d; got: %d", cc.wantCount, len(res))
			}

			if cc.wantOnlyID != "" && len(res) > 0 {
				if cc.wantOnlyID != res[0].ID.Hex() {
					t.Errorf("only_id incorrect; want: %q; got: %q", cc.wantOnlyID, res[0].ID.Hex())
				}
			}
		})
	}
}

func TestUpdateTrip(t *testing.T) {
	c := context.Background()
	mc, err := mongotesting.NewClient(c)
	if err != nil {
		t.Fatalf("cannot connect mongodb: %v", err)
	}
	m := NewMongo(mc.Database("coolcar"))

	tid := id.TripID("6485d7452cf384b4e217ac97")
	aid := id.AccountID("account_for_update")

	var now int64 = 10000
	mgutil.NewObjectIDWithValue(tid)
	mgutil.UpdatedAt = func() int64 {
		return now
	}

	tr, err := m.CreateTrip(c, &rentalpb.Trip{
		AccountId: aid.String(),
		Status:    rentalpb.TripStatus_IN_PROGRESS,
		Start: &rentalpb.LocationStatus{
			PoiName: "start_poi",
		},
	})
	if err != nil {
		t.Fatalf("cannot create trip: %v", err)
	}
	if tr.UpdatedAt != 10000 {
		t.Fatalf("wrong updatedat; want: 10000, got: %d", tr.UpdatedAt)
	}

	update := &rentalpb.Trip{
		AccountId: aid.String(),
		Status:    rentalpb.TripStatus_IN_PROGRESS,
		Start: &rentalpb.LocationStatus{
			PoiName: "start_poi_updated",
		},
	}
	cases := []struct {
		name          string
		now           int64
		withUpdatedAt int64
		wantErr       bool
	}{
		{
			name:          "normal_update",
			now:           20000,
			withUpdatedAt: 10000,
		},
		{
			name:          "update_with_stale_timestamp",
			now:           30000,
			withUpdatedAt: 10000,
			wantErr:       true,
		},
		{
			name:          "update_with_refetch",
			now:           40000,
			withUpdatedAt: 20000,
		},
	}

	for _, cc := range cases {
		now = cc.now
		err := m.UpdateTrip(c, tid, aid, cc.withUpdatedAt, update)
		if cc.wantErr {
			if err == nil {
				t.Errorf("%s: want error; got none", cc.name)
			} else {
				continue
			}
		} else {
			if err != nil {
				t.Errorf("%s: cannot update: %v", cc.name, err)
			}
		}
		updatedTrip, err := m.GetTrip(c, tid, aid)
		if err != nil {
			t.Errorf("%s: cannot get trip after update: %v", cc.name, err)
		}
		if cc.now != updatedTrip.UpdatedAt {
			t.Errorf("%s: incorrect updatedat: want %d, got %d", cc.name, cc.now, updatedTrip.UpdatedAt)
		}
	}
}

func TestMain(m *testing.M) {
	os.Exit(mongotesting.RunWithMongoInDocker(m))
	// os.Exit(mongotesting.RunWithMongoInDocker(m, &mongoURI))
}
