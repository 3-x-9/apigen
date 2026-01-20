package models

type User struct {
	Password string `json:"password"`
	Phone string `json:"phone"`
	UserStatus int `json:"userStatus"`
	Username string `json:"username"`
	Email string `json:"email"`
	FirstName string `json:"firstName"`
	Id int `json:"id"`
	LastName string `json:"lastName"`
}

type Tag struct {
	Name string `json:"name"`
	Id int `json:"id"`
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

type Order struct {
	Id int `json:"id"`
	PetId int `json:"petId"`
	Quantity int `json:"quantity"`
	ShipDate string `json:"shipDate"`
	Status string `json:"status"`
	Complete bool `json:"complete"`
}

type Category struct {
	Name string `json:"name"`
	Id int `json:"id"`
}

