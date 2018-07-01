package main

import (
	"encoding/json"
	"log"
	"sidekick/dataservice_client"

	"github.com/MuggleWei/cascade"
	"github.com/gorilla/websocket"
	"github.com/mitchellh/mapstructure"
)

type ClientCallbackFn func(*cascade.Peer, *ClientReq)

type ClientInfo struct {
	User    string // user name
	Logined bool   // whether already logined
}

type WsHandler struct {
	Hub             *cascade.Hub
	DataService     *dataservice_client.DataServiceClient
	ClientCallbacks map[string]ClientCallbackFn
}

func NewWsHandler() *WsHandler {
	service := &WsHandler{
		Hub:             nil,
		DataService:     nil,
		ClientCallbacks: make(map[string]ClientCallbackFn),
	}

	upgrader := websocket.Upgrader{
		ReadBufferSize:  1024 * 20,
		WriteBufferSize: 1024 * 20,
	}

	hub := cascade.NewHub(service, &upgrader, 10240)
	go hub.Run()

	service.Hub = hub
	service.RegisterClientCallbacks()

	return service
}

// Slot callbacks
func (this *WsHandler) OnActive(peer *cascade.Peer) {
	log.Printf("[Info] OnActive: %v\n", peer.Conn.RemoteAddr().String())
	peer.ExtraInfo = &ClientInfo{
		User:    "",
		Logined: false,
	}
}

func (this *WsHandler) OnInactive(peer *cascade.Peer) {
	log.Printf("[Info] OnInactive: %v\n", peer.Conn.RemoteAddr().String())
}

func (this *WsHandler) OnRead(peer *cascade.Peer, message []byte) {
	var req ClientReq
	err := json.Unmarshal(message, &req)
	if err != nil {
		panic(err)
	}

	callback, ok := this.ClientCallbacks[req.Msg]
	if ok {
		if req.Msg == "login" {
			callback(peer, &req)
		} else {
			clientInfo, _ := peer.ExtraInfo.(*ClientInfo)
			if ok && clientInfo.Logined {
				callback(peer, &req)
			} else {
				log.Printf("[Warning] <%v> %v req message without login\n", peer.Conn.RemoteAddr(), req.Msg)
				RspError(peer, req.Reqid, ErrorCode_PermissionDenied, "request without login")
			}
		}
	} else {
		log.Printf("[Warning] %v message is not handle\n", req.Msg)
		RspError(peer, req.Reqid, ErrorCode_Unhandle, "message without handle function")
	}
}

func (this *WsHandler) OnHubByteMessage(msg *cascade.HubByteMessage) {
	// TODO:
}

func (this *WsHandler) OnHubObjectMessage(msg *cascade.HubObjectMessage) {
}

////////////////////// callback registers //////////////////////
func (this *WsHandler) RegisterClientCallbacks() {
	this.ClientCallbacks["login"] = this.OnLogin
}

////////////////////// client callbacks //////////////////////
func (this *WsHandler) OnLogin(peer *cascade.Peer, req *ClientReq) {
	log.Printf("[Info] OnLogin: %v\n", peer.Conn.RemoteAddr().String())

	var reqLogin ReqLogin
	err := mapstructure.Decode(req.Data, &reqLogin)
	if err != nil {
		panic(err)
	}

	dataReqLogin := dataservice_client.DataReqLogin{
		Name:     reqLogin.User,
		Password: reqLogin.Password,
	}

	loginRsp, err := this.DataService.Login(&dataReqLogin)
	if err != nil {
		rsp := ClientRsp{
			Msg:   "rsplogin",
			Reqid: req.Reqid,
			Err:   NewErrorMsg(ErrorCode_Unknown, err.Error()),
		}
		bytes, _ := json.Marshal(rsp)
		peer.SendChannel <- bytes
	} else {
		rsp := ClientRsp{
			Msg:   "rsplogin",
			Reqid: req.Reqid,
			Err:   NewErrorMsg(loginRsp.Code, loginRsp.Message),
		}
		bytes, _ := json.Marshal(rsp)
		peer.SendChannel <- bytes

		if loginRsp.Success {
			log.Printf("[Info] %v <%v> login success\n", reqLogin.User, peer.Conn.RemoteAddr())
			clientInfo, _ := peer.ExtraInfo.(*ClientInfo)
			clientInfo.User = reqLogin.User
			clientInfo.Logined = true
		}
	}
}
