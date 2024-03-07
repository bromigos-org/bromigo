package run

import (
	"log"
	"net/http"
)

func healthCheckHandler(w http.ResponseWriter, r *http.Request) {
	// Here, you can include checks to verify the bot's health, such as:
	// - Checking if the bot is connected to Discord
	// - Verifying critical components or services the bot relies on are operational

	// If everything is okay, send a 200 OK status
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("OK"))
}

func StartHTTPServer() {
	http.HandleFunc("/health", healthCheckHandler) // Route to handle health check

	go func() {
		// Replace "80" with your preferred port
		if err := http.ListenAndServe(":8080", nil); err != nil {
			log.Fatalf("Failed to start HTTP server: %v", err)
		}
	}()
}
