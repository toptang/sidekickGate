package main

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"sidekick/dataservice_client"
	"time"

	"github.com/go-redis/redis"
	"github.com/gorilla/mux"
	uuid "github.com/satori/go.uuid"
)

const (
	TokenExpireTime = time.Hour * 240
)

type RestHandler struct {
	DataService *dataservice_client.DataServiceClient
	RedisClient *redis.Client
	Router      *mux.Router
}

func NewRestHandler(router *mux.Router) *RestHandler {
	restHandler := &RestHandler{
		DataService: nil,
		RedisClient: nil,
		Router:      router,
	}
	router.HandleFunc("/", restHandler.Index)
	router.HandleFunc("/login", restHandler.Login)
	router.HandleFunc("/order/insert", restHandler.OrderInsert)
	router.HandleFunc("/order/cancel", restHandler.OrderCancel)

	return restHandler
}

func (this *RestHandler) Index(w http.ResponseWriter, r *http.Request) {
	http.Error(w, "404 page not found", http.StatusNotFound)
}

/////////////////////////////////////
// utils
func (this *RestHandler) GenToken() (string, error) {
	uuid, err := uuid.NewV4()
	if err != nil {
		return "", err
	}
	return uuid.String(), nil
}

func (this *RestHandler) CheckToken(user, token string) error {
	result := this.RedisClient.Get(token)
	if result.Err() != nil {
		return result.Err()
	}

	if user != result.Val() {
		return errors.New("token is invalid")
	}

	this.RedisClient.Expire(token, TokenExpireTime)
	return nil
}

func (this *RestHandler) GetUrl(market string) string {
	// TODO:
	return ""
}

/////////////////////////////////////
// login
func (this *RestHandler) Login(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		this.ResponseLogin(&w, ErrorCode_ParseMsg, "failed to read body", "")
		return
	}

	var reqLogin ReqLogin
	err = json.Unmarshal(body, &reqLogin)
	if err != nil {
		this.ResponseLogin(&w, ErrorCode_ParseMsg, "failed to parse login json", "")
		return
	}

	dataReqLogin := dataservice_client.DataReqLogin{
		Name:     reqLogin.User,
		Password: reqLogin.Password,
	}
	dataRsp, err := this.DataService.Login(&dataReqLogin)
	if err != nil {
		this.ResponseLogin(&w, ErrorCode_Unknown, "error password or user", "")
		return
	}

	token := ""
	if dataRsp.Success {
		token, err = this.GenToken()
		if err != nil {
			this.ResponseLogin(&w, ErrorCode_Unknown, "failed to gen token", "")
			return
		}
	} else {
		this.ResponseLogin(&w, dataRsp.Code, dataRsp.Message, "")
		return
	}

	err = this.RedisClient.Set(token, reqLogin.User, TokenExpireTime).Err()
	if err != nil {
		this.ResponseLogin(&w, ErrorCode_Unknown, "failed to save token", "")
		return
	}

	this.ResponseLogin(&w, 0, "", token)
}

func (this *RestHandler) ResponseLogin(w *http.ResponseWriter, code int, what, token string) {
	rsp := RspCommon{
		Code: code,
		What: what,
	}
	if token != "" {
		rsp.Data = RspLogin{Token: token}
	}
	rspmsg, _ := json.Marshal(rsp)
	fmt.Fprint(*w, string(rspmsg))
}

/////////////////////////////////////
// order insert
func (this *RestHandler) OrderInsert(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		panic(err)
	}

	var orderInsert ReqOrderInsert
	err = json.Unmarshal(body, &orderInsert)
	if err != nil {
		panic(err)
	}

	err = this.CheckToken(orderInsert.User, orderInsert.Token)
	if err != nil {
		this.ResponseOrder(&w, ErrorCode_PermissionDenied, "invalid token", nil)
		log.Println(err)
		return
	}

	url := "http://127.0.0.1:8181/strader/okex/future-trade"
	req, err := http.NewRequest("POST", url, bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")
	client := &http.Client{Transport: transport}
	rsp, err := client.Do(req)
	if err != nil {
		this.ResponseOrder(&w, ErrorCode_Unknown, "failed to insert order", nil)
		log.Println(err)
		return
	}
	defer rsp.Body.Close()

	body, _ = ioutil.ReadAll(rsp.Body)
	var rspOrderInsert RspOrder
	err = json.Unmarshal(body, &rspOrderInsert)
	if err != nil {
		this.ResponseOrder(&w, ErrorCode_ParseMsg, "failed to parse rsp insert order json", nil)
		log.Println(err)
		return
	}

	this.ResponseOrder(&w, 0, "", &rspOrderInsert)
}

/////////////////////////////////////
// order cancel
func (this *RestHandler) OrderCancel(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		panic(err)
	}

	var orderCancel ReqOrderCancel
	err = json.Unmarshal(body, &orderCancel)
	if err != nil {
		panic(err)
	}

	err = this.CheckToken(orderCancel.User, orderCancel.Token)
	if err != nil {
		this.ResponseOrder(&w, ErrorCode_PermissionDenied, "invalid token", nil)
		log.Println(err)
		return
	}

	url := "http://127.0.0.1:8181/strader/okex/future-cancel"
	req, err := http.NewRequest("POST", url, bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")
	client := &http.Client{Transport: transport}
	rsp, err := client.Do(req)
	if err != nil {
		this.ResponseOrder(&w, ErrorCode_Unknown, "failed to cancel order", nil)
		log.Println(err)
		return
	}
	defer rsp.Body.Close()

	body, _ = ioutil.ReadAll(rsp.Body)
	var rspOrder RspOrder
	err = json.Unmarshal(body, &rspOrder)
	if err != nil {
		this.ResponseOrder(&w, ErrorCode_ParseMsg, "failed to parse rsp insert order json", nil)
		log.Println(err)
		return
	}

	this.ResponseOrder(&w, 0, "", &rspOrder)
}

func (this *RestHandler) ResponseOrder(w *http.ResponseWriter, code int, what string, rspOrderInsert *RspOrder) {
	rsp := RspCommon{
		Code: code,
		What: what,
	}
	if rspOrderInsert != nil {
		rsp.Data = *rspOrderInsert
	}
	rspmsg, _ := json.Marshal(rsp)
	fmt.Fprint(*w, string(rspmsg))
}
