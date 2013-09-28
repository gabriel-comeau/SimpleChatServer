package main

import (
	"bufio"
	"fmt"
	"github.com/gabriel-comeau/SimpleChatCommon"
	"github.com/gabriel-comeau/tbuikit"
	"net"
	"strconv"
)

var (
	clientHolder ClientHolder
	maxId        uint64 = 0
)

// The port that the server listens on
const PORT = 1337

// Main entry point - start listening and on accept create a new client object
// from the network connection and then fire a goroutine to listen to any incoming
// connections on that connection.
func main() {
	server, err := net.Listen("tcp", ":"+strconv.Itoa(PORT))
	if server == nil || err != nil {
		panic("couldn't start listening: " + err.Error())
	}

	clientHolder.init()

	fmt.Println("SimpleChatServer listening...")

	for {
		client, err := server.Accept()

		if err != nil {
			fmt.Printf("ERROR DURING ACCEPT: %v", err)
		}

		if client != nil {
			fmt.Printf("ACCEPTED: %v <-> %v\n", client.LocalAddr(), client.RemoteAddr())
			chatClient := new(ChatClient)
			chatClient.netClient = client
			chatClient.id = getAndIncId()
			clientHolder.addClient(chatClient)
			chatClient.WriteMessage(SimpleChatCommon.Create(fmt.Sprintf("SERVER: Connected to server at: %v", client.LocalAddr()), "green"))
			go handleClientInput(chatClient)
		}
	}
}

// Meant to be run as a goroutine after accepting a connection
// this will continually read new lines from the client as they
// come in and then send them off for processing.
func handleClientInput(client *ChatClient) {
	b := bufio.NewReader(client.netClient)
	for {
		line, err := b.ReadString('\n')
		if err != nil {
			break
		}

		if line == "" || line == "\n" {
			continue
		}

		if !parseCommand(line, client, &clientHolder) {
			go broadcastOutput(line, client)
		}
	}
	// EOF happened
	clientHolder.removeClient(client.id)
}

// Creates the message object specifically for server-wide
// broadcast.
func broadcastOutput(text string, sender *ChatClient) {
	message := createMessageForBroadCast(text, sender)
	broadcast(message)
}

// Perform the actual broadcast - send the message out
// to every client currently in the client holder.
func broadcast(msg *tbuikit.ColorizedString) {
	for _, c := range clientHolder.getClients() {
		c.WriteMessage(msg)
	}
}

// Keep track of the current max ID of clients.  This function
// could potentially be a concurrency issue but only one blocking
// thread is currently using it as part of a loop - this shouldn't
// cause any problems.
func getAndIncId() uint64 {
	maxId++
	return maxId
}
