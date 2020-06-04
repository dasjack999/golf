// Copyright 2010 The Das Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.
package rpc

import (
	"github.com/wonderivan/logger"
	"io"
	"net"
	"net/http"
	"net/rpc"
	"net/rpc/jsonrpc"
)

//
type BaseService interface {
}

//
type Response struct {
	Data interface{}
}

//
type Request struct {
	Data interface{}
}

const (
	Gob      int = 0
	Json     int = 1
	HttpJson int = 2
)

//
type RpcClient struct {
	//
	Url string
	//
	Name string
	//
	Type int
	//
	client *rpc.Client
}

//
func (r *RpcClient) ConnectGob() error {
	client, err := rpc.Dial("tcp", r.Url)
	if err != nil {

		return err
	}
	r.client = client

	return nil
}

//
func (r *RpcClient) ConnectJson() error {
	conn, err := net.Dial("tcp", r.Url)
	if err != nil {
		return err
	}
	r.client = rpc.NewClientWithCodec(jsonrpc.NewClientCodec(conn))

	return nil
}

//
func (r *RpcClient) Call(service string, param interface{}) (res interface{}, err error) {
	if r.client == nil {
		return
	}
	tmp := Response{}
	//
	err = r.client.Call(r.Name+service, param, &tmp)
	if err == nil {
		res = tmp.Data
	}
	return
}

//
type RpcServer struct {
	//
	Url string
}

//
func (r *RpcServer) Register(srv BaseService) error {
	return rpc.Register(srv)
}

//
func (r *RpcServer) RunGob() error {
	listener, err := net.Listen("tcp", r.Url)
	if err != nil {
		return err
	}
	//
	logger.Debug("server start at", r.Url)
	for {
		conn, err := listener.Accept()
		if err != nil {
			return err
		}
		go rpc.ServeConn(conn)
	}
}

//
func (r *RpcServer) RunJson() error {
	listener, err := net.Listen("tcp", r.Url)
	if err != nil {
		return err
	}
	for {
		conn, err := listener.Accept()
		if err != nil {
			return err
		}
		go rpc.ServeCodec(jsonrpc.NewServerCodec(conn))
	}
}

//
func (r *RpcServer) RunHttpJson(route string) error {
	http.HandleFunc(route, func(writer http.ResponseWriter, request *http.Request) {
		var conn io.ReadWriteCloser = struct {
			io.Writer
			io.ReadCloser
		}{
			ReadCloser: request.Body,
			Writer:     writer,
		}

		_ = rpc.ServeRequest(jsonrpc.NewServerCodec(conn))
	})
	return http.ListenAndServe(r.Url, nil)
}

//
func NewServer(url string, rtype int, srvList ...BaseService) *RpcServer {
	ser := &RpcServer{Url: url}
	var err error
	for _, srv := range srvList {
		err = ser.Register(srv)
		if err != nil {
			logger.Debug("register service fail", srv, err)
		}
	}
	//
	switch rtype {
	case Gob:
		err = ser.RunGob()
		break
	case Json:
		err = ser.RunJson()
		break
	case HttpJson:
		err = ser.RunHttpJson("/service")
	}
	if err != nil {
		logger.Error("run fail", err)
		return nil
	}
	return ser
}

//========================do in server==============================
//type MyService struct {
//	BaseService
//}
////
//func (s *MyService)Hello(word string,res *Response) error  {
//	res.Data = "hi,"+word
//	return nil
//}
//
//func main()  {
//	server :=NewServer(":8095",Gob,&MyService,&MyService2)
//
//}
//==========================do in client================================
//
//type MyServiceClient struct {
//	RpcClient
//}
//
////
//func (m *MyServiceClient)Hello(word string)(res string ,err error)  {
//	var ress interface{}
//	ress,err = m.Call("Hello",word)
//	res = ress.(string)
//	return
//}
//
//func main() {
//
//	client :=&MyServiceClient{}
//	client.Name="MyService.";
//	client.Url=":8095"
//	client.ConnectGob()
//
//	res,err :=client.Hello("jack")
//	logger.Debug("call ",res,err)
//
//}
