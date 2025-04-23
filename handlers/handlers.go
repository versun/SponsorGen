package handlers

import (
        "encoding/json"
        "fmt"
        "log"
        "net/http"
        "os"
        "path/filepath"
        "sync"
        "time"

        "sponsorgen/config"
        "sponsorgen/generator"
        "sponsorgen/sponsors"
)

// Handler manages HTTP handlers for the sponsorkit server
type Handler struct {
        Config         config.Config
        lastGeneration time.Time
        sponsors       []sponsors.Sponsor
        mutex          sync.RWMutex
}

// NewHandler creates a new handler with the given configuration
func NewHandler(cfg config.Config) *Handler {
        return &Handler{
                Config:         cfg,
                lastGeneration: time.Time{},
                sponsors:       []sponsors.Sponsor{},
                mutex:          sync.RWMutex{},
        }
}

// IndexHandler handles the root path
func (h *Handler) IndexHandler(w http.ResponseWriter, r *http.Request) {
        // Get SVG content directly
        var svgContent string
        svgPath := filepath.Join(h.Config.OutputDir, "sponsors.svg")

        // Check if the SVG file exists
        if _, err := os.Stat(svgPath); !os.IsNotExist(err) {
                // Read the SVG file content
                svgBytes, err := os.ReadFile(svgPath)
                if err == nil {
                        svgContent = string(svgBytes)
                }
        }

        // If SVG content is empty (file doesn't exist or read error), use a placeholder
        if svgContent == "" {
                svgContent = `<svg xmlns="http://www.w3.org/2000/svg" width="800" height="100">
                    <text x="50%" y="50%" text-anchor="middle" dominant-baseline="middle">
                        SVG content not available
                    </text>
                </svg>`
        }

        w.Header().Set("Content-Type", "text/html")
        html := `
<!DOCTYPE html>
<html lang="en">
<head>
    <meta charset="UTF-8">
    <meta name="viewport" content="width=device-width, initial-scale=1.0">
    <title>SponsorGen</title>
    <style>
        body {
            font-family: system-ui, -apple-system, 'Segoe UI', Roboto, Ubuntu, Cantarell, 'Noto Sans', sans-serif;
            max-width: 800px;
            margin: 0 auto;
            padding: 2rem;
            line-height: 1.6;
        }
        h1 { color: #333; }
        a { color: #0366d6; text-decoration: none; }
        a:hover { text-decoration: underline; }
        .links { margin-top: 2rem; }
        .links a { display: block; margin-bottom: 0.5rem; }
        pre { 
            background: #f6f8fa;
            border-radius: 6px;
            padding: 1rem;
            overflow: auto;
        }
        .sponsor-svg {
            max-width: 100%;
            height: auto;
            margin-top: 1rem;
            border: 1px solid #eaecef;
            border-radius: 6px;
            overflow: hidden;
        }
    </style>
</head>
<body>
    <h1>SponsorGen</h1>
    <p>Generate sponsor images for your GitHub/OpenCollective/Patreon/Afdian sponsors.</p>
    
    <div class="links">
        <a href="/sponsors.svg">View Sponsors SVG</a>
        <a href="/sponsors.png">View Sponsors PNG (Transparent Background)</a>
        <a href="/sponsors.jpg">View Sponsors JPEG</a>
        <a href="/sponsors.json">View Sponsors JSON</a>
    </div>
    
    <h2>Current Sponsors</h2>
    <div class="sponsor-svg">
        ` + svgContent + `
    </div>
    
    <p>Last updated: ` + h.lastGeneration.Format(time.RFC1123) + `</p>
</body>
</html>
`
        fmt.Fprint(w, html)
}

// SVGHandler serves the generated SVG
func (h *Handler) SVGHandler(w http.ResponseWriter, r *http.Request) {
        h.mutex.RLock()
        defer h.mutex.RUnlock()

        // Check if regeneration is needed
        if h.shouldRegenerate() {
                h.mutex.RUnlock()
                if err := h.GenerateSponsors(); err != nil {
                        h.mutex.RLock()
                        http.Error(w, "Failed to generate sponsor data", http.StatusInternalServerError)
                        return
                }
                h.mutex.RLock()
        }

        // Serve the SVG file
        svgPath := filepath.Join(h.Config.OutputDir, "sponsors.svg")
        if _, err := os.Stat(svgPath); os.IsNotExist(err) {
                http.Error(w, "SVG file not found", http.StatusNotFound)
                return
        }

        w.Header().Set("Content-Type", "image/svg+xml")
        w.Header().Set("Cache-Control", "no-cache, max-age=0")
        http.ServeFile(w, r, svgPath)
}

// JSONHandler serves the generated JSON
func (h *Handler) JSONHandler(w http.ResponseWriter, r *http.Request) {
        h.mutex.RLock()
        defer h.mutex.RUnlock()

        // Check if regeneration is needed
        if h.shouldRegenerate() {
                h.mutex.RUnlock()
                if err := h.GenerateSponsors(); err != nil {
                        h.mutex.RLock()
                        http.Error(w, "Failed to generate sponsor data", http.StatusInternalServerError)
                        return
                }
                h.mutex.RLock()
        }

        // Serve the JSON file
        jsonPath := filepath.Join(h.Config.OutputDir, "sponsors.json")
        if _, err := os.Stat(jsonPath); os.IsNotExist(err) {
                http.Error(w, "JSON file not found", http.StatusNotFound)
                return
        }

        w.Header().Set("Content-Type", "application/json")
        w.Header().Set("Cache-Control", "no-cache, max-age=0")
        http.ServeFile(w, r, jsonPath)
}

// JPEGHandler serves the generated JPEG
func (h *Handler) JPEGHandler(w http.ResponseWriter, r *http.Request) {
        h.mutex.RLock()
        defer h.mutex.RUnlock()

        // Check if regeneration is needed
        if h.shouldRegenerate() {
                h.mutex.RUnlock()
                if err := h.GenerateSponsors(); err != nil {
                        h.mutex.RLock()
                        http.Error(w, "Failed to generate sponsor data", http.StatusInternalServerError)
                        return
                }
                h.mutex.RLock()
        }

        // Check for SVG file
        svgPath := filepath.Join(h.Config.OutputDir, "sponsors.svg")
        if _, err := os.Stat(svgPath); os.IsNotExist(err) {
                http.Error(w, "SVG file not found", http.StatusNotFound)
                return
        }

        // JPEG path
        jpegPath := filepath.Join(h.Config.OutputDir, "sponsors.jpg")

        // Generate JPEG from SVG if needed
        if _, err := os.Stat(jpegPath); os.IsNotExist(err) {
                h.mutex.RUnlock()
                if err := generator.GenerateJPEG(svgPath, jpegPath, 90); err != nil {
                        h.mutex.RLock()
                        http.Error(w, "Failed to generate JPEG: "+err.Error(), http.StatusInternalServerError)
                        return
                }
                h.mutex.RLock()
        }

        // Serve the JPEG file
        w.Header().Set("Content-Type", "image/jpeg")
        w.Header().Set("Cache-Control", "no-cache, max-age=0")
        http.ServeFile(w, r, jpegPath)
}

// PNGHandler serves the generated PNG with transparent background
func (h *Handler) PNGHandler(w http.ResponseWriter, r *http.Request) {
        h.mutex.RLock()
        defer h.mutex.RUnlock()

        // Check if regeneration is needed
        if h.shouldRegenerate() {
                h.mutex.RUnlock()
                if err := h.GenerateSponsors(); err != nil {
                        h.mutex.RLock()
                        http.Error(w, "Failed to generate sponsor data", http.StatusInternalServerError)
                        return
                }
                h.mutex.RLock()
        }

        // Check for SVG file
        svgPath := filepath.Join(h.Config.OutputDir, "sponsors.svg")
        if _, err := os.Stat(svgPath); os.IsNotExist(err) {
                http.Error(w, "SVG file not found", http.StatusNotFound)
                return
        }

        // PNG path
        pngPath := filepath.Join(h.Config.OutputDir, "sponsors.png")

        // Generate PNG from SVG if needed
        if _, err := os.Stat(pngPath); os.IsNotExist(err) {
                h.mutex.RUnlock()
                if err := generator.GeneratePNG(svgPath, pngPath, 90); err != nil {
                        h.mutex.RLock()
                        http.Error(w, "Failed to generate PNG: "+err.Error(), http.StatusInternalServerError)
                        return
                }
                h.mutex.RLock()
        }

        // Serve the PNG file
        w.Header().Set("Content-Type", "image/png")
        w.Header().Set("Cache-Control", "no-cache, max-age=0")
        http.ServeFile(w, r, pngPath)
}

// RefreshHandler forces a regeneration of sponsor data
func (h *Handler) RefreshHandler(w http.ResponseWriter, r *http.Request) {
        // Only allow POST for refreshes
        if r.Method != http.MethodGet && r.Method != http.MethodPost {
                http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
                return
        }

        if err := h.GenerateSponsors(); err != nil {
                http.Error(w, "Failed to refresh sponsor data: "+err.Error(), http.StatusInternalServerError)
                return
        }

        // Redirect back to index
        http.Redirect(w, r, "/", http.StatusSeeOther)
}

// GenerateSponsors fetches sponsor data and generates SVG and JSON files
func (h *Handler) GenerateSponsors() error {
        h.mutex.Lock()
        defer h.mutex.Unlock()

        log.Println("Fetching sponsor data...")

        // Create cache directory if it doesn't exist
        if err := os.MkdirAll(h.Config.CacheDir, 0755); err != nil {
                return fmt.Errorf("failed to create cache directory: %w", err)
        }

        // Create output directory if it doesn't exist
        if err := os.MkdirAll(h.Config.OutputDir, 0755); err != nil {
                return fmt.Errorf("failed to create output directory: %w", err)
        }

        // Collect sponsors from different sources
        var allSponsors []sponsors.Sponsor
        var wg sync.WaitGroup
        var mu sync.Mutex
        var errors []error

        // GitHub sponsors
        if h.Config.GitHubToken != "" && h.Config.GitHubLogin != "" {
                wg.Add(1)
                go func() {
                        defer wg.Done()
                        ghSponsors, err := sponsors.FetchGitHubSponsors(h.Config)
                        if err != nil {
                                mu.Lock()
                                errors = append(errors, fmt.Errorf("GitHub sponsors: %w", err))
                                mu.Unlock()
                                return
                        }
                        mu.Lock()
                        allSponsors = append(allSponsors, ghSponsors...)
                        mu.Unlock()
                }()
        }

        // OpenCollective sponsors
        if h.Config.OpenCollectiveSlug != "" {
                wg.Add(1)
                go func() {
                        defer wg.Done()
                        ocSponsors, err := sponsors.FetchOpenCollectiveSponsors(h.Config)
                        if err != nil {
                                mu.Lock()
                                errors = append(errors, fmt.Errorf("OpenCollective sponsors: %w", err))
                                mu.Unlock()
                                return
                        }
                        mu.Lock()
                        allSponsors = append(allSponsors, ocSponsors...)
                        mu.Unlock()
                }()
        }

        // Patreon sponsors
        if h.Config.PatreonToken != "" && h.Config.PatreonCampaignID != "" {
                wg.Add(1)
                go func() {
                        defer wg.Done()
                        patreonSponsors, err := sponsors.FetchPatreonSponsors(h.Config)
                        if err != nil {
                                mu.Lock()
                                errors = append(errors, fmt.Errorf("Patreon sponsors: %w", err))
                                mu.Unlock()
                                return
                        }
                        mu.Lock()
                        allSponsors = append(allSponsors, patreonSponsors...)
                        mu.Unlock()
                }()
        }

        // Afdian sponsors
        if h.Config.AfdianUserID != "" && h.Config.AfdianToken != "" {
                wg.Add(1)
                go func() {
                        defer wg.Done()
                        afdianSponsors, err := sponsors.FetchAfdianSponsors(h.Config)
                        if err != nil {
                                mu.Lock()
                                errors = append(errors, fmt.Errorf("Afdian sponsors: %w", err))
                                mu.Unlock()
                                return
                        }
                        mu.Lock()
                        allSponsors = append(allSponsors, afdianSponsors...)
                        mu.Unlock()
                }()
        }

        // Wait for all fetchers to complete
        wg.Wait()

        // Check for errors
        if len(errors) > 0 {
                errorMsg := "Errors fetching sponsors:\n"
                for _, err := range errors {
                        errorMsg += "- " + err.Error() + "\n"
                }
                // If we have some sponsors, continue despite errors
                if len(allSponsors) == 0 {
                        return fmt.Errorf(errorMsg)
                }
                log.Println(errorMsg)
        }

        // Apply exclusions and inclusions from config
        allSponsors = sponsors.ApplyFilters(allSponsors, h.Config)

        // Override amounts if specified in config
        for i, sponsor := range allSponsors {
                if amount, ok := h.Config.ForceSponsorAmounts[sponsor.Login]; ok {
                        allSponsors[i].MonthlyAmount = amount
                }
        }

        log.Printf("Found %d sponsors after filtering", len(allSponsors))

        // Generate SVG
        svgPath := filepath.Join(h.Config.OutputDir, "sponsors.svg")
        if err := generator.GenerateSVG(allSponsors, h.Config, svgPath); err != nil {
                return fmt.Errorf("failed to generate SVG: %w", err)
        }
        
        // Remove any existing image files to force regeneration
        jpegPath := filepath.Join(h.Config.OutputDir, "sponsors.jpg")
        if _, err := os.Stat(jpegPath); err == nil {
                if err := os.Remove(jpegPath); err != nil {
                        log.Printf("Warning: Failed to remove existing JPEG file: %v", err)
                }
        }
        
        pngPath := filepath.Join(h.Config.OutputDir, "sponsors.png")
        if _, err := os.Stat(pngPath); err == nil {
                if err := os.Remove(pngPath); err != nil {
                        log.Printf("Warning: Failed to remove existing PNG file: %v", err)
                }
        }

        // Generate JSON
        jsonPath := filepath.Join(h.Config.OutputDir, "sponsors.json")
        jsonFile, err := os.Create(jsonPath)
        if err != nil {
                return fmt.Errorf("failed to create JSON file: %w", err)
        }
        defer jsonFile.Close()

        // Write JSON data
        encoder := json.NewEncoder(jsonFile)
        encoder.SetIndent("", "  ")
        if err := encoder.Encode(allSponsors); err != nil {
                return fmt.Errorf("failed to encode JSON: %w", err)
        }

        // Update state
        h.sponsors = allSponsors
        h.lastGeneration = time.Now()

        log.Printf("Generated sponsors SVG and JSON successfully (PNG and JPEG will be generated on first request)")
        return nil
}

// shouldRegenerate checks if sponsor data should be regenerated
func (h *Handler) shouldRegenerate() bool {
        if h.lastGeneration.IsZero() {
                return true
        }

        refreshInterval := time.Duration(h.Config.RefreshMinutes) * time.Minute
        return time.Since(h.lastGeneration) > refreshInterval
}
