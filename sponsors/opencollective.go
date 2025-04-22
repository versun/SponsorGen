package sponsors

import (
        "bytes"
        "encoding/json"
        "fmt"
        "io"
        "net/http"
        "time"

        "sponsorgen/config"
)

// OpenCollectiveResponse represents the OpenCollective API response
type OpenCollectiveResponse struct {
        Data struct {
                Account struct {
                        Orders struct {
                                Nodes []struct {
                                        FromAccount struct {
                                                ID        string `json:"id"`
                                                Name      string `json:"name"`
                                                Slug      string `json:"slug"`
                                                ImageURL  string `json:"imageUrl"`
                                                Website   string `json:"website"`
                                                Company   string `json:"company"`
                                                IsActive  bool   `json:"isActive"`
                                                CreatedAt string `json:"createdAt"`
                                        } `json:"fromAccount"`
                                        Status     string    `json:"status"`
                                        Amount     struct {
                                                Value    float64 `json:"value"`
                                                Currency string  `json:"currency"`
                                        } `json:"amount"`
                                        Frequency   string `json:"frequency"`
                                        TotalAmount struct {
                                                Value    float64 `json:"value"`
                                                Currency string  `json:"currency"`
                                        } `json:"totalAmount"`
                                        CreatedAt string `json:"createdAt"`
                                        Tier      struct {
                                                Name string `json:"name"`
                                        } `json:"tier"`
                                } `json:"nodes"`
                        } `json:"orders"`
                } `json:"account"`
        } `json:"data"`
}

// FetchOpenCollectiveSponsors fetches sponsors from OpenCollective
func FetchOpenCollectiveSponsors(cfg config.Config) ([]Sponsor, error) {
        sponsors := []Sponsor{}

        query := `
        query($slug: String!) {
                account(slug: $slug) {
                        orders(status: ACTIVE, filter: INCOMING) {
                                nodes {
                                        fromAccount {
                                                id
                                                name
                                                slug
                                                imageUrl
                                                website
                                                company
                                                isActive
                                                createdAt
                                        }
                                        status
                                        amount {
                                                value
                                                currency
                                        }
                                        frequency
                                        totalAmount {
                                                value
                                                currency
                                        }
                                        createdAt
                                        tier {
                                                name
                                        }
                                }
                        }
                }
        }
        `

        client := &http.Client{
                Timeout: 10 * time.Second,
        }

        variables := map[string]interface{}{
                "slug": cfg.OpenCollectiveSlug,
        }

        requestBody, err := json.Marshal(map[string]interface{}{
                "query":     query,
                "variables": variables,
        })

        if err != nil {
                return sponsors, fmt.Errorf("failed to marshal OpenCollective GraphQL request: %w", err)
        }

        req, err := http.NewRequest("POST", "https://api.opencollective.com/graphql/v2", bytes.NewBuffer(requestBody))
        if err != nil {
                return sponsors, fmt.Errorf("failed to create OpenCollective GraphQL request: %w", err)
        }

        if cfg.OpenCollectiveKey != "" {
                req.Header.Set("Api-Key", cfg.OpenCollectiveKey)
        }
        req.Header.Set("Content-Type", "application/json")

        resp, err := client.Do(req)
        if err != nil {
                return sponsors, fmt.Errorf("failed to execute OpenCollective GraphQL request: %w", err)
        }
        defer resp.Body.Close()

        if resp.StatusCode != http.StatusOK {
                body, _ := io.ReadAll(resp.Body)
                return sponsors, fmt.Errorf("OpenCollective GraphQL request failed with status %d: %s", resp.StatusCode, string(body))
        }

        var response OpenCollectiveResponse
        if err := json.NewDecoder(resp.Body).Decode(&response); err != nil {
                return sponsors, fmt.Errorf("failed to decode OpenCollective GraphQL response: %w", err)
        }

        // Process sponsors
        for _, node := range response.Data.Account.Orders.Nodes {
                // Skip inactive accounts
                if !node.FromAccount.IsActive {
                        continue
                }

                // Calculate monthly amount based on frequency
                monthlyAmount := 0.0
                switch node.Frequency {
                case "MONTHLY":
                        monthlyAmount = node.Amount.Value
                case "YEARLY":
                        monthlyAmount = node.Amount.Value / 12
                case "ONE_TIME":
                        // For one-time donations, we'll spread them over a year
                        monthlyAmount = node.Amount.Value / 12
                default:
                        // Skip unknown frequencies
                        continue
                }

                // Create sponsor profile URL
                profileURL := fmt.Sprintf("https://opencollective.com/%s", node.FromAccount.Slug)
                if node.FromAccount.Website != "" {
                        profileURL = node.FromAccount.Website
                }

                // Use company name if available
                name := node.FromAccount.Name
                if node.FromAccount.Company != "" {
                        name = node.FromAccount.Company
                }

                sponsor := Sponsor{
                        ID:            node.FromAccount.ID,
                        Name:          name,
                        Login:         node.FromAccount.Slug,
                        AvatarURL:     node.FromAccount.ImageURL,
                        Link:          profileURL,
                        Platform:      "opencollective",
                        MonthlyAmount: monthlyAmount,
                        CreatedAt:     node.CreatedAt,
                        TierName:      node.Tier.Name,
                }

                sponsors = append(sponsors, sponsor)
        }

        return sponsors, nil
}
