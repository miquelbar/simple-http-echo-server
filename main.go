package main

import (
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"strings"
	"time"
)

type LogRecord struct {
	http.ResponseWriter
	status int
}

func (r *LogRecord) Write(p []byte) (int, error) {
	return r.ResponseWriter.Write(p)
}

func (r *LogRecord) WriteHeader(status int) {
	r.status = status
	r.ResponseWriter.WriteHeader(status)
}

func answer(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
}

func printHeaders(logger *log.Logger, r *http.Request) {
	logger.Println("Recieved headers")
	for header, value := range r.Header {
		logger.Println(`  `, header, strings.Join(value[:], ","))
	}
}

func printBody(logger *log.Logger, r *http.Request) {
	logger.Println("Recieved body")
	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		logger.Printf(err.Error())
		return
	}
	logger.Println(string(body))
	return
}

func logging(logger *log.Logger) func(f http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
			record := &LogRecord{
				ResponseWriter: w,
			}
			defer func() {
				printHeaders(logger, req)
				printBody(logger, req)
				logger.Println(req.Method, req.URL.Path, record.status, req.RemoteAddr, req.UserAgent())
			}()

			next.ServeHTTP(record, req)
		})
	}
}

func main() {
	logger := log.New(os.Stdout, "http: ", log.LstdFlags)
	router := http.NewServeMux()
	port := os.Getenv("APPLICATION_PORT")
	router.HandleFunc("/", answer)

	if len(port) == 0 {
		port = "9999"
	}

	server := &http.Server{
		Addr:         ":" + port,
		Handler:      (logging(logger)(router)),
		ErrorLog:     logger,
		ReadTimeout:  5 * time.Second,
		WriteTimeout: 10 * time.Second,
		IdleTimeout:  15 * time.Second,
	}

	logger.Println("Listening on port", port)
	server.ListenAndServe()
}
