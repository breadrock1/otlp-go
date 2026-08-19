package otlp_go

type OtlpConfig struct {
	AppName string       `mapstructure:"app_name"`
	Logger  LoggerConfig `mapstructure:"logger"`
	Meter   MeterConfig  `mapstructure:"meter"`
	Tracer  TracerConfig `mapstructure:"tracer"`
}

type LoggerConfig struct {
	Level      string             `mapstructure:"level"`
	Attributes []string           `mapstructure:"attributes"`
	Loki       LokiLoggerConfig   `mapstructure:"loki"`
	Syslog     SyslogLoggerConfig `mapstructure:"syslog"`
}

type SyslogLoggerConfig struct {
	Address string `mapstructure:"address"`
	Enable  bool   `mapstructure:"enable"`
}

type LokiLoggerConfig struct {
	Address string `mapstructure:"address"`
	Enable  bool   `mapstructure:"enable"`
}

type MeterConfig struct {
	EndpointURI   string   `mapstructure:"endpoint_uri"`
	ExcludedPaths []string `mapstructure:"excluded_paths"`
}

type TracerConfig struct {
	Address      string `mapstructure:"address"`
	EnableJaeger bool   `mapstructure:"enable_jaeger"`
}
