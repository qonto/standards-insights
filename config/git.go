package config

type PrivateKey struct {
	Path     string
	Password string
}

type BasicAuth struct {
	Username string
	Password string
}

type GitConfig struct {
	PrivateKey PrivateKey `yaml:"private-key" validate:"omitempty,file"`
	BasicAuth  BasicAuth  `yaml:"basic-auth"`
}
