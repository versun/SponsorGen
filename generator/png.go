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

	// Build convert command
	// ImageMagick's convert command will handle the SVG to PNG conversion
	// We use -background transparent to ensure a transparent background
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
		return fmt.Errorf("failed to convert SVG to PNG: %s (%w)", string(output), err)
	}

	// Verify the output file was created
	if _, err := os.Stat(pngPath); os.IsNotExist(err) {
		return fmt.Errorf("conversion completed but PNG file not found at %s", pngPath)
	}

	return nil
}