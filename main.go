package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"time"

	"github.com/gorilla/mux"
	"github.com/javokhirbek1999/microservices-in-go/handlers"
)

func main() {

	l := log.New(os.Stdout, "product-api", log.LstdFlags)

	// create handlers
	ph := handlers.NewProduct(l)

	// create a new serve mux and register the hanlders
	serveMux := mux.NewRouter()

	getRouter := serveMux.Methods(http.MethodGet).Subrouter()
	getRouter.HandleFunc("/", ph.GetProducts)

	putRouter := serveMux.Methods(http.MethodPut).Subrouter()
	putRouter.HandleFunc("/{id:[0-9]+}", ph.UpdateProduct)
	putRouter.Use(ph.MiddlewareProductValidation)

	postRouter := serveMux.Methods(http.MethodPost).Subrouter()
	postRouter.HandleFunc("/", ph.AddProduct)
	postRouter.Use(ph.MiddlewareProductValidation)

	server := &http.Server{
		Addr:         ":9090",
		Handler:      serveMux,
		IdleTimeout:  120 * time.Second, // Keep the connection open with current client for this long
		ReadTimeout:  1 * time.Second,   // Max read-timeout
		WriteTimeout: 1 * time.Second,   // Max write-timeout
	}

	go func() {
		err := server.ListenAndServe()

		if err != nil {
			l.Fatal(err)
		}

	}()

	fmt.Println("Listening to port 9090 on localhost")
	sigChan := make(chan os.Signal)
	signal.Notify(sigChan, os.Interrupt)
	signal.Notify(sigChan, os.Kill)

	sig := <-sigChan // Signal sends the message to channel and channel assigns the message value in this variable
	l.Println("Received terminate, graceful shutdown", sig)

	timeoutContext, _ := context.WithTimeout(context.Background(), 30*time.Second) // context with max wait time to shutdown for incomplete background requests

	server.Shutdown(timeoutContext)

}
