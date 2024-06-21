package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
)

func GetTopUsers() {

}

func GetPostsByLogin() {

}

func main() {
	createUserBody, _ := json.Marshal(map[string]string{
		"login":    "Toby12",
		"password": "Qwerty12345",
	})
	requestCreateUserBody := bytes.NewBuffer(createUserBody)

	respCreateUser, err := http.Post("http://localhost:8080/user/join", "application/json", requestCreateUserBody)
	if err != nil {
		log.Fatalln(err)
	}
	defer respCreateUser.Body.Close()

	bodyCreateUser, err := io.ReadAll(respCreateUser.Body)
	fmt.Println(string(bodyCreateUser))

	createPostBody, _ := json.Marshal(map[string]string{
		"text": "My test",
	})
	requestCreatePostBody := bytes.NewBuffer(createPostBody)

	respCreatePost, err := http.Post("http://localhost:8080/post/create", "application/json", requestCreatePostBody)
	if err != nil {
		log.Fatalln(err)
	}
	defer respCreatePost.Body.Close()

	bodyCreatePost, err := io.ReadAll(respCreatePost.Body)
	fmt.Println(string(bodyCreatePost))
}
