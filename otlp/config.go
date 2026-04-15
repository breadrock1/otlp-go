package otlp_go

type OtlpConfig struct {
	AppName string       `mapstructure:"app_name"`
	Logger  LoggerConfig `mapstructure:"logger"`
	Tracer  TracerConfig `mapstructure:"tracer"`
}

type LoggerConfig struct {
	Level      string `mapstructure:"level"`
	Address    string `mapstructure:"address"`
	EnableLoki bool   `mapstructure:"enable_loki"`
}

type TracerConfig struct {
	Address      string `mapstructure:"address"`
	EnableJaeger bool   `mapstructure:"enable_jaeger"`
}
