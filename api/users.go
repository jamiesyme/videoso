package api

import (
	"encoding/json"
	"log"
	"net/mail"
	"regexp"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/lib/pq"
	"github.com/satori/go.uuid"
	"golang.org/x/crypto/bcrypt"
)

func (server *Server) RegisterUserRoutes() {
	server.Router.Handle(
		"POST",
		"/users",
		newUserHandler(server),
	)
	server.Router.Handle(
		"GET",
		"/users/:userId",
		accessTokenMiddleware(server),
		getUserHandler(server),
	)
}

func isEmailAddressValid(emailAddress string) bool {
	_, err := mail.ParseAddress(emailAddress)
	return err == nil
}

func genCleanEmailAddress(emailAddress string) string {
	return strings.ToLower(emailAddress)
}

var usernameRegex = regexp.MustCompile(`^[\w-]{1,20}$`)

func isUsernameValid(username string) bool {
	return usernameRegex.MatchString(username)
}

func genCleanUsername(username string) string {
	return strings.ToLower(username)
}

func isPasswordValid(password string) bool {
	// 6 <= len(password) <= 1024
	// len() is explicitly used to measure bytes
	return len(password) >= 6 && len(password) <= 1024
}

const bcryptCost = 14

func hashPassword(password string) (string, error) {
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), bcryptCost)
	return string(bytes), err
}

func checkPasswordHash(password, hash string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
	return err == nil
}

type UserInput struct {
	EmailAddress string `json:"emailAddress"`
	Password     string `json:"password"`
}

type UserOutput struct {
	UserId       string `json:"userId"`
	EmailAddress string `json:"emailAddress"`
}

func newUserHandler(server *Server) func(*gin.Context) {
	return func(c *gin.Context) {

		// Extract the user data
		userInput := new(UserInput)
		err := json.NewDecoder(c.Request.Body).Decode(userInput)
		if err != nil {
			c.Status(400)
			log.Println("failed to decode json")
			log.Println(err.Error())
			return
		}
		if !isEmailAddressValid(userInput.EmailAddress) {
			c.JSON(400, gin.H{
				"error": "invalid email address",
			})
			return
		}
		if !isPasswordValid(userInput.Password) {
			c.JSON(400, gin.H{
				"error": "invalid password",
			})
			return
		}

		// Clean the user data
		cleanEmailAddress := genCleanEmailAddress(userInput.EmailAddress)

		// Generate the user id
		userIdUuid, _ := uuid.NewV4()
		userId := userIdUuid.String()

		// Hash the password
		passwordHash, err := hashPassword(userInput.Password)
		if err != nil {
			c.Status(500)
			log.Println("failed to hash password")
			log.Println(err.Error())
			return
		}

		// Try to create the user
		userCreatedAt := time.Now()
		queryStr := "" +
			"INSERT INTO users (" +
			"  user_id," +
			"  email_address," +
			"  password_hash," +
			"  created_at" +
			") VALUES ($1, $2, $3, $4)"
		_, err = server.Db.Exec(
			queryStr,
			userId,
			cleanEmailAddress,
			passwordHash,
			userCreatedAt.UTC(),
		)
		if err != nil {
			pqErr, ok := err.(*pq.Error)
			if ok && pqErr.Code.Name() == "unique_violation" {

				// Check for email conflicts
				if pqErr.Constraint == "users_email_address_key" {
					c.JSON(409, gin.H{
						"error": "email address is taken",
					})
					return
				}
			}
			c.Status(500)
			log.Println("failed to save user to postgres")
			log.Println(err.Error())
			return
		}

		// User creation was a success
		c.Header("location", "/users/"+userId)
		c.JSON(201, &UserOutput{
			UserId:       userId,
			EmailAddress: userInput.EmailAddress,
		})
	}
}

func getUserHandler(server *Server) func(*gin.Context) {
	return func(c *gin.Context) {
		queryStr := "" +
			"SELECT email_address " +
			"FROM users " +
			"WHERE user_id = $1"
		var (
			emailAddress string
		)
		err := server.Db.QueryRow(queryStr, c.Param("userId")).Scan(
			&emailAddress,
		)
		if err != nil {
			c.Status(500)
			log.Println("failed to get user from postgres")
			log.Println(err.Error())
			return
		}
		c.JSON(200, &UserOutput{
			UserId:       c.Param("userId"),
			EmailAddress: emailAddress,
		})
	}
}
