package trip

import (
	"context"
	rentalpb "coolcar/rental/api/gen/v1"
	poi "coolcar/rental/trip/client"
	"coolcar/rental/trip/dao"
	"coolcar/shared/auth"
	"coolcar/shared/id"
	mgutil "coolcar/shared/mongo"
	mongotesting "coolcar/shared/mongo/testing"
	"coolcar/shared/server"
	"encoding/json"
	"fmt"
	"os"
	"testing"
)

func TestCreateTrip(t *testing.T) {
	c := auth.ContextWithAccountID(context.Background(), id.AccountID("account1"))
	mc, err := mongotesting.NewClient(c)
	if err != nil {
		t.Fatalf("cannot create mongo client: %v: ", err)
	}

	logger, err := server.NewZapLogger()
	if err != nil {
		t.Fatalf("cannot create logger: %v", err)
	}

	pm := &profileManager{}
	cm := &carManager{}
	s := &Service{
		ProfileManager: pm,
		CarManager:     cm,
		POIManager:     &poi.Manager{},
		Mongo:          dao.NewMongo(mc.Database("coolcar")),
		Logger:         logger,
	}

	req := &rentalpb.CreateTripRequest{
		CarId: "car1",
		Start: &rentalpb.Location{
			Latitude:  32.13,
			Longitude: 113.27,
		},
	}
	pm.iID = "identity1"
	golden := `{"account_id":"account1","car_id":"car1","start":{"location":{"latitude":32.13,"longitude":113.27},"poi_name":"天安门"},"current":{"location":{"latitude":32.13,"longitude":113.27},"poi_name":"天安门"},"status":1,"identity_id":"identity1"}`
	cases := []struct {
		name         string
		tripID       string
		profileErr   error
		carVerifyErr error
		carUnlockErr error
		want         string
		wantErr      bool
	}{
		{
			name:   "normal_create",
			tripID: "6485d7452cf384b4e217ae91",
			want:   golden,
		},
		{
			name:       "profile_err",
			tripID:     "6485d7452cf384b4e217ae92",
			profileErr: fmt.Errorf("profile"),
			wantErr:    true,
		},
		{
			name:         "car_verify_err",
			tripID:       "6485d7452cf384b4e217ae93",
			carVerifyErr: fmt.Errorf("verify"),
			wantErr:      true,
		},
		{
			name:         "car_unlock_err",
			tripID:       "6485d7452cf384b4e217ae94",
			carUnlockErr: fmt.Errorf("unlock"),
			want:         golden,
		},
	}

	for _, cc := range cases {
		t.Run(cc.name, func(t *testing.T) {
			mgutil.NewObjectIDWithValue(id.TripID(cc.tripID))
			pm.err = cc.profileErr
			cm.UnlockErr = cc.carUnlockErr
			cm.VerifyErr = cc.carVerifyErr
			res, err := s.CreateTrip(c, req)
			if cc.wantErr {
				if err == nil {
					t.Errorf("want error; got none")
				} else {
					return
				}
			}
			if err != nil {
				t.Errorf("error creating trip: %v", err)
				return
			}
			if res.Id != cc.tripID {
				t.Errorf("incorrect id; want %q, got %q", cc.tripID, res.Id)
			}
			b, err := json.Marshal(res.Trip)
			if err != nil {
				t.Errorf("cannot marshall response: %v", err)
			}
			got := string(b)
			if cc.want != got {
				t.Errorf("incorrect response: want %s, got %s", cc.want, got)
			}
		})
	}

}

type profileManager struct {
	iID id.IdentityID
	err error
}

func (p *profileManager) Verify(context.Context, id.AccountID) (id.IdentityID, error) {
	return p.iID, p.err
}

type carManager struct {
	VerifyErr error
	UnlockErr error
}

func (c *carManager) Verify(context.Context, id.CarID, *rentalpb.Location) error {
	return c.VerifyErr
}
func (c *carManager) Unlock(context.Context, id.CarID) error {
	return c.UnlockErr
}

func TestMain(m *testing.M) {
	os.Exit(mongotesting.RunWithMongoInDocker(m))
}
