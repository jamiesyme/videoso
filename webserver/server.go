package webserver

import (
	"log"
	"net/http"

	"github.com/julienschmidt/httprouter"
)

const indexPath = "/opt/videoso/webclient/index.html"

func RunServer(addr string) {
	router := httprouter.New()

	router.GET("/", getIndex)

	log.Fatal(http.ListenAndServe(addr, router))
}

func getIndex(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	http.ServeFile(w, r, indexPath)
}
