package models

import (
	"time"

	"github.com/google/uuid"
)

type (
	Role string
)

type MainDocument struct {
	ID           uuid.UUID      `json:"id"`
	MainTitle    string         `json:"title"`
	Status       uint8          `json:"status"`
	Position     uint8          `json:"position"`
	Date         time.Time      `json:"date"`
	SubDocuments []*SubDocument `json:"sub_documents"`
}


type SubDocument struct {
	ID  	uuid.UUID	`json:"id"`
	MainID 	uuid.UUID	`json:"main_id"`
	SubTitle string 	`json:"sub_title"`
	ProductCode string 	`json:"product_code"`
	SubMessage string 	`json:"sub_message"`//Ürün Özellikleri
	Asset [][]byte		`json:"asset"`
	Position uint8		`json:"position"`
	Status	uint8 		`json:"status"`
	Date 	time.Time 	`json:"date"`
	ContentDocuments []*ContentDocument `json:"content_documents"`
}

type ContentDocument struct {
	ID       uuid.UUID `json:"id"`
	SubID    uuid.UUID `json:"sub_id"`
	ColText  string    `json:"about_collection"`
	JewCare  string    `json:"jewellery_care"`
	Position uint8     `json:"position"`
	Status   uint8     `json:"status"`
	Date     time.Time `json:"date"`
}

type AggregateDocument struct {
	MainDocument     *MainDocument    `json:"main_document"`
	SubDocuments     *SubDocument     `json:"sub_document"`
	ContentDocuments *ContentDocument `json:"content_document"`
}

type User struct {
	UserID        string    `json:"user_id"`
	Name      string    `json:"first_name"`
	Surname   string    `json:"last_name"`
	Email     string    `json:"email"`
	Phone 	  string 	`json:"phone"`
	Password  string    `json:"password"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

