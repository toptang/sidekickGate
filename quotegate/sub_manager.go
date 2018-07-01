package main

import (
	"log"
	"strconv"
	"time"

	"github.com/MuggleWei/cascade"
	uuid "github.com/satori/go.uuid"
)

const (
	TopicStatusNone   = 1 // nobody sub
	TopicStatusSubing = 2 // somebody subed, but have no response
	TopicStatusSubed  = 3 // already subed
)

type SubManager struct {
	// topic status - nil: TopicStatusNone, false: TopicStatusSubing , true: TopicStatusSubed
	TopicStatus map[string]bool

	///////////////// already success subed /////////////////
	// topic - peer set
	TopicPeers map[string]map[*cascade.Peer]bool
	// peer - topic set
	PeerTopics map[*cascade.Peer]map[string]bool

	///////////////// subing /////////////////
	// uuid - subMsgCache set
	UuidMap map[string]map[*SubMsgCache]bool
	// peer - uuid set
	PeerUuidMap map[*cascade.Peer]map[string]bool
	// topic - uuid set
	TopicUuid map[string]string
}

type SubMsgCache struct {
	Client    *cascade.Peer // client peer
	Reqid     string        // client reqid
	Data      interface{}   // client data
	Timestamp int64         // cache timestamp
}

func NewSubManager() *SubManager {
	return &SubManager{
		TopicStatus: make(map[string]bool),
		TopicPeers:  make(map[string]map[*cascade.Peer]bool),
		PeerTopics:  make(map[*cascade.Peer]map[string]bool),
		UuidMap:     make(map[string]map[*SubMsgCache]bool),
		PeerUuidMap: make(map[*cascade.Peer]map[string]bool),
		TopicUuid:   make(map[string]string),
	}
}

func (this *SubManager) GetSubStatus(reqSub *ReqSub) int {
	topic := this.GenTopic(reqSub)
	return this.getTopicStatus(topic)
}

func (this *SubManager) GenTopic(reqSub *ReqSub) string {
	topic := reqSub.Market + "%" + reqSub.Symbol + "%" + reqSub.Type + "%" + reqSub.Table
	depth := strconv.Itoa(reqSub.Optional.Depth)
	topic = topic + "-" + depth + "|" + reqSub.Optional.Period
	return topic
}

func (this *SubManager) GenTopicByUnsub(reqSub *ReqUnsub) string {
	topic := reqSub.Market + "%" + reqSub.Symbol + "%" + reqSub.Type + "%" + reqSub.Table
	depth := strconv.Itoa(reqSub.Optional.Depth)
	topic = topic + "-" + depth + "|" + reqSub.Optional.Period
	return topic
}

func (this *SubManager) GetTopicPeers(reqSub *ReqSub) map[*cascade.Peer]bool {
	topic := this.GenTopic(reqSub)
	return this.TopicPeers[topic]
}

func (this *SubManager) TrySub(peer *cascade.Peer, reqSub *ReqSub, reqId string) (int, string, error) {
	topic := this.GenTopic(reqSub)

	status := this.GetSubStatus(reqSub)
	if status == TopicStatusSubed {
		peerSet, ok := this.TopicPeers[topic]
		if !ok {
			this.TopicPeers[topic] = make(map[*cascade.Peer]bool)
			peerSet = this.TopicPeers[topic]
		}
		peerSet[peer] = true

		topicSet, ok := this.PeerTopics[peer]
		if !ok {
			this.PeerTopics[peer] = make(map[string]bool)
			topicSet = this.PeerTopics[peer]
		}
		topicSet[topic] = true
		return status, "", nil
	} else if status == TopicStatusSubing {
		uuid, ok := this.TopicUuid[topic]
		if !ok {
			what := "can't find subing topic: " + topic
			return 0, "", NewErrorMsg(ErrorCode_Unknown, what)
		}
		ok = this.addSubWait(peer, reqSub, reqId, uuid, topic)
		if !ok {
			what := "failed to add sub wait: " + topic
			return 0, "", NewErrorMsg(ErrorCode_Unknown, what)
		}
		return status, "", nil
	} else {
		uuidV4, err := uuid.NewV4()
		if err != nil {
			return 0, "", err
		}

		uuid := uuidV4.String()
		this.startSubWait(peer, reqSub, reqId, uuid, topic)
		return status, uuid, nil
	}
}

func (this *SubManager) ClearWaitSubed(uuid string, success bool) map[*SubMsgCache]bool {
	//  get topic
	subMsgCaches := this.UuidMap[uuid]
	var reqSub *ReqSub
	for k, _ := range subMsgCaches {
		reqSub, _ = k.Data.(*ReqSub)
		break
	}
	if reqSub == nil {
		return nil
	}
	topic := this.GenTopic(reqSub)

	// clean uuid - subMsgCache set
	delete(this.UuidMap, uuid)

	// clear peer - uuid set
	for k, _ := range subMsgCaches {
		uuidMaps, ok := this.PeerUuidMap[k.Client]
		if ok {
			delete(uuidMaps, uuid)
		}
	}

	// clear topic - uuid
	delete(this.TopicUuid, topic)

	// topic status
	if success {
		this.TopicStatus[topic] = true
	} else {
		delete(this.TopicStatus, topic)
	}

	// topic - peer set and peer - topic set
	if success {
		peerSet, ok := this.TopicPeers[topic]
		if !ok {
			this.TopicPeers[topic] = make(map[*cascade.Peer]bool)
			peerSet = this.TopicPeers[topic]
		}

		for k, _ := range subMsgCaches {
			peer := k.Client
			peerSet[peer] = true

			topicSet, ok := this.PeerTopics[peer]
			if !ok {
				this.PeerTopics[peer] = make(map[string]bool)
				topicSet = this.PeerTopics[peer]
			}
			topicSet[topic] = true
		}
	}

	return subMsgCaches
}

func (this *SubManager) Unsub(peer *cascade.Peer, unsub *ReqUnsub) bool {
	topic := this.GenTopicByUnsub(unsub)
	status := this.getTopicStatus(topic)
	if status != TopicStatusSubed {
		return false
	}

	topics, ok := this.PeerTopics[peer]
	if !ok {
		return false
	}

	_, ok = topics[topic]
	if !ok {
		return false
	}

	delete(topics, topic)

	peers, ok := this.TopicPeers[topic]
	if ok {
		delete(peers, peer)
	}

	return true
}

func (this *SubManager) PeerUnsub(peer *cascade.Peer) {
	topics, ok := this.PeerTopics[peer]
	if !ok {
		return
	}

	for k, _ := range topics {
		peers, ok := this.TopicPeers[k]
		if ok {
			delete(peers, peer)
		}
	}

	delete(this.PeerTopics, peer)
}

func (this *SubManager) PrintSubStatus() {
	for topic, clientPeer := range this.TopicPeers {
		clients := []string{}
		for client, _ := range clientPeer {
			ci, _ := client.ExtraInfo.(*ClientInfo)
			clients = append(clients, ci.User)
		}
		log.Printf("topic: %v -- %+v\n", topic, clients)
	}
}

func (this *SubManager) getTopicStatus(topic string) int {
	status, ok := this.TopicStatus[topic]
	if !ok {
		return TopicStatusNone
	} else {
		if status {
			return TopicStatusSubed
		} else {
			return TopicStatusSubing
		}
	}
}

func (this *SubManager) startSubWait(peer *cascade.Peer, reqSub *ReqSub, reqId, uuid, topic string) bool {
	// topic status
	this.TopicStatus[topic] = false

	// uuid - subMsgCache set
	subMsgCaches, ok := this.UuidMap[uuid]
	if !ok {
		this.UuidMap[uuid] = make(map[*SubMsgCache]bool)
		subMsgCaches = this.UuidMap[uuid]
	}
	subMsgCache := &SubMsgCache{
		Client:    peer,
		Reqid:     reqId,
		Data:      reqSub,
		Timestamp: time.Now().Unix(),
	}
	subMsgCaches[subMsgCache] = true

	// peer - uuid set
	uuidMaps, ok := this.PeerUuidMap[peer]
	if !ok {
		this.PeerUuidMap[peer] = make(map[string]bool)
		uuidMaps = this.PeerUuidMap[peer]
	}
	uuidMaps[uuid] = true

	// topic - uuid
	this.TopicUuid[topic] = uuid

	return true
}

func (this *SubManager) addSubWait(peer *cascade.Peer, reqSub *ReqSub, reqId, uuid, topic string) bool {
	// uuid - subMsgCache set
	subMsgCaches, ok := this.UuidMap[uuid]
	if !ok {
		return false
	}
	subMsgCache := &SubMsgCache{
		Client:    peer,
		Reqid:     reqId,
		Data:      reqSub,
		Timestamp: time.Now().Unix(),
	}
	subMsgCaches[subMsgCache] = true

	// peer - uuid set
	uuidMaps, ok := this.PeerUuidMap[peer]
	if !ok {
		this.PeerUuidMap[peer] = make(map[string]bool)
		uuidMaps = this.PeerUuidMap[peer]
	}
	uuidMaps[uuid] = true

	return true
}
