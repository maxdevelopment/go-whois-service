package main

import (
	"github.com/gorilla/mux"
	"net/http"
	"io/ioutil"
	"log"
	"time"
	"github.com/BurntSushi/toml"
	"fmt"
	"github.com/maxdevelopment/go-whois-service/ws"
)

type server struct {
	IP           string `toml:"server_ip"`
	Port         string `toml:"server_port"`
}

func main() {
	var config server
	if _, err := toml.DecodeFile("config/app.toml", &config); err != nil {
		fmt.Println(err)
		return
	}

	hub := ws.H
	go hub.Run()

	router := mux.NewRouter().StrictSlash(true)
	router.HandleFunc("/", GetIndex).Methods("GET")
	router.PathPrefix("/js/").Handler(http.StripPrefix("/js/", http.FileServer(http.Dir("web/dist/"))))
	router.HandleFunc("/join", ws.Handler).Methods("GET")

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
