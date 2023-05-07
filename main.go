package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	jwt "github.com/dgrijalva/jwt-go"
	"github.com/gorilla/mux"
	log "github.com/sirupsen/logrus"
)

type Notification struct {
	To       string            `json:"to"`
	Data     map[string]string `json:"data"`
	Notification map[string]string `json:"notification"`
}

func main() {
	router := mux.NewRouter()
	router.HandleFunc("/send", sendNotification).Methods("POST")

	srv := &http.Server{
		Handler:      router,
		Addr:         "127.0.0.1:8080",
		WriteTimeout: 15 * time.Second,
		ReadTimeout:  15 * time.Second,
	}

	// Iniciar el servidor HTTP
	log.Info("Iniciando servidor...")
	if err := srv.ListenAndServe(); err != nil {
		log.Error("Error al iniciar el servidor:", err)
	}
}

func sendNotification(w http.ResponseWriter, r *http.Request) {
	var notification Notification
	err := json.NewDecoder(r.Body).Decode(&notification)
	if err != nil {
		log.Error("Error al decodificar el cuerpo de la solicitud:", err)
		http.Error(w, "Error al decodificar el cuerpo de la solicitud", http.StatusBadRequest)
		return
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"iss": "fcm_project_id",
		"sub": "fcm_sender_id",
		"aud": "https://fcm.googleapis.com/",
		"iat": time.Now().Unix(),
		"exp": time.Now().Add(time.Hour * 24).Unix(),
	})

	key := []byte("your_fcm_server_key")
	tokenString, err := token.SignedString(key)
	if err != nil {
		log.Error("Error al firmar el token de autenticación:", err)
		http.Error(w, "Error al firmar el token de autenticación", http.StatusInternalServerError)
		return
	}

	jsonBody, err := json.Marshal(notification)
	if err != nil {
		log.Error("Error al convertir la notificación a JSON:", err)
		http.Error(w, "Error al convertir la notificación a JSON", http.StatusInternalServerError)
		return
	}

	req, err := http.NewRequest("POST", "https://fcm.googleapis.com/fcm/send", bytes.NewBuffer(jsonBody))
	if err != nil {
		log.Error("Error al crear la solicitud HTTP:", err)
		http.Error(w, "Error al crear la solicitud HTTP", http.StatusInternalServerError)
		return
	}

	req.Header.Add("Authorization", "Bearer "+tokenString)
	req.Header.Add("Content-Type", "application/json")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		log.Error("Error al enviar la notificación:", err)
		http.Error(w, "Error al enviar la notificación", http.StatusInternalServerError)
		return
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		log.Error("Error de respuesta de la API de FCM:", resp.Status)
		http.Error(w, "Error de respuesta de la API de FCM", http.StatusInternalServerError)
		return
	}

	respBody, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Error("Error al leer la respuesta de la API de FCM:", err)
		http.Error(w, "Error al leer la respuesta de la API de FCM", http.StatusInternalServerError)
		return
	}
	log.Info("Respuesta de la API de FCM:", string(respBody))

	fmt.Fprint(w, "Notificación enviada correctamente")
}
