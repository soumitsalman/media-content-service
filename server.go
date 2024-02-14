package main

import (
	"github.com/soumitsalman/media-content-service/api"
)

func main() {
	api.NewServer(100, 1000).Run()
}
