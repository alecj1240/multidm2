package main

import (
    "log"
    "net/http"
    "github.com/gorilla/schema"
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
    sendMessage("time to process your message", response.ChannelId, userInfo.UserAccessToken)
  } else {
    sendEphemeralMessage(response.ResponseUrl, "hey, multidm is installed, but we need permission from you, add it here: https://multidm.alecj1240.repl.co/")
  }
}