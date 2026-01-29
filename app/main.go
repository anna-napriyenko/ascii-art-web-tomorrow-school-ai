package main

import (
	"html/template"
	"log"
	"net/http"
)

// --- main: Entry Point ---

func main() {
	// Remove timestamp (date and time) from log output
	log.SetFlags(0)

	// 1. Parse templates on startup
	indexTmpl, err := template.ParseFiles("templates/index.html")
	if err != nil {
		log.Fatalf("Failed to parse templates/index.html: %v", err)
	}
	errorTmpl, err := template.ParseFiles("templates/error.html")
	if err != nil {
		log.Fatalf("Failed to parse templates/error.html: %v", err)
	}

	// 2. Create an application instance with dependencies
	// 'application' struct is defined in web.go (same package)
	app := &application{
		indexTmpl: indexTmpl,
		errorTmpl: errorTmpl,
	}

	// 3. Create a *new* router (ServeMux)
	// This avoids using the global http.DefaultServeMux
	mux := http.NewServeMux()

	// 4. Create a file server for the './static' directory
	// This part is added to serve CSS, JS, and images
	fs := http.FileServer(http.Dir("./static"))
	// Register a handler to serve static files under the '/static/' path
	mux.Handle("/static/", http.StripPrefix("/static/", fs))

	// 5. Register routes on our new mux
	// Handlers are now methods of the 'app' struct
	mux.HandleFunc("/", app.handleIndex)
	mux.HandleFunc("/ascii", app.handleAscii)
	mux.HandleFunc("/favicon.ico", handleFavicon) // This handler is simple, no dependencies needed

	// 6. Create our own http.Server instance
	srv := &http.Server{
		Addr:    ":8080", // Hardcode the port
		Handler: mux,     // Use our custom router
	}

	// 7. Start the server by calling ListenAndServe on *our* instance
	log.Printf("Server starting on http://localhost:8080")
	err = srv.ListenAndServe()
	if err != nil {
		log.Fatal(err)
	}
}
