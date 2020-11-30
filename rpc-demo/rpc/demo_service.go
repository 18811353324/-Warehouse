package rpc

import (
	"context"
	"git.wanpinghui.com/WPH/go_common/wph"
	"git.wanpinghui.com/WPH/go_common/wph/date"
	"git.wanpinghui.com/WPH/go_common/wph/micro"
)

type DemoRQ struct {
	ID wph.Long `json:"id"`
}

type DemoRES struct {
	ID        wph.Long  `json:"id"`
	StartTime date.Time `json:"startTime"`
}

// RPC服务定义结构体
type DemoService struct{}

// 注册RPC服务
func RegisterDemoService(server *micro.RPCServer) {
	s := &DemoService{}
	err := server.Register(s, "")
	if err != nil {
		println("Failed register rpc service")
	}
}

func (s *DemoService) TestFn(ctx context.Context, args *DemoRQ, reply *DemoRES) error {
	if nil == args || args.ID == 0 {
		return nil
	}
	*reply = DemoRES{ID: args.ID, StartTime: date.CurrentTime()}
	return nil
}
