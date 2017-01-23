package main

import stomp "github.com/go-stomp/stomp"
import (
	"fmt"
	"time"
	"encoding/json"
	)

//Connect to ActiveMQ and produce messages
func main() {
	conn, err := stomp.Dial("tcp", "localhost:61613")

	if err != nil {
		fmt.Println(err)
		return;
	}

	Producer(conn)
}

func Producer(conn *stomp.Conn) {
	for {
		time.Sleep(3 * time.Second)

		t := time.Now().Format(time.RFC850)

		mapM := map[string]string{"msg": "Test message", "timestamp": t}

    msg, errjson := json.Marshal(mapM)

		if errjson != nil {
			fmt.Println(errjson)
			return;
		}

		 err := conn.Send("/queue/test-1", "text/plain", msg) // body

		 fmt.Println("Message " + string(msg)  + " sent at " + t)

		if err != nil {
			fmt.Println(err)
			return;
		}
	}
}
