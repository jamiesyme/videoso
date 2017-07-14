package main

import (
	"os"

	"github.com/jamiesyme/videoso/api"
)

func main() {
	config := api.NewServerConfig()
	config.Address = ":8001"
	config.AwsAccessKeyId = os.Getenv("AWS_ACCESS_KEY_ID")
	config.AwsSecretAccessKey = os.Getenv("AWS_SECRET_ACCESS_KEY")
	config.AwsRegion = os.Getenv("AWS_REGION")
	config.AwsBucket = os.Getenv("AWS_BUCKET")

	api.RunServer(config)
}
