// Copyright 2009 The Das Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// Package transport implements net layers.
package transport

//
import (
	"../../base"
	"../../utils"
	"github.com/gorilla/websocket"
	"github.com/wonderivan/logger"
	"net/http"
	"sync"
	"sync/atomic"
	"time"
)

//
type tsWs struct {
	//
	server *http.Server
	//
	msgType int
	//
	pattern string
	//
	name string
	//
	requestHandle base.RequestHandler
	//
	responseHandle base.ResponseHandler
	//
	clients sync.Map
	//
	addr string
	//
	pctl base.Protocoler
}

//
func (ws *tsWs) Start(handle base.RequestHandler, rhandle base.ResponseHandler, pctl base.Protocoler) {
	ws.requestHandle = handle
	ws.pctl = pctl
	ws.responseHandle = rhandle
	http.HandleFunc(ws.pattern, ws.handleClient)
	//
	server := &http.Server{Addr: ws.addr, Handler: nil}
	ws.server = server
	logger.Debug("ws server start at ", ws.addr)
	err := server.ListenAndServe()
	if err != nil {
		logger.Error("ListenAndServe: ", err)
		return
	}
}

//
func (ws *tsWs) Stop() {
	//
	logger.Debug("transport ws stop", ws.clients)
	_ = ws.server.Shutdown(nil)
	//close all clients routines
	ws.clients.Range(func(key interface{}, c interface{}) bool {
		logger.Debug("closing ", c)
		c.(*cClient).close()
		return true
	})
}

//
func (ws *tsWs) Init(cfg map[string]interface{}) {
	ws.pattern = (cfg["pattern"]).(string)
	t := (cfg["msgType"]).(float64)
	ws.msgType = int(t)
	ws.name = ws.pattern
	ws.addr = (cfg["addr"]).(string)
}

//
func (ws *tsWs) Name() string {
	return ws.name
}

//
func (ws *tsWs) handleClient(w http.ResponseWriter, r *http.Request) {
	//
	cid := utils.GetGlobalId()
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		logger.Error("upgrade error", err)
		return
	}
	//
	client := newClient(conn, cid)
	ws.AddClient(client)
	go client.write(ws)
	go client.read(ws)
}

//
func (ws *tsWs) onMessage(c *cClient, message []byte) {
	//
	cmd, _, err := ws.pctl.UnPack(message)
	if err != nil {
		return
	}
	//
	_, _ = ws.requestHandle.HandleRequest(nil, cmd, c)
}

//
func (ws *tsWs) AddClient(c *cClient) {
	logger.Debug("AddClient", c)
	ws.clients.Store(c.id, c)
	c.mgr = ws
	c.param.Store(make(map[string]interface{}, 0))
	_, _ = ws.requestHandle.HandleRequest(nil, &base.Connected{}, c)
}

//
func (ws *tsWs) DelClient(c *cClient) {
	utils.ReStoreGlobalId(c.id)
	ws.clients.Delete(c.id)
	_, _ = ws.requestHandle.HandleRequest(nil, &base.DisConnected{}, c)
}

//
func (ws *tsWs) GetClient(id uint64) *cClient {
	client, ok := ws.clients.Load(id)
	if ok {
		return client.(*cClient)
	}
	return nil
}

//
func (ws *tsWs) Broadcast(cmd base.Cmd, ids []uint64) {
	//
	data, err := ws.pctl.Pack(cmd)
	if err != nil {
		logger.Error("write data fail", data)
		return
	}
	if ids == nil {
		ws.clients.Range(func(key interface{}, c interface{}) bool {
			c.(*cClient).send(data)
			return true
		})
	}
	for _, id := range ids {
		client := ws.GetClient(id)
		if client != nil {
			//
			if cmd == nil {
				client.close()
				return
			}
			client.send(data)
		}
	}
}

//
const (
	// Time allowed to write a message to the peer.
	writeWait = 10 * time.Second
	// Time allowed to read the next pong message from the peer.
	pongWait = 60 * time.Second
	// send pings to peer with this period. Must be less than pongWait.
	pingPeriod = (pongWait * 9) / 10
	// Maximum message size allowed from peer.
	maxMessageSize = 512
)

//
var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
	//EnableCompression: true,
	CheckOrigin: func(r *http.Request) bool { //allow cross origin
		return true
	},
}

//
type cClient struct {
	//
	id uint64
	// The websocket connection.
	socket *websocket.Conn
	// Buffered channel of outbound messages.
	sendBuffer chan []byte
	//
	mgr *tsWs
	//
	param atomic.Value
}

//
func newClient(conn *websocket.Conn, id uint64) *cClient {
	c := &cClient{
		socket:     conn,
		sendBuffer: make(chan []byte, 6),
		param:      atomic.Value{},
		id:         id,
	}
	return c
}

//
func (c *cClient) send(data []byte) {
	c.sendBuffer <- data
}

//
func (c *cClient) close() {
	logger.Debug("client close", c)
	close(c.sendBuffer)
	if c.socket != nil {
		_ = c.socket.Close()
		c.socket = nil
	}
}

//
func (c *cClient) write(ws *tsWs) {
	//for the ping pong loop,sendBuffer ping here,because browser not
	//support sendBuffer ping
	ticker := time.NewTicker(pingPeriod)
	defer func() {
		logger.Debug("socket write routine closed", c)
		ticker.Stop()
	}()

	for {
		select {
		case message, ok := <-c.sendBuffer:
			_ = c.socket.SetWriteDeadline(time.Now().Add(writeWait))
			if !ok {
				_ = c.socket.WriteMessage(websocket.CloseMessage, []byte{})
				logger.Debug("socket sendBuffer chan closed")
				return
			}
			//websocket.TextMessage
			if err := c.socket.WriteMessage(ws.msgType, message); err != nil {
				logger.Debug("socket sendBuffer msg fail", err)
				return
			}
		case <-ticker.C:
			_ = c.socket.SetWriteDeadline(time.Now().Add(writeWait))
			if err := c.socket.WriteMessage(websocket.PingMessage, nil); err != nil {
				logger.Debug("socket pingmessage fail", err)
				return
			}
		}
	}

}

//
func (c *cClient) read(ws *tsWs) {
	//
	c.socket.SetReadLimit(maxMessageSize)
	_ = c.socket.SetReadDeadline(time.Now().Add(pongWait))
	//receive pong here to check alive connect
	c.socket.SetPongHandler(func(string) error {
		//update read wait time to skip pong time
		_ = c.socket.SetReadDeadline(time.Now().Add(pongWait))
		return nil
	})

	defer func() {
		logger.Debug("socket read routine closed")
	}()

	for {
		_, message, err := c.socket.ReadMessage()
		if err != nil {
			if websocket.IsUnexpectedCloseError(err, websocket.CloseGoingAway, websocket.CloseAbnormalClosure) {
				logger.Debug("socket read routine err", err)
			}
			ws.DelClient(c)
			return
		}
		//
		ws.onMessage(c, message)
	}
}

//
func (c *cClient) Write(cmd base.Cmd) {
	//
	if cmd == nil {
		c.close()
		return
	}
	//do the middle
	if c.mgr.responseHandle != nil {
		cmd = c.mgr.responseHandle.HandleResponse(cmd, c)
	}
	//pack
	data, err := c.mgr.pctl.Pack(cmd)
	if err == nil {
		c.send(data)
	}
}

//
func (c *cClient) WriteTo(cmd base.Cmd, whos ...uint64) {
	//do the middle
	if c.mgr.responseHandle != nil {
		cmd = c.mgr.responseHandle.HandleResponse(cmd, c)
	}
	//pack
	if len(whos) > 1 {
		c.mgr.Broadcast(cmd, whos)
	} else {
		//all
		who := whos[0]
		if who == 0 {
			c.mgr.Broadcast(cmd, nil)
		} else if who > 0 {
			//single one
			data, err := c.mgr.pctl.Pack(cmd)
			if err != nil {
				logger.Error("write to ", whos)
				return
			}
			other := c.mgr.GetClient(who)
			if other != nil {
				//
				if cmd == nil {
					other.close()
					return
				}
				other.send(data)
			}
		}
	}
}

//
func (c *cClient) GetSession(key string) (interface{}, bool) {
	val := c.param.Load().(base.Str2Any)
	v, ok := val[key]
	return v, ok
}

//
func (c *cClient) SetSession(key string, v interface{}) {
	val := c.param.Load().(base.Str2Any)
	val[key] = v
}

//
func (c *cClient) Id() uint64 {
	return c.id
}

//
func init() {
	logger.Debug("ws transport init")
	_, _ = base.RegisterClass(&tsWs{})
}
