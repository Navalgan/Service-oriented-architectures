package major

import (
	"Service-oriented-architectures/internal/common"
	"Service-oriented-architectures/internal/common/gen/go/task/proto"
	"Service-oriented-architectures/internal/major/storage"

	"context"
	"encoding/json"
	"errors"
	"log"
	"net/http"
	"strconv"
	"time"

	"github.com/IBM/sarama"
	"github.com/google/uuid"
	"github.com/gorilla/mux"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"golang.org/x/crypto/bcrypt"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

const (
	MinLoginLen    = 4
	MaxLoginLen    = 20
	MinPasswordLen = 7
	MaxPasswordLen = 40
)

type Service struct {
	DB                *storage.DataBase
	GRPCClient        task_v1.TaskClient
	StatisticProducer sarama.SyncProducer
	StatisticConsume  sarama.Consumer
	AnswerConsumer    sarama.PartitionConsumer
}

func NewService(ctx context.Context) (*Service, error) {
	client, err := mongo.Connect(ctx, options.Client().ApplyURI("mongodb://mongo:27017"))
	if err != nil {
		log.Fatalf("Failed to consume partition: %v", err)
	}
	db := storage.NewDataBase(client)

	log.Printf("MongoDB created")

	err = client.Ping(ctx, nil)
	if err != nil {
		log.Fatal(err)
	}

	log.Println("Connected to MongoDB")

	con, err := grpc.Dial("task:9090", grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Fatal(err)
	}

	gRPCClient := task_v1.NewTaskClient(con)

	log.Println("Connected to gRPC")

	producer, err := sarama.NewSyncProducer([]string{"kafka:9092"}, nil)
	if err != nil {
		log.Fatalf("Failed to create producer: %v", err)
	}

	log.Println("Statistic producer created")

	return &Service{
		DB:                db,
		GRPCClient:        gRPCClient,
		StatisticProducer: producer,
	}, nil
}

func (s *Service) UserJoin(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Wrong HTTP method", http.StatusBadRequest)
		return
	}

	var newUser common.UserLogPas
	err := json.NewDecoder(r.Body).Decode(&newUser)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	loginLen := len(newUser.Login)
	passwordLen := len(newUser.Password)

	if loginLen < MinLoginLen || loginLen > MaxLoginLen || passwordLen < MinPasswordLen || passwordLen > MaxPasswordLen {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	hashPassword, err := bcrypt.GenerateFromPassword([]byte(newUser.Password), bcrypt.DefaultCost)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	newUser.Password = string(hashPassword)

	newUser.UserID = uuid.New().String()
	err = s.DB.Join(newUser)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	token, err := common.NewToken(newUser.UserID, newUser.Login, time.Minute*10)
	if err != nil {
		log.Println("failed to generate token")
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	cookie := http.Cookie{
		Name:     "token",
		Value:    token,
		Path:     "/",
		MaxAge:   3600,
		HttpOnly: true,
		Secure:   true,
		SameSite: http.SameSiteLaxMode,
	}
	http.SetCookie(w, &cookie)

	w.WriteHeader(http.StatusCreated)
}

func (s *Service) UserAuth(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Wrong HTTP method", http.StatusBadRequest)
		return
	}

	login := r.URL.Query().Get("login")

	dbUser, err := s.DB.GetUser(login)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	password := r.URL.Query().Get("password")
	err = bcrypt.CompareHashAndPassword([]byte(dbUser.Password), []byte(password))
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	token, err := common.NewToken(dbUser.UserID, dbUser.Login, time.Minute*10)
	if err != nil {
		log.Println("failed to generate token")
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	cookie := http.Cookie{
		Name:     "token",
		Value:    token,
		Path:     "/",
		MaxAge:   3600,
		HttpOnly: true,
		Secure:   true,
		SameSite: http.SameSiteLaxMode,
	}
	http.SetCookie(w, &cookie)

	w.WriteHeader(http.StatusOK)
}

func (s *Service) UserUpdate(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPut {
		log.Println("Wrong http method")
		http.Error(w, "Wrong HTTP method", http.StatusBadRequest)
		return
	}

	cookie, err := r.Cookie("token")
	if err != nil {
		switch {
		case errors.Is(err, http.ErrNoCookie):
			log.Println("Cookie not found")
			http.Error(w, "cookie not found", http.StatusBadRequest)
		default:
			log.Println("server error")
			http.Error(w, "server error", http.StatusInternalServerError)
		}
		return
	}

	token := cookie.Value

	claims, err := common.VerifyToken(token)
	if err != nil {
		log.Println("Wrong token")
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	var newUserInfo common.UserInfo
	err = json.NewDecoder(r.Body).Decode(&newUserInfo)
	if err != nil {
		log.Println("Json decoder error")
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	err = s.DB.Update(claims.UserID, newUserInfo)
	if err != nil {
		log.Println("Mongo error")
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	w.WriteHeader(http.StatusOK)
}

func (s *Service) CreatePost(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		log.Println("Wrong http method")
		http.Error(w, "Wrong HTTP method", http.StatusBadRequest)
		return
	}

	cookie, err := r.Cookie("token")
	if err != nil {
		switch {
		case errors.Is(err, http.ErrNoCookie):
			log.Println("No cookie")
			http.Error(w, "cookie not found", http.StatusBadRequest)
		default:
			log.Println(err)
			http.Error(w, "server error", http.StatusInternalServerError)
		}
		return
	}

	token := cookie.Value

	claims, err := common.VerifyToken(token)
	if err != nil {
		log.Println("Wrong token")
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	var postText common.PostText
	err = json.NewDecoder(r.Body).Decode(&postText)
	if err != nil || len(postText.Text) == 0 {
		log.Println("Empty text")
		http.Error(w, "Empty text", http.StatusBadRequest)
		return
	}

	resp, err := s.GRPCClient.CreatePost(context.Background(), &task_v1.PostRequest{UserID: claims.UserID, Text: postText.Text})
	if err != nil {
		log.Println("Cassandra error")
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	err = json.NewEncoder(w).Encode(resp)
	if err != nil {
		log.Println("Json encoder error")
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	w.WriteHeader(http.StatusOK)
}

func (s *Service) GetPostByID(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Wrong HTTP method", http.StatusBadRequest)
		return
	}

	vars := mux.Vars(r)
	postID, ok := vars["postId"]
	if !ok {
		http.Error(w, "Need post's id", http.StatusNotFound)
		return
	}

	resp, err := s.GRPCClient.GetPostByID(context.Background(), &task_v1.PostIDRequest{PostID: postID})
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	err = json.NewEncoder(w).Encode(resp)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	w.WriteHeader(http.StatusOK)
}

func (s *Service) GetPostsByUser(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Wrong HTTP method", http.StatusBadRequest)
		return
	}

	vars := mux.Vars(r)
	login, ok := vars["login"]
	if !ok {
		http.Error(w, "Need login of post's author", http.StatusNotFound)
		return
	}

	user, err := s.DB.GetUser(login)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	resp, err := s.GRPCClient.GetPostsByUser(context.Background(), &task_v1.UserRequest{UserID: user.UserID})
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	err = json.NewEncoder(w).Encode(resp)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	w.WriteHeader(http.StatusOK)
}

func (s *Service) LikePost(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPut {
		http.Error(w, "Wrong HTTP method", http.StatusBadRequest)
		return
	}

	cookie, err := r.Cookie("token")
	if err != nil {
		switch {
		case errors.Is(err, http.ErrNoCookie):
			log.Println("Cookie not found")
			http.Error(w, "cookie not found", http.StatusBadRequest)
		default:
			log.Println("server error")
			http.Error(w, "server error", http.StatusInternalServerError)
		}
		return
	}

	token := cookie.Value

	claims, err := common.VerifyToken(token)
	if err != nil {
		log.Println("Wrong token")
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	vars := mux.Vars(r)
	postID, ok := vars["postId"]
	if !ok {
		http.Error(w, "Need post's id", http.StatusNotFound)
		return
	}

	requestID := uuid.New().String()

	msg := &sarama.ProducerMessage{
		Topic: "likes",
		Key:   sarama.StringEncoder(requestID),
		Value: sarama.StringEncoder(postID + "," + claims.UserID + "," + strconv.FormatInt(time.Now().Unix(), 10)),
	}

	_, _, err = s.StatisticProducer.SendMessage(msg)
	if err != nil {
		log.Printf("Failed to send message to Kafka: %v", err)
		http.Error(w, "cookie not found", http.StatusBadRequest)
		return
	}

	w.WriteHeader(http.StatusOK)
}

func (s *Service) ViewPost(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPut {
		http.Error(w, "Wrong HTTP method", http.StatusBadRequest)
		return
	}

	cookie, err := r.Cookie("token")
	if err != nil {
		switch {
		case errors.Is(err, http.ErrNoCookie):
			log.Println("Cookie not found")
			http.Error(w, "cookie not found", http.StatusBadRequest)
		default:
			log.Println("server error")
			http.Error(w, "server error", http.StatusInternalServerError)
		}
		return
	}

	token := cookie.Value

	claims, err := common.VerifyToken(token)
	if err != nil {
		log.Println("Wrong token")
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	vars := mux.Vars(r)
	postID, ok := vars["postId"]
	if !ok {
		http.Error(w, "Need post's id", http.StatusNotFound)
		return
	}

	requestID := uuid.New().String()

	msg := &sarama.ProducerMessage{
		Topic: "views",
		Key:   sarama.StringEncoder(requestID),
		Value: sarama.StringEncoder(postID + "," + claims.UserID + "," + strconv.FormatInt(time.Now().Unix(), 10)),
	}

	_, _, err = s.StatisticProducer.SendMessage(msg)
	if err != nil {
		log.Printf("Failed to send message to Kafka: %v", err)
		http.Error(w, "cookie not found", http.StatusBadRequest)
		return
	}

	w.WriteHeader(http.StatusOK)
}

func (s *Service) UpdatePost(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPut {
		http.Error(w, "Wrong HTTP method", http.StatusBadRequest)
		return
	}

	cookie, err := r.Cookie("token")
	if err != nil {
		switch {
		case errors.Is(err, http.ErrNoCookie):
			http.Error(w, "cookie not found", http.StatusBadRequest)
		default:
			log.Println(err)
			http.Error(w, "server error", http.StatusInternalServerError)
		}
		return
	}

	token := cookie.Value

	claims, err := common.VerifyToken(token)
	if err != nil {
		log.Println("Wrong token")
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	vars := mux.Vars(r)
	postID, ok := vars["postId"]
	if !ok {
		http.Error(w, "Need post's id", http.StatusNotFound)
		return
	}

	var updatePost common.PostText
	err = json.NewDecoder(r.Body).Decode(&updatePost)
	if err != nil || len(updatePost.Text) == 0 {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	_, err = s.GRPCClient.UpdatePost(context.Background(), &task_v1.UpdatePostRequest{PostID: postID, UserID: claims.UserID, Text: updatePost.Text})
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	w.WriteHeader(http.StatusOK)
}

func (s *Service) DeletePost(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodDelete {
		http.Error(w, "Wrong HTTP method", http.StatusBadRequest)
		return
	}

	cookie, err := r.Cookie("token")
	if err != nil {
		switch {
		case errors.Is(err, http.ErrNoCookie):
			http.Error(w, "cookie not found", http.StatusBadRequest)
		default:
			log.Println(err)
			http.Error(w, "server error", http.StatusInternalServerError)
		}
		return
	}

	token := cookie.Value

	claims, err := common.VerifyToken(token)
	if err != nil {
		log.Println("Wrong token")
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	vars := mux.Vars(r)
	postID, ok := vars["postId"]
	if !ok {
		http.Error(w, "Need post's id", http.StatusNotFound)
		return
	}

	_, err = s.GRPCClient.DeletePost(context.Background(), &task_v1.DeletePostRequest{PostID: postID, UserID: claims.UserID})
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	w.WriteHeader(http.StatusOK)
}
