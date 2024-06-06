package config

func GetOtelCollectorConfig() string {
	return getEnv("OTEL_COLLECTOR", "")
}
