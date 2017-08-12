package api

import (
	"crypto/rand"
	"database/sql"
	"encoding/hex"
	"encoding/json"
	"log"
	"regexp"
	"strconv"
	"strings"
	"time"

	"github.com/dgrijalva/jwt-go"
	"github.com/gin-gonic/gin"
)

func (server *Server) RegisterAuthenticationRoutes() {
	server.Router.Handle(
		"DELETE",
		"/refresh-tokens/-",
		refreshTokenMiddleware(server),
		deleteRefreshTokenHandler(server),
	)
	server.Router.Handle(
		"POST",
		"/refresh-tokens",
		newRefreshTokenHandler(server),
	)
	server.Router.Handle(
		"POST",
		"/access-tokens",
		refreshTokenMiddleware(server),
		newAccessTokenHandler(server),
	)
}

const refreshTokenIdLen = 48

var refreshTokenIdRegex = regexp.MustCompile(
	`^[0-9a-fA-F]{` + strconv.Itoa(refreshTokenIdLen) + `}$`,
)

func generateRefreshTokenId() (string, error) {
	bytes := make([]byte, refreshTokenIdLen/2)
	_, err := rand.Read(bytes)
	if err != nil {
		return "", err
	}
	return hex.EncodeToString(bytes), nil
}

func validateRefreshTokenId(refreshTokenId string) bool {
	return refreshTokenIdRegex.MatchString(refreshTokenId)
}

func refreshTokenMiddleware(server *Server) gin.HandlerFunc {
	return func(c *gin.Context) {
		fail := func() {
			c.Header("www-authenticate", "Bearer")
			c.AbortWithStatus(401)
		}

		authHeader := c.Request.Header.Get("authorization")
		authParts := strings.SplitN(authHeader, " ", 2)
		if len(authParts) != 2 || authParts[0] != "Bearer" {
			fail()
			return
		}

		refreshTokenId := authParts[1]
		if !validateRefreshTokenId(refreshTokenId) {
			fail()
			return
		}

		queryStr := "" +
			"SELECT user_id " +
			"FROM refresh_tokens " +
			"WHERE refresh_token_id = $1"
		userId := ""
		err := server.Db.QueryRow(queryStr, refreshTokenId).Scan(&userId)
		if err == sql.ErrNoRows {
			fail()
			return
		}
		if err != nil {
			c.Status(500)
			log.Println("failed to get refresh token from postgres")
			return
		}

		c.Set("refreshTokenId", refreshTokenId)
		c.Set("userId", userId)
	}
}

func deleteRefreshTokenHandler(server *Server) func(*gin.Context) {
	return func(c *gin.Context) {
		queryStr := "" +
			"DELETE FROM refresh_tokens " +
			"WHERE refresh_token_id = $1"
		result, err := server.Db.Exec(queryStr, c.GetString("refreshTokenId"))
		if err != nil {
			c.Status(500)
			log.Println("failed to delete refresh token from postgres")
			log.Println(err.Error())
			return
		}
		if rows, _ := result.RowsAffected(); rows == 0 {
			// The refresh token is verified in middleware before hitting this
			// route, so the only way this can happen is if the token is deleted
			// in between the middleware and here.
			// TODO: Perhaps it's better to return 404
			c.Status(500)
			log.Println("failed to delete refresh token from postgres")
			log.Println("token not found: " + c.GetString("refreshTokenId"))
			return
		}
		c.Status(204)
	}
}

func newRefreshTokenHandler(server *Server) func(*gin.Context) {
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
				"error": "bad email address",
			})
			return
		}
		if !isPasswordValid(userInput.Password) {
			c.JSON(400, gin.H{
				"error": "bad password",
			})
			return
		}

		// Get the user data
		cleanEmailAddress := genCleanEmailAddress(userInput.EmailAddress)
		queryStr := "" +
			"SELECT user_id, password_hash " +
			"FROM users " +
			"WHERE email_address = $1 " +
			"LIMIT 1"
		var (
			userId       string
			passwordHash string
		)
		err = server.Db.QueryRow(queryStr, cleanEmailAddress).Scan(
			&userId,
			&passwordHash,
		)
		if err == sql.ErrNoRows {
			c.JSON(404, gin.H{
				"error": "invalid email address or password",
			})
			return
		}
		if err != nil {
			c.Status(500)
			log.Println("failed to query postgres for user")
			log.Println(err.Error())
			return
		}

		// Check the password
		if !checkPasswordHash(userInput.Password, passwordHash) {
			c.JSON(404, gin.H{
				"error": "invalid email address or password",
			})
			return
		}

		// Successfully authenticated; create the refresh token
		refreshTokenId, err := generateRefreshTokenId()
		if err != nil {
			c.Status(500)
			log.Println("failed to generate refresh token id")
			log.Println(err.Error())
			return
		}
		refreshTokenCreatedAt := time.Now().UTC()
		queryStr = "" +
			" INSERT INTO refresh_tokens (" +
			"   refresh_token_id," +
			"   user_id," +
			"   created_at" +
			" ) VALUES ($1, $2, $3)"
		_, err = server.Db.Exec(
			queryStr,
			refreshTokenId,
			userId,
			refreshTokenCreatedAt,
		)
		if err != nil {
			c.Status(500)
			log.Println("failed to save refresh token to postgres")
			log.Println(err.Error())
			return
		}

		// Refresh token creation was a success
		c.JSON(201, gin.H{
			"refreshTokenId": refreshTokenId,
		})
	}
}

func accessTokenMiddleware(server *Server) gin.HandlerFunc {
	return func(c *gin.Context) {
		fail := func() {
			c.Header("www-authenticate", "Bearer")
			c.AbortWithStatus(401)
		}

		authHeader := c.Request.Header.Get("authorization")
		authParts := strings.SplitN(authHeader, " ", 2)
		if len(authParts) != 2 || authParts[0] != "Bearer" {
			fail()
			return
		}

		accessTokenStr := authParts[1]
		jwtParser := jwt.Parser{}
		jwtParser.ValidMethods = append(
			jwtParser.ValidMethods,
			jwt.SigningMethodHS256.Alg(),
		)
		accessToken, err := jwtParser.ParseWithClaims(
			accessTokenStr,
			&jwt.StandardClaims{},
			func(token *jwt.Token) (interface{}, error) {
				return server.AccessTokenSecret, nil
			},
		)
		if err != nil {
			fail()
			return
		}

		claims := accessToken.Claims.(*jwt.StandardClaims)
		c.Set("userId", claims.Subject)
	}
}

func newAccessTokenHandler(server *Server) func(*gin.Context) {
	return func(c *gin.Context) {

		// Generate the access token
		expireTime := time.Now().UTC().Add(time.Duration(time.Minute * 20))
		claims := &jwt.StandardClaims{
			ExpiresAt: expireTime.Unix(),
			Subject:   c.GetString("userId"),
		}
		token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
		signedStr, err := token.SignedString(server.AccessTokenSecret)
		if err != nil {
			c.Status(500)
			log.Println("failed to sign access token")
			log.Println(err.Error())
			return
		}

		// Access token creation was a success
		c.JSON(201, gin.H{
			"accessToken": signedStr,
		})
	}
}
