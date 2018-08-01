package nan0chat

import (
	"github.com/yomiji/nan0"
	"time"
	"net"
	"io"
	"fmt"
)

type ChatServer struct {
	users    map[int64]*ConnectedUser
	internal *nan0.Service
}

type ConnectedUser struct {
	conn net.Conn
}

func Serve(port int) (err error) {
	service := &ChatServer{
		internal: &nan0.Service{
			ServiceName: "Nan0 Chat",
			Port:        int32(port),
			HostName:    "localhost",
			StartTime:   time.Now().Unix(),
			ServiceType: "Chat",
		},
		users: make(map[int64]*ConnectedUser),
	}
	listener, err := service.internal.Start()
	if err != nil {
		return err
	}
	defer listener.Close()
	fmt.Println("Secure nan0chat server started. Use interrupt command (ctrl+c) to exit.")
	for ; ; {
		conn, err := listener.Accept()
		if err != nil {
			return err
		}
		//io.Copy(conn, conn)
		newUserId := random.Int63()
		service.users[newUserId] = &ConnectedUser{
			conn: conn,
		}

		fmt.Printf("New user %v connected.\n", newUserId)
		go service.startDistributor(newUserId, conn)
	}
	return
}

// Starts a handler for the given user that will distribute the user's messages to each other connected user
func (s *ChatServer) startDistributor(userId int64, conn net.Conn) {
	defer conn.Close()
	for ; ; {
		buff := make([]byte, 1024)
		message := make([]byte, 0)
		conn.SetReadDeadline(time.Now().Add(500 * time.Millisecond))
		for n, err := conn.Read(buff); !(n == 0 || err == io.EOF); n, err = conn.Read(buff) {
			message = append(message, buff[:n]...)
		}
		if len(message) > 0 {
			for id, user := range s.users {
				if id != userId {
					err := user.conn.SetWriteDeadline(time.Now().Add(15 * time.Second))
					handleErr(err, nil)
					_, err = user.conn.Write(message)
					handleErr(err, nil)
				}
			}
		}
		time.Sleep(30 * time.Millisecond)
	}
}
