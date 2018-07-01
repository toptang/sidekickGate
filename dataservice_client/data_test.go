package dataservice_client

import (
	"log"
	"os"
	"testing"
)

var dataServiceAddr = "http://127.0.0.1:6001"

func init() {
	log.SetOutput(os.Stdout)
	log.SetFlags(log.LstdFlags | log.LUTC | log.Lmicroseconds | log.Lshortfile)
}

func Test_Add_User(t *testing.T) {
	client := NewDataService(dataServiceAddr)

	req := DataReqAddUser{
		Name:     "foo",
		Password: "bar",
	}
	rsp, err := client.AddUser(&req)

	if err != nil {
		t.Error(err)
	}

	if !rsp.Success {
		t.Error("failed to add user")
	}

	log.Printf("Success add user: %+v, get response: %+v\n", req, *rsp)
}

func Test_Login(t *testing.T) {
	client := NewDataService(dataServiceAddr)

	req := DataReqLogin{
		Name:     "foo",
		Password: "bar",
	}
	rsp, err := client.Login(&req)

	if err != nil {
		t.Error(err)
	}

	if !rsp.Success {
		t.Error("failed to login")
	}

	log.Printf("Success login: %+v, get response: %+v\n", req, *rsp)
}

func Test_Del_User(t *testing.T) {
	client := NewDataService(dataServiceAddr)

	req := DataReqDelUser{
		Name: "foo",
	}
	rsp, err := client.DelUser(&req)

	if err != nil {
		t.Error(err)
	}

	if !rsp.Success {
		t.Error("failed to delete user")
	}

	log.Printf("Success delete user: %+v, get response: %+v\n", req, *rsp)
}
