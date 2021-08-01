package main

import (
  "log"
  "net/http"
  "github.com/gorilla/schema"
  "strings"
  "unicode/utf8"
  "fmt"
)

type IncomingMessage struct {
  Token string `schema:"token"`
  TeamId string `schema:"team_id"`
  TeamDomain string `schema:"team_domain"`
  ChannelId string `schema:"channel_id"`
  ChannelName string `schema:"channel_name"`
  UserId string `schema:"user_id"`
  UserName string `schema:"user_name"`
  Command string `schema:"command"`
  Text string `schema:"text"`
  IsEnterpriseInstall string `schema:"is_enterprise_install"`
  ResponseUrl string `schema:"response_url"`
  TriggerId string `schema:"trigger_id"`
  ApiAppId string `schema:"api_app_id"`
}

func receiveMessage(w http.ResponseWriter, r *http.Request){
  err := r.ParseForm()
  if err != nil {
    log.Fatal(err)
  }

  var response IncomingMessage
  decoder := schema.NewDecoder()

  err = decoder.Decode(&response, r.Form)
  if err != nil {
    log.Fatal(err)
  }

  isUser, userInfo := checkUserExists(response.UserId)

  if isUser == true {
    recipients, message := splitMessage(response.Text)
    teamUsers := getUsers(userInfo.UserAccessToken)

    for _, recipient := range recipients {
      for _, teamUser := range teamUsers.Members {
        if recipient == teamUser.Name {
          fmt.Println("sending to: " + teamUser.Name)
          sendMessage(message, teamUser.Id, userInfo.UserAccessToken)
        }
      }
    }
  } else {
    sendEphemeralMessage(response.ResponseUrl, "hey, multidm is installed for your team, but we need also permission directly from you:  https://multidm.carrd.co/#install")
  }
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