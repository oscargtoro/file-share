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

// Send creates a buffer for the body, then using that buffer creates a multipart
// writer where the file will be loaded then sent in a request.
func send(path string) {

	client := &http.Client{
		Timeout: time.Second * 10,
	}

	body := &bytes.Buffer{}
	writer := multipart.NewWriter(body)
	fw, err := writer.CreateFormFile("file", path)
	if err != nil {
		fmt.Println("Error openning the file " + path)
		return
	}

	file, err := os.Open(path)
	if err != nil {
		fmt.Println("Error openning the file " + path)
		return
	}

	_, err = io.Copy(fw, file)
	if err != nil {
		fmt.Print("There was an error copying the file\n", err)
		return
	}

	writer.Close()
	req, err := http.NewRequest("POST", "http://localhost:8080/send", bytes.NewReader(body.Bytes()))
	if err != nil {
		fmt.Println(err)
		return
	}

	req.Header.Set("Content-Type", writer.FormDataContentType())
	rsp, err := client.Do(req)
	if err != nil {
		fmt.Println("Could not connect to server")
		return
	}

	if rsp.StatusCode != http.StatusOK {
		log.Printf("Request failed with response code: %d", rsp.StatusCode)
		return
	}

	defer rsp.Body.Close()
	response, _ := io.ReadAll(rsp.Body)
	fmt.Println(string(response))
}

func main() {
	if len(os.Args) < 3 {
		fmt.Println("Not enough arguments")
		os.Exit(0)
	}
	if os.Args[1] == "send" {
		send(os.Args[2])
		os.Exit(0)
	}
}
