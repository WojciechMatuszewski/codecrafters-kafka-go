package main

import (
	"context"
	"fmt"
	"net"
	"os/signal"
	"syscall"
)

func main() {
	const (
		host string = "localhost"
		port string = "9092"
	)

	server := newServer(host, port)
	err := server.Start()
	if err != nil {
		panic(err)
	}
}

type Server struct {
	Host string
	Port string

	listener *net.Listener
}

func newServer(host string, port string) Server {
	return Server{
		Host: host,
		Port: port,
	}
}

func (s Server) Start() error {
	ctx, stop := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer stop()

	address := fmt.Sprintf("%s:%s", s.Host, s.Port)

	lc := net.ListenConfig{}
	listener, err := lc.Listen(ctx, "tcp", address)
	if err != nil {
		return fmt.Errorf("failed to bind to port %s: %w", s.Port, err)
	}
	defer listener.Close()

	go acceptLoop(ctx, listener)

	<-ctx.Done()

	fmt.Println("Shutting down")

	return nil
}

func acceptLoop(ctx context.Context, listener net.Listener) {
	for {
		select {
		case <-ctx.Done():
			fmt.Println("Stopping to accept connections")
			return
		default:
			connection, err := listener.Accept()
			if err != nil {
				if ctx.Err() != nil {
					fmt.Println("Listener closed, stopping the accept loop")
					return
				}

				fmt.Println("Error accepting the connection:", err)
				continue
			}

			fmt.Println("New connection to the server:", connection.RemoteAddr())

			go connectionLoop(connection)
		}
	}
}

func connectionLoop(connection net.Conn) {
	defer connection.Close()

	received := make([]byte, 1024)
	for {
		_, err := connection.Read(received)
		if err != nil {
			fmt.Println("Failed to read:", err)
			continue
		}

		fmt.Println(string(received))

		header := make([]byte, 4)
		message := []byte{0x00, 0x00, 0x00, 0x07}

		var res []byte
		res = append(res, header...)
		res = append(res, message...)

		_, err = connection.Write(res)
		if err != nil {
			fmt.Println("Failed to write:", err)
			continue
		}
	}
}
