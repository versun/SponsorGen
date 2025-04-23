package main

import (
        "flag"
        "fmt"
        "log"
        "net/http"
        "os"
        "time"

        "sponsorgen/config"
        "sponsorgen/handlers"
)

// scheduleMidnightRefresh sets up a scheduler to refresh sponsor data at midnight (00:00) every day
func scheduleMidnightRefresh(handler *handlers.Handler) {
        for {
                now := time.Now()
                // Calculate time until next midnight
                midnight := time.Date(now.Year(), now.Month(), now.Day(), 0, 0, 0, 0, now.Location()).Add(24 * time.Hour)
                duration := midnight.Sub(now)
                
                // Sleep until next midnight
                log.Printf("Next scheduled refresh in %s at %s", duration.Round(time.Second), midnight.Format("2006-01-02 15:04:05"))
                time.Sleep(duration)
                
                // Refresh the sponsors at midnight
                log.Println("Executing scheduled midnight refresh...")
                if err := handler.GenerateSponsors(); err != nil {
                        log.Printf("Warning: Scheduled refresh failed: %v", err)
                } else {
                        log.Println("Scheduled refresh completed successfully")
                }
        }
}

func main() {
        // Define command line flags
        port := flag.Int("port", 8000, "Port to serve on")
        flag.Parse()

        // Load configuration from environment variables
        cfg, err := config.LoadConfig()
        if err != nil {
                log.Fatalf("Failed to load configuration: %v", err)
        }

        // Create output directory if it doesn't exist
        if err := os.MkdirAll(cfg.OutputDir, 0755); err != nil {
                log.Fatalf("Failed to create output directory: %v", err)
        }

        // Setup HTTP handler
        handler := handlers.NewHandler(cfg)

        // Register handlers
        http.HandleFunc("/", handler.IndexHandler)
        http.HandleFunc("/sponsors.svg", handler.SVGHandler)
        http.HandleFunc("/sponsors.json", handler.JSONHandler)
        http.HandleFunc("/sponsors.jpg", handler.JPEGHandler)
        http.HandleFunc("/refresh", handler.RefreshHandler)

        // Serve static files
        fs := http.FileServer(http.Dir(cfg.OutputDir))
        http.Handle("/static/", http.StripPrefix("/static/", fs))

        // Start server
        addr := fmt.Sprintf("0.0.0.0:%d", *port)
        log.Printf("SponsorGen server starting on %s", addr)
        log.Printf("Configuration loaded from environment variables")
        log.Printf("Serving SVG at http://localhost:%d/sponsors.svg", *port)
        log.Printf("Serving JSON at http://localhost:%d/sponsors.json", *port)
        log.Printf("Serving JPEG at http://localhost:%d/sponsors.jpg", *port)
        log.Printf("Force refresh with http://localhost:%d/refresh", *port)

        // Generate initial sponsor data
        if err := handler.GenerateSponsors(); err != nil {
                log.Printf("Warning: Failed to generate initial sponsor data: %v", err)
        }

        // Setup daily refresh at midnight (00:00)
        go scheduleMidnightRefresh(handler)
        log.Println("Scheduled daily refresh at 00:00")

        // Start HTTP server
        if err := http.ListenAndServe(addr, nil); err != nil {
                log.Fatalf("Server failed: %v", err)
        }
}
