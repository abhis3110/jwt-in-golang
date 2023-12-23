package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"time"
	"github.com/golang-jwt/jwt"
)

  type Message struct {
	Status string `json:"status"`
	Info   string `json:"info"`
}

var secretKey =[]byte("alphadelta")

func generateJWT()(string, error){

	// create new token
	token := jwt.New(jwt.SigningMethodEdDSA)

	// use claim method of token to modify jwt
	claims := token.Claims.(jwt.MapClaims)
    claims["exp"] = time.Now().Add(10 * time.Minute)
    claims["authorized"] = true
    claims["user"] = "username"

	tokenString, err := token.SignedString(secretKey)
    if err != nil {
		return "Signing Error", err
	}
    return tokenString, nil
}


func verifyJWT(endpointHandler func(writer http.ResponseWriter, request *http.Request)) http.HandlerFunc {
	return http.HandlerFunc(func(writer http.ResponseWriter, request *http.Request) {
		if request.Header["Token"] != nil {
			token, err := jwt.Parse(request.Header["Token"][0], func(token *jwt.Token) (interface{}, error) {
				_, ok := token.Method.(*jwt.SigningMethodECDSA)
				if !ok {
					writer.WriteHeader(http.StatusUnauthorized)
					_, err := writer.Write([]byte("You're Unauthorized"))
					if err != nil {
						return nil, err
					}
				}
				return "", nil

			})
			// parsing errors result
			if err != nil {
				writer.WriteHeader(http.StatusUnauthorized)
				_, err2 := writer.Write([]byte("You're Unauthorized due to error parsing the JWT"))
				if err2 != nil {
					return
				}

			}
			// if there's a token
			if token.Valid {
				endpointHandler(writer, request)
			} else {
				writer.WriteHeader(http.StatusUnauthorized)
				_, err := writer.Write([]byte("You're Unauthorized due to invalid token"))
				if err != nil {
					return
				}
			}
		} else {
			writer.WriteHeader(http.StatusUnauthorized)
			_, err := writer.Write([]byte("You're Unauthorized due to No token in the header"))
			if err != nil {
				return
			}
		}
		// response for if there's no token header
	})
}

func extractClaims(_ http.ResponseWriter, request *http.Request) (string, error) {
	if request.Header["Token"] != nil {
		tokenString := request.Header["Token"][0]
		token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {

			if _, ok := token.Method.(*jwt.SigningMethodECDSA); !ok {
				return nil, fmt.Errorf("there's an error with the signing method")
			}
			return secretKey, nil
		})
		if err != nil {
			return "Error Parsing Token: ", err
		}
		claims, ok := token.Claims.(jwt.MapClaims)
		if ok && token.Valid {
			username := claims["username"].(string)
			return username, nil
		}
	}

	return "unable to extract claims", nil
}

func authPage() {
	token, _ := generateJWT()
	client := &http.Client{}
	req, _ := http.NewRequest("GET", "http://localhost:8080/", nil)
	req.Header.Set("Token", token)
	_, _ = client.Do(req)
}



func handlePage(writer http.ResponseWriter, request *http.Request) {
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



	_, err := generateJWT()
	if err != nil {
		log.Fatalln("Error generating JWT", err)
	}
	writer.Header().Set("Content-Type", "application/json")
	type_ := "application/json"
	writer.Header().Set("Content-Type", type_)

	var message Message
	err = json.NewDecoder(request.Body).Decode(&message)
	if err != nil {
		fmt.Println(err)
		return
	}
	err = json.NewEncoder(writer).Encode(message)
	if err != nil {
		fmt.Println(err)
		return
	}

	// Send a response
	//writer.WriteHeader(http.StatusOK)
	//writer.Write([]byte("End Point Return Values"))
}



 // creating a simple web server with an endpoint that will be secured with a JWT.
 func main() {
	http.HandleFunc("/home", verifyJWT(handlePage))
	err := http.ListenAndServe(":8080", nil)
	if err != nil {
	   log.Println("There was an error listening on port :8080", err)
	}
 }