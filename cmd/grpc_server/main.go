package main

import (
	"context"
	"log"

	"github.com/merynayr/auth/internal/app"
)

// var configPath string

// func init() {
// 	flag.StringVar(&configPath, "config-path", ".env", "path to config file")
// }

func main() {
	// flag.Parse()
	ctx := context.Background()

	a, err := app.NewApp(ctx)
	if err != nil {
		log.Fatalf("failed to init app: %s", err.Error())
	}

	err = a.Run()
	if err != nil {
		log.Fatalf("failed to run app: %s", err.Error())
	}
}
