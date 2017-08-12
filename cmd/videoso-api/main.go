package main

import (
	"encoding/hex"
	"log"
	"os"

	"github.com/jamiesyme/videoso/api"
)

func main() {
	server := api.NewServer()
	server.S3BucketVideos = os.Getenv("S3_BUCKET_VIDEOS")

	accessTokenSecretHex := os.Getenv("ACCESS_TOKEN_SECRET")
	accessTokenSecret, err := hex.DecodeString(accessTokenSecretHex)
	if err != nil {
		log.Println("invalid access token secret")
		return
	}
	server.AccessTokenSecret = accessTokenSecret

	server.Db, err = api.NewPostgresDatabase(api.PostgresParams{
		Database: os.Getenv("PG_DB_NAME"),
		Password: os.Getenv("PG_PASSWORD"),
		User:     os.Getenv("PG_USER"),
	})
	if err != nil {
		log.Println("failed to connect to postgres")
		log.Println(err.Error())
		return
	}

	server.S3Uploader, err = api.NewS3Uploader(api.AwsParams{
		AccessKeyId:     os.Getenv("AWS_ACCESS_KEY_ID"),
		Region:          os.Getenv("AWS_REGION"),
		SecretAccessKey: os.Getenv("AWS_SECRET_ACCESS_KEY"),
	})
	if err != nil {
		log.Println("failed to connect to s3")
		log.Println(err.Error())
		return
	}

	server.RegisterAuthenticationRoutes()
	server.RegisterUserRoutes()
	server.RegisterVideoRoutes()

	server.Run()
}
