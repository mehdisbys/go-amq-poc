package main

import stomp "github.com/go-stomp/stomp"
import (
	"encoding/json"
	"fmt"
	"time"
	"strconv"
	"os"
)

//Connect to ActiveMQ and produce messages
func main() {
	conn, err := stomp.Dial("tcp", "localhost:61613")

	if err != nil {
		fmt.Println(err)
		return
	}

  messages:= 100000
	startTime := time.Now()
	Producer(messages,conn)
	endTime := time.Since(startTime)

	fmt.Printf("%d messages sent in %s \n", messages, endTime)

}

func Producer(messages int, conn *stomp.Conn) {

	for i := 0; i < messages; i++ {
		t := time.Now().Format(time.RFC850)

		mapM := map[string]string{"msg": "Test message " + strconv.Itoa(i), "orderID": strconv.Itoa(i), "timestamp": t}

		msg, errjson := json.Marshal(mapM)

		if errjson != nil {
			fmt.Println(errjson)
			return
		}

		err := conn.Send("/queue/orders-received", "text/plain", msg, stomp.SendOpt.Receipt, stomp.SendOpt.Header("persistent", "true"))

		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}
	}
}
