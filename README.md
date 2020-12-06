[![PkgGoDev](https://pkg.go.dev/badge/github.com/korableg/flow)](https://pkg.go.dev/github.com/korableg/flow) [![Go Report Card](https://goreportcard.com/badge/github.com/korableg/flow)](https://goreportcard.com/report/github.com/korableg/flow) ![Go](https://github.com/korableg/flow/workflows/Go/badge.svg)

# Flow
The simplest message queue framework

RESTful implementation - https://github.com/korableg/Flow-Gin

# Usage:
```go
// Create database factory. Flow supports leveldb, but you can write your driver by implementing repo.DB interface
db := leveldb.New("./")

// Creating the flow instance
f := flow.New(db)

// Creating a node producer
producer, err := f.NewNode("producer", true)
if err != nil {
	panic(err)
}

// Creating a node consumer
consumer, err := f.NewNode("consumer", true)
if err != nil {
	panic(err)
}

// Creating a another one node consumer
consumer1, err := f.NewNode("consumer1", true)
if err != nil {
	panic(err)
}

// Creating a hub
hub, err := f.NewHub("mainhub")
if err != nil {
	panic(err)
}

// Adding the nods into the hub
err = f.AddNodeToHub(hub.Name(), consumer.Name())
if err != nil { panic(err) }
err = f.AddNodeToHub(hub.Name(), consumer1.Name())
if err != nil { panic(err) }

dummyData := make([]byte, 0)
// Sending message to hub, message will be deliver to each node in hub
f.SendMessageToHub(producer.Name(), hub.Name(), dummyData)

// or sending message to node directly
f.SendMessageToNode(producer.Name(), consumer.Name(), dummyData)

// The Flow guarantees messages order, and you get them in the same order
m := consumer.FrontMessage()
_ = m

// After receiving the message, you must confirm message, else FrontMessage return same message
consumer.RemoveFrontMessage()

// Closes Flow db before exit application
f.Close()
```
