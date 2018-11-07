package components

import (
	"framework/messages"
	"shared/shared"
	"fmt"
	"os"
	"shared/parameters"
	"strings"
	"strconv"
	"net"
	"shared/errors"
	"encoding/json"
)

type NotificationConsumer struct{}

var Queues = map[string]chan messages.MessageMOM{}

func (NC NotificationConsumer) I_PosInvP(msg *messages.SAMessage, r *bool) {
	inv := msg.Payload.(messages.Invocation)

	switch inv.Op {
	case "Publish":
		_args := inv.Args.([]interface{})
		_topic := _args[0].(string)
		_msg := _args[1].(map[string]interface{})
		_msgHeader := _msg["Header"].(map[string]interface{})
		_headerDestination := _msgHeader["Destination"].(string)
		_msgPayload := _msg["PayLoad"].(string)
		_msgPub := messages.MessageMOM{Header: messages.Header{Destination: _headerDestination}, PayLoad: _msgPayload}
		_r := NC.Publish(_topic, _msgPub)

		go NC.NotifySubscribers(_msgPub)

		_ter := shared.QueueingTermination{_r}
		*msg = messages.SAMessage{_ter}
	case "Consume":
		_args := inv.Args.([]interface{})
		_topic := _args[0].(string)
		_r := NC.Consume(_topic)

		_ter := shared.QueueingTermination{_r}
		*msg = messages.SAMessage{_ter}

	default:
		fmt.Println("NotificationEngine:: Operation " + inv.Op + " is not implemented by NotificationEngine")
		os.Exit(0)
	}
}

func (NC NotificationConsumer) Publish(topic string, msg messages.MessageMOM) bool {
	r := false

	if _, ok := Queues[topic]; !ok {
		Queues[topic] = make(chan messages.MessageMOM, parameters.QUEUE_SIZE)
	}

	if len(Queues[topic]) < parameters.QUEUE_SIZE {
		Queues[topic] <- msg
		r = true
	} else {
		r = false
	}

	return r
}

func (NotificationConsumer) Consume(topic string) messages.MessageMOM {
	r := messages.MessageMOM{}
	if _, ok := Queues[topic]; !ok {
		Queues[topic] = make(chan messages.MessageMOM, parameters.QUEUE_SIZE)
	}
	if len(Queues[topic]) == 0 {
		r = messages.MessageMOM{Header: messages.Header{Destination: topic}, PayLoad: "QUEUE EMPTY"} // TODO
	} else {
		r = <-Queues[topic]
	}
	return r
}

func (NotificationConsumer) NotifySubscribers(msg messages.MessageMOM){

	// Notify Subscribers
	host := "192.168.0.15"
	port := 1313
	addr := strings.Join([]string{host, strconv.Itoa(port)}, ":")
	conn, err = net.Dial("tcp", addr)

	//defer conn.Close()

	portTmp = port
	if err != nil {
		fmt.Println(err)
		myError := errors.MyError{Source: "Notification Consumer", Message: "Unable to open connection to " + host + " : " + strconv.Itoa(port)}
		myError.ERROR()
	}

	encoder := json.NewEncoder(conn)
	err = encoder.Encode(msg)
	if err != nil {
		fmt.Println(err)
		myError := errors.MyError{Source: "Notification Consumer", Message: "Unable to send data to " + host + ":" + strconv.Itoa(port)}
		myError.ERROR()
	}
	return
}