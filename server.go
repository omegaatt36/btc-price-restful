package main

import (
	"BTC-price-restful/routes"
	"fmt"
	"net/http"

	negronilogrus "github.com/meatballhat/negroni-logrus"
	"github.com/sirupsen/logrus"
	"github.com/urfave/negroni"
)

func main() {

	host := "127.0.0.1"
	port := "8080"
	mux := routes.NewRouter()

	/* Create the logger for the web application. */
	l := logrus.New()
	n := negroni.New()
	n.Use(negronilogrus.NewMiddlewareFromLogger(l, "web"))

	n.UseHandler(mux)

	/* Create the main server object */
	server := http.Server{
		Addr:    fmt.Sprintf("%s:%s", host, port),
		Handler: n,
	}

	logrus.Println(fmt.Sprintf("Run the web server at %s:%s", host, port))
	logrus.Fatal(server.ListenAndServe())
}
