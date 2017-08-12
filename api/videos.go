package api

import (
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
	"github.com/aws/aws-sdk-go/service/s3/s3manager"
	"github.com/gin-gonic/gin"
	"github.com/satori/go.uuid"
)

func (server *Server) RegisterVideoRoutes() {
	server.Router.Handle(
		"GET",
		"/videos",
		getVideosHandler(server),
	)
	server.Router.Handle(
		"POST",
		"/videos",
		newVideoHandler(server),
	)
}

func getVideosHandler(server *Server) func(*gin.Context) {
	return func(c *gin.Context) {
		queryStr := "" +
			"SELECT video_id, title, video_mpd_s3_url, created_at " +
			"FROM videos " +
			"ORDER BY created_at DESC"
		rows, err := server.Db.Query(queryStr)
		if err != nil {
			c.Status(500)
			log.Println("failed to query videos")
			log.Println(err.Error())
			return
		}
		defer rows.Close()

		videosJson := make([]gin.H, 0)
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
			videosJson = append(videosJson, gin.H{
				"videoId":     videoId,
				"title":       title,
				"videoMpdUrl": videoMpdUrl,
				"createdAt":   createdAt.Unix(),
			})
		}
		err = rows.Err()
		if err != nil {
			log.Println("failed to iterate over video rows")
			log.Println(err.Error())
		}

		c.Header("access-control-allow-origin", "*")
		c.JSON(200, gin.H{
			"videos": videosJson,
		})
	}
}

func newVideoHandler(server *Server) func(*gin.Context) {
	return func(c *gin.Context) {

		// Apply the file size limit
		// We add an extra MB to account for the other request body parameters
		maxReqBytes := int64((100 + 1) * math.Pow(2, 20))
		c.Request.Body = http.MaxBytesReader(c.Writer, c.Request.Body, maxReqBytes)

		// Read the input file
		in, _, err := c.Request.FormFile("videoFile")
		if err != nil {
			c.Status(400)
			log.Println("new upload failed - bad request - file")
			log.Println(err.Error())
			return
		}
		defer in.Close()

		// Read the input title
		videoTitle := c.Request.FormValue("title")
		if utf8.RuneCountInString(videoTitle) == 0 {
			c.Status(400)
			log.Println("new upload failed - bad request - title")
			log.Println("title required")
			return
		}
		if utf8.RuneCountInString(videoTitle) > 100 {
			c.Status(400)
			log.Println("new upload failed - bad request - title")
			log.Println("title too long")
			return
		}

		// Create a temporary directory for us to work in
		tempDir, err := ioutil.TempDir("", "upload")
		if err != nil {
			c.Status(500)
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
			c.Status(500)
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
		const genSamplesPath = "/opt/videoso/api/gen-samples.sh"
		log.Println("generating samples...")
		samplesCmd := exec.Command(genSamplesPath, tempDir, videoId)
		samplesCmd.Dir = tempDir
		err = samplesCmd.Run()
		if err != nil {
			c.Status(500)
			log.Println("failed generating samples")
			log.Println(err.Error())
			return
		}
		log.Println("generated samples")

		// Generate the mpeg dash data
		const genDashPath = "/opt/videoso/api/gen-mpeg-dash.sh"
		log.Println("generating dash info...")
		dashCmd := exec.Command(genDashPath, tempDir, videoId)
		dashCmd.Dir = tempDir
		err = dashCmd.Run()
		if err != nil {
			c.Status(500)
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
			uploadInfo, err := server.S3Uploader.Upload(
				&s3manager.UploadInput{
					Bucket: aws.String(server.S3BucketVideos),
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
		queryStr := "" +
			"INSERT INTO videos (" +
			"  video_id," +
			"  title," +
			"  original_video_s3_url," +
			"  video_mpd_s3_url," +
			"  created_at" +
			") VALUES ($1, $2, $3, $4, $5)"
		_, err = server.Db.Exec(
			queryStr,
			videoId,
			videoTitle,
			originalVideoUrl,
			videoMpdUrl,
			videoCreatedAt.UTC(),
		)
		if err != nil {
			c.Status(500)
			log.Println("failed to save to postgres")
			log.Println(err.Error())
			return
		}
		log.Println("saved to postgres")

		log.Println("new upload successful")
		c.Header("access-control-allow-origin", "*")
		c.JSON(201, gin.H{
			"videoId":     videoId,
			"title":       videoTitle,
			"videoMpdUrl": videoMpdUrl,
			"createdAt":   videoCreatedAt.Unix(),
		})
	}
}
