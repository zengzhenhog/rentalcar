package trip

import (
	"context"
	trippb "coolcar/proto/gen/go"
)

// type TripServiceServer interface {
// 	GetTrip(context.Context, *GetTripRequest) (*GetTripRespons, error)
// 	mustEmbedUnimplementedTripServiceServer()
// }

type Service struct{}

func (*Service) GetTrip(c context.Context, req *trippb.GetTripRequest) (*trippb.GetTripRespons, error) {
	return &trippb.GetTripRespons{
		Id: req.Id,
		Trip: &trippb.Trip{
			Start:       "abc",
			End:         "def",
			DurationSec: 3000,
			FeeCent:     300000,
			StartPos: &trippb.Location{
				Latitude:  30,
				Longitude: 120,
			},
			EndPos: &trippb.Location{
				Latitude:  35,
				Longitude: 115,
			},
			PathLocations: []*trippb.Location{
				{
					Latitude:  31,
					Longitude: 119,
				},
				{
					Latitude:  32,
					Longitude: 118,
				},
			},
			Status: trippb.TripStatus_FINISHED,
		},
	}, nil
}

func (*Service) mustEmbedUnimplementedTripServiceServer() {}
