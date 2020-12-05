package node

import (
	"encoding/json"
	"github.com/korableg/mini-gin/flow/cmn"
	"github.com/korableg/mini-gin/flow/errs"
	"github.com/korableg/mini-gin/flow/msgs"
	"github.com/korableg/mini-gin/flow/repo"
)

type Node struct {
	name     string
	messages *msgs.MessageQueue
}

func New(name string, careful bool, db repo.DB) (n *Node, err error) {
	if err = checkName(name); err != nil {
		return
	}

	var mqdb repo.MQDB
	if db != nil {
		mqdb = db.NewMQRepository(name)
	}

	messages, err := msgs.NewMessageQueue(mqdb, careful)
	if err != nil {
		return
	}

	n = new(Node)
	n.name = name
	n.messages = messages

	return
}

func (n *Node) Name() string {
	return n.name
}

func (n *Node) IsCareful() bool {
	return n.messages.IsCareful()
}

func (n *Node) PushMessage(m *msgs.Message) error {
	return n.messages.Push(m)
}

func (n *Node) FrontMessage() *msgs.Message {
	return n.messages.Front()
}

func (n *Node) RemoveFrontMessage() error {
	return n.messages.RemoveFront()
}

func (n *Node) Len() int {
	return n.messages.Len()
}

func (n *Node) DeleteMessageDB() error {
	return n.messages.DeleteDB()
}

func (n *Node) MarshalJSON() ([]byte, error) {

	nodeMap := make(map[string]interface{})
	nodeMap["name"] = n.Name()
	nodeMap["careful"] = n.IsCareful()
	nodeMap["messages"] = n.Len()

	return json.Marshal(nodeMap)

}

func checkName(name string) error {
	if len(name) == 0 {
		return errs.ErrNodeNameIsempty
	}
	if len([]rune(name)) > 100 {
		return errs.ErrNodeNameOver100
	}
	if !cmn.NameMatchedPattern(name) {
		return errs.ErrNodeNameNotMatchedPattern
	}
	return nil
}
