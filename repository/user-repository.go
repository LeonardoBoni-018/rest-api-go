package repository

import (
	"database/sql"
	"fmt"

	"go-api/model"
)

type UserRepository struct {
	connection *sql.DB
}

func NewUserRepository(conn *sql.DB) *UserRepository {
	return &UserRepository{connection: conn}
}

func (ur *UserRepository) CreateUser(u model.User, passwordHash string) (int, error) {
	var id int
	err := ur.connection.QueryRow(
		"INSERT INTO users (username, password_hash) VALUES ($1, $2) RETURNING id",
		u.Username, passwordHash,
	).Scan(&id)

	if err != nil {
		return 0, err
	}

	return id, nil
}

func (ur *UserRepository) GetuserByName(username string) (*model.User, string, error) {
	query := "SELECT id, username, password_hash FROM users WHERE username = $1"
	row := ur.connection.QueryRow(query, username)

	var u model.User
	var hash string
	err := row.Scan(&u.ID, &u.Username, &hash)

	if err != nil {
		if err == sql.ErrNoRows {
			return nil, "", nil // usuário não encontrado
		}
		fmt.Println("Erro ao buscar usuário:", err)
		return nil, "", err
	}
	return &u, hash, nil
}
