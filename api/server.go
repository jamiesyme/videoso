package api

import (
	"database/sql"
	"encoding/json"
	"io"
	"io/ioutil"
	"log"
	"math"
	"net/http"
	"os"
	"os/exec"
	"path/filepath"
	"time"
	"unicode/utf8"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3/s3manager"
	"github.com/julienschmidt/httprouter"
	_ "github.com/lib/pq"
	"github.com/satori/go.uuid"
)

const genSamplesPath = "/opt/videoso/api/gen-samples.sh"
const genDashPath = "/opt/videoso/api/gen-mpeg-dash.sh"

type ServerConfig struct {
	Address            string
	AwsAccessKeyId     string
	AwsBucket          string
	AwsRegion          string
	AwsSecretAccessKey string
	AwsToken           string
	PgDbName           string
	PgUser             string
	PgPassword         string
	PgHost             string
	PgPort             string
}

func NewServerConfig() *ServerConfig {
	config := new(ServerConfig)
	config.PgHost = "localhost"
	config.PgPort = "5432"
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
	db         *sql.DB
}

func newServerContext(config *ServerConfig) *serverContext {
	ctx := new(serverContext)
	ctx.config = config
	awsConfig := aws.NewConfig().
		WithRegion(config.AwsRegion).
		WithCredentials(config.newAwsCreds())
	ctx.awsSess = session.Must(session.NewSession(awsConfig))
	ctx.s3Uploader = s3manager.NewUploader(ctx.awsSess)
	dbParams := "" +
		" dbname=" + config.PgDbName +
		" user=" + config.PgUser +
		" password=" + config.PgPassword +
		" host=" + config.PgHost +
		" port=" + config.PgPort
	var err error
	ctx.db, err = sql.Open("postgres", dbParams)
	if err != nil {
		log.Println("failed to connect to database")
		log.Println(err.Error())
	} else {
		err = ctx.db.Ping()
		if err != nil {
			log.Println("failed to connect to database")
			log.Println(err.Error())
		}
	}
	return ctx
}

func makeHandler(ctx *serverContext, h func(*serverContext, http.ResponseWriter, *http.Request, httprouter.Params)) httprouter.Handle {
	return func(w http.ResponseWriter, r *http.Request, params httprouter.Params) {
		h(ctx, w, r, params)
	}
}

func RunServer(config *ServerConfig) {
	ctx := newServerContext(config)
	defer ctx.db.Close()

	router := httprouter.New()
	router.POST("/videos", makeHandler(ctx, uploadVideo))
	router.GET("/videos", makeHandler(ctx, getVideos))
	log.Fatal(http.ListenAndServe(config.Address, router))
}

func uploadVideo(ctx *serverContext, w http.ResponseWriter, r *http.Request, _ httprouter.Params) {

	// Apply the file size limit
	// We add an extra MB to account for the other request body parameters
	maxReqBytes := int64((100 + 1) * math.Pow(2, 20))
	r.Body = http.MaxBytesReader(w, r.Body, maxReqBytes)

	// Read the input file
	in, _, err := r.FormFile("videoFile")
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		log.Println("new upload failed - bad request - file")
		log.Println(err.Error())
		return
	}
	defer in.Close()

	// Read the input title
	videoTitle := r.FormValue("title")
	if utf8.RuneCountInString(videoTitle) == 0 {
		w.WriteHeader(http.StatusBadRequest)
		log.Println("new upload failed - bad request - title")
		log.Println("title required")
		return
	}
	if utf8.RuneCountInString(videoTitle) > 100 {
		w.WriteHeader(http.StatusBadRequest)
		log.Println("new upload failed - bad request - title")
		log.Println("title too long")
		return
	}

	// Create a temporary directory for us to work in
	tempDir, err := ioutil.TempDir("", "upload")
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		log.Println("new upload failed - failed to create temp dir")
		log.Println(err.Error)
		return
	}
	defer os.RemoveAll(tempDir)

	// Open the output file
	videoId := uuid.NewV4().String()
	videoPath := filepath.Join(tempDir, videoId+".mp4")
	out, err := os.OpenFile(videoPath, os.O_WRONLY|os.O_CREATE, 0644)
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
	samplesCmd := exec.Command(genSamplesPath, tempDir, videoId)
	samplesCmd.Dir = tempDir
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
	dashCmd := exec.Command(genDashPath, tempDir, videoId)
	dashCmd.Dir = tempDir
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
	originalVideoUrl := ""
	videoMpdUrl := ""
	filenameGlob, _ := filepath.Glob(filepath.Join(tempDir, "*_dash*"))
	filenamesToUpload := []string{videoPath}
	filenamesToUpload = append(filenamesToUpload, filenameGlob...)
	for _, filename := range filenamesToUpload {
		in, err = os.Open(filename)
		if err != nil {
			log.Println("failed to read file for upload")
			log.Println(err.Error())
			continue
		}
		defer in.Close()
		uploadInfo, err := ctx.s3Uploader.Upload(
			&s3manager.UploadInput{
				Bucket: aws.String(ctx.config.AwsBucket),
				Key:    aws.String(filepath.Base(filename)),
				Body:   in,
			},
		)
		if err != nil {
			log.Println("failed to upload file")
			log.Println(err.Error())
			continue
		}
		if filename == videoPath {
			originalVideoUrl = uploadInfo.Location
		} else if filepath.Base(filename) == videoId+"_dash.mpd" {
			videoMpdUrl = uploadInfo.Location
		}
	}
	log.Println("uploaded to s3")

	log.Println("saving to postgres...")
	videoCreatedAt := time.Now()
	queryStr := "INSERT INTO videos(video_id, title, original_video_s3_url, video_mpd_s3_url, created_at) VALUES($1, $2, $3, $4, $5);"
	_, err = ctx.db.Exec(
		queryStr,
		videoId,
		videoTitle,
		originalVideoUrl,
		videoMpdUrl,
		videoCreatedAt.UTC(),
	)
	if err != nil {
		log.Println("failed to save to postgres")
		log.Println(err.Error())
		return
	}
	log.Println("saved to postgres")

	log.Println("new upload successful")
	videoJson := make(map[string]interface{})
	videoJson["videoId"] = videoId
	videoJson["title"] = videoTitle
	videoJson["videoMpdUrl"] = videoMpdUrl
	videoJson["createdAt"] = videoCreatedAt.Unix()
	resJson := make(map[string]interface{})
	resJson["video"] = videoJson
	jsonStr, err := json.Marshal(resJson)
	if err != nil {
		log.Println("failed to marshal json response")
		log.Println(err.Error())
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Content-Type", "application/json")
	w.Write(jsonStr)
}

func getVideos(ctx *serverContext, w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	queryStr := "" +
		" SELECT video_id, title, video_mpd_s3_url, created_at" +
		" FROM videos" +
		" ORDER BY created_at DESC;"
	rows, err := ctx.db.Query(queryStr)
	if err != nil {
		log.Println("failed to query videos")
		log.Println(err.Error())
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	defer rows.Close()

	videosJson := make(map[string]interface{})
	videosJson["videos"] = make([]interface{}, 0)
	for rows.Next() {
		var (
			videoId     string
			title       string
			videoMpdUrl string
			createdAt   time.Time
		)
		err = rows.Scan(
			&videoId,
			&title,
			&videoMpdUrl,
			&createdAt,
		)
		if err != nil {
			log.Println("failed to scan video row")
			log.Println(err.Error())
			continue
		}
		videoJson := make(map[string]interface{})
		videoJson["videoId"] = videoId
		videoJson["title"] = title
		videoJson["videoMpdUrl"] = videoMpdUrl
		videoJson["createdAt"] = createdAt.Unix()
		videosJson["videos"] = append(videosJson["videos"].([]interface{}), videoJson)
	}
	err = rows.Err()
	if err != nil {
		log.Println("failed to iterate over video rows")
		log.Println(err.Error())
	}

	jsonStr, err := json.Marshal(videosJson)
	if err != nil {
		log.Println("failed to marshal json response")
		log.Println(err.Error())
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Content-Type", "application/json")
	w.Write(jsonStr)
}
