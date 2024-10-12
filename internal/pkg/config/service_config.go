package config

type ServiceConfig struct {
	ServiceHost string `env:"GATEWAY_SERVICE_HOST"`
	ServicePort string `env:"GATEWAY_SERVICE_PORT"`
	DBConfig    DBConfig
}

type DBConfig struct {
	Host     string `env:"GATEWAY_SERVICE_DB_HOST"`
	Port     string `env:"GATEWAY_SERVICE_DB_PORT"`
	Database string `env:"GATEWAY_SERVICE_DB_NAME"`
	Username string `env:"GATEWAY_SERVICE_DB_USER"`
	Password string `env:"GATEWAY_SERVICE_DB_PASSWORD"`
}
