package generator

import (
        "bytes"
        "fmt"
        "log"
        "math"
        "os"
        "text/template"

        "sponsorgen/config"
        "sponsorgen/sponsors"
        "sponsorgen/utils"
)

// SVGData represents the data to be passed to the SVG template
type SVGData struct {
        Width           int
        Height          int
        AvatarSize      int
        FontSize        int
        FontFamily      string
        ShowAmount      bool
        ShowName        bool
        BackgroundColor string
        PaddingX        int
        PaddingY        int
        Sponsors        []SponsorData
}

// SponsorData represents a sponsor in the SVG
type SponsorData struct {
        Name    string
        Avatar  string
        Link    string
        Amount  string
        X       int
        Y       int
        Size    int
        NameX   int
        NameY   int
        AmountX int
        AmountY int
}

// GenerateSVG generates an SVG file for the sponsors
func GenerateSVG(allSponsors []sponsors.Sponsor, cfg config.Config, outputPath string) error {
        // Ensure default avatar exists
        if _, err := os.Stat(cfg.DefaultAvatar); os.IsNotExist(err) {
                if err := createDefaultAvatar(cfg.DefaultAvatar); err != nil {
                        return fmt.Errorf("failed to create default avatar: %w", err)
                }
        }
        
        // Ensure cache directory exists
        if err := os.MkdirAll(cfg.CacheDir, 0755); err != nil {
                return fmt.Errorf("failed to create cache directory: %w", err)
        }

        // Sort sponsors by creation date (newer first)
        sortedSponsors := sponsors.SortSponsors(allSponsors)

        // Calculate SVG dimensions and sponsor positions
        svgData, err := calculateSVGLayout(sortedSponsors, cfg)
        if err != nil {
                return fmt.Errorf("failed to calculate SVG layout: %w", err)
        }

        // Parse template
        tmpl, err := template.New("svg").Parse(cfg.SVGTemplate)
        if err != nil {
                return fmt.Errorf("failed to parse SVG template: %w", err)
        }

        // Generate SVG
        var svgBuffer bytes.Buffer
        if err := tmpl.Execute(&svgBuffer, svgData); err != nil {
                return fmt.Errorf("failed to execute SVG template: %w", err)
        }

        // Write to file
        if err := os.WriteFile(outputPath, svgBuffer.Bytes(), 0644); err != nil {
                return fmt.Errorf("failed to write SVG file: %w", err)
        }

        return nil
}

// calculateSVGLayout calculates the positions of sponsors in the SVG
func calculateSVGLayout(sortedSponsors []sponsors.Sponsor, cfg config.Config) (SVGData, error) {
        svgData := SVGData{
                Width:           cfg.SVGWidth,
                Height:          100, // Initial height, will be updated
                AvatarSize:      cfg.AvatarSize,
                FontSize:        cfg.FontSize,
                FontFamily:      cfg.FontFamily,
                ShowAmount:      cfg.ShowAmount,
                ShowName:        cfg.ShowName,
                BackgroundColor: cfg.BackgroundColor,
                PaddingX:        cfg.PaddingX,
                PaddingY:        cfg.PaddingY,
                Sponsors:        []SponsorData{},
        }

        // Calculate positions for sponsors
        maxY := 0
        currentX := cfg.PaddingX
        rowY := cfg.PaddingY + 10 // Small padding from top

        avatarSize := cfg.AvatarSize // Use default avatar size for all sponsors

        for _, sponsor := range sortedSponsors {
                // Skip to next row if this sponsor doesn't fit
                if currentX + avatarSize > cfg.SVGWidth - cfg.PaddingX {
                        currentX = cfg.PaddingX
                        rowY = maxY + cfg.AvatarMargin
                }

                // Prepare avatar URL
                avatarURL := sponsor.AvatarURL
                if avatarURL == "" {
                        avatarURL = cfg.DefaultAvatar
                }
                
                // Download and embed the avatar image
                embeddedAvatar, err := utils.DownloadImage(avatarURL, cfg.CacheDir)
                if err != nil {
                        log.Printf("Failed to download avatar for %s: %v, using default avatar", sponsor.Name, err)
                        embeddedAvatar, _ = utils.DownloadImage(cfg.DefaultAvatar, cfg.CacheDir)
                }

                // Format amount string
                amountStr := fmt.Sprintf("%.2f", sponsor.MonthlyAmount)

                // Create sponsor data
                sponsorData := SponsorData{
                        Name:    sponsor.Name,
                        Avatar:  embeddedAvatar, // Use the embedded avatar instead of the URL
                        Link:    sponsor.Link,
                        Amount:  amountStr,
                        X:       currentX,
                        Y:       rowY,
                        Size:    avatarSize,
                        NameX:   currentX + avatarSize + 5,
                        NameY:   rowY + avatarSize/2,
                        AmountX: currentX + avatarSize + 5,
                        AmountY: rowY + avatarSize/2 + cfg.FontSize + 2,
                }

                svgData.Sponsors = append(svgData.Sponsors, sponsorData)

                // Update position for next sponsor
                currentX += avatarSize + cfg.AvatarMargin
                maxY = int(math.Max(float64(maxY), float64(rowY+avatarSize)))
        }

        // Update SVG height
        svgData.Height = maxY + cfg.PaddingY + avatarSize

        return svgData, nil
}

// createDefaultAvatar creates a default SVG avatar
func createDefaultAvatar(path string) error {
        defaultAvatar := `<svg xmlns="http://www.w3.org/2000/svg" viewBox="0 0 100 100">
  <rect width="100" height="100" fill="#f2f2f2" />
  <text x="50" y="50" font-family="sans-serif" font-size="20" text-anchor="middle" dominant-baseline="middle" fill="#666">?</text>
</svg>`

        return os.WriteFile(path, []byte(defaultAvatar), 0644)
}
