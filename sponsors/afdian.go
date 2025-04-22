package sponsors

import (
        "bytes"
        "crypto/md5"
        "encoding/hex"
        "encoding/json"
        "fmt"
        "io"
        "log"
        "net/http"
        "strconv"
        "time"

        "sponsorgen/config"
)

// AfdianSponsor represents data returned from Afdian API
type AfdianSponsor struct {
        SponsorPlans []AfdianPlan   `json:"sponsor_plans"`
        CurrentPlan  AfdianPlan     `json:"current_plan"`
        AllSumAmount string         `json:"all_sum_amount"`
        CreateTime   int64          `json:"create_time"`
        LastPayTime  int64          `json:"last_pay_time"`
        User         AfdianUserInfo `json:"user"`
}

// AfdianPlan represents a plan in Afdian API
type AfdianPlan struct {
        PlanID string `json:"plan_id"`
        Name   string `json:"name"`
        Price  string `json:"price,omitempty"`
}

// AfdianUserInfo represents user information in Afdian API
type AfdianUserInfo struct {
        UserID string `json:"user_id"`
        Name   string `json:"name"`
        Avatar string `json:"avatar"`
}

// AfdianResponse represents the response from Afdian API
type AfdianResponse struct {
        EC   int    `json:"ec"`
        EM   string `json:"em"`
        Data struct {
                TotalCount int             `json:"total_count"`
                TotalPage  int             `json:"total_page"`
                List       []AfdianSponsor `json:"list"`
        } `json:"data"`
}

// FetchAfdianSponsors retrieves sponsors from Afdian
func FetchAfdianSponsors(cfg config.Config) ([]Sponsor, error) {
        if cfg.AfdianUserID == "" || cfg.AfdianToken == "" {
                return nil, fmt.Errorf("Afdian user ID or token not provided")
        }

        log.Println("Fetching Afdian sponsors...")

        var allSponsors []Sponsor
        page := 1
        totalPages := 1

        // Loop through pages
        for page <= totalPages {
                // Create parameters for this page
                params := map[string]interface{}{
                        "page":    page,
                        "per_page": 50, // Maximum allowed
                }

                // Convert params to JSON string
                paramsJSON, err := json.Marshal(params)
                if err != nil {
                        return nil, fmt.Errorf("encoding params: %w", err)
                }

                // Generate timestamp
                ts := strconv.FormatInt(time.Now().Unix(), 10)

                // Generate signature
                signStr := fmt.Sprintf("%sparams%sts%suser_id%s", 
                        cfg.AfdianToken, 
                        string(paramsJSON), 
                        ts, 
                        cfg.AfdianUserID)
                
                hash := md5.Sum([]byte(signStr))
                sign := hex.EncodeToString(hash[:])

                // Create request body
                reqBody := map[string]string{
                        "user_id": cfg.AfdianUserID,
                        "params":  string(paramsJSON),
                        "ts":      ts,
                        "sign":    sign,
                }

                reqJSON, err := json.Marshal(reqBody)
                if err != nil {
                        return nil, fmt.Errorf("encoding request: %w", err)
                }

                // Create HTTP request
                req, err := http.NewRequest("POST", "https://afdian.com/api/open/query-sponsor", bytes.NewBuffer(reqJSON))
                if err != nil {
                        return nil, fmt.Errorf("creating request: %w", err)
                }

                req.Header.Set("Content-Type", "application/json")

                // Execute request
                client := &http.Client{Timeout: 30 * time.Second}
                resp, err := client.Do(req)
                if err != nil {
                        return nil, fmt.Errorf("sending request: %w", err)
                }
                defer resp.Body.Close()

                // Read response
                body, err := io.ReadAll(resp.Body)
                if err != nil {
                        return nil, fmt.Errorf("reading response: %w", err)
                }

                // Parse response
                var afdianResp AfdianResponse
                if err := json.Unmarshal(body, &afdianResp); err != nil {
                        return nil, fmt.Errorf("parsing response: %w", err)
                }

                // Check response status
                if afdianResp.EC != 200 {
                        return nil, fmt.Errorf("API error: %s", afdianResp.EM)
                }

                // Process sponsors
                for _, afdianSponsor := range afdianResp.Data.List {
                        // Calculate monthly amount - convert Afdian's total to monthly average
                        monthlyAmount := 0.0
                        if afdianSponsor.AllSumAmount != "" {
                                totalAmount, err := strconv.ParseFloat(afdianSponsor.AllSumAmount, 64)
                                if err == nil && totalAmount > 0 {
                                        // If we have last pay time and create time, we can calculate a monthly average
                                        if afdianSponsor.LastPayTime > 0 && afdianSponsor.CreateTime > 0 {
                                                // Calculate the number of months between first and last payment
                                                months := float64(afdianSponsor.LastPayTime-afdianSponsor.CreateTime) / (30 * 24 * 60 * 60)
                                                if months < 1 {
                                                        months = 1 // Minimum 1 month to avoid division by zero
                                                }
                                                monthlyAmount = totalAmount / months
                                        } else {
                                                // Fallback: assume it's a one-time donation
                                                monthlyAmount = totalAmount
                                        }
                                }
                        }
                        
                        // Determine tier name
                        tierName := ""
                        if afdianSponsor.CurrentPlan.Name != "" {
                                tierName = afdianSponsor.CurrentPlan.Name
                        }

                        // Create sponsor
                        sponsor := Sponsor{
                                ID:            afdianSponsor.User.UserID,
                                Name:          afdianSponsor.User.Name,
                                Login:         afdianSponsor.User.Name, // Use name as login since Afdian doesn't have a separate login field
                                AvatarURL:     afdianSponsor.User.Avatar,
                                Link:          fmt.Sprintf("https://afdian.com/@%s", afdianSponsor.User.UserID),
                                Platform:      "afdian",
                                MonthlyAmount: monthlyAmount,
                                CreatedAt:     time.Unix(afdianSponsor.CreateTime, 0).Format(time.RFC3339),
                                TierName:      tierName,
                        }

                        allSponsors = append(allSponsors, sponsor)
                }

                // Update total pages
                totalPages = afdianResp.Data.TotalPage
                
                // Move to next page
                page++
        }

        log.Printf("Found %d Afdian sponsors", len(allSponsors))
        return allSponsors, nil
}