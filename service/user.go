package service

import (
	"errors"
	"fmt"
	"log"

	"github.com/dieklingel/core/internal/api"
	"github.com/ostafen/clover/v2"
	"github.com/ostafen/clover/v2/document"
	"github.com/ostafen/clover/v2/query"
	"golang.org/x/crypto/bcrypt"
)

type user struct {
	username string
	password string
	role     string
}

func (user *user) Username() string {
	return user.username
}

func (user *user) Password() string {
	return user.password
}

func (user *user) Role() string {
	return user.role
}

type UserService struct{}

func NewUserService() *UserService {
	return &UserService{}
}

func (service *UserService) Create(username string, password string, role string) (api.User, error) {
	db, err := clover.Open("")
	if err != nil {
		log.Printf("an error occured while open the database: %s", err.Error())
		return nil, errors.New("an error occoured while opening the database")
	}
	defer db.Close()

	if succes, _ := db.HasCollection("users"); !succes {
		db.CreateCollection("users")
	}

	if exists, _ := db.Exists(query.NewQuery("users").Where(query.Field("username").Eq(username))); exists {
		return nil, errors.New("a user with this username already exists")
	}

	passwordHash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return nil, fmt.Errorf("cannot process the given password: %s", err.Error())
	}

	doc := document.NewDocument()
	doc.Set("username", username)
	doc.Set("password", string(passwordHash))
	doc.Set("role", role)

	if err := db.Insert("users", doc); err != nil {
		log.Println("an error occured while store the newly created user")
		return nil, err
	}

	return service.GetByUsername(username), nil
}

func (service *UserService) GetByUsername(username string) api.User {
	db, err := clover.Open("")
	if err != nil {
		log.Printf("an error occured while open the database: %s", err.Error())
		return nil
	}
	defer db.Close()

	if succes, _ := db.HasCollection("users"); !succes {
		return nil
	}

	doc, err := db.FindFirst(query.NewQuery("users").Where(query.Field("username").Eq(username)))
	if err != nil {
		return nil
	}

	if doc == nil {
		return nil
	}

	return &user{
		username: doc.Get("username").(string),
		password: doc.Get("password").(string),
		role:     doc.Get("role").(string),
	}
}
