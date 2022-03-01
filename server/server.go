package main

import (
	"bytes"
	"fmt"
	"io"
	"log"
	"mime/multipart"
	"net/http"
	"os"
	"time"
)

const (
	HOST = "localhost"
	PORT = "8080"
)

func start(channels map[string]map[string]string) {

	// handler function for /join, logic for when channel registration is issued
	joinHandler := func(w http.ResponseWriter, r *http.Request) {
		channel := r.FormValue("name")
		log.Println("Request to join channel " + channel + " recieved...")
		if _, ok := channels[channel]; !ok {
			channels[channel] = make(map[string]string)
			channels[channel]["host"] = r.FormValue("host")
			channels[channel]["port"] = r.FormValue("port")
		}
		io.WriteString(w, "Registered in channel "+channel)
	}

	// handler function for /send, logic for when a file share request is issued
	sendHandler := func(w http.ResponseWriter, r *http.Request) {
		// Limits body to a size of 32MB this way it will only use the 32MB limit on mem used bellow
		r.Body = http.MaxBytesReader(w, r.Body, 32<<20)
		client := &http.Client{
			Timeout: time.Second * 10,
		}
		body := &bytes.Buffer{}
		writer := multipart.NewWriter(body)

		//Limits use of memory to 32.5MB
		err := r.ParseMultipartForm(32<<20 + 512) // Max 32MB
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			return
		}

		file, fHeader, err := r.FormFile("file")
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			return
		}

		fw, err := writer.CreateFormFile("file", "./"+fHeader.Filename)
		if err != nil {
			log.Println(err)
			return
		}

		_, err = io.Copy(fw, file)
		if err != nil {
			log.Println(err)
			return
		}
		writer.Close()

		channel := r.FormValue("channel")

		if _, ok := channels[channel]; ok {
			for _, v := range channels {
				req, err := http.NewRequest("POST", fmt.Sprintf("http://%s:%s/send", v["host"], v["port"]), bytes.NewReader(body.Bytes()))
				if err != nil {
					fmt.Println(err)
					return
				}
				fmt.Println(writer.FormDataContentType())
				req.Header.Set("Content-Type", writer.FormDataContentType())
				rsp, err := client.Do(req)
				if err != nil {
					fmt.Println(err)
					return
				}
				defer rsp.Body.Close()
				if rsp.StatusCode != http.StatusOK {
					log.Printf("Request failed with response code: %d", rsp.StatusCode)
					return
				}
				response, _ := io.ReadAll(rsp.Body)
				fmt.Println(string(response) + "Host" + v["host"])
			}
		}

		w.Write([]byte("File sent"))
	}

	http.HandleFunc("/join", joinHandler)
	http.HandleFunc("/send", sendHandler)
	log.Println("Listening on port " + PORT)
	log.Fatal(http.ListenAndServe(":"+PORT, nil))
}

func main() {
	channels := make(map[string]map[string]string)
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
