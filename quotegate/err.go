package main

import (
	"encoding/json"
	"fmt"

	"github.com/MuggleWei/cascade"
)

const (
	ErrorCode_ParseMsg         = 0x01 // 解析消息错误
	ErrorCode_PermissionDenied = 0x02 // 没有权限执行
	ErrorCode_Unhandle         = 0x03 // 没有对应的消息处理
	ErrorCode_Passwd           = 0x04 // 密码错误
	ErrorCode_Unknown          = 0x05 // 未知错误
	ErrorCode_ServerNotConnect = 0x06 // 服务未连接
)

type ErrorMsg struct {
	Code int    `json:"code"`
	What string `json:"what"`
}

func (e ErrorMsg) Error() string {
	return fmt.Sprintf("code: %v,  what: %v", e.Code, e.What)
}
func NewErrorMsg(code int, what string) ErrorMsg {
	return ErrorMsg{
		Code: code,
		What: what,
	}
}

func RspError(peer *cascade.Peer, reqid string, code int, what string) {
	rsp := ClientRsp{
		Msg:   "err",
		Reqid: reqid,
		Err:   NewErrorMsg(code, what),
	}
	bytes, _ := json.Marshal(rsp)
	peer.SendChannel <- bytes
}
