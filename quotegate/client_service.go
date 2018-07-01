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
type ObjectCallbackFn func(*cascade.Peer, interface{})
type ServerCallbackFn func(*ServerRsp)

type ClientInfo struct {
	User    string // user name
	Logined bool   // whether already logined
}

type SubObj struct {
	Reqid  string
	Uuid   string
	ReqSub ReqSub
}

type ClientService struct {
	Hub *cascade.Hub

	DataService   *dataservice_client.DataServiceClient
	SubManager    *SubManager
	QuoteServices map[string]*cascade.Hub

	ClientCallbacks map[string]ClientCallbackFn
	ObjectCallbacks map[string]ObjectCallbackFn
	ServerCallbacks map[string]ServerCallbackFn
}

func NewClientService() *ClientService {
	service := &ClientService{
		Hub: nil,

		DataService:   nil,
		SubManager:    nil,
		QuoteServices: make(map[string]*cascade.Hub),

		ClientCallbacks: make(map[string]ClientCallbackFn),
		ObjectCallbacks: make(map[string]ObjectCallbackFn),
		ServerCallbacks: make(map[string]ServerCallbackFn),
	}

	upgrader := websocket.Upgrader{
		ReadBufferSize:  1024 * 20,
		WriteBufferSize: 1024 * 20,
	}

	hub := cascade.NewHub(service, &upgrader, 10240)
	go hub.Run()

	service.Hub = hub
	service.RegisterClientCallbacks()
	service.RegisterObjectCallbacks()
	service.RegisterServerMsgCallbacks()

	return service
}

// Slot callbacks
func (this *ClientService) OnActive(peer *cascade.Peer) {
	log.Printf("[Info] OnActive: %v\n", peer.Conn.RemoteAddr().String())
	peer.ExtraInfo = &ClientInfo{
		User:    "",
		Logined: false,
	}
}

func (this *ClientService) OnInactive(peer *cascade.Peer) {
	log.Printf("[Info] OnInactive: %v\n", peer.Conn.RemoteAddr().String())
	this.SubManager.PeerUnsub(peer)
}

func (this *ClientService) OnRead(peer *cascade.Peer, message []byte) {
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

func (this *ClientService) OnHubByteMessage(msg *cascade.HubByteMessage) {
	var msgIdentify MessageIdentify
	err := json.Unmarshal(msg.Message, &msgIdentify)
	if err != nil {
		panic(err)
	}

	if msgIdentify.Msg == "quote" {
		var quote Quote
		err := json.Unmarshal(msg.Message, &quote)
		if err != nil {
			panic(err)
		}
		this.OnQuote(&quote)
	} else {
		var rsp ServerRsp
		err := json.Unmarshal(msg.Message, &rsp)
		if err != nil {
			panic(err)
		}

		callback, ok := this.ServerCallbacks[rsp.Msg]
		if ok {
			callback(&rsp)
		} else {
			log.Printf("[Warning] %v message is not handle\n", rsp.Msg)
		}
	}
}

func (this *ClientService) OnHubObjectMessage(msg *cascade.HubObjectMessage) {
	callback, ok := this.ObjectCallbacks[msg.ObjectName]
	if ok {
		callback(msg.Peer, msg.ObjectPtr)
	} else {
		log.Printf("[Error] client manager: %v object message is not handle\n", msg.ObjectName)
	}
}

////////////////////// callback registers //////////////////////
func (this *ClientService) RegisterClientCallbacks() {
	this.ClientCallbacks["login"] = this.OnLogin
	this.ClientCallbacks["sub"] = this.OnSub
	this.ClientCallbacks["unsub"] = this.OnUnsub
}

func (this *ClientService) RegisterObjectCallbacks() {
	this.ObjectCallbacks["subobj"] = func(peer *cascade.Peer, pointer interface{}) {
		this.OnSubObj(peer, pointer.(*SubObj))
	}
}

func (this *ClientService) RegisterServerMsgCallbacks() {
	this.ServerCallbacks["rspsub"] = this.OnRspSub
}

////////////////////// client callbacks //////////////////////
func (this *ClientService) OnLogin(peer *cascade.Peer, req *ClientReq) {
	log.Printf("[Info] OnLogin: %v\n", peer.Conn.RemoteAddr().String())

	var reqLogin ReqLogin
	err := mapstructure.Decode(req.Data, &reqLogin)
	if err != nil {
		panic(err)
	}

	dataReqLogin := dataservice_client.DataReqLogin{
		Name:     reqLogin.User,
		Password: reqLogin.Passwd,
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
func (this *ClientService) OnSub(peer *cascade.Peer, req *ClientReq) {
	clientInfo, _ := peer.ExtraInfo.(*ClientInfo)
	log.Printf("[Info] OnSub: %+v <%v> sub %+v\n", clientInfo.User, peer.Conn.RemoteAddr().String(), *req)

	var subObj SubObj
	err := mapstructure.Decode(req.Data, &subObj.ReqSub)
	if err != nil {
		panic(err)
	}

	subObj.Reqid = req.Reqid

	// default value
	if subObj.ReqSub.Table == "orderbook" {
		if subObj.ReqSub.Optional.Depth == 0 {
			subObj.ReqSub.Optional.Depth = 5
		}
	} else if subObj.ReqSub.Table == "kline" {
		if subObj.ReqSub.Optional.Period == "" {
			subObj.ReqSub.Optional.Period = "1m"
		}
	}

	obj := &cascade.HubObjectMessage{
		Peer:       peer,
		ObjectName: "subobj",
		ObjectPtr:  &subObj,
	}
	this.Hub.ObjectMessageChannel <- obj
}
func (this *ClientService) OnUnsub(peer *cascade.Peer, req *ClientReq) {
	clientInfo, _ := peer.ExtraInfo.(*ClientInfo)
	log.Printf("[Info] OnUnsub: %+v <%v> sub %+v\n", clientInfo.User, peer.Conn.RemoteAddr().String(), *req)

	var reqUnsub ReqUnsub
	err := mapstructure.Decode(req.Data, &reqUnsub)
	if err != nil {
		panic(err)
	}

	ok := this.SubManager.Unsub(peer, &reqUnsub)
	var errMsg ErrorMsg
	if ok {
		errMsg = NewErrorMsg(0, "")
	} else {
		errMsg = NewErrorMsg(ErrorCode_Unknown, "failed to unsub topic")
	}

	rsp := ClientRsp{
		Msg:   "rspunsub",
		Reqid: req.Reqid,
		Err:   errMsg,
		Data:  reqUnsub,
	}
	bytes, _ := json.Marshal(rsp)
	peer.SendChannel <- bytes
}

////////////////////// object callbacks //////////////////////
func (this *ClientService) OnSubObj(peer *cascade.Peer, subObj *SubObj) {
	clientInfo, _ := peer.ExtraInfo.(*ClientInfo)
	log.Printf("[Info] OnSubObj: %+v <%v> sub %+v\n", clientInfo.User, peer.Conn.RemoteAddr().String(), *subObj)

	var rsp *ClientRsp

	quoteService, ok := this.QuoteServices[subObj.ReqSub.Market]
	if ok {
		status, uuid, err := this.SubManager.TrySub(peer, &subObj.ReqSub, subObj.Reqid)
		if err != nil {
			rsp = &ClientRsp{
				Msg:   "rspsub",
				Reqid: subObj.Reqid,
				Data:  subObj.ReqSub,
				Err:   NewErrorMsg(ErrorCode_Unknown, err.Error()),
			}
		} else if status == TopicStatusSubed {
			rsp = &ClientRsp{
				Msg:   "rspsub",
				Reqid: subObj.Reqid,
				Err:   ErrorMsg{Code: 0, What: ""},
				Data:  subObj.ReqSub,
			}
			this.SubManager.PrintSubStatus()
		} else if status == TopicStatusNone {
			subObj.Uuid = uuid
			quoteService.ObjectMessageChannel <- &cascade.HubObjectMessage{Peer: peer, ObjectName: "subobj", ObjectPtr: subObj}
		}
	} else {
		what := "failed to connect server: " + subObj.ReqSub.Market
		rsp = &ClientRsp{
			Msg:   "rspsub",
			Reqid: subObj.Reqid,
			Err:   NewErrorMsg(ErrorCode_ServerNotConnect, what),
			Data:  subObj.ReqSub,
		}
	}

	if rsp != nil {
		bytes, _ := json.Marshal(*rsp)
		peer.SendChannel <- bytes
	}
}

////////////////////// server message callbacks //////////////////////
func (this *ClientService) OnRspSub(rsp *ServerRsp) {
	var success bool
	if rsp.Err == nil {
		success = true
	} else {
		var errmsg ErrorMsg
		err := mapstructure.Decode(rsp.Err, &errmsg)
		if err == nil && errmsg.Code == 0 {
			success = true
		} else {
			success = false
		}
	}

	subMsgCaches := this.SubManager.ClearWaitSubed(rsp.Uuid, success)
	for subMsgCache, _ := range subMsgCaches {
		clientRsp := ClientRsp{
			Msg:   rsp.Msg,
			Reqid: subMsgCache.Reqid,
			Data:  subMsgCache.Data,
			Err:   rsp.Err,
		}
		bytes, _ := json.Marshal(clientRsp)
		subMsgCache.Client.SendChannel <- bytes
	}

	this.SubManager.PrintSubStatus()
}

func (this *ClientService) OnQuote(quote *Quote) {
	// TODO: 这里等文档改完之后，要优化一波
	if quote.Table == "orderbook" {
		var orderbooks []Orderbook
		err := mapstructure.Decode(quote.Data, &orderbooks)
		if err != nil {
			log.Printf("failed to parse orderbook\n")
			log.Printf("%+v\n", quote.Data)
			return
		}

		for _, orderbook := range orderbooks {
			reqSub := ReqSub{
				Market:   quote.Market,
				Symbol:   orderbook.Symbol,
				Table:    quote.Table,
				Type:     orderbook.Type,
				Optional: quote.Optional,
			}

			peers := this.SubManager.GetTopicPeers(&reqSub)
			if peers == nil {
				log.Printf("failed to get topic peers: %+v\n", reqSub)
				return
			}

			bytes, _ := json.Marshal(quote)
			for peer, _ := range peers {
				peer.SendChannel <- bytes
			}
		}
	} else if quote.Table == "ticker" {
		var tickers []Ticker
		err := mapstructure.Decode(quote.Data, &tickers)
		if err != nil {
			log.Printf("failed to parse ticker: %v\n", err.Error())
			log.Printf("%+v\n", quote.Data)
			return
		}

		for _, ticker := range tickers {
			ticker.Symbol = quote.Symbol
			ticker.Type = quote.Type

			reqSub := ReqSub{
				Market:   quote.Market,
				Symbol:   ticker.Symbol,
				Table:    quote.Table,
				Type:     ticker.Type,
				Optional: quote.Optional,
			}

			peers := this.SubManager.GetTopicPeers(&reqSub)
			if peers == nil {
				log.Printf("failed to get topic peers: %+v\n", reqSub)
				return
			}

			bytes, _ := json.Marshal(quote)
			for peer, _ := range peers {
				peer.SendChannel <- bytes
			}
		}
	} else if quote.Table == "login" {
		bytes, _ := json.Marshal(quote)

		log.Println(bytes)

		for peer, _ := range this.Hub.Peers {
			peer.SendChannel <- bytes
		}
	}
}
