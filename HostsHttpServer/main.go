package main

import (
	"fmt"
	"github.com/gpmgo/gopm/modules/log"
	"io/ioutil"
	"net/http"
	"os"
	"time"
)

func main() {
	go StartListen()

	os.Mkdir("/tmp/hosts/",0755)

	UpdateHosts()
}

func UpdateHosts() {
	ticker := time.NewTicker(48*time.Hour)
	for {
		r := GetPage("https://cdn.jsdelivr.net/gh/neoFelhz/neohosts@gh-pages/full/hosts")
		WriteWithIoutil("/tmp/hosts/adhosts",&r)
		<-ticker.C

	}
}

func WriteWithIoutil(name string ,content *string) {
	data :=  []byte(*content)
	if ioutil.WriteFile(name,data,0755) != nil {
		fmt.Println("Write Error")
	}
}

func GetPage(url string) string {
	resp, err := http.Get(url)
	if err != nil {
		log.Fatal("i/o timeout")
	}

	if resp.StatusCode != 200 {
		log.Fatal("status code error: %d %s", resp.StatusCode, resp.Status)
	}

	defer resp.Body.Close()
	html, _ := ioutil.ReadAll(resp.Body)
	return string(html)
}

func StartListen() {
	mux := http.NewServeMux()
	mux.Handle("/",http.FileServer(http.Dir("/tmp/hosts")))

	server := http.Server{
		Addr:    "0.0.0.0:85",
		Handler: mux,
	}

	if err := server.ListenAndServe();err != nil {
		log.Fatal(err.Error())
	}
}