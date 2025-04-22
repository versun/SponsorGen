package sponsors

import (
        "encoding/json"
        "fmt"
        "io"
        "net/http"
        "time"

        "sponsorgen/config"
)

// PatreonResponse represents the Patreon API response
type PatreonResponse struct {
        Data []struct {
                ID         string `json:"id"`
                Attributes struct {
                        FullName     string    `json:"full_name"`
                        Email        string    `json:"email"`
                        CreatedAt    time.Time `json:"created"`
                        PatronStatus string    `json:"patron_status"`
                        LastChargeDate   *time.Time `json:"last_charge_date"`
                        LastChargeStatus string     `json:"last_charge_status"`
                        LifetimeSupportCents int    `json:"lifetime_support_cents"`
                        CurrentlyEntitledAmountCents int `json:"currently_entitled_amount_cents"`
                        PledgeRelationshipStart *time.Time `json:"pledge_relationship_start"`
                } `json:"attributes"`
                Relationships struct {
                        Currently_entitled_tiers struct {
                                Data []struct {
                                        ID   string `json:"id"`
                                        Type string `json:"type"`
                                } `json:"data"`
                        } `json:"currently_entitled_tiers"`
                } `json:"relationships"`
                Type string `json:"type"`
        } `json:"data"`
        Included []struct {
                ID         string `json:"id"`
                Attributes struct {
                        Title       string `json:"title"`
                        Description string `json:"description"`
                        AmountCents int    `json:"amount_cents"`
                } `json:"attributes"`
                Type string `json:"type"`
        } `json:"included"`
        Links struct {
                Next string `json:"next"`
        } `json:"links"`
}

// FetchPatreonSponsors fetches sponsors from Patreon
func FetchPatreonSponsors(cfg config.Config) ([]Sponsor, error) {
        sponsors := []Sponsor{}

        client := &http.Client{
                Timeout: 10 * time.Second,
        }

        // Build URL with campaign ID
        url := fmt.Sprintf("https://www.patreon.com/api/oauth2/v2/campaigns/%s/members?include=currently_entitled_tiers&fields[member]=full_name,email,patron_status,last_charge_date,last_charge_status,lifetime_support_cents,currently_entitled_amount_cents,pledge_relationship_start&fields[tier]=title,description,amount_cents", cfg.PatreonCampaignID)

        hasNextPage := true
        for hasNextPage {
                req, err := http.NewRequest("GET", url, nil)
                if err != nil {
                        return sponsors, fmt.Errorf("failed to create Patreon API request: %w", err)
                }

                req.Header.Set("Authorization", "Bearer "+cfg.PatreonToken)

                resp, err := client.Do(req)
                if err != nil {
                        return sponsors, fmt.Errorf("failed to execute Patreon API request: %w", err)
                }
                defer resp.Body.Close()

                if resp.StatusCode != http.StatusOK {
                        body, _ := io.ReadAll(resp.Body)
                        return sponsors, fmt.Errorf("Patreon API request failed with status %d: %s", resp.StatusCode, string(body))
                }

                var response PatreonResponse
                if err := json.NewDecoder(resp.Body).Decode(&response); err != nil {
                        return sponsors, fmt.Errorf("failed to decode Patreon API response: %w", err)
                }

                // Process sponsors
                for _, patron := range response.Data {
                        // Skip patrons who are not active
                        if patron.Attributes.PatronStatus != "active_patron" {
                                continue
                        }

                        // Get monthly amount in dollars
                        monthlyAmount := float64(patron.Attributes.CurrentlyEntitledAmountCents) / 100.0

                        // Extract tier name
                        tierName := ""
                        if len(patron.Relationships.Currently_entitled_tiers.Data) > 0 {
                                tierID := patron.Relationships.Currently_entitled_tiers.Data[0].ID
                                for _, tier := range response.Included {
                                        if tier.ID == tierID && tier.Type == "tier" {
                                                tierName = tier.Attributes.Title
                                                break
                                        }
                                }
                        }

                        // Use member ID as login if it's a public profile
                        login := patron.ID
                        // Name defaults to "Anonymous" if not provided
                        name := patron.Attributes.FullName
                        if name == "" {
                                name = "Anonymous Patron"
                        }

                        // For avatar URL, we use a placeholder since Patreon doesn't provide avatars in the API
                        avatarURL := "https://c8.patreon.com/2/200/0"

                        // For link, use a Patreon profile URL if available (this is simplified)
                        link := fmt.Sprintf("https://www.patreon.com/user?u=%s", patron.ID)

                        createdAt := time.Now().Format(time.RFC3339)
                        if patron.Attributes.PledgeRelationshipStart != nil {
                                createdAt = patron.Attributes.PledgeRelationshipStart.Format(time.RFC3339)
                        }

                        sponsor := Sponsor{
                                ID:            patron.ID,
                                Name:          name,
                                Login:         login,
                                AvatarURL:     avatarURL,
                                Link:          link,
                                Platform:      "patreon",
                                MonthlyAmount: monthlyAmount,
                                CreatedAt:     createdAt,
                                TierName:      tierName,
                        }

                        sponsors = append(sponsors, sponsor)
                }

                // Check if there are more pages
                hasNextPage = response.Links.Next != ""
                if hasNextPage {
                        url = response.Links.Next
                }
        }

        return sponsors, nil
}
