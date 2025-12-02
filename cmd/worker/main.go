package main

import (
	"log"
	"os"
	"strconv"

	"github.com/Luzin7/alert-service/internal/infra/smtp"

	"github.com/Luzin7/alert-service/internal/infra/database"
	"github.com/Luzin7/alert-service/internal/infra/messenger"
	"github.com/Luzin7/alert-service/internal/infra/providers"
	"github.com/Luzin7/alert-service/internal/transport/consumer"
	"github.com/Luzin7/alert-service/internal/usecases"
	"github.com/joho/godotenv"
)

func main() {
	if os.Getenv("ENV") != "production" {
		err := godotenv.Load(".env")

		if err != nil {
			log.Fatal("Error loading .env file")
		}
	}

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	databaseConnectionString := os.Getenv("DATABASE_URL")
	if databaseConnectionString == "" {
		log.Fatal("DATABASE_URL is not set")
	}

	db, err := database.DatabaseConnection(databaseConnectionString)
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}

	// cacheAddr := os.Getenv("CACHE_ADDR")
	// cachePassword := os.Getenv("CACHE_PASSWORD")
	// cacheUsername := os.Getenv("CACHE_USERNAME")
	// cacheDB := os.Getenv("CACHE_DB")

	// cacheConn, err := cache.CacheConnection(cacheAddr, cachePassword, cacheUsername, atoi(cacheDB))
	// if err != nil {
	// 	log.Fatalf("Failed to connect to cache: %v", err)
	// }

	messengerUsername := os.Getenv("MESSENGER_USERNAME")
	messengerPassword := os.Getenv("MESSENGER_PASSWORD")
	messengerHost := os.Getenv("MESSENGER_HOST")
	messengerPort := os.Getenv("MESSENGER_PORT")
	if messengerUsername == "" || messengerPassword == "" || messengerHost == "" || messengerPort == "" {
		log.Fatal("Messenger configuration is not set properly")
	}

	messengerConn, err := messenger.MessengerConnection(messengerUsername, messengerPassword, messengerHost, messengerPort)
	if err != nil {
		log.Fatalf("Failed to connect to messenger: %v", err)
	}

	repo := database.NewRepository(db)
	smtpServer := os.Getenv("SMTP_SERVER")
	smtpPortStr := os.Getenv("SMTP_PORT")
	smtpPort, err := strconv.Atoi(smtpPortStr)
	if err != nil {
		log.Fatalf("Invalid SMTP_PORT: %v", err)
	}
	smtpUsername := os.Getenv("SMTP_USERNAME")
	smtpPassword := os.Getenv("SMTP_PASSWORD")
	senderConn, err := smtp.SMTPConnection(smtpServer, smtpPort, smtpUsername, smtpPassword)
	if err != nil {
		log.Fatalf("Failed to connect to SMTP: %v", err)
	}

	linkGenerator := providers.GoogleFlightsGenerator{
		BaseURL: "https://www.google.com/travel/flights",
	}

	processAlertUseCase := usecases.NewProcessAlert(linkGenerator, repo, senderConn)

	handler := consumer.NewHandler(processAlertUseCase)

	worker := consumer.NewWorker(messengerConn, handler)

	queueName := os.Getenv("QUEUE_NAME")
	if queueName == "" {
		queueName = "price-alerts"
	}

	log.Printf("Starting worker on queue: %s", queueName)
	worker.Start(queueName)
}
