package generator

import (
        "fmt"
        "os"
        "os/exec"
        "strconv"
)

// GeneratePNG generates a PNG file with transparent background from the SVG sponsors image using ImageMagick
func GeneratePNG(svgPath, pngPath string, quality int) error {
        // Check if the SVG file exists
        if _, err := os.Stat(svgPath); os.IsNotExist(err) {
                return fmt.Errorf("SVG file does not exist at path %s: %w", svgPath, err)
        }

        // First try with 'magick' command (newer ImageMagick v7)
        magickCmd := exec.Command(
                "magick",
                "convert",
                "-background", "transparent",
                "-quality", strconv.Itoa(quality),
                svgPath,
                pngPath,
        )

        magickOutput, magickErr := magickCmd.CombinedOutput()
        if magickErr == nil {
                // Check if file was created successfully
                if _, err := os.Stat(pngPath); err == nil {
                        return nil
                }
        }

        // If 'magick' command fails, try with legacy 'convert' command (IMv6)
        cmd := exec.Command(
                "convert",
                "-background", "transparent",
                "-quality", strconv.Itoa(quality),
                svgPath,
                pngPath,
        )

        // Run the command
        output, err := cmd.CombinedOutput()
        if err != nil {
                // Include both error outputs in the error message
                return fmt.Errorf("failed to convert SVG to PNG: %s (%w). Magick error: %s", 
                        string(output), err, string(magickOutput))
        }

        // Verify the output file was created
        if _, err := os.Stat(pngPath); os.IsNotExist(err) {
                return fmt.Errorf("conversion completed but PNG file not found at %s", pngPath)
        }

        return nil
}