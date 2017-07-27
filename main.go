package main

import (
	"flag"
	"fmt"
	"log"
	"net/http"
	"os"
	"io/ioutil"
	"strconv"
)

var (
	host = flag.String("host", "0.0.0.0", "Bind address for HTTP server")
	port = flag.Int("port", 8080, "Bind port for HTTP server")
)

func main() {
	log.SetOutput(os.Stdout) // All the logs belong to stdout
	flag.Parse()             // Parse command line arguments

	gmHandler := NewGoodnightMoonHandler(
		"resources/flesh.html",
		"resources/flesh.css",
		"resources/flesh.js",
	)
	imgHandler := NewImageHandler("resources/goodnight_moon.jpg")

	http.Handle("/", gmHandler)
	http.Handle("/goodnight_moon.jpg", imgHandler)
	http.HandleFunc("/flesch", fleschHandler)
	http.HandleFunc("/-/health", healthHandler)

	listen := fmt.Sprintf("%s:%d", *host, *port)
	log.Printf("james:web listening on %s", listen)
	log.Fatal(http.ListenAndServe(listen, Log(http.DefaultServeMux))) // Return a rc != 0 on failure
}

func healthHandler(w http.ResponseWriter, r *http.Request) {
	status := http.StatusOK
	w.WriteHeader(status)
	w.Header().Set("content-type", "text/plain")
	w.Write([]byte("It's OK!"))
}

func fleschHandler(w http.ResponseWriter, r *http.Request) {
	status := http.StatusOK
	w.WriteHeader(status)
	w.Header().Set("content-type", "text/plain")
	w.Write([]byte("It's OK!"))
}

type imageHandler struct {
	imagePath string
	imageBytes []byte
}

func check(e error) {
	if e != nil {
		panic(e)
	}
}

func NewImageHandler(imagePath string) *imageHandler {
	bytes, err := ioutil.ReadFile(imagePath)
	check(err)
	return &imageHandler{
		imagePath: imagePath,
		imageBytes: bytes,
	}
}

func (handler *imageHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "image/jpeg")
	w.Header().Set("Content-Length", strconv.Itoa(len(handler.imageBytes)))

	w.WriteHeader(200)
	w.Write(handler.imageBytes)
}

type goodnightMoonHandler struct {
	htmlPath string
	cssPath  string
	jsPath   string
	html string
	css  string
	js   string
}

func NewGoodnightMoonHandler(htmlPath string, cssPath string, jsPath string) *goodnightMoonHandler {
	htmlBytes, err := ioutil.ReadFile(htmlPath)
	check(err)
	cssBytes, err := ioutil.ReadFile(cssPath)
	check(err)
	jsBytes, err := ioutil.ReadFile(jsPath)
	check(err)
	return &goodnightMoonHandler{
		htmlPath: htmlPath,
		cssPath: cssPath,
		jsPath: jsPath,
		html: string(htmlBytes),
		css: string(cssBytes),
		js: string(jsBytes),
	}
}

func (handler *goodnightMoonHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/html")

	w.WriteHeader(200)
	w.Write(handler.BodyFmt())
}

func Log(handler http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		log.Printf("%s %s %s", r.RemoteAddr, r.Method, r.URL)
		handler.ServeHTTP(w, r)
	})
}

func (handler *goodnightMoonHandler) BodyFmt() []byte {
	return []byte(fmt.Sprintf(indexBody, handler.css, handler.js, handler.html))
}

const indexBody = `
<html>
<head>
<script src="https://ajax.googleapis.com/ajax/libs/jquery/3.1.0/jquery.min.js"></script>
<link rel="stylesheet" href="https://maxcdn.bootstrapcdn.com/bootstrap/3.3.7/css/bootstrap.min.css" integrity="sha384-BVYiiSIFeK1dGmJRAkycuHAHRg32OmUcww7on3RYdg4Va+PmSTsz/K68vbdEjh4u" crossorigin="anonymous">
<style>%s</style>
<script>%s</script>
</head>
<body>
%s
</audio>
</body>
</html>
`
