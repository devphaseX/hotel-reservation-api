package types

import (
	"fmt"
	"regexp"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"golang.org/x/crypto/bcrypt"
)

const (
	bcryptCost         = 12
	minFirstNameLength = 2
	minLastNameLength  = 2
	minPasswordLength  = 7
)

type UpdateUserParams struct {
	FirstName string `json:"firstName"`
	LastName  string `json:"lastName"`
}

func (params UpdateUserParams) ToBSON() bson.M {
	m := bson.M{}

	if len(params.FirstName) > 0 {
		m["first_name"] = params.FirstName
	}

	if len(params.LastName) > 0 {
		m["last_name"] = params.LastName
	}

	return m
}

type CreateUserParams struct {
	FirstName string `json:"firstName"`
	LastName  string `json:"lastName"`
	Email     string `json:"email"`
	Password  string `json:"password"`
}

func (params CreateUserParams) Validate() map[string]string {
	errors := map[string]string{}
	if len(params.FirstName) < minFirstNameLength {
		errors["firstName"] = fmt.Sprintf("firstname length should contain a minimum of %d characters", minFirstNameLength)
	}

	if len(params.LastName) < minLastNameLength {
		errors["lastName"] = fmt.Sprintf("lastname length should contain a minimum of %d characters", minLastNameLength)
	}

	if len(params.Password) < minPasswordLength {
		errors["password"] = fmt.Sprintf("lastname length should contain a minimum of %d characters", minLastNameLength)
	}

	if !IsValidEmail(params.Email) {
		errors["email"] = "email is invalid"
	}

	return errors

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

func IsValidPassword(epw, pw string) bool {
	return bcrypt.CompareHashAndPassword([]byte(epw), []byte(pw)) == nil
}
