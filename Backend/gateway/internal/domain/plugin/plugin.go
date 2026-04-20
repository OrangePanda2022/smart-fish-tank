package plugin

type PluginType string

const (
	PluginJWT       PluginType = "JWT"
	PluginRateLimit PluginType = "REDISLIMITER"
	PluginIPFilter  PluginType = "IP_FILTER"
)

type Plugin struct {
	Name   string
	Type   PluginType
	Config map[string]interface{}
}
