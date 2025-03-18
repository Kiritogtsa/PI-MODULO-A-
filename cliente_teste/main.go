package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
)

var arduino_response = make(chan string, 10) // Canal com buffer para evitar bloqueios

type Data struct {
	Message string `json:"reposta"`
}

func read_input() {
	for {
		fmt.Println("digite 1 para sim,2 para nao,3 para nao sei")
		var input string
		fmt.Scanln(&input)
		arduino_response <- input
	}
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
	go read_input()
	go client()

	select {}
}
