package flow

import (
	"fmt"
	"github.com/korableg/mini-gin/flow/errs"
	"github.com/korableg/mini-gin/flow/msgs"
	mockDB2 "github.com/korableg/mini-gin/flow/repo/mockDB"
	"math/rand"
	"strconv"
	"testing"
	"time"
)

func TestFlow(t *testing.T) {

	rand.Seed(time.Now().UnixNano())

	mockDB := new(mockDB2.MockDB)

	m := New(mockDB)

	hub, err := m.NewHub("testHub")
	if err != nil {
		t.Fatal(err)
	}
	_, err = m.NewHub("testHub")
	if err != errs.ErrHubIsAlreadyExists {
		t.Error(err)
	}

	hub = nil
	hub = m.GetHub("testHub")
	if hub == nil {
		t.Fatal("Err get hub")
	}

	nodeProducer, _ := m.NewNode("node_producer", true)
	nodeConsumer, _ := m.NewNode("node_consumer", true)
	nodeConsumerDirectly, _ := m.NewNode("node_consumer_directly", true)

	for nodeConsumerDirectly.Len() > 0 {
		if err := nodeConsumerDirectly.RemoveFrontMessage(); err != nil {
			t.Error(err)
		}
	}

	_, err = m.SendMessageToNode("node_producer", "node_consumer_directly", nil)
	if err != nil {
		t.Error(err)
	}

	err = m.AddNodeToHub(hub.Name(), nodeConsumer.Name())
	if err != nil {
		t.Error(err)
	}
	NameNodeConsumer := nodeConsumer.Name()
	nodeConsumer = nil
	nodeConsumer = m.GetNode(NameNodeConsumer)
	if nodeConsumer == nil {
		t.Fatal("node_consumer not found")
	}
	if nodeConsumer.Len() != 4 {
		t.Error("message queue len error")
	}
	for nodeConsumer.Len() > 0 {
		err := nodeConsumer.RemoveFrontMessage()
		if err != nil {
			t.Error(err)
		}
	}
	nodesCount := rand.Intn(100) + 1

	for i := 0; i < nodesCount; i++ {
		n, _ := m.NewNode("testNode"+strconv.Itoa(i), true)
		err := m.AddNodeToHub(hub.Name(), n.Name())
		if err != nil {
			t.Error(err)
		}
	}

	messageCount := rand.Intn(100) + 1
	mSent := make([]*msgs.Message, 0, messageCount*3)
	mSentChan := make(chan interface{})
	mSentChan2 := make(chan interface{})
	mSentChan3 := make(chan interface{})
	mReceivedChan := make(chan []*msgs.Message)

	funcSend := func(out chan interface{}) {
		for i := 0; i < messageCount; i++ {
			data := make([]byte, rand.Intn(1024*1024*10))
			rand.Read(data)
			mes, _ := m.SendMessageToHub(nodeProducer.Name(), hub.Name(), data)
			mSent = append(mSent, mes)
		}
		out <- 1
	}

	funcReceive := func() {
		mReceived := make([]*msgs.Message, messageCount*3, messageCount*3)
		for i := 0; i < messageCount*3; {
			mes, _ := m.GetMessage(nodeConsumer.Name())
			if mes != nil {
				mReceived[i] = mes
				err := m.RemoveMessage(nodeConsumer.Name())
				if err != nil {
					t.Error(err)
				}
				i++
			} else {
				time.Sleep(time.Millisecond * 500)
			}
		}
		mReceivedChan <- mReceived
	}

	go funcSend(mSentChan)
	go funcSend(mSentChan2)
	go funcSend(mSentChan3)
	go funcReceive()

	<-mSentChan
	<-mSentChan2
	<-mSentChan3

	mReceived := <-mReceivedChan

	for i := 0; i < messageCount*3; i++ {
		if mSent[i] != mReceived[i] {
			t.Fatal("Sent != Received")
		}
	}

	fmt.Printf("Nodes count %d\n", nodesCount)
	fmt.Printf("msgs received %d\n", len(mReceived))

	err = m.AddNodeToHub(hub.Name(), "NodeNotFound")
	if err != errs.ErrNodeNotFound {
		t.Error("must be error ERR_NODE_NOT_FOUND")
	}

	err = m.AddNodeToHub("HubNotFound", nodeConsumer.Name())
	if err != errs.ErrHubNotFound {
		t.Error("must be error ERR_HUB_NOT_FOUND")
	}

	err = m.DeleteNodeFromHub(hub.Name(), "NodeNotFound")
	if err != errs.ErrNodeNotFound {
		t.Error("must be error ERR_NODE_NOT_FOUND")
	}

	err = m.DeleteNodeFromHub("HubNotFound", nodeConsumer.Name())
	if err != errs.ErrHubNotFound {
		t.Error("must be error ERR_HUB_NOT_FOUND")
	}

	err = m.DeleteNodeFromHub(hub.Name(), nodeConsumer.Name())
	if err != nil {
		t.Error(err)
	}

	err = m.DeleteNode(nodeConsumer.Name())
	if err != nil {
		t.Error(err)
	}

	hubs := m.GetAllHubs()
	if len(hubs) != 2 {
		t.Errorf("hubs len have %d, must 2", len(hubs))
	}

	err = m.DeleteHub(hub.Name())
	if err != nil {
		t.Error(err)
	}

	err = m.Close()
	if err != nil {
		t.Fatal(err)
	}

}
