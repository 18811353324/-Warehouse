package main

import (
	"context"
	"flag"
	"fmt"
	"git.wanpinghui.com/WPH/go_common/wph"
	mw "git.wanpinghui.com/WPH/go_common/wph/http_middleware"
	"git.wanpinghui.com/WPH/go_common/wph/logger"
	"git.wanpinghui.com/WPH/go_common/wph/micro"
	"git.wanpinghui.com/WPH/rpc-demo/controller"
	"git.wanpinghui.com/WPH/rpc-demo/rpc"
	"github.com/gin-contrib/gzip"
	"github.com/gin-gonic/gin"
	"github.com/spf13/viper"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"
)

// 全局常量定义（作用域整个package）
var (
	hostPrefix = "/rpc" // 部署到线上时添加的子系统路径前缀
)

var configPath string

func init() {
	flag.StringVar(&configPath, "c", "config", "config file path command line param")
}

// 程序启动入口
// 日常编码，只需要关心"注册路由"部分即可
func main() {
	flag.Parse() // 解析命令行参数定义的flag
	// 加载配置文件
	currentPath, err := wph.GetCurrentDir()
	if err != nil {
		panic(fmt.Errorf("Get current directory failed: %s \n", err))
	}
	viper.SetConfigName(configPath)
	if len(currentPath) > 0 {
		viper.AddConfigPath(currentPath)
	}
	viper.AddConfigPath("./")
	viper.SetDefault("service.port", ":8087")
	err = viper.ReadInConfig()
	if err != nil {
		panic(fmt.Errorf("Fatal error config file: %s \n", err))
	}
	LOG := logger.New()

	router := gin.New()

	// 添加 GIN 请求处理中间件
	router.Use(
		gzip.Gzip(gzip.DefaultCompression),
		mw.AddCrossOriginHeaders(),                      // 添加跨域头
		mw.HandleOptionsMethod(),                        // Options和Head请求处理
		mw.OpenTracer(),                                 // 服务调用链路追踪
		mw.Logger(hostPrefix, LOG, "/test/imgood", "/"), // 统一日志记录
		mw.Recovery(hostPrefix, LOG),                    // 统一错误处理
		// jws.NeedLogin(),               // 全局JWT权限校验
	)

	// 设置404、405响应
	router.NoMethod(func(c *gin.Context) {
		resp := wph.E(-http.StatusMethodNotAllowed, "Method Not Allowed", nil)
		resp.HttpStatus = http.StatusMethodNotAllowed
		panic(resp)
	})
	router.NoRoute(func(c *gin.Context) {
		resp := wph.E(-http.StatusMethodNotAllowed, "Endpoint Not Found", nil)
		resp.HttpStatus = http.StatusMethodNotAllowed
		panic(resp)
	})

	// 注册路由
	controller.RegisterTestController(router)

	// 新建 RPCX 对象实例，NewRPCServer返回包装好的 *micro.RPCServer类型
	rpcServer := micro.NewRPCServer(micro.SetRPCServerLogger(LOG))

	// 注册RPC服务
	rpc.RegisterDemoService(rpcServer)

	// 开始监听RPC请求
	go func(rpcServer *micro.RPCServer, log *logger.Logger) {
		if err = rpcServer.Serve("tcp", ":10000"); err != nil {
			log.Warn("Failed serve rpc", err.Error())
		}
	}(rpcServer, LOG)

	// 开始监听HTTP请求
	addr := viper.GetString("service.port")
	// router.Run(addr)
	srv := &http.Server{Addr: addr, Handler: &mw.PrefixCut{Handler: router, HostPrefix: hostPrefix}}
	go func(httpServer *http.Server, log *logger.Logger) {
		if err := httpServer.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Warn(err)
		}
	}(srv, LOG)

	signalHandler(rpcServer, srv)
}

func signalHandler(rpcServer *micro.RPCServer, httpServer *http.Server) {
	var (
		c chan os.Signal
		s os.Signal
	)
	c = make(chan os.Signal, 1)
	signal.Notify(c, syscall.SIGHUP, syscall.SIGQUIT, syscall.SIGTERM, syscall.SIGINT)
	for {
		s = <-c
		switch s {
		case syscall.SIGQUIT, syscall.SIGTERM, syscall.SIGINT:
			println("Shutdown Server ...")
			ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
			//if err := rpcServer.Shutdown(ctx); err != nil {
			//	LOG.Fatal("RPC Server Shutdown:", err)
			//}
			if err := httpServer.Shutdown(ctx); err != nil {
				println("HTTP Server Shutdown:", err)
			}
			cancel()
			println("Server Exited")
			return
		default:
			return
		}
	}
}
