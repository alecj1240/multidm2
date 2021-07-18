package main

import (
    "fmt"
    "log"
    "net/http"
    "io/ioutil"
    "os"
    "encoding/json"
    "github.com/m3o/m3o-go/client"
)

type AuthorizedResponse struct {
  Ok bool `json:"ok"`
  AppId string `json:"app_id"`
  BotUserId string `json:"bot_user_id"`
  IsEnterpriseInstall bool `json:"is_enterprise_install"`
  BotAccessToken string `json:"access_token"`

  AuthedUser struct {
    Id string `json:"id"`
    AccessToken string `json:"access_token"`
  } `json:"authed_user"`

  Team struct {
    Id string `json:"id"`
    Name string `json:"name"`
  } `json:"team"`
}

type DBUser struct {
  BotAccessToken string
  UserAccessToken string
  AuthedUser string
  BotUserId string
  IsEnterpriseInstall bool
  AppId string
  TeamId string
  TeamName string
}

func appAuthorized(w http.ResponseWriter, r *http.Request){
  code, ok := r.URL.Query()["code"]

  if !ok || len(code[0]) < 1 {
    fmt.Printf("Url Param 'code' is missing")
    return // best to redirect to an error page
  }

  url := "https://slack.com/api/oauth.v2.access?client_id=" + os.Getenv("CLIENT_ID") + "&client_secret=" + os.Getenv("CLIENT_SECRET") + "&code=" + code[0]

  resp, err := http.Get(url)
   if err != nil {
      log.Fatalln(err)
   }

  body, err := ioutil.ReadAll(resp.Body)
  if err != nil {
    log.Fatalln(err)
  }

  var response AuthorizedResponse	
  json.Unmarshal(body, &response)

  c := client.NewClient(&client.Options{Token: os.Getenv("MICRO_TOKEN")})

	req := map[string]interface{}{
    "record": map[string]interface{} {
      "botAccessToken": response.BotAccessToken,
      "userAccessToken": response.AuthedUser.AccessToken,
      "authedUser": response.AuthedUser.Id,
      "botUserId": response.BotUserId,
      "isEnterpriseInstall": response.IsEnterpriseInstall,
      "appId": response.AppId,
      "teamId": response.Team.Id,
      "teamName": response.Team.Name,
    },
    "table" : "users",
  }

	var rsp map[string]interface{}

	if err := c.Call("db", "Create", req, &rsp); err != nil {
		log.Fatalln(err)
    // redirect to an error page here
	}

  http.Redirect(w, r, "https://multidm.alecj1240.repl.co/done.html", 301)
}

func checkUserExists(userId string) (bool, DBUser) {
  c := client.NewClient(&client.Options{Token: os.Getenv("MICRO_TOKEN")})
  query := "authedUser == '" + userId + "'"

	req := map[string]interface{}{
    "table" : "users",
    "query" : query,
  }

	var rsp map[string][]DBUser

	if err := c.Call("db", "Read", req, &rsp); err != nil {
		fmt.Println(err)
	}

  records, ok := rsp["records"]
  if ok {
    return true, records[0]
  }

  return false, DBUser{}
}