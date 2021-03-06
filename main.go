package main

import (
	"log"
	"net/http"

	c "github.com/soyuka/incolore/config"
	"github.com/soyuka/incolore/handlers"
	t "github.com/soyuka/incolore/transports"
)

func main() {
	config := c.GetConfig()

	transport, err := t.NewTransport(&config)
	if err != nil {
		log.Fatal(err)
	}

	env := &handlers.Env{
		Transport: transport,
		Config:    config,
	}

	http.Handle("/favicon.ico", handlers.Handler{Env: env, Handler: handlers.Favicon})
	http.Handle("/", handlers.Handler{Env: env, Handler: handlers.GetIndex})

	log.Fatal(http.ListenAndServe(":"+config.Port, nil))
}
