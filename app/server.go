package main

import (
	"context"
	"encoding/binary"
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
	_, err := connection.Read(received)
	if err != nil {
		fmt.Println("Failed to read the data:", err)
		return
	}

	// requestLength := binary.BigEndian.Uint32(received[:4])
	// requestApiKey := binary.BigEndian.Uint16(received[4:6])
	// requestApiVersion := binary.BigEndian.Uint16(received[6:8])
	requestCorrelationId := received[8:12]

	// Prepare the response body
	errorCode := []byte{0, 0}                // No Error
	apiVersions := []byte{0, 18, 0, 0, 0, 4} // API key 18, MaxVersion 4
	body := append(errorCode, apiVersions...)

	// Calculate the total length (header + body)
	totalLength := 4 + len(requestCorrelationId) + len(body)
	responseLength := make([]byte, 4)
	binary.BigEndian.PutUint32(responseLength, uint32(totalLength))

	// Prepare the response
	var response []byte
	response = append(response, responseLength...)
	response = append(response, requestCorrelationId...)
	response = append(response, body...)

	_, err = connection.Write(response)
	if err != nil {
		fmt.Println("Failed to write:", err)
	}
}
