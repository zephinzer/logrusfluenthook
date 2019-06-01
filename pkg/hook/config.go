package hook

// Config is for storing configurations for the hook to
// instantiate a fluent logger instance
type Config struct {
	// Host defines the hostname of the fluent service
	Host string
	// Port defines the port the fluent service is listening on
	Port uint64
	// Base Tag defines the tag which will appear in fluent
	BaseTag string
	// Levels defines the levels for which the hook will be fired
	Levels []string
	// FieldMap defines the property labels to use for the various log fields
	FieldMap map[string]string
}
