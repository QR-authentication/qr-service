package main

import (
	"log"

	"github.com/QR-authentication/qr-service/internal/config"
)

func main() {
	cfg := config.MustLoad()
	log.Println(cfg)
}
