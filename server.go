package main

import (
	"BTC-price-restful/db"
	"BTC-price-restful/routes"
	"context"
	"fmt"
	"net/http"
	"os"
	"time"

	negronilogrus "github.com/meatballhat/negroni-logrus"
	"github.com/sirupsen/logrus"
	"github.com/urfave/negroni"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func main() {
	port := "8080"

	/* create logger */
	l := logrus.New()

	/* init mongo db */
	ctx, cancle := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancle()
	client, err := mongo.Connect(ctx, clientOptions())
	if err != nil {
		l.Fatal(err)
	}
	db.SetClint(client)

	/* init router */
	mux := routes.NewRouter()

	/* use logrus as middleware logger */
	n := negroni.New()
	n.Use(negronilogrus.NewMiddlewareFromLogger(l, "web"))

	n.UseHandler(mux)

	/* Create the main server object */
	server := http.Server{
		Addr:    fmt.Sprintf(":%s", port),
		Handler: n,
	}

	logrus.Println(fmt.Sprintf("Run the web server at :%s", port))
	logrus.Fatal(server.ListenAndServe())
}

func clientOptions() *options.ClientOptions {
	host := "db"
	if os.Getenv("profile") != "prod" {
		host = "localhost"
	}
	return options.Client().ApplyURI(
		"mongodb://" + host + ":27017",
	)
}
