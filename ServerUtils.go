package main

import (
	"fmt"
	"github.com/gabriel-comeau/SimpleChatCommon"
	"github.com/gabriel-comeau/tbuikit"
	"strconv"
	"strings"
)

// Creates a new message object from a string (the message body) and a chat client object,
// which is used to get the message's color
func createMessageForBroadCast(text string, sender *ChatClient) *tbuikit.ColorizedString {
	formatted := formatBroadCastMessage(text, sender)
	msg := SimpleChatCommon.Create(formatted, sender.color)
	return msg
}

// Formats a message's test in the "nickname: message" format
func formatBroadCastMessage(message string, sender *ChatClient) string {
	message = strings.Trim(message, " \n")
	var output string
	if sender.nick != "" {
		output = fmt.Sprintf("%v: %v", sender.nick, message)
	} else {
		output = fmt.Sprintf("%v: %v", sender.id, message)
	}

	return output
}

// Formats a message in format suitable for private messages
func formatWhisperMessage(message string, sender *ChatClient) string {
	message = strings.Trim(message, " \n")
	var output string
	if sender.nick != "" {
		output = fmt.Sprintf("<PRIVATE MESSSAGE> %v: %v", sender.nick, message)
	} else {
		output = fmt.Sprintf("<PRIVATE MESSSAGE> %v: %v", sender.id, message)
	}

	return output
}

// Checks if a given nickname is currently being occupied
func nickTaken(n string, holder *ClientHolder) bool {
	for _, c := range clientHolder.getClients() {
		if c.nick == n {
			return true
		} else if n == strconv.FormatUint(c.id, 10) {
			return true
		}
	}
	return false
}
