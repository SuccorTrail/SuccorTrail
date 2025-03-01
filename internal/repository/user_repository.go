package repository

import (
	"github.com/SuccorTrail/SuccorTrail/internal/db"
	"github.com/SuccorTrail/SuccorTrail/internal/model"
)

type UserRepository interface {
	Create(user *model.User) error
	Update(user *model.User) error
	GetByEmail(email string) (*model.User, error)
	UserExists(email string) (bool, error)
}

type userRepository struct{}

func NewUserRepository() UserRepository {
	return &userRepository{}
}

func (r *userRepository) Create(user *model.User) error {
	_, err := db.GetDB().Exec(
		"INSERT INTO users (id, name, email, password, user_type, created_at) VALUES (?, ?, ?, ?, ?, ?)",
		user.ID, user.Name, user.Email, user.Password, user.UserType, user.CreatedAt)
	return err
}

func (r *userRepository) Update(user *model.User) error {
	_, err := db.GetDB().Exec(
		"UPDATE users SET name = ?, email = ?, password = ?, user_type = ? WHERE id = ?",
		user.Name, user.Email, user.Password, user.UserType, user.ID)
	return err
}

func (r *userRepository) GetByEmail(email string) (*model.User, error) {
	user := &model.User{}
	err := db.GetDB().QueryRow(
		"SELECT id, name, email, password, user_type, created_at FROM users WHERE email = ?",
		email).Scan(&user.ID, &user.Name, &user.Email, &user.Password, &user.UserType, &user.CreatedAt)
	if err != nil {
		return nil, err
	}
	return user, nil
}

func (r *userRepository) UserExists(email string) (bool, error) {
	var count int
	err := db.GetDB().QueryRow("SELECT COUNT(*) FROM users WHERE email = ?", email).Scan(&count)
	if err != nil {
		return false, err
	}
	return count > 0, nil
}
