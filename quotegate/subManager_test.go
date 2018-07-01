package main

import (
	"log"
	"os"
	"testing"
)

func init() {
	log.SetOutput(os.Stdout)
	log.SetFlags(log.LstdFlags | log.LUTC | log.Lmicroseconds | log.Lshortfile)
}

func Test_SubManager(t *testing.T) {
	manager := NewSubManager()
	reqSub := ReqSub{
		Market: "simulation",
		Symbol: "BCH",
		Type:   "this_week",
		Table:  "orderbook",
		Optional: ReqSubOptional{
			Depth:  5,
			Period: "",
		},
	}

	reqUnsub := ReqUnsub{
		Market: "simulation",
		Symbol: "BCH",
		Type:   "this_week",
		Table:  "orderbook",
		Optional: ReqSubOptional{
			Depth:  5,
			Period: "",
		},
	}

	// gen topic
	topic := manager.GenTopic(&reqSub)
	log.Printf("gen topic: %+v -> %v\n", reqSub, topic)

	// try sub
	status, uuid, err := manager.TrySub(nil, &reqSub, "0")
	if err != nil {
		t.Error(err)
	}
	if status != TopicStatusNone {
		t.Error("failed topic sub status")
	}
	log.Printf("gen uuid: %v\n", uuid)
	if uuid == "" {
		t.Error("failed generate uuid")
	}

	// try sub again
	status, uuid2, err := manager.TrySub(nil, &reqSub, "0")
	if err != nil {
		t.Error(err)
	}
	if status != TopicStatusSubing {
		t.Error("failed topic sub status")
	}
	if uuid2 != "" {
		t.Error("twice generate uuid")
	}

	// success sub
	subMsgCaches := manager.ClearWaitSubed(uuid, true)
	for k, _ := range subMsgCaches {
		if k.Client != nil {
			t.Error("failed get client")
		}
		if k.Reqid != "0" {
			t.Error("failed get reqid")
		}
		log.Printf("get message cache: %+v\n", *k)
	}

	subMsgCaches = manager.ClearWaitSubed(uuid, true)
	if subMsgCaches != nil {
		t.Error("get sub msg caches twice")
	}

	// try sub again
	status, uuid3, err := manager.TrySub(nil, &reqSub, "0")
	if err != nil {
		t.Error(err)
	}
	if status != TopicStatusSubed {
		t.Error("failed topic sub status")
	}
	if uuid3 != "" {
		t.Error("twice generate uuid")
	}

	// unsub
	ok := manager.Unsub(nil, &reqUnsub)
	if !ok {
		t.Error("failed unsub")
	}
}
