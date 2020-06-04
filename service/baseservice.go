package service

//
import (
	"../base"
	"../middle"
	"context"
	"github.com/wonderivan/logger"
)

//
type BaseService struct {
	//
	roomMiddler *middle.MdRoom
	//
	Name string
}

//
func (s *BaseService) GetName() string {
	return s.Name
}

//
func (s *BaseService) GetMiddleFromCtx(ctx context.Context, name string) base.Middler {
	value := ctx.Value("mdmgr")
	if value == nil {
		return nil
	}
	mgr := value.(base.MiddleMger)
	return mgr.GetMiddle(name)
}

//
func (s *BaseService) GetRoomMiddler(ctx context.Context) *middle.MdRoom {
	if s.roomMiddler == nil {
		//
		md := s.GetMiddleFromCtx(ctx, "MdRoom")
		if md == nil {
			logger.Error("room middleware not setup")
			return nil
		}

		s.roomMiddler = md.(*middle.MdRoom)
	}
	return s.roomMiddler
}
