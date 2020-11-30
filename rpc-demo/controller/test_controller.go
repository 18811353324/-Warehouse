package controller

import (
	"fmt"
	"git.wanpinghui.com/WPH/go_common/wph"
	"git.wanpinghui.com/WPH/go_common/wph/micro"
	"git.wanpinghui.com/WPH/rpc-demo/rpc"
	"github.com/gin-gonic/gin"
	"net/http"
)

type TestController struct{}

func RegisterTestController(router *gin.Engine) {
	c := &TestController{}
	test := router.Group("/test")
	{
		test.GET("/imgood", c.imGood)
		test.GET("/rpc_call1", c.rpcCall1)
		test.GET("/rpc_call2", c.rpcCall2)
	}
}

// 节点健康检查
func (c *TestController) imGood(ctxt *gin.Context) {
	ctxt.JSON(http.StatusOK, wph.R(nil))
}

func (c *TestController) rpcCall1(ctxt *gin.Context) {
	c.rpcCall(ctxt, "localhost:10000")
}

func (c *TestController) rpcCall2(ctxt *gin.Context) {
	c.rpcCall(ctxt, "test:10000")
}

func (c *TestController) rpcCall(ctxt *gin.Context, addr string) {
	println("call ", addr)
	cc := micro.NewRPCClient(ctxt, addr, "DemoService")
	defer func(c *micro.RPCClient) {
		_ = c.Close()
	}(cc)

	var result interface{}
	rpcRQ := &rpc.DemoRQ{ID: wph.NextId()}
	rpcRES := &rpc.DemoRES{}
	err := cc.Call("TestFn", rpcRQ, rpcRES)
	if err != nil {
		result = fmt.Sprintf("rpc client call err=%s", err.Error())
		println(result)
		ctxt.JSON(http.StatusBadRequest, wph.E(-1, err.Error(), err.Error()))
		return
	} else {
		result = rpcRES
	}
	ctxt.JSON(http.StatusOK, wph.R(result))
}
