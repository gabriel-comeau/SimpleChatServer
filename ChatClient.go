package main

import (
	"github.com/gabriel-comeau/SimpleChatCommon"
	"github.com/gabriel-comeau/termbox-uikit"
	"net"
)

// Represents a connected chat client - holds on to various
// meta data for the client
type ChatClient struct {
	netClient net.Conn
	id        uint64
	nick      string
	color     string
}

// Sends a message to the client by serializing a message object into
// a string and then converting this string to a byte slice and writing
// it over the network.
func (this *ChatClient) WriteMessage(message *termbox-uikit.ColorizedString) {
	bytesToWrite := []byte(SimpleChatCommon.Pack(message))
	this.netClient.Write(bytesToWrite)
}
