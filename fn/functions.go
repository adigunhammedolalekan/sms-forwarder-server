package fn

import (
	"errors"
	"fmt"
	"github.com/adigunhammedolalekan/sms-forwarder/types"
	"github.com/dgrijalva/jwt-go"
	"golang.org/x/crypto/bcrypt"
	"log"
	"math/rand"
	"regexp"
	"strings"
	"time"
)

var letters = []rune("abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ")
var ErrInvalidEmail = errors.New("invalid email address")
var userRegexp = regexp.MustCompile("^[a-zA-Z0-9!#$%&'*+/=?^_`{|}~.-]+$")
var hostRegexp = regexp.MustCompile("^[^\\s]+\\.[^\\s]+$")

func init() { rand.Seed(time.Now().UnixNano()) }

// GenerateRandomString returns a randomly generated string
func GenerateRandomString(n int) string {
	b := make([]rune, n)
	for i := range b {
		b[i] = letters[rand.Intn(len(letters))]
	}
	return strings.ToUpper(string(b))
}

func ValidateEmail(email string) error {
	email = strings.TrimSpace(email)
	if len(email) < 6 || len(email) > 254 {
		return ErrInvalidEmail
	}

	at := strings.LastIndex(email, "@")
	if at <= 0 || at > len(email)-3 {
		return ErrInvalidEmail
	}

	user := email[:at]
	host := email[at+1:]

	if len(user) > 64 {
		return ErrInvalidEmail
	}

	if !userRegexp.MatchString(user) || !hostRegexp.MatchString(host) {
		return ErrInvalidEmail
	}

	return nil
}

func ReadableTime(ti time.Time) string {
	return fmt.Sprintf("%d %s %d, %d:%d", ti.Day(), ti.Month().String(), ti.Year(), ti.Hour(), ti.Minute())
}

func HashPassword(password string) string {
	hashed, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return ""
	}
	return string(hashed)
}

func VerifyPassword(hashedPassword, password string) bool {
	if err := bcrypt.CompareHashAndPassword([]byte(hashedPassword), []byte(password)); err != nil {
		return false
	}
	return true
}

type JwtTokenGenerator struct {
	jwtSecret []byte
}

func NewJwtTokenGenerator(secret []byte) *JwtTokenGenerator {
	return &JwtTokenGenerator{jwtSecret:secret}
}

func (j *JwtTokenGenerator) SignJwtToken(accountId uint, email string) string {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, &types.Token{
		UserId: accountId, Email: email,
	})
	tokenString, err := token.SignedString(j.jwtSecret)
	if err != nil {
		log.Printf("failed to generate JWT token; %v", err)
		return ""
	}
	return tokenString
}


func (j *JwtTokenGenerator) ParseJwtToken(headerValue string) (*types.Token, error) {
	parts := strings.Split(headerValue, " ")
	if len(parts) != 2 || parts[1] == "" {
		return nil, errors.New("authorization token is missing")
	}
	tokenString := parts[1]
	tk := &types.Token{}
	token, err := jwt.ParseWithClaims(tokenString, tk, func(token *jwt.Token) (interface{}, error) {
		if token.Method != jwt.GetSigningMethod("HS256") {
			return nil, errors.New("invalid signing method")
		}
		return j.jwtSecret, nil
	})

	if err != nil {
		return nil, errors.New("malformed authorization token")
	}
	if !token.Valid {
		return nil, errors.New("malformed or invalid authorization token")
	}
	return tk, nil
}
