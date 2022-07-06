package main

import (
	"bufio"
	"flag"
	"fmt"
	"log"
	"net"
)

//we will use this type to send messages through the chat
type Client chan<- string

var (
	//clients that are connecting to our chat
	incomingClients = make(chan Client)
	//clients that are disconnecting from our chat
	leavingClients = make(chan Client)
	//messages of our chat
	messages = make(chan string)
)

var (
	host = flag.String("h", "localhost", "host")
	port = flag.String("p", "3090", "port")
)

//every team a client connects to our server, it will be assigned to HandleConnection
func HandleConnection(conn net.Conn) {
	defer conn.Close()
	//this channel will be used to send messages to the client
	message := make(chan string)
	//this function prints every message through the conn variable
	go MessageWrite(conn, message)

	clientName := conn.RemoteAddr().String()
	//this message will be only sent to the client that just connected
	message <- "Welcome to the chat " + clientName
	//this message will be sent to all the clients
	messages <- "New client " + clientName + " has joined the chat"

	//add the client channel to the list of channels
	incomingClients <- message

	//this loop will run until the client disconnects and will scan for messages
	inputMessage := bufio.NewScanner(conn)
	for inputMessage.Scan() {
		messages <- clientName + ": " + inputMessage.Text()
	}
	//remove the client channel from the list of channels
	leavingClients <- message
	//this message will be sent to all the clients
	messages <- "Client " + clientName + " has left the chat"
}

func MessageWrite(conn net.Conn, messages <-chan string) {
	for message := range messages {
		fmt.Fprintln(conn, message)
	}
}

func Broadcast() {
	clients := make(map[Client]bool)
	for {
		select {
		case message := <-messages:
			for client := range clients {
				client <- message
			}
		case newClient := <-incomingClients:
			clients[newClient] = true
		case leavingClient := <-leavingClients:
			delete(clients, leavingClient)
			close(leavingClient)
		}
	}
}

func main() {
	listener, err := net.Listen("tcp", *host+":"+*port)
	if err != nil {
		log.Fatal(err)
	}
	go Broadcast()
	for {
		conn, err := listener.Accept()
		if err != nil {
			log.Print(err)
			continue
		}
		go HandleConnection(conn)
	}
}
