package main

import (
	"log"
	"net/http"
	"os"

	"github.com/coreos/go-etcd/etcd"
	"github.com/jkakar/switchboard"
)

func main() {
	// Initialize the exchange.
	client := etcd.NewClient([]string{"http://127.0.0.1:4001"})
	mux := switchboard.NewExchangeServeMux()
	exchange := switchboard.NewExchange("example", client, mux)
	exchange.Init()

	// Watch for service changes in etcd.
	go func() {
		log.Print("Watching for service configuration changes in etcd")
		stop := make(chan bool)
		exchange.Watch(stop)
	}()

	// Listen for HTTP requests.
	port := os.Getenv("PORT")
	log.Printf("Listening for HTTP requests on port %v", port)
	err := http.ListenAndServe("localhost:"+port, Log(mux))
	if err != nil {
		log.Print(err)
	}
	log.Print("Shutting down")
}

func Log(handler http.Handler) http.Handler {
	wrapper := func(writer http.ResponseWriter, request *http.Request) {
		log.Printf("%s %s %s", request.RemoteAddr, request.Method, request.URL)
		handler.ServeHTTP(writer, request)
	}
	return http.HandlerFunc(wrapper)
}
