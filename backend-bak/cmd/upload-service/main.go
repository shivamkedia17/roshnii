package main

import (
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/gin-gonic/gin"

	"github.com/shivamkedia17/roshnii-backend/internal/config"
)

func main() {
	log.Println("Starting Upload Service...")

	// Load configuration
	cfg, err := config.LoadConfig(".")
	if err != nil {
		log.Printf("FATAL: Error loading configuration: %v", err)
		os.Exit(1)
	}

	log.Printf("Configuration loaded successfully.")
	log.Printf("Service Port: %s", cfg.ServerPort)
	log.Printf("Log Level: %s", cfg.LogLevel)
	log.Printf("Storage Type: %s", cfg.BlobStorageType)
	if cfg.BlobStorageType == "local" {
		log.Printf("Local Storage Path: %s", cfg.LocalstoragePath)
	}

	// Set Gin mode based on environment
	if cfg.Environment == "production" {
		gin.SetMode(gin.ReleaseMode)
	} else {
		gin.SetMode(gin.DebugMode)
	}

	router := gin.Default() // Default includes logger and recovery middleware

	// --- Define Basic Routes ---
	// Simple health check / ping route
	router.GET("/ping", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"message": "pong",
			"service": "upload-service",
		})
	})

	// --- Placeholder for future steps ---
	// TODO: Initialize database connections (PostgreSQL)
	// TODO: Initialize blob storage client
	// TODO: Initialize message queue publisher (optional)
	// TODO: Define *actual* API routes and handlers (e.g., POST /api/v1/upload)
	// TODO: Inject dependencies (config, db, storage) into handlers

	// --- Start the HTTP server ---
	listenAddr := fmt.Sprintf(":%s", cfg.ServerPort)
	log.Printf("Upload Service starting to listen on %s", listenAddr)

	// Start listening for incoming connections
	// router.Run blocks the main goroutine until the server is shut down
	err = router.Run(listenAddr)
	if err != nil {
		log.Fatalf("FATAL: Failed to start Gin server: %v", err)
		os.Exit(1) // Exit if server fails to start
	}

	// Code below router.Run() will likely not execute unless the server stops gracefully
	log.Println("Upload Service stopped.")
}
