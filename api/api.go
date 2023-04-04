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
	FPing uint8 = iota

	FUserCreate
	FLogIn
	FLogOut
	FUserInfo
)

type Request struct {
	Func uint8
	Args []byte
}

type ArgsFUserCreate struct {
	Login      string
	Password   string
	FirstName  string
	LastName   string
	Patronymic string
}

type Response struct {
	Code uint8
	Data []byte
}

const (
	EExists uint8 = iota + 1 // record exists

	EArgsInval uint8 = 253 // invalid arguments
	ENoFun     uint8 = 254 // function does not exist
	EUnknown   uint8 = 255 // unknown error
)

// TODO: remove conn (?)

type RequestHandler func(conn net.Conn, r *Request) (*Response, error)

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

			// read the request
			buf := make([]byte, 4096)
			n, err := conn.Read(buf)
			if err != nil {
				log.Fatal(err)
			}

			// deserialize the request
			var msg Request
			err = msgpack.Unmarshal(buf[:n], &msg)
			if err != nil {
				log.Fatal(err)
			}

			// handle the request
			response, err := handler(conn, &msg)
			if err != nil {
				log.Fatal(err)
			}

			// serialize the response
			data, err := msgpack.Marshal(response)
			if err != nil {
				log.Fatal(err)
			}

			// send the response
			_, err = conn.Write(data)
			if err != nil {
				log.Fatal(err)
			}
		}(conn)
	}

	return nil
}
