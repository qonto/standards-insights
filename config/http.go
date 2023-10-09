package config

type HTTPConfig struct {
	ExposeConfiguration bool   `yaml:"expose-configuration"`
	Host                string `validate:"required"`
	Port                uint32 `validate:"required,gt=1024,lt=65535"`
	WriteTimeout        int    `yaml:"write-timeout" validate:"gt=-1,lt=60"`
	ReadTimeout         int    `yaml:"read-timeout" validate:"gt=-1,lt=60"`
	ReadHeaderTimeout   int    `yaml:"read-header-timeout" validate:"gt=-1,lt=60"`
	CertPath            string `yaml:"cert-path" validate:"omitempty,file"`
	KeyPath             string `yaml:"key-path" validate:"omitempty,file"`
	CacertPath          string `yaml:"ca-cert-path" validate:"omitempty,file"`
	InsecureSkipVerify  bool   `yaml:"insecure-skip-verify"`
	ClientAuthType      string `yaml:"client-auth-type" validate:"omitempty,oneof=NoClientCert RequestClientCert RequireAnyClientCert VerifyClientCertIfGiven RequireAndVerifyClientCert"`
}
