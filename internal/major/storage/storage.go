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
	PhoneNumber string `json:"phone_number,omitempty"`
}

type DataBase struct {
	c *mongo.Collection
}

func NewDataBase(client *mongo.Client) *DataBase {
	return &DataBase{c: client.Database("users").Collection("info")}
}

func (db *DataBase) Join(newUser common.UserLogPas) error {
	filter := bson.D{{"login", newUser.Login}}

	count, err := db.c.CountDocuments(context.Background(), filter)
	if count != 0 {
		log.Printf("Join error: user already exist")
		return errors.UserAlreadyExist
	}

	newUserDB := UserDB{
		Login:    newUser.Login,
		Password: newUser.Password,
	}

	insertResult, err := db.c.InsertOne(context.Background(), newUserDB)
	if err != nil {
		return err
	}
	fmt.Println("Inserted a new user: ", insertResult.InsertedID)

	return nil
}

func (db *DataBase) Update(login string, newInfo common.UserInfo) error {
	filter := bson.D{{"login", login}}
	update := bson.D{{"$set", bson.D{
		{"name", newInfo.Name},
		{"surname", newInfo.Surname},
		{"date_of_birth", newInfo.DateOfBirth},
		{"mail", newInfo.Mail},
		{"phone_number", newInfo.PhoneNumber},
	}}}

	updateResult, err := db.c.UpdateOne(context.TODO(), filter, update)
	if err != nil {
		return err
	}

	if updateResult.ModifiedCount == 0 {
		return errors.UserNotFound
	}

	return nil
}

func (db *DataBase) GetUser(userLogin string) (*common.UserLogPas, error) {
	filter := bson.D{{"login", userLogin}}

	var result UserDB
	err := db.c.FindOne(context.TODO(), filter).Decode(&result)
	if err != nil {
		return nil, err
	}

	return &common.UserLogPas{
		Login:    result.Login,
		Password: result.Password,
	}, nil
}
