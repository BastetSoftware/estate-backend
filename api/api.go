package api

import (
	"fmt"
	"github.com/vmihailenco/msgpack/v5"
	"log"
	"net"
	"os"
	"os/signal"
	"syscall"
)

const (
	FPing = iota

	FUserCreate
	FLogin
	FLogout
	FUserInfo
)

type Request struct {
	Func int8
	Args []byte
}

type RequestHandler func(conn *net.Conn, r *Request)

// TODO: add cleanup func

func Listen(address string, handler RequestHandler) error {
	socket, err := net.Listen("unix", address)
	if err != nil {
		return fmt.Errorf("unable to create a socket: %v", err)
	}

	channel := make(chan os.Signal, 1)
	signal.Notify(channel, os.Interrupt, syscall.SIGTERM)
	go func() {
		<-channel
		err := os.Remove(address)
		if err != nil {
			log.Fatalf("unable to remove the socket: %v", err)
		}

		os.Exit(1)
	}()

	for {
		conn, err := socket.Accept()
		if err != nil {
			log.Fatal(err)
		}

		go func(conn net.Conn) {
			defer conn.Close()

			buf := make([]byte, 4096)
			n, err := conn.Read(buf)
			if err != nil {
				log.Fatal(err)
			}

			var msg Request
			err = msgpack.Unmarshal(buf[:n], &msg)
			if err != nil {
				log.Fatal(err)
			}
			handler(&conn, &msg)
		}(conn)
	}

	return nil
}
