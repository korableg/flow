package hub

import (
	"encoding/json"
	"github.com/korableg/mini-gin/flow/errs"
	"github.com/korableg/mini-gin/flow/msgs"
	"github.com/korableg/mini-gin/flow/node"
	"github.com/korableg/mini-gin/flow/repo/mockDB"
	"testing"
)

func TestHub(t *testing.T) {

	nameHub := "TestHub1"
	nameNode := "TestNode1"

	db := new(MockDB.MockDB)

	_, err := New("   ", nil)
	if err != errs.ErrHubNameNotMatchedPattern {
		t.Error(err)
	}
	_, err = New("", nil)
	if err != errs.ErrHubNameIsempty {
		t.Error(err)
	}
	_, err = New("123456789012345678901234567890123456789012345678901234567890123456789012345678901234567890123456789012345678901234567890", nil)
	if err != errs.ErrHubNameOver100 {
		t.Error(err)
	}
	hub, err := New(nameHub, db)
	if err != nil {
		t.Fatal(err)
	}
	if nameHub != hub.Name() {
		t.Error(nameHub + " != " + hub.Name())
	}

	n, err := node.New(nameNode, true, db)
	if err != nil {
		t.Fatal(err)
	}

	err = hub.AddNode(n)
	if err != nil {
		t.Error(err)
	}

	err = hub.PushMessage(msgs.NewMessage(nameNode, nil))
	if err != nil {
		t.Error(err)
	}

	_, err = json.Marshal(hub)

	err = hub.DeleteNode(n)
	if err != nil {
		t.Error(err)
	}

	err = hub.AddNode(n)
	if err != nil {
		t.Error(err)
	}

	err = hub.deleteAllNodes()
	if err != nil {
		t.Error(err)
	}
	err = hub.deleteNodeDB()
	if err != nil {
		t.Error(err)
	}

}
