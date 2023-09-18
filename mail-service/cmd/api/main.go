package main

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"strconv"
)

type Config struct {
	Mailer Mail
}

const webPort = "80"

func main() {

	app := Config{
		Mailer: createMailer(),
	}

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

func createMailer() Mail {
	port, _ := strconv.Atoi(os.Getenv("MAIL_PORT"))

	m := Mail{
		Domain:       os.Getenv("MAIL_DOMAIN"),
		Host:         os.Getenv("MAIL_HOST"),
		Port:         port,
		Username:     os.Getenv("MAIL_USERNAME"),
		Password:     os.Getenv("MAIL_PASSWORD"),
		Encryptation: os.Getenv("MAIL_ENCRYPTION"),
		FromName:     os.Getenv("MAIL_FROM_NAME"),
		FromAdress:   os.Getenv("MAIL_FROM_ADDRESS"),
	}

	return m

}
