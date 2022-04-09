package internal

import (
	"encoding/json"
	"fmt"
	"log"

	"github.com/GeorgeHN666/chat-app-gorilla/internal/models"
	"github.com/GeorgeHN666/chat-app-gorilla/internal/utils"
	"github.com/gorilla/websocket"
)

type UserChat struct {
	Channels   *Channel
	UserName   string
	Connection *websocket.Conn
}

func NewUserChat(c *Channel, username string, conn *websocket.Conn) *UserChat {
	return &UserChat{
		Channels:   c,
		UserName:   username,
		Connection: conn,
	}
}

// ListeningMessages Its gonna be listening for message
func (u *UserChat) ListeningMessages() {

	for {
		if _, message, err := u.Connection.ReadMessage(); err != nil {
			log.Println("error trying to read message:::", err.Error())
			break
		} else {
			msg := &models.Message{}

			err := json.Unmarshal(message, msg)
			if err != nil {
				log.Println("Couldn´t Parse:::", err.Error())
			} else {
				u.Channels.MessageChan <- msg
			}

		}
	}

	u.Channels.LeaveChan <- u

}

func (u *UserChat) SendMessageToClient(m *models.Message) error {
	m.ID = utils.GetID()

	data, err := json.Marshal(m)
	if err != nil {
		return fmt.Errorf("error trying to send message to client::%v", err)
	} else {
		err = u.Connection.WriteMessage(websocket.TextMessage, data)
		log.Printf("Message sended:%s to %s", m.Sender, m.Target)
		return err
	}

}
