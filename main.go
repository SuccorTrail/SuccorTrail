package main

import (
	"log"
	"net/http"
	"os"

	"github.com/SuccorTrail/SuccorTrail/internal/db"
	"github.com/SuccorTrail/SuccorTrail/internal/router"
	"github.com/SuccorTrail/SuccorTrail/internal/util"
	"github.com/joho/godotenv"
	"github.com/sirupsen/logrus"
)

func main() {
	// Load environment variables first
	if err := godotenv.Load(); err != nil {
		logrus.Info("No .env file found, using system environment variables")
	}

	// Init logger
	logrus.SetFormatter(&logrus.JSONFormatter{})
	logrus.SetOutput(os.Stdout)
	logrus.SetLevel(logrus.InfoLevel)

	// Init application
	if err := util.InitApp(); err != nil {
		logrus.Fatal(err)
	}

	// Init db
	err := db.InitDB()
	if err != nil {
		logrus.Fatal(err)
	}
	defer db.CloseDB()

	r := router.InitRouter()

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	logrus.Infof("Server starting on port %s...", port)
	log.Fatal(http.ListenAndServe(":"+port, r))
}
