package config

import (
        "fmt"
        "os"
        "strings"

        "gopkg.in/yaml.v2"
)

// Config represents the application configuration
type Config struct {
        // Output settings
        OutputDir      string `yaml:"outputDir"`
        CacheDir       string `yaml:"cacheDir"`
        SVGTemplate    string `yaml:"svgTemplate"`
        DefaultAvatar  string `yaml:"defaultAvatar"`
        RefreshMinutes int    `yaml:"refreshMinutes"`

        // GitHub sponsor settings
        GitHubToken      string   `yaml:"githubToken"`
        GitHubLogin      string   `yaml:"githubLogin"`
        IncludePrivate   bool     `yaml:"includePrivate"`
        GitHubOrgs       []string `yaml:"githubOrgs"`
        ExcludeSponsors  []string `yaml:"excludeSponsors"`
        IncludeSponsors  []string `yaml:"includeSponsors"`
        ForceSponsorAmounts map[string]float64 `yaml:"forceSponsorAmounts"`

        // OpenCollective settings
        OpenCollectiveSlug string `yaml:"openCollectiveSlug"`
        OpenCollectiveKey  string `yaml:"openCollectiveKey"`

        // Patreon settings
        PatreonToken      string `yaml:"patreonToken"`
        PatreonCampaignID string `yaml:"patreonCampaignId"`

        // Afdian settings
        AfdianUserID      string `yaml:"afdianUserId"`
        AfdianToken       string `yaml:"afdianToken"`

        // Rendering settings
        AvatarSize       int    `yaml:"avatarSize"`
        AvatarMargin     int    `yaml:"avatarMargin"`
        SVGWidth         int    `yaml:"svgWidth"`
        FontSize         int    `yaml:"fontSize"`
        FontFamily       string `yaml:"fontFamily"`
        ShowAmount       bool   `yaml:"showAmount"`
        ShowName         bool   `yaml:"showName"`
        BackgroundColor  string `yaml:"backgroundColor"`
        PaddingX         int    `yaml:"paddingX"`
        PaddingY         int    `yaml:"paddingY"`
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

// LoadConfig loads the configuration from a file
func LoadConfig(filename string) (Config, error) {
        config := DefaultConfig()

        data, err := os.ReadFile(filename)
        if err != nil {
                return config, fmt.Errorf("reading config file: %w", err)
        }

        if err := yaml.Unmarshal(data, &config); err != nil {
                return config, fmt.Errorf("parsing config file: %w", err)
        }

        // Check for environment variables to override config
        if env := os.Getenv("GITHUB_TOKEN"); env != "" {
                config.GitHubToken = env
        }

        if env := os.Getenv("OPENCOLLECTIVE_KEY"); env != "" {
                config.OpenCollectiveKey = env
        }

        if env := os.Getenv("PATREON_TOKEN"); env != "" {
                config.PatreonToken = env
        }
        
        if env := os.Getenv("AFDIAN_USER_ID"); env != "" {
                config.AfdianUserID = env
        }
        
        if env := os.Getenv("AFDIAN_TOKEN"); env != "" {
                config.AfdianToken = env
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
