package repository

import (
	"database/sql"
	"fmt"
	"log/slog"
)

const (
	limit = 10 // limit users per page
)

type Users struct {
	db *sql.DB
}

func New(con *sql.DB) *Users {
	return &Users{
		db: con,
	}
}

type User struct {
	ID          int
	Name        string
	Surname     string
	Patronymic  *string
	Age         *int
	Gender      *string
	Nationality *string
}

func (s *Users) Delete(id int) error {
	result, err := s.db.Exec(`DELETE FROM users WHERE id = $1`, id)
	if err != nil {
		return err
	}
	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return err
	}
	if rowsAffected == 0 {
		return sql.ErrNoRows
	}
	return nil
}

func (s *Users) Get(id int) (*User, error) {
	var user User

	err := s.db.
		QueryRow(`SELECT id, name, surname, patronymic, age, gender, nationality FROM users WHERE id = $1`, id).
		Scan(&user.ID, &user.Name, &user.Surname, &user.Patronymic, &user.Age, &user.Gender, &user.Nationality)
	if err != nil {
		return nil, err
	}

	return &user, nil
}

func (s *Users) Insert(name string, surname string, patronymic string, age int, gender string, nationality string) (int, error) {
	var id int

	err := s.db.
		QueryRow(`
			INSERT INTO users (name, surname, patronymic, age, gender, nationality) 
			VALUES ($1, $2, $3, $4, $5, $6) 
			RETURNING id`, name, surname, patronymic, age, gender, nationality).
		Scan(&id)
	if err != nil {
		return 0, err
	}

	return id, nil
}

func (s *Users) Update(id int, name string, surname string, patronymic string, age int, gender string, nationality string) (int, error) {
	result, err := s.db.Exec(`
		UPDATE users 
		SET 	name = $1, 
			surname = $2, 
			patronymic = $3, 
			age = $4, 
			gender = $5, 
			nationality = $6 
		WHERE id = $7`, name, surname, patronymic, age, gender, nationality, id)
	if err != nil {
		return 0, err
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return 0, err
	}

	if rowsAffected == 0 {
		return 0, sql.ErrNoRows
	}

	return id, nil
}

type UserList struct {
	Users []User
	Count int
	Pages int
	Page  int
}

func (s *Users) List(page int) (*UserList, error) {
	var count int
	if err := s.db.QueryRow(`SELECT COUNT(*) FROM users`).Scan(&count); err != nil {
		return nil, fmt.Errorf("count: %w", err)
	}

	offset := (page - 1) * limit
	if offset < 0 {
		offset = 0
		page = 1
	}
	if page > count/limit {
		page = count/limit + 1
		offset = (page - 1) * limit
	}

	rows, err := s.db.Query(`
		SELECT id, name, surname, patronymic, age, gender, nationality 
		FROM users 
		OFFSET $1 
		LIMIT $2`, offset, limit)
	if err != nil {
		return nil, err
	}
	defer func() {
		if err := rows.Close(); err != nil {
			slog.Debug("db rows close", "err", err)
		}
	}()

	userList := UserList{
		Users: make([]User, 0, limit),
		Count: count,
		Pages: count/limit + 1,
		Page:  page,
	}

	for rows.Next() {
		var user User

		err := rows.Scan(
			&user.ID,
			&user.Name,
			&user.Surname,
			&user.Patronymic,
			&user.Age,
			&user.Gender,
			&user.Nationality,
		)
		if err != nil {
			return nil, err
		}

		userList.Users = append(userList.Users, user)
	}

	return &userList, nil
}
