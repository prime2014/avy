package accounts

import (
	"db"
	"fmt"
	"log"
	"time"

	validation "github.com/go-ozzo/ozzo-validation"
	"github.com/go-ozzo/ozzo-validation/is"
	"golang.org/x/crypto/bcrypt"
)

type ValidationError struct{}

type Authentication struct {
	AccessToken  string `json:"access_token"`
	RefreshToken string `json:"refresh_token"`
}

func (v *ValidationError) Error() string {
	return "Invalid Credentials!"
}

// This functionality hashes the password
func HashPassword(password string) string {
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), 14)

	if err != nil {
		panic(err.Error())
	}

	hash := string(bytes)
	return hash
}

// Validates the plaintext password is the same as the hash password
func ValidatePassword(hash, password string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))

	return err == nil
}

type Users struct {
	ID           int       `json:"id"`
	Firstname    string    `json:"firstname"`
	Lastname     string    `json:"lastname"`
	Email        string    `json:"email"`
	Password     string    `json:"password"`
	RefreshToken string    `json:"refresh_token"`
	IsActive     bool      `json:"is_active"`
	DateJoined   time.Time `json:"date_joined"`
}

func (u Users) Save() Users {
	mydb := db.DBconnect()
	Lastinsertid := 0
	hash := HashPassword(u.Password)
	err := mydb.QueryRow("INSERT INTO users(firstname, lastname, email, password) VALUES($1, $2, $3, $4) RETURNING id",
		u.Firstname,
		u.Lastname,
		u.Email,
		hash,
	).Scan(&Lastinsertid)

	if err != nil {
		log.Fatal(err.Error())
	}
	u.ID = Lastinsertid
	u.IsActive = false
	u.DateJoined = time.Now().Local()

	return u
}

func (u Users) Validate() error {
	return validation.ValidateStruct(&u,
		//email cannot be empty
		validation.Field(&u.Email, validation.Required, is.Email),
		validation.Field(&u.Firstname, validation.Required),
		validation.Field(&u.Lastname, validation.Required),
		validation.Field(&u.Password, validation.Required, validation.Length(8, 12)),
	)
}

// Validates that the login credentials are legitimate
func (u Users) Authenticate() (*Authentication, error) {
	mydb := db.DBconnect()

	user := Users{}
	fmt.Printf("Parsed user: %v", u.Email)
	fmt.Printf("Parsed password: %v", u.Password)

	// check if the email exists in the database
	err := mydb.QueryRow("SELECT id, firstname, lastname, email, password, is_active, date_joined FROM users WHERE email=$1", u.Email).Scan(
		&user.ID,
		&user.Firstname,
		&user.Lastname,
		&user.Email,
		&user.Password,
		&user.IsActive,
		&user.DateJoined,
	)

	fmt.Printf("The user is: %v", user)

	if err != nil {
		return &Authentication{}, &ValidationError{}
	}

	// check if the two passwords match
	isTruePassword := ValidatePassword(user.Password, u.Password)

	if !isTruePassword {
		return &Authentication{}, &ValidationError{}
	}

	access_token, err := GenerateJwt(user.Email, user.ID)
	refresh_token, err := GenerateRefreshToken(user.Email, user.ID)

	if err != nil {
		return &Authentication{}, &ValidationError{}
	}

	dbdriver := db.DBconnect()

	_, err = dbdriver.Exec("UPDATE users SET refresh_token=$1 WHERE email=$2", refresh_token, user.Email)

	if err != nil {
		return &Authentication{}, err
	}

	return &Authentication{
		AccessToken:  access_token,
		RefreshToken: refresh_token,
	}, nil

}
