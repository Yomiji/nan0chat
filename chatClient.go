package nan0chat

import (
	"time"
	"fmt"
	"github.com/Yomiji/nan0"
)

type ChatClient struct {
	internal *nan0.Service
	user     *User
}

func NewChatClient() (client *ChatClient) {
	// turn off nan0 logging
	nan0.NoLogging()
	client = &ChatClient{
	}
	return
}

// Connects to the target service (defined in the application flags)
func (client *ChatClient) Connect() {
	// create the initial client connection descriptor targeting the chat server
	client.internal = &nan0.Service{
		HostName:    *Host,
		Port:        int32(*Port),
		ServiceType: "Chat",
		ServiceName: "ChatServer",
		StartTime:   time.Now().Unix(),
		Expired:     false,
	}

	// create another random user id
	newUserId := random.Int63()
	// use this auto-generated username unless a custom username has been assigned
	client.user = &User{
		UserName: fmt.Sprintf("Connected_User#%v\n", newUserId),
	}
	if *CustomUsername != "" {
		client.user.SetUserName(*CustomUsername)
	}

	// convert the base64 keys to usable byte arrays
	encKey, authKey := KeysToNan0Bytes(*EncryptKey, *Signature)

	// connect to the server securely
	nan0chat, err := client.internal.DialNan0Secure(encKey, authKey).
		ReceiveBuffer(1).
		SendBuffer(0).
		AddMessageIdentity(new(ChatMessage)).
	Build()
	// close the connection when this application closes
	defer nan0chat.Close()

	if err != nil {
		panic(err)
	}

	// get the channels used to communicate with the server
	serviceReceiver := nan0chat.GetReceiver()
	serviceSender := nan0chat.GetSender()

	// create and start a new UI
	var chatClientUI ChatClientUI
	// create a message channel for passing ui message to backend and to server
	messageChannel := make(chan string)
	go chatClientUI.Start(fmt.Sprintf("@%v: ", client.user.UserName), messageChannel)

	for {
		select {
		// when a new message comes in, handle it
		case m := <-serviceReceiver:
			if message, ok := m.(*ChatMessage); ok {
				chatClientUI.outputBox.addMessage(message.Message)
			}
		// when a new message is generated in the UI, broadcast it
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
