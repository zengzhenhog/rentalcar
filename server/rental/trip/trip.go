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
	ProfileManage ProfileManage
	Mongo         *dao.Mongo
	Logger        *zap.Logger
}

// ProfileManage defines the ACL (Anti Corruption Layer)
// for profile verification logic
// id.IdentityID防止驾驶者身份在中途修改，存储行程创建时的身份快照
type ProfileManage interface {
	Verify(context.Context, id.AccountID) (id.IdentityID, error)
}

func (s *Service) CreateTrip(c context.Context, req *rentalpb.CreateTripRequest) (*rentalpb.TripEntity, error) {
	aid, err := auth.AccountIDFromContext(c)
	if err != nil {
		return nil, err
	}
	// 验证驾驶者身份，虽然在同个微服务，不同领域使用ACL(防止入侵层)进行隔离
	iID, err := s.ProfileManage.Verify(c, aid)
	if err != nil {
		// err.Error()返回错误的string
		return nil, status.Error(codes.FailedPrecondition, err.Error())
	}

	// 按以下流程操作，除非两个人同时完成操作，很难出错
	// 检查车辆状态
	// 创建行程，写入数据库，开始计费
	// 车辆开锁

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
