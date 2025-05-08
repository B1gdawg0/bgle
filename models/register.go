package models

type Register struct {
	Profiles map[string]ProfileEntry `yaml:"profiles"`
}

type ProfileEntry struct {
	Profile string `yaml:"profile"`
	Dir     string `yaml:"dir"`
}
