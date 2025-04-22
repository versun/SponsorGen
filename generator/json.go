package generator

import (
        "encoding/json"
        "os"

        "sponsorgen/sponsors"
)

// JSONSponsor represents a sponsor in the JSON output
type JSONSponsor struct {
        Name          string  `json:"name"`
        Login         string  `json:"login"`
        AvatarURL     string  `json:"avatarUrl"`
        Link          string  `json:"link"`
        Platform      string  `json:"platform"`
        MonthlyAmount float64 `json:"monthlyAmount"`
        CreatedAt     string  `json:"createdAt"`
        TierName      string  `json:"tierName"`
}

// GenerateJSON generates a JSON file with sponsor data
func GenerateJSON(sponsors []sponsors.Sponsor, outputPath string) error {
        // Create the output file
        file, err := os.Create(outputPath)
        if err != nil {
                return err
        }
        defer file.Close()

        // Create a JSON encoder
        encoder := json.NewEncoder(file)
        encoder.SetIndent("", "  ")

        // Encode the sponsors
        return encoder.Encode(sponsors)
}
