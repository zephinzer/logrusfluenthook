package hook

// Config is for storing configurations for the hook to
// instantiate a fluent logger instance
type Config struct {
	// Host defines the hostname of the fluent service, defaults to 0.0.0.0
	Host string
	// Port defines the port the fluent service is listening on, defaults to 24224
	Port uint16
	// Base Tag defines the tag which will appear in fluent, defaults to "app"
	BaseTag string
	// Levels defines the levels for which the hook will be fired, defaults to all
	Levels []string
	// FieldMap defines the property labels to use for the various log fields
	FieldMap map[string]string
	// TimestampFormat defines the format of the timestamp attached to the log entry
	TimestampFormat string
}

func (config *Config) assignDefaults() {
	if config.Host == "" {
		config.Host = DefaultFluentHost
	}

	if config.Port == 0 {
		config.Port = DefaultFluentPort
	}

	if config.BaseTag == "" {
		config.BaseTag = DefaultFluentTagPrefix
	}

	if config.TimestampFormat == "" {
		config.TimestampFormat = DefaultTimestampFormat
	}

	if len(config.Levels) == 0 {
		for level, _ := range LevelMap {
			config.Levels = append(config.Levels, level)
		}
	}
}