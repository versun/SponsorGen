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

        // First try with 'magick' command (newer ImageMagick v7)
        magickCmd := exec.Command(
                "magick",
                "convert",
                "-quality", strconv.Itoa(quality),
                svgPath,
                jpegPath,
        )

        magickOutput, magickErr := magickCmd.CombinedOutput()
        if magickErr == nil {
                // Check if file was created successfully
                if _, err := os.Stat(jpegPath); err == nil {
                        return nil
                }
        }

        // If 'magick' command fails, try with legacy 'convert' command (IMv6)
        cmd := exec.Command(
                "convert",
                "-quality", strconv.Itoa(quality),
                svgPath,
                jpegPath,
        )

        // Run the command
        output, err := cmd.CombinedOutput()
        if err != nil {
                // Include both error outputs in the error message
                return fmt.Errorf("failed to convert SVG to JPEG: %s (%w). Magick error: %s", 
                        string(output), err, string(magickOutput))
        }

        // Verify the output file was created
        if _, err := os.Stat(jpegPath); os.IsNotExist(err) {
                return fmt.Errorf("conversion completed but JPEG file not found at %s", jpegPath)
        }

        return nil
}