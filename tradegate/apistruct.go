package main

/////////////////////////////
// restful message

type RspCommon struct {
	Code int         `json:"code"`
	What string      `json:"what"`
	Data interface{} `json:"data,omitempty"`
}

// login
type ReqLogin struct {
	User     string `json:"user"`
	Password string `json:"password"`
}

type RspLogin struct {
	Token string `json:"token"`
}

// order insert
type ReqOrderInsert struct {
	User  string `json:"user"`
	Token string `json:"token"`

	Market       string `json:"market"`
	Symbol       string `json:"symbol"`
	ContractType string `json:"contract_type"`
	Price        string `json:"price"`
	Amount       string `json:"amount"`
	Type         string `json:"type"`
	MatchPrice   string `json:"match_price"`
	LeverRate    string `json:"lever_rate"`
}

type ReqOrderCancel struct {
	User  string `json:"user"`
	Token string `json:"token"`

	Market       string `json:"market"`
	Symbol       string `json:"symbol"`
	OrderId      string `json:"order_id"`
	ContractType string `json:"contract_type"`
}

type RspOrder struct {
	RetCode int          `json:"retcode"`
	Message string       `json:"message"`
	Data    RspDataOrder `json:"data,omitempty"`
}

type RspDataOrder struct {
	OrderId   int  `json:"order_id,omitempty"`
	Result    bool `json:"result,omitempty"`
	ErrorCode int  `json:"error_code,omitempty"`
}

// ws message
type ClientReq struct {
	Msg   string      `json:"msg"`   // message identify
	Reqid string      `json:"reqid"` // client request id
	Data  interface{} `json:"data"`  // client request data
}

type ClientRsp struct {
	Msg   string      `json:"msg"`   // message identify
	Reqid string      `json:"reqid"` // client request id
	Err   interface{} `json:"err"`   // error
	Data  interface{} `json:"data"`  // message data
}

// server common message
type ServerReq struct {
	Msg  string      `json:"msg"`  // message identify
	Uuid string      `json:"uuid"` // uuid
	Data interface{} `json:"data"` // message data
}

type ServerRsp struct {
	Msg  string      `json:"msg"`  // message identify
	Uuid string      `json:"uuid"` // uuid
	Err  interface{} `json:"err"`  // error
	Data interface{} `json:"data"` // message data
}
