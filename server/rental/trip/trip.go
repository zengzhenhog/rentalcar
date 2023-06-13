package trip

import (
	"context"
	rentalpb "coolcar/rental/api/gen/v1"
	"coolcar/rental/trip/dao"
	"coolcar/shared/auth"
	"coolcar/shared/id"

	// "coolcar/shared/auth"

	"go.uber.org/zap"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type Service struct {
	Mongo  *dao.Mongo
	Logger *zap.Logger
}

func (s *Service) CreateTrip(c context.Context, req *rentalpb.CreateTripRequest) (*rentalpb.TripEntity, error) {
	return nil, status.Error(codes.Unimplemented, "")
}

func (s *Service) GetTrip(c context.Context, req *rentalpb.GetTripRequest) (*rentalpb.Trip, error) {
	return nil, status.Error(codes.Unimplemented, "")
}

func (s *Service) GetTipps(c context.Context, req *rentalpb.GetTripsRequest) (*rentalpb.GetTripsResponse, error) {
	return nil, status.Error(codes.Unimplemented, "")
}

func (s *Service) UpdateTrip(c context.Context, req *rentalpb.UpdateTripRequest) (*rentalpb.Trip, error) {
	aid, err := auth.AccountIDFromContext(c)
	if err != nil {
		return nil, status.Error(codes.Unauthenticated, "")
	}
	tid := id.TripID(req.Id)
	tr, err := s.Mongo.GetTrip(c, tid, aid)
	// if err != nil {
	// 	return nil, status.Error(codes.Unauthenticated, "")
	// }
	if req.Current != nil {
		tr.Trip.Current = s.calcCurrentStatus(tr.Trip, req.Current)
	}

	if req.EndTrip {
		tr.Trip.End = tr.Trip.Current
		tr.Trip.Status = rentalpb.TripStatus_FINISHED
	}

	s.Mongo.UpdateTrip(c, tid, aid, tr.UpdatedAt, tr.Trip)
	return nil, status.Error(codes.Unimplemented, "")
}

func (s *Service) calcCurrentStatus(trip *rentalpb.Trip, cur *rentalpb.Location) *rentalpb.LocationStatus {
	return nil
}
