package handlers

import (
	"context"
	"fmt"
	"net/http"
)

var ServerInstance *http.Server

func StopServer(w http.ResponseWriter, r *http.Request) {
	fmt.Printf("\"Stop server starts\": %v\n", "Stop server starts")
	resp := map[string]string{
		"result": "ServerShutdown success",
	}
	encodeAndSendJSON(&resp, &w)
	go ServerInstance.Shutdown(context.TODO())
}
