package main

import (
	"encoding/json"
	"log"

	"github.com/MuggleWei/cascade"
	uuid "github.com/satori/go.uuid"
)

type QuoteService struct {
	Name   string
	Hub    *cascade.Hub
	Server *cascade.Peer

	ClientService *ClientService

	ObjectCallbacks map[string]ObjectCallbackFn
	Topics          map[ReqSub]bool
}

func NewQuoteService(name string) *QuoteService {
	service := &QuoteService{
		Name: name,
		Hub:  nil,

		ClientService: nil,

		ObjectCallbacks: make(map[string]ObjectCallbackFn),
		Topics:          make(map[ReqSub]bool),
	}

	hub := cascade.NewHub(service, nil, 10240)
	service.Hub = hub
	service.RegisterObjectCallbacks()

	return service
}

// Slot callbacks
func (this *QuoteService) OnActive(peer *cascade.Peer) {
	log.Printf("%v connected: %v\n", this.Name, peer.Conn.RemoteAddr().String())
	this.Server = peer

	for reqSub, _ := range this.Topics {
		uuid, err := uuid.NewV4()
		if err != nil {
			log.Printf("[warning] failed to create uuid!\n")
			return
		}
		subObj := SubObj{
			Uuid:   uuid.String(),
			ReqSub: reqSub,
		}
		this.OnSubObj(nil, &subObj)
	}

	login := ServerReq{
		Msg: "sub",
		Data: ReqSub{
			Market: "okex",
			Symbol: "rocking",
			Type:   "roll",
			Table:  "login",
		},
	}
	b, err := json.Marshal(login)
	if err == nil {
		log.Println("login")
		peer.SendChannel <- b
	}
}

func (this *QuoteService) OnInactive(peer *cascade.Peer) {
	log.Printf("%v disconnected: %v\n", this.Name, peer.Conn.RemoteAddr().String())
	this.Server = nil
}

func (this *QuoteService) OnRead(peer *cascade.Peer, message []byte) {
	this.ClientService.Hub.ByteMessageChannel <- &cascade.HubByteMessage{Peer: peer, Message: message}
}

func (this *QuoteService) OnHubByteMessage(msg *cascade.HubByteMessage) {
}

func (this *QuoteService) OnHubObjectMessage(msg *cascade.HubObjectMessage) {
	callback, ok := this.ObjectCallbacks[msg.ObjectName]
	if ok {
		callback(msg.Peer, msg.ObjectPtr)
	} else {
		log.Printf("[Error] %v: %v object message is not handle\n", this.Name, msg.ObjectName)
	}
}

// callbacks
func (this *QuoteService) RegisterObjectCallbacks() {
	this.ObjectCallbacks["subobj"] = func(peer *cascade.Peer, pointer interface{}) {
		this.OnSubObj(peer, pointer.(*SubObj))
	}
}

func (this *QuoteService) OnSubObj(peer *cascade.Peer, subObj *SubObj) {
	if this.Server == nil {
		what := "failed to connect server: " + subObj.ReqSub.Market
		rsp := ServerRsp{
			Msg:  "rspsub",
			Uuid: subObj.Uuid,
			Err:  NewErrorMsg(ErrorCode_ServerNotConnect, what),
		}
		bytes, _ := json.Marshal(rsp)
		this.ClientService.Hub.ByteMessageChannel <- &cascade.HubByteMessage{Peer: nil, Message: bytes}
		return
	}

	req := ServerReq{
		Msg:  "sub",
		Uuid: subObj.Uuid,
		Data: subObj.ReqSub,
	}
	bytes, _ := json.Marshal(req)
	this.Server.SendChannel <- bytes

	log.Printf("to %v, sub %v\n", subObj.ReqSub.Market, string(bytes))

	this.Topics[subObj.ReqSub] = true
}
