package generator

import (
        "bytes"
        "fmt"
        "image"
        "image/jpeg"
        "os"

        "github.com/srwiley/oksvg"
        "github.com/srwiley/rasterx"
)

// GenerateJPEG generates a JPEG file from the SVG sponsors image
func GenerateJPEG(svgPath, jpegPath string, quality int) error {
        // Read the SVG file
        svgBytes, err := os.ReadFile(svgPath)
        if err != nil {
                return fmt.Errorf("failed to read SVG file: %w", err)
        }

        // Parse SVG bytes
        icon, err := oksvg.ReadIconStream(bytes.NewReader(svgBytes))
        if err != nil {
                return fmt.Errorf("failed to parse SVG: %w", err)
        }

        // Get icon dimensions
        w, h := int(icon.ViewBox.W), int(icon.ViewBox.H)
        img := image.NewRGBA(image.Rect(0, 0, w, h))

        // Create a rasterizer for SVG
        scanner := rasterx.NewScannerGV(w, h, img, img.Bounds())
        raster := rasterx.NewDasher(w, h, scanner)

        // Rasterize the SVG to the image
        icon.Draw(raster, 1.0)

        // Create the output file
        file, err := os.Create(jpegPath)
        if err != nil {
                return fmt.Errorf("failed to create JPEG file: %w", err)
        }
        defer file.Close()

        // Set JPEG quality (1-100)
        if quality < 1 {
                quality = 75 // default quality
        } else if quality > 100 {
                quality = 100 // max quality
        }

        // Encode as JPEG
        if err := jpeg.Encode(file, img, &jpeg.Options{Quality: quality}); err != nil {
                return fmt.Errorf("failed to encode JPEG: %w", err)
        }

        return nil
}