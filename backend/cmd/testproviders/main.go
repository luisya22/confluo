package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os/exec"
	"time"

	"github.com/luisya22/confluo/backend/internal/executor"
	"github.com/luisya22/confluo/backend/internal/providers/github"
)

type GithubDeviceFlowResponse struct {
	DeviceCode      string `json:"device_code"`
	UserCode        string `json:"user_code"`
	VerificationUri string `json:"verification_uri"`
	ExpiresIn       int    `json:"expires_in"`
	Interval        int    `json:"interval"`
}

func main() {
	clientId := ""

	posturl := fmt.Sprintf("https://github.com/login/device/code?client_id=%v&scope=repo", clientId)
	r, err := http.NewRequest(http.MethodPost, posturl, bytes.NewBuffer([]byte{}))
	if err != nil {
		panic(err)
	}

	r.Header.Add("Content-Type", "application/json")
	r.Header.Add("Accept", "application/json")

	client := &http.Client{}
	res, err := client.Do(r)
	if err != nil {
		panic(err)
	}
	defer res.Body.Close()

	var deviceFlowRes GithubDeviceFlowResponse
	err = json.NewDecoder(res.Body).Decode(&deviceFlowRes)
	if err != nil {
		panic(err)
	}

	if res.StatusCode != http.StatusOK {
		panic(res.Status)
	}

	fmt.Println("Go to: ", deviceFlowRes.VerificationUri)
	fmt.Println("Use this code on your browser: ", deviceFlowRes.UserCode)

	exec.Command("cmd.exe", "/c", "start", deviceFlowRes.VerificationUri).Start()

	pollTokenUrl := fmt.Sprintf(
		"https://github.com/login/oauth/access_token?client_id=%v&device_code=%v&grant_type=%v",
		clientId,
		deviceFlowRes.DeviceCode,
		"urn:ietf:params:oauth:grant-type:device_code",
	)

	r, err = http.NewRequest(http.MethodPost, pollTokenUrl, bytes.NewBuffer([]byte{}))

	r.Header.Add("Content-Type", "application/json")
	r.Header.Add("Accept", "application/json")

	var token string
	var ok bool

	for {
		res, err := client.Do(r)
		if err != nil {
			log.Printf("Failed to send request: %v", err)
			continue
		}

		defer res.Body.Close()
		var data map[string]interface{}
		if err := json.NewDecoder(res.Body).Decode(&data); err != nil {
			log.Printf("Failed to decode response: %v", err)
		} else {
			log.Println("Response data:", data)
		}

		if res.StatusCode != http.StatusOK {
			log.Printf("Non-OK HTTP status: %v", res.Status)
			continue
		}

		token, ok = data["access_token"].(string)
		if ok {
			break
		}

		time.Sleep(time.Duration(deviceFlowRes.Interval) * time.Second)
	}

	extr := executor.NewExecutor()
	github.Initialize(extr)

	params := make(map[string]interface{})

	params["token"] = token
	params["owner"] = "luisya22"
	params["repo"] = "galactic-exchange"
	params["lastIssue"] = 11

	params, err = extr.Execute("Github", "New Issue", params)
	if err != nil {
		log.Panic(err)
	}

	fmt.Println(params)

}
