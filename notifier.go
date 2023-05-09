package main

import (
	"bytes"
	"encoding/json"
	"flag"
	"fmt"
	"io/ioutil"
	"net/http"
	"strings"

	"golang.org/x/oauth2"
	"golang.org/x/oauth2/google"
	"google.golang.org/api/option"
)

type FcmMessage struct {
	To           string            `json:"to"`
	Data         map[string]string `json:"data"`
	Notification *FcmNotification  `json:"notification,omitempty"`
}

type FcmNotification struct {
	Title string `json:"title"`
	Body  string `json:"body"`
	Lang  string `json:"lang,omitempty"`
}

func main() {

	titleFlag := flag.String("title", "", "Title of the notification")
	bodyFlag := flag.String("body", "", "Body of the notification")
	flag.Parse()

	title := *titleFlag
	body := *bodyFlag


	configData, err := ioutil.ReadFile("config.txt")
	if err != nil {
		fmt.Println("Error reading config file:", err)
		return
	}

	configLines := strings.Split(string(configData), "\n")
	if len(configLines) < 2 {
		fmt.Println("Invalid config file")
		return
	}


	creds, err := google.CredentialsFromJSON(
		oauth2.NoContext, []byte(configData), "https://www.googleapis.com/auth/cloud-platform",
	)
	if err != nil {
		fmt.Println("Error creating credentials:", err)
		return
	}


	data := map[string]string{
		"message": body,
	}

	notification := &FcmNotification{
		Title: title,
		Body:  body,
	}

	message := &FcmMessage{
		To:           "/topics/all",
		Data:         data,
		Notification: notification,
	}

	jsonMessage, err := json.Marshal(message)
	if err != nil {
		fmt.Println("Error encoding message:", err)
		return
	}

	// Send FCM message
	projectID := strings.TrimSpace(configLines[0])
	url := fmt.Sprintf("https://fcm.googleapis.com/v1/projects/%s/messages:send", projectID)

	request, err := http.NewRequest("POST", url, bytes.NewBuffer(jsonMessage))
	if err != nil {
		fmt.Println("Error creating request:", err)
		return
	}

	request.Header.Set("Authorization", "Bearer "+creds.TokenSource.Token().AccessToken)
	request.Header.Set("Content-Type", "application/json")
	request.Header.Set("X-GOOG-API-FORMAT-VERSION", "2")

	client := &http.Client{}
	response, err := client.Do(request)
	if err != nil {
		fmt.Println("Error sending request:", err)
		return
	}
	defer response.Body.Close()

	responseBody, err := ioutil.ReadAll(response.Body)
	if err != nil {
		fmt.Println("Error reading response:", err)
		return
	}

	if response.StatusCode == http.StatusOK {
		fmt.Println("Notification sent successfully:", string(responseBody))
	} else {
		fmt.Println("Error sending notification:", string(responseBody))
	}
}