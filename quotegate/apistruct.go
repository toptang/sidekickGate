package main

// client common message
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

type MessageIdentify struct {
	Msg string `json:"msg"`
}

type ServerRsp struct {
	Msg  string      `json:"msg"`  // message identify
	Uuid string      `json:"uuid"` // uuid
	Err  interface{} `json:"err"`  // error
	Data interface{} `json:"data"` // message data
}

// client request
type ReqLogin struct {
	User   string `json:"user"`
	Passwd string `json:"passwd"`
}

type ReqSubOptional struct {
	Depth  int    `json:"depth,omitempty"`
	Period string `json:"period,omitempty"`
}

type ReqSub struct {
	Market   string         `json:"market"`
	Symbol   string         `json:"symbol"`
	Type     string         `json:"type"`
	Table    string         `json:"table"`
	Optional ReqSubOptional `json:"optional"`
}

type ReqUnsub ReqSub

// quote
type Quote struct {
	Msg      string         `json:"msg"`    // message type
	Market   string         `json:"market"` // market name - bitmex, okex
	Table    string         `json:"table"`  // table name - orderbook, ticker
	Symbol   string         `json:"symbol"` // symbol name - btc, eos
	Type     string         `json:"type"`   // type name - this_week, next_week
	Channel  string         `json:"channel"`
	Optional ReqSubOptional `json:"optional"` // optional
	Data     interface{}    `json:"data"`     // data
}

type Orderbook struct {
	Symbol    string      `json:"symbol"` // symbol
	Type      string      `json:"type"`   // type
	Bids      interface{} `json:"bids"`   // bids
	Asks      interface{} `json:"asks"`   // asks
	Timestamp int64       `json:"ts"`     // timestamp
}

type Ticker struct {
	Symbol     string `json:"symbol"`
	Type       string `json:"type"`
	LimitHigh  string `json:"limitHigh"`
	LimitLow   string `json:"limitLow"`
	Vol        string `json:"vol"`
	Last       string `json:"last"`
	Sell       string `json:"sell"`
	Buy        string `json:"buy"`
	UnitAmount string `json:"unitAmount"`
	HoldAmount string `json:"hold_amount"`
	ContractId int64  `json:"contractId"`
	High       string `json:"high"`
	Low        string `json:"low"`
}
