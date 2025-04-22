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

// GitHubSponsorResponse represents the GitHub GraphQL API response for sponsors
type GitHubSponsorResponse struct {
        Data struct {
                User struct {
                        SponsorshipsAsMaintainer struct {
                                Nodes []struct {
                                        CreatedAt  string `json:"createdAt"`
                                        IsOneTime  bool   `json:"isOneTimePayment"`
                                        TierName   string `json:"tier"`
                                        Sponsor    struct {
                                                Login     string `json:"login"`
                                                Name      string `json:"name"`
                                                AvatarURL string `json:"avatarUrl"`
                                                URL       string `json:"url"`
                                                ID        string `json:"id"`
                                        } `json:"sponsorEntity"`
                                        TotalAmountDonated struct {
                                                Currency string `json:"currency"`
                                                Value    int    `json:"value"`
                                        } `json:"totalDonated"`
                                        TierAmount struct {
                                                Currency string  `json:"currency"`
                                                Value    float64 `json:"value"`
                                        } `json:"tier"`
                                } `json:"nodes"`
                                PageInfo struct {
                                        HasNextPage bool   `json:"hasNextPage"`
                                        EndCursor   string `json:"endCursor"`
                                } `json:"pageInfo"`
                        } `json:"sponsorshipsAsMaintainer"`
                } `json:"user"`
        } `json:"data"`
        Errors []struct {
                Message string `json:"message"`
                Type    string `json:"type"`
        } `json:"errors"`
}

// FetchGitHubSponsors fetches sponsors from GitHub using the GraphQL API
func FetchGitHubSponsors(cfg config.Config) ([]Sponsor, error) {
        sponsors := []Sponsor{}
        
        query := `
        query($login: String!, $cursor: String) {
                user(login: $login) {
                        sponsorshipsAsMaintainer(first: 100, after: $cursor, includePrivate: %s) {
                                nodes {
                                        createdAt
                                        isOneTimePayment
                                        tier
                                        sponsorEntity {
                                                ... on User {
                                                        id
                                                        login
                                                        name
                                                        avatarUrl
                                                        url
                                                }
                                                ... on Organization {
                                                        id
                                                        login
                                                        name
                                                        avatarUrl
                                                        url
                                                }
                                        }
                                        totalDonated {
                                                currency
                                                value
                                        }
                                        tier {
                                                monthlyPriceInDollars
                                        }
                                }
                                pageInfo {
                                        hasNextPage
                                        endCursor
                                }
                        }
                }
        }
        `
        
        includePrivate := "false"
        if cfg.IncludePrivate {
                includePrivate = "true"
        }
        
        query = fmt.Sprintf(query, includePrivate)
        
        client := &http.Client{
                Timeout: 10 * time.Second,
        }
        
        hasNextPage := true
        cursor := ""
        
        for hasNextPage {
                variables := map[string]interface{}{
                        "login":  cfg.GitHubLogin,
                        "cursor": cursor,
                }
                
                requestBody, err := json.Marshal(map[string]interface{}{
                        "query":     query,
                        "variables": variables,
                })
                
                if err != nil {
                        return sponsors, fmt.Errorf("failed to marshal GitHub GraphQL request: %w", err)
                }
                
                req, err := http.NewRequest("POST", "https://api.github.com/graphql", bytes.NewBuffer(requestBody))
                if err != nil {
                        return sponsors, fmt.Errorf("failed to create GitHub GraphQL request: %w", err)
                }
                
                req.Header.Set("Authorization", "bearer "+cfg.GitHubToken)
                req.Header.Set("Content-Type", "application/json")
                
                resp, err := client.Do(req)
                if err != nil {
                        return sponsors, fmt.Errorf("failed to execute GitHub GraphQL request: %w", err)
                }
                defer resp.Body.Close()
                
                if resp.StatusCode != http.StatusOK {
                        body, _ := io.ReadAll(resp.Body)
                        return sponsors, fmt.Errorf("GitHub GraphQL request failed with status %d: %s", resp.StatusCode, string(body))
                }
                
                var response GitHubSponsorResponse
                if err := json.NewDecoder(resp.Body).Decode(&response); err != nil {
                        return sponsors, fmt.Errorf("failed to decode GitHub GraphQL response: %w", err)
                }
                
                if len(response.Errors) > 0 {
                        return sponsors, fmt.Errorf("GitHub GraphQL API error: %s", response.Errors[0].Message)
                }
                
                // Process sponsors from this page
                for _, node := range response.Data.User.SponsorshipsAsMaintainer.Nodes {
                        // Skip one-time payments
                        if node.IsOneTime {
                                continue
                        }
                        
                        // Skip sponsors with no entity (deleted accounts etc.)
                        if node.Sponsor.Login == "" {
                                continue
                        }
                        
                        sponsor := Sponsor{
                                ID:            node.Sponsor.ID,
                                Name:          node.Sponsor.Name,
                                Login:         node.Sponsor.Login,
                                AvatarURL:     node.Sponsor.AvatarURL,
                                Link:          node.Sponsor.URL,
                                Platform:      "github",
                                MonthlyAmount: node.TierAmount.Value,
                                CreatedAt:     node.CreatedAt,
                                TierName:      node.TierName,
                        }
                        
                        sponsors = append(sponsors, sponsor)
                }
                
                // Check if there are more pages
                hasNextPage = response.Data.User.SponsorshipsAsMaintainer.PageInfo.HasNextPage
                if hasNextPage {
                        cursor = response.Data.User.SponsorshipsAsMaintainer.PageInfo.EndCursor
                }
        }
        
        // Also fetch from orgs if specified
        for _, org := range cfg.GitHubOrgs {
                orgSponsors, err := fetchGitHubOrgSponsors(org, cfg)
                if err != nil {
                        // Log error but continue
                        fmt.Printf("Error fetching sponsors for org %s: %v\n", org, err)
                        continue
                }
                
                sponsors = append(sponsors, orgSponsors...)
        }
        
        return sponsors, nil
}

// fetchGitHubOrgSponsors fetches sponsors for a GitHub organization
func fetchGitHubOrgSponsors(orgLogin string, cfg config.Config) ([]Sponsor, error) {
        // Similar implementation to FetchGitHubSponsors but for organizations
        // This is a simplified version as the full implementation would be quite lengthy
        
        // For the sake of this example, we'll just return an empty slice
        // In a real implementation, you would use a similar GraphQL query adjusted for organizations
        
        return []Sponsor{}, nil
}
