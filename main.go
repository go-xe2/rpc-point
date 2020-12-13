/*****************************************************************
* Copyright©,2020-2022, email: 279197148@qq.com
* Version: 1.0.0
* @Author: yangtxiang
* @Date: 2020-08-19 16:32
* Description:
*****************************************************************/

package main

import (
	"context"
	"fmt"
	"github.com/go-xe2/x/os/xfile"
	"github.com/go-xe2/x/os/xlog"
	"github.com/go-xe2/x/utils/xconfig"
	"github.com/go-xe2/xthrift/gateway"
	"github.com/go-xe2/xthrift/rpcPoint"
	"log"
	"net"
	"os"
)


var server *rpcPoint.TEndPointServer
//
//func appInit(ctx context.Context) ([]net.Listener, error)  {
//	options, ok := ctx.Value("options").(*rpcPoint.TOptions)
//	if !ok {
//		return nil, errors.New("参数options为nil")
//	}
//	httpLst, err := net.Listen("tcp", options.HttpAddr)
//	if err != nil {
//		return nil, err
//	}
//	thriftLst, err := net.Listen("tcp", options.ThriftAddr)
//	if err != nil {
//		return nil, err
//	}
//	return []net.Listener{httpLst, thriftLst}, nil
//}
//
//func appRun(ctx context.Context, listeners []net.Listener) error  {
//	options := ctx.Value("options").(*rpcPoint.TOptions)
//	server = rpcPoint.NewEndPointServer(options)
//	server.SetHttpListener(listeners[0])
//	server.SetThriftListener(listeners[1])
//	// 创建网关
//	gatewayHandle := gateway.NewSvrHttpHandler(server.GetPdlQuery(), server, options.BaseRouter, options.BaseNamespace)
//	server.PathPrefixHandle("", options.BaseRouter, gatewayHandle, "服务调用", "调用已注册的内部服务")
//	xlog.Info("协议存放目录:", options.PDLPath)
//	xlog.Info("服务器地址存放目录:", options.HostPath)
//	return server.Serve()
//}
//
//func appStop(ctx context.Context) {
//	if server != nil {
//		if err := server.Stop(); err != nil {
//			fmt.Println("stop error:", err)
//		}
//	}
//}


func main() {
	//loader := hotLoad.NewHotLoadService("rpc-point.pid")
	c := xconfig.Config()
	c.SetFileName("config.yaml")

	mp := c.GetMap("options")

	options, err := rpcPoint.NewOptionsFromMap(mp)
	if err != nil {
		xlog.Info(err)
		return
	}
	if !xfile.Exists("./logs") {
		_ = xfile.Mkdir("./logs")
	}
	logFile, logErr := os.OpenFile(xfile.Join("./logs", "rpc-point.log"), os.O_CREATE|os.O_RDWR|os.O_APPEND, 0666)
	if logErr != nil {
		fmt.Println("设置日志输出文件出错:", logErr)
	} else {
		defer logFile.Close()
		log.SetOutput(logFile)
	}
	if err := xlog.SetPath("./logs"); err != nil {
		fmt.Println("设置日志输出目录出错:", err)
	}
	cxt := context.Background()

	cxt = context.WithValue(cxt, "options", options)

	httpLst, err := net.Listen("tcp", options.HttpAddr)
	if err != nil {
		panic(err)
	}
	thriftLst, err := net.Listen("tcp", options.ThriftAddr)
	if err != nil {
		panic(err)
	}

	server = rpcPoint.NewEndPointServer(options)
	server.SetHttpListener(httpLst)
	server.SetThriftListener(thriftLst)

	// 创建网关
	gatewayHandle := gateway.NewSvrHttpHandler(server.GetPdlQuery(), server, options.BaseRouter, options.BaseNamespace)
	server.PathPrefixHandle("", options.BaseRouter, gatewayHandle, "服务调用", "调用已注册的内部服务")
	xlog.Info("协议存放目录:", options.PDLPath)
	xlog.Info("服务器地址存放目录:", options.HostPath)

	if err := server.Serve(); err != nil {
		xlog.Info(err)
	}
	select {
	}
	//
	//
	//err = loader.Load(cxt, appInit, appRun, appStop)
	//if err != nil {
	//	fmt.Println("error:", err)
	//}
}