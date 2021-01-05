package main

import (
	"fmt"
	"net/http"
	"os"
)

const serviceNameHeader = "Svc_name"
const destinationIPHeader = "Destination_ip"

func main() {
	fmt.Println("Starting Client ...")
	fmt.Println("Server IP: " + os.Getenv("server-ip"))
	url := "http://" + os.Getenv("server-ip") + ":80/hello"
	// url := "http://127.0.0.1:8000/hello"
	req, _ := http.NewRequest("GET", url, nil)
	req.Header.Add(serviceNameHeader, os.Getenv("SVC_NAME"))
	req.Header.Add(destinationIPHeader, os.Getenv("server-ip"))
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		fmt.Println(err.Error())
	} else {
		fmt.Println(resp)
		fmt.Println(resp.Header)
	}
}
