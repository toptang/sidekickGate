package dataservice_client

import (
	"bytes"
	"encoding/json"
	"io/ioutil"
	"log"
	"net/http"
	"time"
)

type DataServiceClient struct {
	Addr      string
	Transport *http.Transport
}

func NewDataService(addr string) *DataServiceClient {
	transport := &http.Transport{
		MaxIdleConns:        0,
		MaxIdleConnsPerHost: 10,
		IdleConnTimeout:     30 * time.Second,
	}
	return &DataServiceClient{
		Addr:      addr,
		Transport: transport,
	}
}

func (this *DataServiceClient) RestPost(url string, body []byte) (*DataRspCommon, error) {
	client := &http.Client{Transport: this.Transport}

	rsp, err := client.Post(url, "application/json", bytes.NewBuffer(body))
	if err != nil {
		log.Printf("Error handle message: %v", err)
		return nil, err
	}
	defer rsp.Body.Close()

	body, err = ioutil.ReadAll(rsp.Body)
	if err != nil {
		return nil, err
	}

	var rspLogin DataRspCommon
	err = json.Unmarshal(body, &rspLogin)
	if err != nil {
		return nil, err
	}

	return &rspLogin, nil
}

func (this *DataServiceClient) Login(req *DataReqLogin) (*DataRspCommon, error) {
	url := this.Addr + "/v1.0/user/login"
	body, err := json.Marshal(*req)
	if err != nil {
		return nil, err
	}

	return this.RestPost(url, body)
}

func (this *DataServiceClient) AddUser(req *DataReqAddUser) (*DataRspCommon, error) {
	url := this.Addr + "/v1.0/user/add"
	body, err := json.Marshal(*req)
	if err != nil {
		return nil, err
	}

	return this.RestPost(url, body)
}

func (this *DataServiceClient) DelUser(req *DataReqDelUser) (*DataRspCommon, error) {
	url := this.Addr + "/v1.0/user/del"
	body, err := json.Marshal(*req)
	if err != nil {
		return nil, err
	}

	return this.RestPost(url, body)
}

func (this *DataServiceClient) UpdateUserPassword(req *DataReqUpdateUserPassword) (*DataRspCommon, error) {
	url := this.Addr + "/v1.0/user/update_password"
	body, err := json.Marshal(*req)
	if err != nil {
		return nil, err
	}

	return this.RestPost(url, body)
}
