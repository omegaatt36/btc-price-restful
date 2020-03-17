package main

import (
	"BTC-price-restful/db"
	"BTC-price-restful/remote"
	"BTC-price-restful/routes"
	"context"
	"fmt"
	"net/http"
	"os"
	"time"

	"github.com/garyburd/redigo/redis"
	"github.com/gorilla/handlers"
	negronilogrus "github.com/meatballhat/negroni-logrus"
	"github.com/sirupsen/logrus"
	"github.com/urfave/negroni"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

var l *logrus.Logger

func main() {
	port := "8080"

	/* create logger */
	l = logrus.New()

	/* init mongo db */
	ctx, cancle := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancle()
	client, err := mongo.Connect(ctx, clientOptions())
	if err != nil {
		l.Fatal(err)
	}
	db.SetMongoClint(client)
	SetRedisClint, err := redis.Dial("tcp", "127.0.0.1:6379")
	if err != nil {
		l.Fatal("Connect to redis error", err)
	}
	db.SetRedisClint(SetRedisClint)
	defer SetRedisClint.Close()

	/* init remote API */
	file, err := os.Open("APIconf.json")
	if err != nil {
		l.Info("can't open local APIconfig.json file, make sure the file has been created.")
		defaultFile, err := os.Open("APIconfDefault.json")
		if err != nil {
			file.Close()
			l.Fatal(err)
		}
		file.Close()
		file = defaultFile
	}
	defer file.Close()
	if err = remote.ParseConfig(file); err != nil {
		l.Info("can't parse local APIconfig.json file")
		l.Fatal(err)
	}
	remote.InitAPIs()

	/* init router */
	mux := routes.NewRouter()

	/* use logrus as middleware logger */
	n := negroni.New()
	n.Use(negronilogrus.NewMiddlewareFromLogger(l, "web"))
	allowedHeaders := handlers.AllowedHeaders([]string{"Content-Type", "Authorization"})
	allowedOrigins := handlers.AllowedOrigins([]string{"*"})
	allowedMethods := handlers.AllowedMethods([]string{"GET", "POST", "OPTIONS"})
	n.UseHandler(handlers.CORS(allowedHeaders, allowedOrigins, allowedMethods)(mux))

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
