package models

type Profile struct {
    Project     string            `yaml:"project"`
    Profile     string            `yaml:"profile"`
    Dir         string            `yaml:"directory"`
    Branch      string            `yaml:"branch,omitempty"`
    Docker      Docker            `yaml:"docker"`
    EnvFile     string            `yaml:"env_file,omitempty"`
    EnvVars     map[string]string `yaml:"env_vars,omitempty"`
    PreScripts  []string          `yaml:"pre_scripts,omitempty"`
    Scripts     []string          `yaml:"scripts"`
    PostScripts []string          `yaml:"post_scripts,omitempty"`
}
