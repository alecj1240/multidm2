package main

import (
  "log"
  "net/http"
  "github.com/gorilla/schema"
  "strings"
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
    if (strings.Split(response.Text, " "))[0] == "help" {
      sendHelp(response)
    } else if (strings.Split(response.Text, " "))[0] == "schedule" {
      scheduleMultiDM(response, userInfo)
    } else {
      multiDM(response, userInfo)
    }
    return;
  } else {
    sendEphemeralMessage(response.ResponseUrl, "hey, it looks like multidm is installed for your team, but we need also permission directly from you: https://multidm.carrd.co/#install")
    return;
  }
}