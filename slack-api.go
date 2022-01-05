package main

import (
    "log"
    "net/http"
    "io/ioutil"
    "encoding/json"
    "bytes"
)

type GetUserResponse struct {
  Ok string `json:"okay"`
  Members []User `json:"members"`
}

type User struct {
  Id string `json:"id"`
  Team_id string `json:"team_id"`
  Name string `json:"name"`
  Real_Name string `json:"real_name"`
  AccessToken string `json:"access_token"`
  Is_Admin bool `json:"is_admin"`
  Is_Owner bool `json:"is_owner"`
  Is_Bot bool `json:"is_bot"`
}

type UserInfoResponse struct {
	Ok   bool `json:"ok"`
	UserInfo struct {
		ID       string `json:"id"`
		TeamID   string `json:"team_id"`
		Name     string `json:"name"`
		RealName string `json:"real_name"`
		Tz       string `json:"tz"`
		TzLabel  string `json:"tz_label"`
		TzOffset int    `json:"tz_offset"`
	} `json:"user"`
}

func getUsers(accessToken string) GetUserResponse {
  url := "https://slack.com/api/users.list"

  req, err := http.NewRequest("GET", url, nil)
  req.Header.Add("Authorization", "Bearer " + accessToken)

  client := &http.Client{}
  resp, err := client.Do(req)

  if err != nil {
    log.Println("Error on response.\n[ERROR] -", err)
  }

  defer resp.Body.Close()

  body, err := ioutil.ReadAll(resp.Body)
  if err != nil {
    log.Println("Error while reading the response bytes:", err)
  }

  var response GetUserResponse	
  json.Unmarshal(body, &response)

  return response
}

func getUserInfo(accessToken string, userId string) UserInfoResponse {
  url := "https://slack.com/api/users.info?user=" + userId

  req, err := http.NewRequest("GET", url, nil)
  req.Header.Add("Authorization", "Bearer " + accessToken)

  client := &http.Client{}
  resp, err := client.Do(req)

  if err != nil {
    log.Println("Error on response.\n[ERROR] -", err)
  }

  defer resp.Body.Close()

  body, err := ioutil.ReadAll(resp.Body)
  if err != nil {
    log.Println("Error while reading the response bytes:", err)
  }

  var response UserInfoResponse	
  json.Unmarshal(body, &response)

  return response
}

func sendMessage(text string, channel string, accessToken string) {
  authStr := "Bearer " + accessToken

  values := map[string]string{
    "channel": channel, 
    "text": text,
    "as_user": "true",
  }
  
  jsonData, err := json.Marshal(values)

  var jsonStr = []byte(jsonData)

	req, err := http.NewRequest("POST", "https://slack.com/api/chat.postMessage", bytes.NewBuffer(jsonStr))

	req.Header.Set("Content-Type", "application/json")
  req.Header.Set("Authorization", authStr)

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		panic(err)
	}
	defer resp.Body.Close()
}

func sendEphemeralMessage(url string, text string) {
  values := map[string]string{
    "response_type": "ephemeral", 
    "text": text,
  }
  json_data, err := json.Marshal(values)

  if err != nil {
    log.Fatal(err)
  }

  resp, err := http.Post(url, "application/json", bytes.NewBuffer(json_data))

  if err != nil {
    log.Fatal(err)
  }

  var res map[string]interface{}

  json.NewDecoder(resp.Body).Decode(&res)
}