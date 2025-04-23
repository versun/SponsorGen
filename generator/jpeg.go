package generator

import (
        "fmt"
        "os"
        "os/exec"
        "strconv"
)

// GenerateJPEG generates a JPEG file from the SVG sponsors image using ImageMagick
func GenerateJPEG(svgPath, jpegPath string, quality int) error {
        // Check if the SVG file exists
        if _, err := os.Stat(svgPath); os.IsNotExist(err) {
                return fmt.Errorf("SVG file does not exist at path %s: %w", svgPath, err)
        }

        // Set JPEG quality (1-100)
        if quality < 1 {
                quality = 75 // default quality
        } else if quality > 100 {
                quality = 100 // max quality
        }

        // Build convert command
        // ImageMagick's convert command will handle the SVG to JPEG conversion
        cmd := exec.Command(
                "convert",
                "-quality", strconv.Itoa(quality),
                svgPath,
                jpegPath,
        )

        // Run the command
        output, err := cmd.CombinedOutput()
        if err != nil {
                return fmt.Errorf("failed to convert SVG to JPEG: %s (%w)", string(output), err)
        }

        // Verify the output file was created
        if _, err := os.Stat(jpegPath); os.IsNotExist(err) {
                return fmt.Errorf("conversion completed but JPEG file not found at %s", jpegPath)
        }

        return nil
}