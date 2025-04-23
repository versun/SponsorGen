package config

import (
        "fmt"
        "os"
        "strconv"
        "strings"
)

// Config represents the application configuration
type Config struct {
        // Output settings
        OutputDir      string
        CacheDir       string
        SVGTemplate    string
        DefaultAvatar  string
        RefreshMinutes int

        // GitHub sponsor settings
        GitHubToken          string
        GitHubLogin          string
        IncludePrivate       bool
        GitHubOrgs           []string
        ExcludeSponsors      []string
        IncludeSponsors      []string
        ForceSponsorAmounts  map[string]float64

        // OpenCollective settings
        OpenCollectiveSlug   string
        OpenCollectiveKey    string

        // Patreon settings
        PatreonToken         string
        PatreonCampaignID    string

        // Afdian settings
        AfdianUserID         string
        AfdianToken          string

        // Rendering settings
        AvatarSize           int
        AvatarMargin         int
        SVGWidth             int
        FontSize             int
        FontFamily           string
        ShowAmount           bool
        ShowName             bool
        BackgroundColor      string
        PaddingX             int
        PaddingY             int
}

// DefaultConfig returns a default configuration
func DefaultConfig() Config {
        return Config{
                OutputDir:      "./output",
                CacheDir:       "./cache",
                RefreshMinutes: 60,
                DefaultAvatar:  "./assets/default_avatar.svg",
                AvatarSize:     45,
                AvatarMargin:   5,
                SVGWidth:       800,
                FontSize:       14,
                FontFamily:     "system-ui, -apple-system, 'Segoe UI', Roboto, Ubuntu, Cantarell, 'Noto Sans', sans-serif",
                ShowAmount:     false,
                ShowName:       false,
                BackgroundColor: "transparent",
                PaddingX:       10,
                PaddingY:       10,
                GitHubToken:      "",
                GitHubLogin:      "",
                IncludePrivate:   false,
                GitHubOrgs:       []string{},
                ExcludeSponsors:  []string{},
                IncludeSponsors:  []string{},
                ForceSponsorAmounts: map[string]float64{},
                OpenCollectiveSlug: "",
                OpenCollectiveKey:  "",
                PatreonToken:      "",
                PatreonCampaignID: "",
                AfdianUserID:      "",
                AfdianToken:       "",
        }
}

// LoadConfig loads the configuration from environment variables
func LoadConfig() (Config, error) {
        config := DefaultConfig()

        // Output settings
        if env := os.Getenv("OUTPUT_DIR"); env != "" {
                config.OutputDir = env
        }
        
        if env := os.Getenv("CACHE_DIR"); env != "" {
                config.CacheDir = env
        }
        
        if env := os.Getenv("DEFAULT_AVATAR"); env != "" {
                config.DefaultAvatar = env
        }
        
        if env := os.Getenv("REFRESH_MINUTES"); env != "" {
                if val, err := strconv.Atoi(env); err == nil {
                        config.RefreshMinutes = val
                }
        }

        // GitHub settings
        if env := os.Getenv("GITHUB_TOKEN"); env != "" {
                config.GitHubToken = env
        }
        
        if env := os.Getenv("GITHUB_LOGIN"); env != "" {
                config.GitHubLogin = env
        }
        
        if env := os.Getenv("INCLUDE_PRIVATE"); env != "" {
                config.IncludePrivate = (strings.ToLower(env) == "true")
        }
        
        if env := os.Getenv("GITHUB_ORGS"); env != "" {
                config.GitHubOrgs = strings.Split(env, ",")
        }
        
        if env := os.Getenv("EXCLUDE_SPONSORS"); env != "" {
                config.ExcludeSponsors = strings.Split(env, ",")
        }
        
        if env := os.Getenv("INCLUDE_SPONSORS"); env != "" {
                config.IncludeSponsors = strings.Split(env, ",")
        }

        // OpenCollective settings
        if env := os.Getenv("OPENCOLLECTIVE_SLUG"); env != "" {
                config.OpenCollectiveSlug = env
        }
        
        if env := os.Getenv("OPENCOLLECTIVE_KEY"); env != "" {
                config.OpenCollectiveKey = env
        }

        // Patreon settings
        if env := os.Getenv("PATREON_TOKEN"); env != "" {
                config.PatreonToken = env
        }
        
        if env := os.Getenv("PATREON_CAMPAIGN_ID"); env != "" {
                config.PatreonCampaignID = env
        }

        // Afdian settings
        if env := os.Getenv("AFDIAN_USER_ID"); env != "" {
                config.AfdianUserID = env
        }
        
        if env := os.Getenv("AFDIAN_TOKEN"); env != "" {
                config.AfdianToken = env
        }

        // Rendering settings
        if env := os.Getenv("AVATAR_SIZE"); env != "" {
                if val, err := strconv.Atoi(env); err == nil {
                        config.AvatarSize = val
                }
        }
        
        if env := os.Getenv("AVATAR_MARGIN"); env != "" {
                if val, err := strconv.Atoi(env); err == nil {
                        config.AvatarMargin = val
                }
        }
        
        if env := os.Getenv("SVG_WIDTH"); env != "" {
                if val, err := strconv.Atoi(env); err == nil {
                        config.SVGWidth = val
                }
        }
        
        if env := os.Getenv("FONT_SIZE"); env != "" {
                if val, err := strconv.Atoi(env); err == nil {
                        config.FontSize = val
                }
        }
        
        if env := os.Getenv("FONT_FAMILY"); env != "" {
                config.FontFamily = env
        }
        
        if env := os.Getenv("SHOW_AMOUNT"); env != "" {
                config.ShowAmount = (strings.ToLower(env) == "true")
        }
        
        if env := os.Getenv("SHOW_NAME"); env != "" {
                config.ShowName = (strings.ToLower(env) == "true")
        }
        
        if env := os.Getenv("BACKGROUND_COLOR"); env != "" {
                config.BackgroundColor = env
        }
        
        if env := os.Getenv("PADDING_X"); env != "" {
                if val, err := strconv.Atoi(env); err == nil {
                        config.PaddingX = val
                }
        }
        
        if env := os.Getenv("PADDING_Y"); env != "" {
                if val, err := strconv.Atoi(env); err == nil {
                        config.PaddingY = val
                }
        }

        // Create SVG template if not provided
        if config.SVGTemplate == "" {
                config.SVGTemplate = DefaultSVGTemplate()
        }

        return config, nil
}

// DefaultSVGTemplate returns a default SVG template
func DefaultSVGTemplate() string {
        return `<svg xmlns="http://www.w3.org/2000/svg" xmlns:xlink="http://www.w3.org/1999/xlink" width="{{.Width}}" height="{{.Height}}" viewBox="0 0 {{.Width}} {{.Height}}">
  <style>
    .avatar {
      width: {{.AvatarSize}}px;
      height: {{.AvatarSize}}px;
      border-radius: 50%;
      overflow: hidden;
    }
  </style>
  <rect width="100%" height="100%" fill="{{.BackgroundColor}}" />
  <g transform="translate({{.PaddingX}}, {{.PaddingY}})">
    {{range .Sponsors}}
    <g transform="translate({{.X}}, {{.Y}})">
      <title>{{.Name}}</title>
      <image xlink:href="{{.Avatar}}" class="avatar" width="{{.Size}}" height="{{.Size}}" x="0" y="0" />
    </g>
    {{end}}
  </g>
</svg>`
}



// ValidateConfig ensures the configuration is valid
func (c *Config) ValidateConfig() error {
        var errors []string

        // Check if we have any source of sponsors
        if c.GitHubToken == "" && c.OpenCollectiveSlug == "" && c.PatreonToken == "" && c.AfdianUserID == "" {
                errors = append(errors, "No sponsor source configured (GitHub, OpenCollective, Patreon, or Afdian)")
        }

        // Check GitHub configuration
        if c.GitHubToken != "" && c.GitHubLogin == "" {
                errors = append(errors, "GitHub token provided but GitHub login is missing")
        }

        // Check OpenCollective configuration
        if c.OpenCollectiveSlug != "" && c.OpenCollectiveKey == "" {
                errors = append(errors, "OpenCollective slug provided but API key is missing")
        }

        // Check Patreon configuration
        if c.PatreonToken != "" && c.PatreonCampaignID == "" {
                errors = append(errors, "Patreon token provided but campaign ID is missing")
        }
        
        // Check Afdian configuration
        if c.AfdianUserID != "" && c.AfdianToken == "" {
                errors = append(errors, "Afdian user ID provided but token is missing")
        }

        // Return combined errors if any
        if len(errors) > 0 {
                return fmt.Errorf("configuration validation failed:\n- %s", strings.Join(errors, "\n- "))
        }

        return nil
}
