package types

import (
	"fmt"
	"regexp"

	"go.mongodb.org/mongo-driver/bson/primitive"
	"golang.org/x/crypto/bcrypt"
)

const (
	bcryptCost         = 12
	minFirstNameLength = 2
	minLastNameLength  = 2
	minPasswordLength  = 7
)

type CreateUserParams struct {
	FirstName string `json:"firstName"`
	LastName  string `json:"lastName"`
	Email     string `json:"email"`
	Password  string `json:"password"`
}

func (params CreateUserParams) Validate() error {
	if len(params.FirstName) < minFirstNameLength {
		return fmt.Errorf("firstname length should contain a minimum of %d characters", minFirstNameLength)
	}

	if len(params.LastName) < minLastNameLength {
		return fmt.Errorf("lastname length should contain a minimum of %d characters", minLastNameLength)
	}

	if len(params.Password) < minPasswordLength {
		return fmt.Errorf("lastname length should contain a minimum of %d characters", minLastNameLength)
	}

	if !IsValidEmail(params.Email) {
		return fmt.Errorf("email is invalid")
	}

	return nil

}

func IsValidEmail(email string) bool {
	reg := regexp.MustCompile(`^[a-z0-9._%+\-]+@[a-z0-9.\-]+\.[a-z]{2,4}$`)
	return reg.MatchString(email)
}

type User struct {
	ID                primitive.ObjectID `bson:"_id,omitempty" json:"id,omitempty"`
	FirstName         string             `bson:"first_name" json:"firstName"`
	LastName          string             `bson:"last_name" json:"lastName"`
	Email             string             `bson:"email" json:"email"`
	EncryptedPassword string             `bson:"encrpyted_password" json:"-"`
}

func NewUserFromParams(params CreateUserParams) (*User, error) {
	encp, err := bcrypt.GenerateFromPassword([]byte(params.Password), bcryptCost)

	if err != nil {
		return nil, err
	}

	return &User{FirstName: params.FirstName,
		LastName:          params.LastName,
		Email:             params.Email,
		EncryptedPassword: string(encp),
	}, nil
}
