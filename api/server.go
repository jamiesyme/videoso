package api

import (
	"database/sql"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3/s3manager"
	"github.com/gin-gonic/gin"
	_ "github.com/lib/pq"
)

type Server struct {
	AccessTokenSecret []byte
	Address           string
	Db                *sql.DB
	Router            *gin.Engine
	S3Uploader        *s3manager.Uploader
	S3BucketVideos    string
}

func NewServer() *Server {
	server := new(Server)
	server.Address = ":8001"
	server.Router = gin.Default()
	return server
}

func (server *Server) Run() {
	server.Router.Run(server.Address)
}

type AwsParams struct {
	AccessKeyId     string
	Region          string
	SecretAccessKey string
	Token           string
}

func NewS3Uploader(p AwsParams) (*s3manager.Uploader, error) {

	// Build aws config
	creds := credentials.NewStaticCredentials(
		p.AccessKeyId,
		p.SecretAccessKey,
		p.Token,
	)
	config := aws.NewConfig().WithRegion(p.Region).WithCredentials(creds)

	// Connect to S3
	awsSession, err := session.NewSession(config)
	if err != nil {
		return nil, err
	}
	return s3manager.NewUploader(awsSession), nil
}

type PostgresParams struct {
	Database string
	User     string
	Password string
	Host     string
	Port     string
}

func NewPostgresDatabase(p PostgresParams) (*sql.DB, error) {

	// Set parameter defaults
	if p.Host == "" {
		p.Host = "localhost"
	}
	if p.Port == "" {
		p.Port = "5432"
	}

	// Build the connection string
	pStr := "" +
		"dbname=" + p.Database +
		" user=" + p.User +
		" password=" + p.Password +
		" host=" + p.Host +
		" port=" + p.Port

	// Create connection
	db, err := sql.Open("postgres", pStr)
	if err != nil {
		return nil, err
	}

	// Verify connection
	err = db.Ping()
	if err != nil {
		db.Close()
		return nil, err
	}

	return db, nil
}
