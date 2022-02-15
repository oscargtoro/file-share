package main

import (
	"io"
	"log"
	"net/http"
	"os"
)

const (
	HOST = "localhost"
	PORT = "8080"
)

func main() {
	sendHandler := func(w http.ResponseWriter, r *http.Request) {
		err := r.ParseMultipartForm(32 << 20) // Max 32MB
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			return
		}

		file, fHeader, err := r.FormFile("file")
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			return
		}
		tmpfile, err := os.Create("./" + fHeader.Filename)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
		defer tmpfile.Close()
		_, err = io.Copy(tmpfile, file)

		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
		w.Write([]byte("File Recieved"))
	}

	http.HandleFunc("/send", sendHandler)
	log.Println("Listening on port " + PORT)
	log.Fatal(http.ListenAndServe(":"+PORT, nil))
}
