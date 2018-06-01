package main

import (
	"github.com/gorilla/mux"
	"net/http"
	"io/ioutil"
	"log"
	"time"
	"github.com/BurntSushi/toml"
	"fmt"
)

type server struct {
	IP           string `toml:"server_ip"`
	Port         string `toml:"server_port"`
	WriteTimeout int    `toml:"server_write_timeout"`
	ReadTimeout  int    `toml:"server_read_timeout"`
}

func main() {
	var config server
	if _, err := toml.DecodeFile("config/app.toml", &config); err != nil {
		fmt.Println(err)
		return
	}

	router := mux.NewRouter().StrictSlash(true)
	router.HandleFunc("/", GetIndex).Methods("GET")
	router.PathPrefix("/js/").Handler(http.StripPrefix("/js/", http.FileServer(http.Dir("web/dist/"))))

	srv := &http.Server{
		Handler:      router,
		Addr:         config.IP + ":" + config.Port,
		WriteTimeout: 10 * time.Second,
		ReadTimeout:  10 * time.Second,
	}

	log.Fatal(srv.ListenAndServe())
}

func GetIndex(writer http.ResponseWriter, request *http.Request) {
	indexFile, err := ioutil.ReadFile("web/index.html")
	if err != nil {
		log.Println(err)
		writer.WriteHeader(http.StatusNotFound)
		return
	}

	writer.Write(indexFile)
}
