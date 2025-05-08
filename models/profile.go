package models

type Profile struct {
    Project     string            `yaml:"projec,omitempty"`
    Profile     string            `yaml:"profile,omitempty"`
    Dir         string            `yaml:"directory,omitempty"`
    Branch      string            `yaml:"branch,omitempty"`
    Docker      Docker            `yaml:"docker"`
    EnvFile     string            `yaml:"env_file,omitempty"`
    EnvVars     map[string]string `yaml:"env_vars,omitempty"`
    Bootstrap   Bootstrap         `yaml:"bootstrap"`
    PreScripts  []string          `yaml:"pre_scripts,omitempty"`
    Scripts     []string          `yaml:"scripts"`
    PostScripts []string          `yaml:"post_scripts,omitempty"`
}

type OutputProfile struct {
    Branch      string            `yaml:"branch,omitempty"`
    Docker      Docker            `yaml:"docker"`
    EnvFile     string            `yaml:"env_file,omitempty"`
    EnvVars     map[string]string `yaml:"env_vars,omitempty"`
    Bootstrap   Bootstrap         `yaml:"bootstrap"`
    PreScripts  []string          `yaml:"pre_scripts,omitempty"`
    Scripts     []string          `yaml:"scripts"`
    PostScripts []string          `yaml:"post_scripts,omitempty"`
}