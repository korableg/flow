package msgs

import (
	"container/list"
	"github.com/korableg/mini-gin/flow/repo"
	"sort"
	"sync"
)

type MessageQueue struct {
	l       *list.List
	mutex   sync.Mutex
	db      repo.MQDB
	careful bool
}

func NewMessageQueue(db repo.MQDB, careful bool) (*MessageQueue, error) {

	mq := new(MessageQueue)
	mq.l = list.New()
	mq.db = db
	mq.careful = careful

	if mq.db != nil {

		messages, err := mq.db.All()
		if err != nil {
			return nil, err
		}
		if !sort.IsSorted(messages) {
			sort.Sort(messages)
		}

		for _, messageRepo := range messages {
			message := new(Message)
			message.id = messageRepo.ID
			message.data = messageRepo.Data
			message.from = messageRepo.From
			if err := mq.Push(message); err != nil {
				return nil, err
			}
		}
	}

	return mq, nil
}

func (mq *MessageQueue) IsCareful() bool {
	return mq.careful
}

func (mq *MessageQueue) Front() (mes *Message) {
	mq.mutex.Lock()
	defer mq.mutex.Unlock()

	e := mq.l.Front()
	if e != nil {
		mes = e.Value.(*Message)
	}

	return mes
}

func (mq *MessageQueue) Push(mes *Message) (err error) {
	mq.mutex.Lock()
	defer mq.mutex.Unlock()

	if mq.careful && mq.db != nil {
		messageRepo := new(repo.Message)
		messageRepo.ID = mes.ID()
		messageRepo.Data = mes.Data()
		messageRepo.From = mes.From()
		err = mq.db.Store(messageRepo)
		if err != nil {
			return
		}
	}

	mq.l.PushBack(mes)

	return
}

func (mq *MessageQueue) RemoveFront() (err error) {
	mq.mutex.Lock()
	defer mq.mutex.Unlock()

	e := mq.l.Front()
	if e != nil {
		if mq.careful && mq.db != nil {
			err = mq.db.Delete(e.Value.(*Message).id)
			if err != nil {
				return
			}
		}
		mq.l.Remove(e)
	}

	return
}

func (mq *MessageQueue) Len() int {
	mq.mutex.Lock()
	defer mq.mutex.Unlock()

	return mq.l.Len()
}

func (mq *MessageQueue) Close() error {
	if mq.db == nil {
		return nil
	}

	if !mq.careful {

		messageCount := mq.Len()

		mq.mutex.Lock()
		defer mq.mutex.Unlock()

		messagesRepo := make(repo.MessageSlice, messageCount, messageCount)
		i := 0

		for element := mq.l.Front(); element != nil; element = element.Next() {
			message := element.Value.(*Message)
			messageRepo := new(repo.Message)
			messageRepo.ID = message.ID()
			messageRepo.From = message.From()
			messageRepo.Data = message.Data()
			messagesRepo[i] = messageRepo
			i++
		}

		if err := mq.db.StoreAll(messagesRepo); err != nil {
			return err
		}

	}

	return mq.db.Close()
}

func (mq *MessageQueue) DeleteDB() error {

	mq.mutex.Lock()
	defer mq.mutex.Unlock()

	err := mq.db.DeleteDB()
	if err == nil {
		mq.db = nil
	}

	return err
}
