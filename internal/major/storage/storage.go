package storage

import (
	"Service-oriented-architectures/internal/common"
	"Service-oriented-architectures/internal/errors"
	"context"
	"fmt"
	"log"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

type UserDB struct {
	// The user's id
	UserID string `bson:"user_id,omitempty"`
	// Must be unique in the system
	Login string `bson:"login,omitempty"`
	// Can't be empty and will not change
	Password string `bson:"password,omitempty"`
	// The user's name
	Name string `bson:"name,omitempty"`
	// The user's surname
	Surname string `bson:"surname,omitempty"`
	// The user's date of birth
	DateOfBirth string `bson:"date_of_birth,omitempty"`
	// The user's mail
	Mail string `bson:"mail,omitempty"`
	// The user's phone number
	PhoneNumber string `bson:"phone_number,omitempty"`
}

type DataBase struct {
	c *mongo.Collection
}

func NewDataBase(client *mongo.Client) *DataBase {
	return &DataBase{c: client.Database("users").Collection("info")}
}

func (db *DataBase) Join(newUser common.UserLogPas) error {
	filter := bson.D{{"login", newUser.Login}}

	ctx := context.Background()

	count, err := db.c.CountDocuments(ctx, filter)
	if count != 0 {
		log.Printf("Join error: user already exist")
		return errors.UserAlreadyExist
	}

	newUserDB := UserDB{
		UserID:   newUser.UserID,
		Login:    newUser.Login,
		Password: newUser.Password,
	}

	insertResult, err := db.c.InsertOne(ctx, newUserDB)
	if err != nil {
		return err
	}
	fmt.Println("Inserted a new user: ", insertResult.InsertedID)

	return nil
}

func (db *DataBase) Update(userID string, newInfo common.UserInfo) error {
	filter := bson.D{{"user_id", userID}}
	update := bson.D{{"$set", bson.D{
		{"name", newInfo.Name},
		{"surname", newInfo.Surname},
		{"date_of_birth", newInfo.DateOfBirth},
		{"mail", newInfo.Mail},
		{"phone_number", newInfo.PhoneNumber},
	}}}

	ctx := context.Background()

	updateResult, err := db.c.UpdateOne(ctx, filter, update)
	if err != nil {
		log.Println("Can not update user")
		return err
	}

	if updateResult.ModifiedCount == 0 {
		log.Println("User not found or the data already meets the requirements")
		return errors.UserNotFound
	}

	return nil
}

func (db *DataBase) GetUser(userLogin string) (*common.UserLogPas, error) {
	filter := bson.D{{"login", userLogin}}

	var result UserDB
	err := db.c.FindOne(context.TODO(), filter).Decode(&result)
	if err != nil {
		log.Println(err)
		return nil, err
	}

	return &common.UserLogPas{
		UserID:   result.UserID,
		Login:    result.Login,
		Password: result.Password,
	}, nil
}

func (db *DataBase) GetUserByID(userID string) (*common.UserLogPas, error) {
	filter := bson.D{{"user_id", userID}}

	var result UserDB
	err := db.c.FindOne(context.TODO(), filter).Decode(&result)
	if err != nil {
		log.Println("error on get user by id: " + userID)
		log.Println(err.Error())
		return nil, err
	}

	return &common.UserLogPas{
		UserID:   result.UserID,
		Login:    result.Login,
		Password: result.Password,
	}, nil
}
