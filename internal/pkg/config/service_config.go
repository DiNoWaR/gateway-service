package config

type ServiceConfig struct {
	ServiceHost             string `env:"GATEWAY_SERVICE_HOST"`
	ServicePort             string `env:"GATEWAY_SERVICE_PORT"`
	ServiceCallbackEndpoint string `env:"GATEWAY_SERVICE_CALLBACK_ENDPOINT"`
	SoapGatewayConfig       SoapGatewayConfig
	RestGatewayConfig       RestGatewayConfig
	DBConfig                DBConfig
	RetryInterval           int `env:"GATEWAY_SERVICE_INTERVAL"`
	RetryElapseTime         int `env:"GATEWAY_SERVICE_ELAPSE_TIME"`
}

type DBConfig struct {
	Host     string `env:"GATEWAY_SERVICE_DB_HOST"`
	Port     string `env:"GATEWAY_SERVICE_DB_PORT"`
	Database string `env:"GATEWAY_SERVICE_DB_NAME"`
	Username string `env:"GATEWAY_SERVICE_DB_USER"`
	Password string `env:"GATEWAY_SERVICE_DB_PASSWORD"`
}

type SoapGatewayConfig struct {
	GatewayId    string `env:"SOAP_GATEWAY_ID"`
	Endpoint     string `env:"SOAP_GATEWAY_ENDPOINT"`
	EndpointHost string `env:"SOAP_GATEWAY_ENDPOINT_HOST"`
	EndpointPort string `env:"SOAP_GATEWAY_ENDPOINT_PORT"`
}

type RestGatewayConfig struct {
	GatewayId string `env:"REST_GATEWAY_ID"`
	Host      string `env:"REST_GATEWAY_HOST"`
	Port      string `env:"REST_GATEWAY_PORT"`
}
