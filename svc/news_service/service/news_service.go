package service

import (
	"fmt"
	httptransport "github.com/go-kit/kit/transport/http"
	"github.com/gorilla/mux"
	nm "github.com/the-gigi/delinkcious/pkg/news_manager"
	"log"
	"net/http"
	"os"
)

func Run() {
	log.Println("Service started...")

	port := os.Getenv("PORT")
	if port == "" {
		port = "6060"
	}

	redisHostname := os.Getenv("NEWS_MANAGER_REDIS_SERVICE_HOST")
	redisPort := os.Getenv("NEWS_MANAGER_REDIS_SERVICE_PORT")

	var store nm.Store
	if redisHostname == "" {
		store = nm.NewInMemoryNewsStore()
	} else {
		address := fmt.Sprintf("%s:%s", redisHostname, redisPort)
		store, err := nm.NewRedisNewsStore(address)
		if err != nil {
			log.Fatal(err)
		}
	}

	natsHostname := os.Getenv("NATS_CLUSTER_SERVICE_HOST")
	natsPort := os.Getenv("NATS_CLUSTER_SERVICE_PORT")

	svc, err := nm.NewNewsManager(store, natsHostname, natsPort)
	if err != nil {
		log.Fatal(err)
	}

	getNewsHandler := httptransport.NewServer(
		makeGetNewsEndpoint(svc),
		decodeGetNewsRequest,
		encodeGetNewsResponse,
	)

	r := mux.NewRouter()
	r.Methods("GET").Path("/news/{username}").Handler(getNewsHandler)

	log.Printf("Listening on port %s...\n", port)
	log.Fatal(http.ListenAndServe(":"+port, r))
}
