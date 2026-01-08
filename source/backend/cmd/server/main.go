package main

import (
	"log"
	"net/http"
	"os"

	"github.com/joho/godotenv"
	"github.com/school-monitoring/backend/internal/api"
	"github.com/school-monitoring/backend/internal/database"
	"github.com/school-monitoring/backend/internal/services/maintenance"
	"github.com/school-monitoring/backend/internal/services/notifications"
	"github.com/school-monitoring/backend/internal/websocket"
)

func main() {
	// Cargar variables de entorno desde .env
	if err := godotenv.Load(); err != nil {
		// Intentar desde la raiz del proyecto
		if err := godotenv.Load("../../.env"); err != nil {
			log.Println("No .env file found, using environment variables")
		}
	}

	// Inicializar base de datos
	db, err := database.Initialize()
	if err != nil {
		log.Fatal("Failed to connect to database:", err)
	}

	log.Println("Database initialized successfully")

	// Inicializar WebSocket hub
	hub := websocket.NewHub()
	go hub.Run()

	// Worker de notificaciones (outbox)
	stop := make(chan struct{})
	notifWorker := notifications.NewWorker(db, hub)
	go notifWorker.Run(stop)

	// Retención (limpieza periódica)
	go maintenance.RunRetention(db, stop)

	// Crear router
	router := api.NewRouter(db, hub)

	// Obtener puerto
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	log.Printf("Server starting on port %s", port)
	log.Printf("API available at http://localhost:%s/api/v1", port)
	log.Printf("WebSocket available at ws://localhost:%s/ws", port)

	if err := http.ListenAndServe(":"+port, router); err != nil {
		log.Fatal("Server error:", err)
	}
}
