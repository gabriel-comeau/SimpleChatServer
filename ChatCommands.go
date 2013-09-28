package main

import (
	"fmt"
	"github.com/gabriel-comeau/SimpleChatCommon"
	"regexp"
	"strconv"
	"strings"
)

// Monolithic (for now) function to parse validate and parse incoming message from users.
//
// It checks if the message are commands, if they are it performs the related action
// and returns true.  If not, it returns false and the message can be broadcast
// as a regular chat message.
func parseCommand(line string, client *ChatClient, clientHolder *ClientHolder) bool {
	helpCmdRegex := regexp.MustCompile(`^\/help$`)
	whoCmdRegex := regexp.MustCompile(`^\/who$`)
	nickCmdRegex := regexp.MustCompile(`^\/nick.+$`)
	colorCmdRegex := regexp.MustCompile(`^\/color.+$`)
	whisperCmdRegex := regexp.MustCompile(`^\/w(|hisper) .+ .+$`)

	line = strings.Trim(line, "\n ")

	if helpCmdRegex.MatchString(line) {
		helpSl := make([]string, 0)
		helpSl = append(helpSl, "")
		helpSl = append(helpSl, "SERVER: COMMANDS LIST:")
		helpSl = append(helpSl, "SERVER: /help - Print this message")
		helpSl = append(helpSl, "SERVER: /who - Displays connected users")
		helpSl = append(helpSl, "SERVER: /nick <desirednickname> - Change nickname")
		helpSl = append(helpSl, "SERVER: /color <desiredcolor> - Change text color")
		helpSl = append(helpSl, "SERVER: /w OR /whisper <nickname> <message> - Sends private message to user")
		helpSl = append(helpSl, "")

		helpMessages := SimpleChatCommon.ColStringSliceFromStringSlice(helpSl, "blue")
		for _, msg := range helpMessages {
			client.WriteMessage(msg)
		}
		return true
	}

	if whisperCmdRegex.MatchString(line) {
		splitCmd := strings.Split(line, ` `)
		sendTo := splitCmd[1]

		// This approach only removes the first occurrence of the command and the nick,
		// so they can be put into the message later.
		body := strings.Replace(line, fmt.Sprintf("%v %v", splitCmd[0], sendTo), "", 1)

		// First try to lookup client by nickname
		recip := clientHolder.getClientByNick(sendTo)
		if recip == nil {
			// fall back to ID
			num, _ := strconv.ParseUint(sendTo, 10, 64)
			recip = clientHolder.getClientById(num)
		}

		// If we still don't have a recipient, print an error and bail
		if recip == nil {
			client.WriteMessage(SimpleChatCommon.Create(fmt.Sprintf("SERVER: Could not send private message to: %v - no such user", sendTo), "red"))
		} else {
			recip.WriteMessage(SimpleChatCommon.Create(formatWhisperMessage(body, client), client.color))
			client.WriteMessage(SimpleChatCommon.Create(fmt.Sprintf("SERVER: Sent"), "blue"))
		}

		return true
	}

	if whoCmdRegex.MatchString(line) {
		whoSl := make([]string, 0)
		whoSl = append(whoSl, "SERVER: Connected Users")
		for _, c := range clientHolder.getClients() {
			if c.nick != "" {
				whoSl = append(whoSl, fmt.Sprintf("nickname: %v, id: %v", c.nick, c.id))
			} else {
				whoSl = append(whoSl, fmt.Sprintf("id: %v, no nickname set", c.id))
			}
		}
		whoSl = append(whoSl, "")

		whoMessages := SimpleChatCommon.ColStringSliceFromStringSlice(whoSl, "blue")
		for _, msg := range whoMessages {
			client.WriteMessage(msg)
		}
		return true
	}

	if nickCmdRegex.MatchString(line) {
		splitCmd := strings.Split(line, ` `)
		newNick := splitCmd[1]
		nickValidRegex := regexp.MustCompile(`^[a-zA-Z]`)
		if !nickValidRegex.MatchString(newNick) {
			client.WriteMessage(SimpleChatCommon.Create(fmt.Sprintf("SERVER: Could not change nickname to: %v - nicknames must start with a letter!", newNick), "red"))
		} else if nickTaken(newNick, clientHolder) {
			client.WriteMessage(SimpleChatCommon.Create(fmt.Sprintf("SERVER: Could not change nickname to: %v - already in use", splitCmd[1]), "red"))
		} else if strings.ToLower(newNick) == "server" {
			client.WriteMessage(SimpleChatCommon.Create(fmt.Sprintf("SERVER: Could not change nickname to: %v - already in use", splitCmd[1]), "red"))
		} else {
			msgText := ""
			oldNick := client.nick
			if oldNick != "" {
				msgText = fmt.Sprintf("SERVER: User formerly known as: %v has changed their nickname to: %v", oldNick, newNick)
			} else {
				msgText = fmt.Sprintf("SERVER: User with id: %v has set their nickname to: %v", client.id, newNick)
			}
			client.nick = newNick
			fmt.Printf("Accepted nickname change command for user id: %v to change to nickname: %v\n", client.id, client.nick)
			broadcast(SimpleChatCommon.Create(msgText, "blue"))
		}
		return true
	}

	if colorCmdRegex.MatchString(line) {
		colors := new(SimpleChatCommon.ColorList)
		splitCmd := strings.Split(line, ` `)
		if colors.IsColor(strings.ToLower(splitCmd[1])) {
			client.color = strings.ToLower(splitCmd[1])
			client.WriteMessage(SimpleChatCommon.Create(fmt.Sprintf("SERVER: Color changed to: %v", strings.ToLower(splitCmd[1])), "blue"))
		} else {
			client.WriteMessage(SimpleChatCommon.Create(fmt.Sprintf("SERVER: Sorry: %v was not a valid color choice", splitCmd[1]), "red"))
		}
		return true
	}

	return false
}
