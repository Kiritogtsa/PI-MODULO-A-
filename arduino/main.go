package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	"github.com/sirupsen/logrus"
	"go.bug.st/serial"
)

var arduino_response = make(chan string, 10) // Canal com buffer para evitar bloqueios

func openArduino() (serial.Port, error) {
	ports, err := serial.GetPortsList()
	var port serial.Port
	if err != nil {
		log.Fatal(err)
	}
	if len(ports) == 0 {
		log.Fatal("No Serial ports found!")
	}
	for _, port := range ports {
		fmt.Printf("Found port: %v\n", port)
	}
	mode := &serial.Mode{BaudRate: 9600}
	port, err = serial.Open(ports[0], mode)
	// for _, a := range ports {
	// 	if a == "/dev/ttyUSB" {
	// 		port, err = serial.Open(a, mode)
	// 	}
	// }
	if err != nil {
		return nil, err
	}
	return port, nil
}

func readArduino(port serial.Port) {
	buff := make([]byte, 100)
	for {
		n, err := port.Read(buff)
		if err != nil {
			logrus.Error("Error reading from Arduino:", err)
			break
		}
		if n == 0 {
			fmt.Println("\nEOF")
			break
		}
		message := string(buff[:n])
		fmt.Println("Received from Arduino:", message)
		arduino_response <- message
	}
}

type Data struct {
	Message string `json:"reposta"`
}

func client() {
	for message := range arduino_response {
		data := Data{
			Message: message,
		}

		jsonData, err := json.Marshal(data)
		if err != nil {
			fmt.Println("Error encoding JSON:", err)
			continue
		}

		url := "http://localhost:8080/arduino"
		req, err := http.NewRequest("POST", url, bytes.NewBuffer(jsonData)) // ALTERADO PARA POST
		if err != nil {
			fmt.Println("Error creating request:", err)
			continue
		}
		req.Header.Set("Content-Type", "application/json")

		client := &http.Client{}
		resp, err := client.Do(req)
		if err != nil {
			fmt.Println("Error sending request:", err)
			continue
		}
		defer resp.Body.Close()

		fmt.Println("Server response:", resp.Status)
	}
}

func main() {
	port, err := openArduino()
	if err != nil {
		log.Fatal(err)
	}

	go readArduino(port)
	go client()

	select {}
}
