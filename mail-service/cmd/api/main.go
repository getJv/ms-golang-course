package main

import (
	"fmt"
	"log"
	"net/http"
)

type Config struct {
}

const webPort = "80"

func main() {

	app := Config{}

	log.Println("Starting mailing server on port: ", webPort)

	srv := &http.Server{
		Addr:    fmt.Sprintf(":%s", webPort),
		Handler: app.routes(),
	}

	//Start server
	err := srv.ListenAndServe()
	if err != nil {
		log.Panic(err)
	}

}
