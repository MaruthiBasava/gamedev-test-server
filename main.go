package main 

import (
	"flag"
	"log"
	"net/http"
	ws "github.com/gorilla/websocket"
)

var upgrader = ws.Upgrader{}
var addr = flag.String("addr", "localhost:8080", "http service address")

func main() {
	http.HandleFunc("/stream", StreamHandler())
	log.Fatal(http.ListenAndServe(*addr, nil))
}

func StreamHandler() http.HandlerFunc {
	return http.HandlerFunc(func(res http.ResponseWriter, req *http.Request) { 
		c, err := upgrader.Upgrade(res, req, nil)
		if err != nil {
			log.Print("upgrade:", err)
			return
		}

		defer c.Close()

		for {

			mt, message, err := c.ReadMessage()
			if err != nil {
				log.Println("read:", err)
				break
			}

			log.Printf("recv: %s", message)
			
			payload := json.Unmarshal(message)

		}


	})
}