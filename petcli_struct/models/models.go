package models

type Order struct {
	PetId int `json:"petId"`
	Quantity int `json:"quantity"`
	ShipDate string `json:"shipDate"`
	Status string `json:"status"`
	Complete bool `json:"complete"`
	Id int `json:"id"`
}

type Category struct {
	Id int `json:"id"`
	Name string `json:"name"`
}

type User struct {
	UserStatus int `json:"userStatus"`
	Username string `json:"username"`
	Email string `json:"email"`
	FirstName string `json:"firstName"`
	Id int `json:"id"`
	LastName string `json:"lastName"`
	Password string `json:"password"`
	Phone string `json:"phone"`
}

type Tag struct {
	Id int `json:"id"`
	Name string `json:"name"`
}

type Pet struct {
	Id int `json:"id"`
	Name string `json:"name"`
	PhotoUrls []string `json:"photoUrls"`
	Status string `json:"status"`
	Tags []Tag `json:"tags"`
	Category Category `json:"category"`
}

type ApiResponse struct {
	Type string `json:"type"`
	Code int `json:"code"`
	Message string `json:"message"`
}

