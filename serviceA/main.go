package main

import (
	"fmt"
	"net/http"
	"os"
)

const serviceName = "Svc_name"

func main() {
	fmt.Println("Starting Service A ...")
	fmt.Println("Service B IP: " + os.Getenv("service-b-ip"))
	//url := "http://" + os.Getenv("service-b-ip") + ":80/hello"
	url := "http://127.0.0.1:8000/hello"
	req, _ := http.NewRequest("GET", url, nil)
	req.Header.Add("Svc_name", os.Getenv("SVC_NAME"))
	req.Header.Add("destination-ip", os.Getenv("service-b-ip"))
	resp, _ := http.DefaultClient.Do(req)
	fmt.Println(resp)
}
