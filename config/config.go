package config

var defaultModuleConfigs = make(map[string]ModuleConfig)

func AddDefaultConfig(name string, config ModuleConfig) {
	defaultModuleConfigs[name] = config
}

type Config struct {
	AuthorisationToken string                  `json:"authorisation_token"`
	ApplicationID      string                  `json:"application_id"`
	ModulePriorities   map[string]int          `json:"module_priorities"`
	ModuleConfigs      map[string]ModuleConfig `json:"module_configs"`
}

type ModuleConfig interface{}

func GetDefaultConfig() *Config {
	return &Config{
		AuthorisationToken: "",
		ApplicationID:      "",
		ModulePriorities: map[string]int{
			"Spotify": 0,
			"Last FM": 1,
		},
		ModuleConfigs: defaultModuleConfigs,
	}
}
