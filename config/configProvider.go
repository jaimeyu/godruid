package config

// Provider - taken from https://github.com/gohugoio/hugo/blob/master/config/configProvider.go
type Provider interface {
	GetString(key string) string
	GetStringSlice(key string) []string
	GetInt(key string) int
	GetFloat64(key string) float64
	GetBool(key string) bool
	GetStringMap(key string) map[string]interface{}
	GetStringMapString(key string) map[string]string
	Get(key string) interface{}
	Set(key string, value interface{})
	IsSet(key string) bool
}
