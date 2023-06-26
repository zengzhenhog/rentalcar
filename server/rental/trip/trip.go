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
	ProfileManager ProfileManager
	CarManager     CarManager
	POIManager     POIManager
	Mongo          *dao.Mongo
	Logger         *zap.Logger
}

// ProfileManager defines the ACL (Anti Corruption Layer)
// for profile verification logic
// id.IdentityID防止驾驶者身份在中途修改，存储行程创建时的身份快照
type ProfileManager interface {
	Verify(context.Context, id.AccountID) (id.IdentityID, error)
}

type CarManager interface {
	Verify(context.Context, id.CarID, *rentalpb.Location) error
	Unlock(context.Context, id.CarID) error
}

type POIManager interface {
	Resolve(context.Context, *rentalpb.Location) (string, error)
}

func (s *Service) CreateTrip(c context.Context, req *rentalpb.CreateTripRequest) (*rentalpb.TripEntity, error) {
	aid, err := auth.AccountIDFromContext(c)
	if err != nil {
		return nil, err
	}
	// 验证驾驶者身份，虽然在同个微服务，不同领域使用ACL(防止入侵层)进行隔离
	iID, err := s.ProfileManager.Verify(c, aid)
	if err != nil {
		// err.Error()返回错误的string
		return nil, status.Error(codes.FailedPrecondition, err.Error())
	}

	// 按以下流程操作，除非两个人同时完成操作，很难出错

	// 检查车辆状态
	carID := id.CarID(req.CarId)
	err = s.CarManager.Verify(c, carID, req.Start)
	if err != nil {
		return nil, status.Error(codes.FailedPrecondition, err.Error())
	}

	poi, err := s.POIManager.Resolve(c, req.Start)
	if err != nil {
		s.Logger.Info("cannot resolve poi", zap.Stringer("location", req.Start), zap.Error(err)) // zap.Stringer指参数带有string方法req.Start.string()
	}

	// 创建行程，写入数据库，开始计费
	ls := &rentalpb.LocationStatus{
		Location: req.Start,
		PoiName:  poi,
	}
	tr, err := s.Mongo.CreateTrip(c, &rentalpb.Trip{
		AccountId:  aid.String(),
		CarId:      carID.String(),
		IdentityId: iID.String(),
		Status:     rentalpb.TripStatus_IN_PROGRESS,
		Start:      ls,
		Current:    ls,
	})
	if err != nil {
		s.Logger.Warn("cannot create trip", zap.Error(err))
		// 行程创建失败最可能的原因是形成已经存在，使用code.Code_ALREADY_EXISTS
		return nil, status.Error(codes.AlreadyExists, "")
	}

	// 车辆开锁，开锁是后台任务，不包括在当前请求，先返回行程创建成功
	go func() {
		err := s.CarManager.Unlock(context.Background(), carID)
		if err != nil {
			s.Logger.Error("cannot unlock car", zap.Error(err))
		}
	}()

	return &rentalpb.TripEntity{
		Id:   tr.ID.Hex(),
		Trip: tr.Trip,
	}, nil
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
