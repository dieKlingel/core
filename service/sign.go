package service

import (
	"log"

	"github.com/dieklingel/core/internal/api"
	"github.com/ostafen/clover/v2"
	"github.com/ostafen/clover/v2/document"
	"github.com/ostafen/clover/v2/query"
)

type sign struct {
	id     string
	name   string
	script string
}

func (sign *sign) Id() string {
	return sign.id
}

func (sign *sign) Name() string {
	return sign.name
}

func (sign *sign) Script() string {
	return sign.script
}

type SignService struct{}

func NewSignService() *SignService {
	return &SignService{}
}

func (service *SignService) List() []api.Sign {
	db, err := clover.Open("")
	if err != nil {
		log.Printf("an error occured while open the database: %s", err.Error())
		return make([]api.Sign, 0)
	}
	defer db.Close()

	if succes, _ := db.HasCollection("signs"); !succes {
		return make([]api.Sign, 0)
	}

	docs, err := db.FindAll(query.NewQuery("signs"))
	if err != nil {
		log.Printf("an error occured while fetching signs: %s", err.Error())
		return make([]api.Sign, 0)
	}

	signs := make([]api.Sign, len(docs))
	for index, doc := range docs {
		id := doc.ObjectId()
		name := doc.Get("name").(string)
		script := doc.Get("script").(string)

		signs[index] = &sign{
			id:     id,
			name:   name,
			script: script,
		}
	}

	return signs
}

func (service *SignService) Add(name string, script string) (string, error) {
	db, err := clover.Open("")
	if err != nil {
		log.Printf("an error occured while open the database: %s", err.Error())
		return "", err
	}
	defer db.Close()

	if succes, _ := db.HasCollection("signs"); !succes {
		db.CreateCollection("signs")
	}

	doc := document.NewDocument()
	doc.Set("name", name)
	doc.Set("script", script)

	if err := db.Insert("signs", doc); err != nil {
		log.Printf("an error occured while inserting into the signs database: %s", err.Error())
		return "", err
	}

	return doc.ObjectId(), nil
}

func (service *SignService) Remove(id string) error {
	db, err := clover.Open("")
	if err != nil {
		log.Printf("an error occured while open the database: %s", err.Error())
		return err
	}
	defer db.Close()

	if succes, _ := db.HasCollection("signs"); !succes {
		return nil
	}

	return db.DeleteById("signs", id)
}

func (service *SignService) GetById(id string) api.Sign {
	db, err := clover.Open("")
	if err != nil {
		log.Printf("an error occured while open the database: %s", err.Error())
		return nil
	}
	defer db.Close()

	if succes, _ := db.HasCollection("signs"); !succes {
		return nil
	}

	doc, err := db.FindById("signs", id)
	if err != nil {
		log.Printf("an error occured while fetching a sign %s", err.Error())
		return nil
	}

	if doc == nil {
		return nil
	}

	sign := &sign{
		id:     doc.ObjectId(),
		name:   doc.Get("name").(string),
		script: doc.Get("script").(string),
	}
	return sign
}
