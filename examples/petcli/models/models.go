package models

type Order struct {
	ShipDate string `json:"shipDate"`
	Status string `json:"status"`
	Complete bool `json:"complete"`
	Id int `json:"id"`
	PetId int `json:"petId"`
	Quantity int `json:"quantity"`
}

type Category struct {
	Id int `json:"id"`
	Name string `json:"name"`
}

type User struct {
	Phone string `json:"phone"`
	UserStatus int `json:"userStatus"`
	Username string `json:"username"`
	Email string `json:"email"`
	FirstName string `json:"firstName"`
	Id int `json:"id"`
	LastName string `json:"lastName"`
	Password string `json:"password"`
}

type Tag struct {
	Id int `json:"id"`
	Name string `json:"name"`
}

type Pet struct {
	Category Category `json:"category"`
	Id int `json:"id"`
	Name string `json:"name"`
	PhotoUrls []string `json:"photoUrls"`
	Status string `json:"status"`
	Tags []Tag `json:"tags"`
}

type ApiResponse struct {
	Code int `json:"code"`
	Message string `json:"message"`
	Type string `json:"type"`
}

