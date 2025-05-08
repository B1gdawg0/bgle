package models

type Docker struct {
    Enabled      bool     `yaml:"enabled"`
    ComposeFiles []string `yaml:"compose_files,omitempty"`
    Up           bool     `yaml:"up"`
    // Build        bool     `yaml:"build"`
}
