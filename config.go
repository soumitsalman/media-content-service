package main

import "os"

func getPort() string {
	if port := os.Getenv("SERVER_PORT"); port != "" {
		return ":" + port
	}
	return ":8080"
}

func getInternalAuthToken() string {
	return os.Getenv("INTERNAL_AUTH_TOKEN")
}
