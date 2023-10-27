package models

type CreateUserRequest struct {
	Name       string `json:"name"`
	Surname    string `json:"surname"`
	Patronymic string `json:"patronymic"`
}

type CreateUserResponse struct {
	ID int `json:"id"`
}
