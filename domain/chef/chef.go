package chef

import (
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

type Chef struct {
	gorm.Model
	FirstName string `json:"first_name"`
	LastName  string `json:"last_name"`
	Email     string `json:"email" gorm:"unique"`
	Password  string `json:"password"`
}

func NewChef(firstName, lastName, email, password string) (*Chef, error) {
	if firstName == "" || lastName == "" || email == "" || password == "" {
		return nil, ErrInvalidInput
	}
	hashed, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return nil, err
	}
	return &Chef{FirstName: firstName, LastName: lastName, Email: email, Password: string(hashed)}, nil
}

func (c *Chef) VerifyPassword(raw string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(c.Password), []byte(raw))
	return err == nil
}

func (c *Chef) UpdateProfile(firstName, lastName, email string) error {
	if firstName == "" && lastName == "" && email == "" {
		return ErrInvalidInput
	}
	if firstName != "" {
		c.FirstName = firstName
	}
	if lastName != "" {
		c.LastName = lastName
	}
	if email != "" {
		c.Email = email
	}
	return nil
}

func (c *Chef) SetPassword(raw string) error {
	if raw == "" {
		return ErrInvalidInput
	}
	hashed, err := bcrypt.GenerateFromPassword([]byte(raw), bcrypt.DefaultCost)
	if err != nil {
		return err
	}
	c.Password = string(hashed)
	return nil
}
