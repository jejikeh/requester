package main

import (
	"fmt"
	"log"
	"net/http"

	"github.com/jejikeh/requester/internal"
	"github.com/jejikeh/requester/routers"
)

func main() {
	client := internal.NewClient()

	taskManager := internal.NewInMemoryTaskManager(client)

	router := routers.NewRouter(taskManager)

	host := fmt.Sprintf(":%d", 8080)

	log.Printf("Listening on port %s...\n", host)

	log.Fatal(http.ListenAndServe(host, router))
}
