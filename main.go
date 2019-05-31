package main

import (
	"flag"
	"log"
	"net/http"
	"net/url"
	"os"
)

const (
	baseUrl = "http://thecatapi.com"
	getPath = "/api/images/get"
)

const (
	portKey    = "port"
	logFileKey = "log"

	// Can be:
	// jpg, png, gif
	imgTypeKey = "type"

	// Can be:
	// small, med, full
	imgSizeKey = "size"
)

func buildServeFunction(urlPath, imgType, imgSize string) http.HandlerFunc {
	u, err := url.Parse(urlPath)
	if err != nil {
		log.Fatal(err)
	}
	q := u.Query()
	q.Set(imgTypeKey, imgType)
	q.Set(imgSizeKey, imgSize)
	u.RawQuery = q.Encode()
	fullUrlPath := u.String()
	log.Printf("Rendered cat-pic URL: %s\n", fullUrlPath)
	return func(w http.ResponseWriter, r *http.Request) {
		log.Printf("Got user connection: %s\n", r.Header.Get("User-Agent"))
		http.Redirect(w, r, fullUrlPath, http.StatusTemporaryRedirect)
	}
}

func main() {
	// Add command line parameter
	flag.String(portKey, "80", "Specify port to serve on")
	flag.String(logFileKey, "log.log", "the output log file")
	flag.Parse()

	f, err := os.OpenFile("log.log", os.O_RDWR|os.O_CREATE|os.O_APPEND, 0666)
	if err != nil {
		log.Fatalf("error opening file: %v", err)
	}
	defer f.Close()

	log.SetOutput(f)

	port, imgType, imgSize := "80", "jpg", "small"

	sv := buildServeFunction(baseUrl+getPath, imgType, imgSize)
	http.HandleFunc("/", sv)
	addr := ":" + port
	log.Printf("Start serving on %s ...", addr)
	log.Fatal(http.ListenAndServe(addr, http.DefaultServeMux))
}
