package main

import (
	"log"
	"net/http"
	"tiktok-oauth2/config"
	"tiktok-oauth2/handlers"
	"tiktok-oauth2/models"
	"tiktok-oauth2/utils"

	"github.com/gorilla/mux"
)

func main() {
	// Load configuration
	config.LoadConfig()

	// Create router
	router := mux.NewRouter()

	// Add CORS middleware
	router.Use(corsMiddleware)

	// Health check endpoint
	router.HandleFunc("/health", healthHandler).Methods("GET")

	// OAuth endpoints
	router.HandleFunc("/auth", handlers.AuthHandler).Methods("GET")
	router.HandleFunc("/callback", handlers.CallbackHandler).Methods("GET")
	router.HandleFunc("/refresh", handlers.RefreshTokenHandler).Methods("POST")
	router.HandleFunc("/user", handlers.UserInfoHandler).Methods("GET")

	// Start server
	port := ":" + config.ServerPort
	log.Printf("🚀 TikTok OAuth2 Server starting on port %s", config.ServerPort)
	log.Printf("📱 Auth URL: http://localhost%s/auth", port)
	log.Printf("🔄 Callback URL: %s", config.RedirectURI)

	if err := http.ListenAndServe(port, router); err != nil {
		log.Fatal("❌ Server failed to start:", err)
	}
}

// healthHandler provides a simple health check endpoint
func healthHandler(w http.ResponseWriter, r *http.Request) {
	utils.WriteJSONResponse(w, http.StatusOK, models.APIResponse{
		Success: true,
		Message: "TikTok OAuth2 Server is running",
		Data: map[string]interface{}{
			"status":  "healthy",
			"version": "1.0.0",
		},
	})
}

// corsMiddleware adds CORS headers to responses
func corsMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization")

		if r.Method == "OPTIONS" {
			w.WriteHeader(http.StatusOK)
			return
		}

		next.ServeHTTP(w, r)
	})
}
