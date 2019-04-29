package main 

import (
	"flag"
	"log"
	"net/http"
	"github.com/gorilla/websocket"
	"fmt"
	"github.com/golang/protobuf/proto"
	"sync"
)

var upgrader = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool {
		return true	
	},
}
var addr = flag.String("addr", "172.16.0.7:8080", "http service address")
var broadcast = make(chan []byte)
var mapLock = &sync.RWMutex{}

var connectedPlayers struct {
	Players map[*websocket.Conn]PlayerData
}

func main() {

	connectedPlayers.Players = make( map[*websocket.Conn]PlayerData)

	http.HandleFunc("/stream", StreamHandler())
	go echo()


	log.Fatal(http.ListenAndServe(*addr, nil))
}

func writer(data []byte) {
	broadcast <- data
}

func echo() {
	for {
		val := <- broadcast
		fmt.Println("\n Sent: ", val)

		for c := range connectedPlayers.Players {
			err := c.WriteMessage(websocket.BinaryMessage, val)
			if err != nil {
				c.Close()
				delete(connectedPlayers.Players, c)
				fmt.Print("bye")
			}
		}

	}
}

func StreamHandler() http.HandlerFunc {
	return http.HandlerFunc(func(res http.ResponseWriter, req *http.Request) { 

		fmt.Print("hello")

		c, err := upgrader.Upgrade(res, req, nil)
		if err != nil {
			log.Print("upgrade:", err)
			return
		}

		defer c.Close()
		

		for {

			_, message, err := c.ReadMessage()
			if err != nil {
				log.Println("read:", err)
				break
			}

			player := PlayerData{}

			err = proto.Unmarshal(message, &player)
			if err != nil {
				panic(err)
			}

			fmt.Println("\n Received: ", player)
		
			mapLock.RLock()
			connectedPlayers.Players[c] = player
			mapLock.RUnlock()

			go writer(message)

		}
	})
}