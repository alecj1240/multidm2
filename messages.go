package main

import(
  "fmt"
  "strings"
  "unicode/utf8"
)

func multiDM(response IncomingMessage, userInfo DBUser) {
  recipients, message := splitMessage(response.Text)
  if len(recipients) < 1 {
    sendEphemeralMessage(response.ResponseUrl, "Hey, it looks like you didn't tell us who to DM that to!")
    return;
  }

  teamUsers := getUsers(userInfo.UserAccessToken)
  if len(teamUsers.Members) < 1 {
    sendEphemeralMessage(response.ResponseUrl, "Hey, we can't find any slack team members for you -- you might want to try reinstalling multiDM, https://multidm.carrd.co/#install");
    return;
  }

  for _, recipient := range recipients {
    for _, teamUser := range teamUsers.Members {
      if recipient == teamUser.Name {
        fmt.Println("sending to: " + teamUser.Name)
        sendMessage(message, teamUser.Id, userInfo.UserAccessToken)
      }
    }
  }
  
  sendEphemeralMessage(response.ResponseUrl, "Messages sent successfully! You just saved time using MultiDM!");
  return;
}


func splitMessage(message string) ([]string, string) {
  words := strings.Split(message, " ")
  recipients := make([]string, 0)

  for _, word := range words {
    characters := strings.Split(word, "")
    if characters[0] == "@" {
      recipients = append(recipients, trimFirstCharacter(word))
    } else {
      break
    }
	}

  for j := 1; j <= len(recipients); j++ {
    words = words[1:]
  }

  joinedText := strings.Join(words, " ")

  return recipients, joinedText
}

func trimFirstCharacter(s string) string {
	_, i := utf8.DecodeRuneInString(s)
	return s[i:]
}

func sendHelp(response IncomingMessage) {
  sendEphemeralMessage(response.ResponseUrl, "Hey, you need some help with MultiDM? No problem! \n The format to send a message is: */multidm @bob @jane this is my message to bob and jane* \n \n if you continue to have trouble, reach out to alecjones@hey.com");
}