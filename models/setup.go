package models

type Bootstrap struct{
	Enabled      bool             `yaml:"enabled"`
	Repo_URL     string           `yaml:"repository_url"`
	Scripts     []string          `yaml:"init_scripts"`
}