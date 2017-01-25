package main

import stomp "github.com/go-stomp/stomp"
import (
	"encoding/json"
	"fmt"
	"sync"
	"os"
)

func main() {
	conn, err := stomp.Dial("tcp", "localhost:61613")

	if err != nil {
		fmt.Println(err)
	}

	connSend, err := stomp.Dial("tcp", "localhost:61613")

	if err != nil {
		fmt.Println(err)
	}

	received := make(chan string)

	var wg sync.WaitGroup

	wg.Add(2)
	go Consumer(received, conn)
	go OrderPicked(received, connSend)
	wg.Wait()
}

func Consumer(received chan string, conn *stomp.Conn) {

	sub, err := conn.Subscribe("/queue/orders-received", stomp.AckAuto)

	if err != nil {
		fmt.Println("Error : ", err)
	}

	for {
		msg := <-sub.C
		fmt.Println("Success : ", string(msg.Body))
		received <- string(msg.Body)
		<-received
	}

	err = sub.Unsubscribe()
	if err != nil {
		fmt.Println(err)
	}
	defer conn.Disconnect()
}

func OrderPicked(received chan string, conn *stomp.Conn) {
	for {
		select {
		case msg := <-received:

			var mapMsg map[string]string
			json.Unmarshal([]byte(msg), &mapMsg)
			mapM := map[string]string{"msg": "Order picked " + mapMsg["orderID"], "orderID": mapMsg["orderID"]}
			picked, errjson := json.Marshal(mapM)

			if errjson != nil {
				fmt.Println(errjson)
				return
			}

			err := conn.Send("/queue/orders-picked", "text/plain", picked, stomp.SendOpt.Receipt)

			if errjson != nil {
				fmt.Println(err)
				os.Exit(1)
			}
			received <- "y"
		}
	}
}
