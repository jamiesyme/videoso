package api

import (
	"io"
	"log"
	"net/http"
	"os"
	"os/exec"
	"path/filepath"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3/s3manager"
	"github.com/julienschmidt/httprouter"
)

const genSamplesPath = "/opt/videoso/api/gen-samples.sh"
const genDashPath = "/opt/videoso/api/gen-mpeg-dash.sh"
const storePath = "/var/www/videoso/"

type ServerConfig struct {
	Address            string
	AwsAccessKeyId     string
	AwsBucket          string
	AwsRegion          string
	AwsSecretAccessKey string
	AwsToken           string
}

func NewServerConfig() *ServerConfig {
	config := new(ServerConfig)
	return config
}

func (config *ServerConfig) newAwsCreds() *credentials.Credentials {
	return credentials.NewStaticCredentials(
		config.AwsAccessKeyId,
		config.AwsSecretAccessKey,
		config.AwsToken,
	)
}

type serverContext struct {
	config     *ServerConfig
	awsSess    *session.Session
	s3Uploader *s3manager.Uploader
}

func newServerContext(config *ServerConfig) *serverContext {
	ctx := new(serverContext)
	ctx.config = config
	awsConfig := aws.NewConfig().
		WithRegion(config.AwsRegion).
		WithCredentials(config.newAwsCreds())
	ctx.awsSess = session.Must(session.NewSession(awsConfig))
	ctx.s3Uploader = s3manager.NewUploader(ctx.awsSess)
	return ctx
}

func makeHandler(ctx *serverContext, h func(*serverContext, http.ResponseWriter, *http.Request, httprouter.Params)) httprouter.Handle {
	return func(w http.ResponseWriter, r *http.Request, params httprouter.Params) {
		h(ctx, w, r, params)
	}
}

func RunServer(config *ServerConfig) {
	ctx := newServerContext(config)

	router := httprouter.New()
	router.POST("/upload-video", makeHandler(ctx, uploadVideo))
	log.Fatal(http.ListenAndServe(config.Address, router))
}

func uploadVideo(ctx *serverContext, w http.ResponseWriter, r *http.Request, _ httprouter.Params) {

	// Read the input file
	in, _, err := r.FormFile("videoFile")
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		log.Println("new upload failed - bad request")
		log.Println(err.Error)
		return
	}
	defer in.Close()

	// Open the output file
	out, err := os.OpenFile(storePath + "sample.mp4", os.O_WRONLY|os.O_CREATE, 0644)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		log.Println("new upload failed - failed to save")
		log.Println(err.Error())
		return
	}
	defer out.Close()

	// Save the input file
	log.Println("saving input file...")
	io.Copy(out, in)
	log.Println("saved input file")

	// Generate the samples
	log.Println("generating samples...")
	samplesCmd := exec.Command(genSamplesPath)
	samplesCmd.Dir = storePath
	err = samplesCmd.Run()
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		log.Println("failed generating samples")
		log.Println(err.Error())
		return
	}
	log.Println("generated samples")

	// Generate the mpeg dash data
	log.Println("generating dash info...")
	dashCmd := exec.Command(genDashPath)
	dashCmd.Dir = storePath
	err = dashCmd.Run()
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		log.Println("failed generating dash info")
		log.Println(err.Error())
		return
	}
	log.Println("generated dash info")

	// Upload the files to s3
	log.Println("uploading to s3...")
	filenameGlob, _ := filepath.Glob(storePath + "sample*_dashinit.mp4")
	filenamesToUpload := []string{
		storePath + "sample.mp4",
		storePath + "sample_dash.mpd",
	}
	filenamesToUpload = append(filenamesToUpload, filenameGlob...)
	for _, filename := range filenamesToUpload {
		in, err = os.Open(filename)
		if err != nil {
			log.Println("failed to read file for upload")
			log.Println(err.Error())
			continue
		}
		defer in.Close()
		_, err = ctx.s3Uploader.Upload(
			&s3manager.UploadInput{
				Bucket: aws.String(ctx.config.AwsBucket),
				Key: aws.String(filepath.Base(filename)),
				Body: in,
			},
		)
		if err != nil {
			log.Println("failed to upload file")
			log.Println(err.Error())
			continue
		}
	}
	log.Println("uploaded to s3")

	log.Println("new upload successful")
	w.WriteHeader(http.StatusOK)
}
