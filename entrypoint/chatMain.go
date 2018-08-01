package main

import (
	"fmt"
	"nan0chat"
	"flag"
)

func main() {
	flag.Parse()
	if *nan0chat.EncryptKey == "" || *nan0chat.Signature == "" {
		fmt.Println("Please use --key and --sig flags to add security keys")
	} else if *nan0chat.IsServer {
		startServer()
	} else {
		startClient()
	}
}

func startServer() {
	err := nan0chat.Serve(*nan0chat.Port)
	if err != nil {
		panic(err)
	}
}

func startClient() {
	nan0chat.NewChatClient().Connect()
}
