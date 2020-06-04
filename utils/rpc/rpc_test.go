package rpc

import (
	"github.com/wonderivan/logger"
	"testing"
)

//
type MyServiceClient struct {
	RpcClient
}

//
func (m *MyServiceClient) Hello(word string) (res string, err error) {
	var r interface{}
	r, err = m.Call("Hello", word)
	res = r.(string)
	return
}

func TestRpcClient_ConnectGob(t *testing.T) {

	client := &MyServiceClient{RpcClient{Url: ":8095"}}

	res, err := client.Hello("jack")
	logger.Debug("call ", res, err)

}
