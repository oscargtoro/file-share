package main

import (
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
)

const (
	HOST = "localhost"
	PORT = "8080"
)

func start(channels map[string][]string) {

	// handler function for /join, logic for when channel registration is issued
	joinHandler := func(w http.ResponseWriter, r *http.Request) {
		channel := r.FormValue("name")
		log.Println("Request to join channel " + channel + " recieved...")
		if _, ok := channels[channel]; !ok {
			channels[channel] = append(channels[channel], r.FormValue("host"))
			channels[channel] = append(channels[channel], r.FormValue("port"))
		}
		/* client := new(strings.Builder)
		log.Fprint(client, clientPort)
		log.Println(clientPort)
		clients[string(clientPort)] = name */
		io.WriteString(w, "Registered in channel "+channel)
		log.Println(channels)
	}

	// handler function for /send, logic for when a file share request is issued
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

	http.HandleFunc("/join", joinHandler)
	http.HandleFunc("/send", sendHandler)
	log.Println("Listening on port " + PORT)
	log.Fatal(http.ListenAndServe(":"+PORT, nil))
}

func main() {
	channels := make(map[string][]string)
	if len(os.Args) == 1 {
		fmt.Println("Not enough arguments")
		os.Exit(0)
	}
	if os.Args[1] == "start" {
		start(channels)
	} else {
		fmt.Println("Invalid Argument")
	}
}
