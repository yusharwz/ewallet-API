package userRepository

import (
	"database/sql"
	"errors"
	"final-project-enigma/src/user"
	"fmt"
)

type userRepository struct {
	db *sql.DB
}

func NewUserRepository(db *sql.DB) user.UserRepository {
	return &userRepository{db}
}

func (repo *userRepository) CekEmail(email string) (bool, error) {
	var resp string
	query := "SELECT email FROM users WHERE email = $1"
	if err := repo.db.QueryRow(query, email).Scan(&resp); err != nil {
		fmt.Println(email)
		fmt.Println(err)
		return false, errors.New("invalid email")
	}

	return true, nil
}

func (repo *userRepository) InsertCode(code, email string) (bool, error) {
	var resp string
	query := "UPDATE users SET verification_code = $1 WHERE email = $2 RETURNING email;"
	if err := repo.db.QueryRow(query, code, email).Scan(&resp); err != nil {
		fmt.Println(err)
		return false, errors.New("fail Insert Code")
	}

	return true, nil
}
