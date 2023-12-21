package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
)

 type Message struct {
	Status string `json:"status"`
	Info   string `json:"info"`
}


func handlePage(writer http.ResponseWriter, request *http.Request) {
	writer.Header().Set("Content-Type", "application/json")

	// Accessing information from the *http.Request
	method := request.Method
	url := request.URL
	headers := request.Header
	//body := request.Body

 
	// Do something with the request information
	fmt.Printf("Method: %s\n", method)
	fmt.Printf("URL: %s\n", url)
	fmt.Printf("Headers: %v\n", headers)

	// Read the request body
	// this code cause error EOF 
	// 1. Reading the Body Twice: google bard, whats the concept ??

	// buf := make([]byte, 1024)
	// n, _ := body.Read(buf)
	// bodyContent := buf[:n]
	// fmt.Printf("Body: %s\n", bodyContent)


	var message Message
	err := json.NewDecoder(request.Body).Decode(&message)
	if err != nil {
		fmt.Println(err)
		return
	}
	err = json.NewEncoder(writer).Encode(message)
	if err != nil {
		fmt.Println(err)
		return
	}

	//fmt.Println(message)

	// Send a response
	//writer.WriteHeader(http.StatusOK)
	//writer.Write([]byte("End Point Return Values"))
}



 // creating a simple web server with an endpoint that will be secured with a JWT.
 func main() {
	http.HandleFunc("/home", handlePage)
	err := http.ListenAndServe(":8080", nil)
	if err != nil {
	   log.Println("There was an error listening on port :8080", err)
	}
 
 }