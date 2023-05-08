package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
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

	secretKey, err := ioutil.ReadFile("fcm_secret.txt")
	if err != nil {
		fmt.Println("Error reading secret key file:", err)
		return
	}

	data := map[string]string{
		"message": "Test message",
	}

	notification := &FcmNotification{
		Title: "Notification title",
		Body:  "Notification body",
		Lang:  "en", // Set notification language to English.
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

	request, err := http.NewRequest("POST", "https://fcm.googleapis.com/fcm/send", bytes.NewBuffer(jsonMessage))
	if err != nil {
		fmt.Println("Error creating request:", err)
		return
	}

	request.Header.Set("Authorization", "key="+string(secretKey))
	request.Header.Set("Content-Type", "application/json")

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
