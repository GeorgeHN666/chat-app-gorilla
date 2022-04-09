package internal

import (
	"fmt"
	"log"
	"net/http"
	"strings"

	"github.com/GeorgeHN666/chat-app-gorilla/internal/models"
	"github.com/GeorgeHN666/chat-app-gorilla/internal/utils"
	"github.com/gorilla/websocket"
)

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
	CheckOrigin:     originCheck,
}

func originCheck(r *http.Request) bool {

	log.Printf("Method:%s \n Host:%s \n ReqeustURL:%s \n Protocol:%s ", r.Method, r.Host, r.RequestURI, r.Proto)

	return r.Method == http.MethodGet

}

type MessageChannel chan *models.Message

type UserChan chan *UserChat

// Create communication with the Global Chat and the other users
type Channel struct {
	MessageChan MessageChannel
	LeaveChan   UserChan
}

// Chat Global
type WSChat struct {
	UsersList map[string]*UserChat
	// This Tell If A User Connect
	JoinChan UserChan
	Channels *Channel
}

func NewWsChat() *WSChat {
	return &WSChat{
		UsersList: make(map[string]*UserChat),
		JoinChan:  make(UserChan),
		Channels: &Channel{
			MessageChan: make(MessageChannel),
			LeaveChan:   make(UserChan),
		},
	}
}

func (w *WSChat) HandleUserConn(rw http.ResponseWriter, r *http.Request) {
	ws, err := upgrader.Upgrade(rw, r, nil)
	if err != nil {
		log.Println("Couldnt connect", r.Host, err)
		rw.WriteHeader(http.StatusInternalServerError)
	}

	key := r.URL.Query().Get("username")

	user := strings.TrimSpace(key)

	if user == "" {
		user = fmt.Sprintf("user-%d", utils.GetID())
	}

	u := NewUserChat(w.Channels, user, ws)

	w.JoinChan <- u
	u.ListeningMessages()

}

// This function its gonna handle all the chats
func (w *WSChat) UsersChatHandler() {

	for {
		select {
		case UserChat := <-w.JoinChan:
			w.AddUser(UserChat)
		case message := <-w.Channels.MessageChan:
			w.SendMessage(message)

		case UserChat := <-w.Channels.LeaveChan:
			w.LeaveChat(UserChat.UserName)

		}
	}

}

func (w *WSChat) AddUser(user *UserChat) {
	u, ok := w.UsersList[user.UserName]
	if ok {
		u.Connection = user.Connection
	} else {
		w.UsersList[user.UserName] = user
		log.Printf("New user:%s\n", user.UserName)
	}
}

func (w *WSChat) SendMessage(m *models.Message) {
	user, ok := w.UsersList[m.Target]
	if ok {
		err := user.SendMessageToClient(m)
		if err != nil {
			log.Printf("Couldnt Send message to clien::%s, Error::%v", m.Target, err.Error())
		}
	}
}

func (w *WSChat) LeaveChat(username string) {
	user, ok := w.UsersList[username]
	if ok {
		defer user.Connection.Close()
		delete(w.UsersList, username)
		log.Printf("User::: %s Leave chat\n", username)
	}
}
