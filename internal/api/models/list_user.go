package models

type User struct {
	ID          int     `json:"id"`
	Name        string  `json:"name"`
	Surname     string  `json:"surname"`
	Patronymic  *string `json:"patronymic"`
	Age         *int    `json:"age"`
	Gender      *string `json:"gender"`
	Nationality *string `json:"nationality"`
}

type ListUsersResponse struct {
	Users []User `json:"users"`
	Page  int    `json:"page"`
	Pages int    `json:"pages"`
}
