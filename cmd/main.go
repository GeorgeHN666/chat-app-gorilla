package main

import (
	"log"
	"net/http"

	"github.com/GeorgeHN666/chat-app-gorilla/internal"
)

func main() {

	log.Println("Websocket Listening on LocalHost:8080")

	ws := internal.NewWsChat()

	http.HandleFunc("/chat", ws.HandleUserConn)

	go ws.UsersChatHandler()

	log.Fatalln(http.ListenAndServe("localhost:8080", nil))

}
