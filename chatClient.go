package nan0chat

import (
	"github.com/yomiji/nan0"
	"time"
	"fmt"
	"math/rand"
)

type ChatClient struct {
	internal *nan0.Service
	user     *User
}

func NewChatClient() (client *ChatClient) {
	nan0.NoLogging()
	client = &ChatClient{
	}
	return
}

func (client *ChatClient) Connect() {
	client.internal = &nan0.Service{
		HostName:    *Host,
		Port:        int32(*Port),
		ServiceType: "Chat",
		ServiceName: "ChatServer",
		StartTime:   time.Now().Unix(),
		Expired:     false,
	}
	newUserId := rand.Int63()
	client.user = &User{
		UserName: fmt.Sprintf("Connected_User#%v\n", newUserId),
	}
	if *CustomUsername != "" {
		client.user.SetUserName(*CustomUsername)
	}
	encKey, authKey := KeysToNan0Bytes(*EncryptKey, *Signature)
	nan0chat, err := client.internal.DialNan0Secure(encKey, authKey).
		ReceiveBuffer(1).
		SendBuffer(0).
		MessageIdentity(new(ChatMessage)).
		BuildNan0()
	defer nan0chat.Close()

	if err != nil {
		panic(err)
	}
	serviceReceiver := nan0chat.GetReceiver()
	serviceSender := nan0chat.GetSender()

	var chatClientUI ChatClientUI

	messageChannel := make(chan string)
	go chatClientUI.Start(fmt.Sprintf("@%v: ", client.user.UserName), messageChannel)

	for {
		select {
		case m := <-serviceReceiver:
			if message, ok := m.(*ChatMessage); ok {
				chatClientUI.outputBox.addMessage(message.Message)
			}
		case newmsg := <-messageChannel:
			serviceSender <- &ChatMessage{
				Message:   newmsg,
				Time:      time.Now().Unix(),
				MessageId: random.Int63(),
				UserId:    random.Int63(),
			}
		}
	}
}
