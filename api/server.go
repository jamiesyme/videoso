package api

import (
	"io"
	"log"
	"net/http"
	"os"
	"os/exec"

	"github.com/julienschmidt/httprouter"
)

const genSamplesPath = "/opt/videoso/api/gen-samples.sh"
const genDashPath = "/opt/videoso/api/gen-mpeg-dash.sh"
const storePath = "/var/www/videoso"
const sampleVideoPath = storePath + "/sample.mp4"

func RunServer(addr string) {
	router := httprouter.New()

	router.POST("/upload-video", uploadVideo)

	log.Fatal(http.ListenAndServe(addr, router))
}

func uploadVideo(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {

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
	out, err := os.OpenFile(sampleVideoPath, os.O_WRONLY|os.O_CREATE, 0644)
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

	log.Println("new upload successful")
	w.WriteHeader(http.StatusOK)
}
