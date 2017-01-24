package main

import stomp "github.com/go-stomp/stomp"
import (
	"encoding/json"
	"fmt"
)

//Connect to ActiveMQ and listen for messages
func main() {
	conn, err := stomp.Dial("tcp", "localhost:61613")

	if err != nil {
		fmt.Println(err)
	}

	sub, err := conn.Subscribe("/queue/orders-received", stomp.AckAuto)

	if err != nil {
		fmt.Println("Error : ", err)
	}

	for {
		msg := <-sub.C

		fmt.Println("Success : ", string(msg.Body))

		var mapMsg map[string]string

		errjson := json.Unmarshal(msg.Body, &mapMsg)

		if errjson != nil {
			fmt.Println(errjson)
			return
		}

		mapM := map[string]string{"msg": "Order picked " + mapMsg["orderID"], "orderID": mapMsg["orderID"]}

		picked, errjson2 := json.Marshal(mapM)

		err := conn.Send("/queue/orders-picked", "text/plain", picked) // body

		if errjson2 != nil {
			fmt.Println(errjson)
			return
		}

		if err != nil {
			fmt.Println(err)
			return
		}
	}

	err = sub.Unsubscribe()

	if err != nil {
		fmt.Println(err)
	}
	defer conn.Disconnect()
}
