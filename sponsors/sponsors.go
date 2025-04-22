package sponsors

import (
        "strings"

        "sponsorgen/config"
)

// Sponsor represents a sponsor from any platform
type Sponsor struct {
        ID            string  `json:"id"`
        Name          string  `json:"name"`
        Login         string  `json:"login"`
        AvatarURL     string  `json:"avatarUrl"`
        Link          string  `json:"link"`
        Platform      string  `json:"platform"` // github, opencollective, patreon, afdian
        MonthlyAmount float64 `json:"monthlyAmount"`
        CreatedAt     string  `json:"createdAt"`
        TierName      string  `json:"tierName,omitempty"`
}

// ApplyFilters applies exclusion and inclusion filters from the config
func ApplyFilters(sponsors []Sponsor, cfg config.Config) []Sponsor {
        var filtered []Sponsor

        // Create lookup maps for faster checking
        excludeMap := make(map[string]bool)
        includeMap := make(map[string]bool)

        for _, login := range cfg.ExcludeSponsors {
                excludeMap[strings.ToLower(login)] = true
        }

        for _, login := range cfg.IncludeSponsors {
                includeMap[strings.ToLower(login)] = true
        }

        for _, sponsor := range sponsors {
                lowerLogin := strings.ToLower(sponsor.Login)

                // Skip if in exclude list, unless also in include list (include takes precedence)
                if excludeMap[lowerLogin] && !includeMap[lowerLogin] {
                        continue
                }

                filtered = append(filtered, sponsor)
        }

        return filtered
}

// MergeDuplicates combines sponsors with the same login across platforms
func MergeDuplicates(sponsors []Sponsor) []Sponsor {
        merged := make(map[string]Sponsor)

        for _, sponsor := range sponsors {
                key := strings.ToLower(sponsor.Login)
                
                if existing, found := merged[key]; found {
                        // Add the monthly amounts
                        existing.MonthlyAmount += sponsor.MonthlyAmount
                        
                        // Keep the earliest creation date
                        if sponsor.CreatedAt < existing.CreatedAt {
                                existing.CreatedAt = sponsor.CreatedAt
                        }
                        
                        // Update platforms
                        if !strings.Contains(existing.Platform, sponsor.Platform) {
                                existing.Platform = existing.Platform + "," + sponsor.Platform
                        }
                        
                        merged[key] = existing
                } else {
                        merged[key] = sponsor
                }
        }

        // Convert map back to slice
        result := make([]Sponsor, 0, len(merged))
        for _, sponsor := range merged {
                result = append(result, sponsor)
        }

        return result
}

// SortSponsors sorts sponsors by amount (descending) and then by creation date
func SortSponsors(sponsors []Sponsor) []Sponsor {
        // Implementation using a simple bubble sort for clarity
        result := make([]Sponsor, len(sponsors))
        copy(result, sponsors)

        for i := 0; i < len(result); i++ {
                for j := i + 1; j < len(result); j++ {
                        // Sort by creation date (descending - newer sponsors first)
                        if result[i].CreatedAt < result[j].CreatedAt {
                                result[i], result[j] = result[j], result[i]
                        }
                }
        }

        return result
}


