package main

import (
	statistic_v1 "Service-oriented-architectures/internal/common/gen/go/statistic/proto"
	task_v1 "Service-oriented-architectures/internal/common/gen/go/task/proto"
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"math/rand"
	"net/http"
	"net/url"
	"strconv"
	"strings"
	"sync"
	"time"
)

type Jar struct {
	lk      sync.Mutex
	cookies map[string][]*http.Cookie
}

func NewJar() *Jar {
	jar := new(Jar)
	jar.cookies = make(map[string][]*http.Cookie)
	return jar
}

func (jar *Jar) SetCookies(u *url.URL, cookies []*http.Cookie) {
	jar.lk.Lock()
	jar.cookies[u.Host] = cookies
	jar.lk.Unlock()
}

func (jar *Jar) Cookies(u *url.URL) []*http.Cookie {
	return jar.cookies[u.Host]
}

type User struct {
	Login string
	ID    string
	Posts map[string]*task_v1.PostResponse
}

func GetPostsStat(client *http.Client, users []User) bool {
	for i, user := range users {
		for _, post := range user.Posts {
			resp, err := client.Get("http://localhost:8080/post/" + post.PostID + "/stat")
			if err != nil {
				return false
			}
			defer resp.Body.Close()

			postStat := statistic_v1.PostStatResponse{}
			if err = json.NewDecoder(resp.Body).Decode(&postStat); err != nil {
				log.Fatal(err)
			}

			if i == 0 {
				if postStat.Likes != 3 {
					return false
				}
			} else if i == 1 {
				if postStat.Likes != 2 {
					return false
				}
			} else if i == 2 {
				if postStat.Likes != 0 {
					return false
				}
			}
		}
	}

	return true
}

func GetPostsByLogin(client *http.Client, users []User) bool {
	for _, user := range users {
		resp, err := client.Get("http://localhost:8080/posts/" + user.Login)
		if err != nil {
			return false
		}
		defer resp.Body.Close()

		allPosts := task_v1.PostsResponse{}

		if err = json.NewDecoder(resp.Body).Decode(&allPosts); err != nil {
			log.Fatal(err)
		}

		if len(allPosts.Posts) != len(user.Posts) {
			return false
		}

		for _, post := range allPosts.Posts {
			userPost, ok := user.Posts[post.PostID]
			if !ok {
				return false
			}

			if strings.Compare(userPost.PostID, post.PostID) != 0 || strings.Compare(userPost.AuthorID, post.AuthorID) != 0 ||
				strings.Compare(userPost.Text, post.Text) != 0 || userPost.Date != post.Date {
				return false
			}
		}
	}

	return true
}

func main() {
	jar := NewJar()
	client := http.Client{Jar: jar}

	users := make([]User, 0)

	randTest := int(rand.Uint32() % 100000)

	for i := 0; i < 3; i++ {
		login := "MyTestUser" + strconv.Itoa(randTest+i)

		users = append(users, User{Login: login, Posts: make(map[string]*task_v1.PostResponse)})

		createUserBody, _ := json.Marshal(map[string]string{
			"login":    login,
			"password": "Qwerty12345",
		})
		requestCreateUserBody := bytes.NewBuffer(createUserBody)

		respCreateUser, err := client.Post("http://localhost:8080/user/join", "application/json", requestCreateUserBody)
		if err != nil {
			log.Fatalln(err)
		}
		defer respCreateUser.Body.Close()

		bodyCreateUser, err := io.ReadAll(respCreateUser.Body)
		fmt.Println(string(bodyCreateUser))

		for j := 0; j < 3; j++ {
			createPostBody, _ := json.Marshal(map[string]string{
				"text": "My test",
			})
			requestCreatePostBody := bytes.NewBuffer(createPostBody)

			respCreatePost, err := client.Post("http://localhost:8080/post/create", "application/json", requestCreatePostBody)
			if err != nil {
				log.Fatalln(err)
			}
			defer respCreatePost.Body.Close()

			var jsonPostResp task_v1.PostResponse
			if err = json.NewDecoder(respCreatePost.Body).Decode(&jsonPostResp); err != nil {
				log.Fatal(err)
			}

			users[i].ID = jsonPostResp.AuthorID

			users[i].Posts[jsonPostResp.PostID] = &jsonPostResp
		}
	}

	for i, user := range users {
		_, err := client.Get("http://localhost:8080/user/auth?login=" + user.Login + "&password=Qwerty12345")
		if err != nil {
			log.Fatal(err)
		}

		for j := 0; j < 2; j++ {
			if i == 2 && j == 1 {
				break
			}
			for _, post := range users[j].Posts {
				req, err := http.NewRequest("PUT", "http://localhost:8080/post/"+post.PostID+"/like", nil)
				if err != nil {
					log.Fatalln(err)
				}

				_, err = client.Do(req)
				if err != nil {
					log.Fatal(err)
				}
			}
		}
	}

	time.Sleep(time.Second * 5)

	fmt.Println("Running GetPostsStat")
	if !GetPostsStat(&client, users) {
		fmt.Println("FAIL")
		log.Fatal("FAIL")
	}
	fmt.Println("OK")

	fmt.Println("Running GetPostsByLogin")
	if !GetPostsByLogin(&client, users) {
		fmt.Println("FAIL")
		log.Fatal("FAIL")
	}
	fmt.Println("OK")
}
