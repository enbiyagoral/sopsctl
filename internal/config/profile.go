package config

// Profile represents a SOPS encryption profile.
type Profile struct {
	Name        string `yaml:"-"` // Populated from map key
	Description string `yaml:"description,omitempty"`

	// Encryption backend
	Age *AgeConfig `yaml:"age,omitempty"`

	// SOPS-specific options
	SOPS SOPSOptions `yaml:"sops,omitempty"`
}

// AgeConfig represents age encryption configuration.
type AgeConfig struct {
	Recipients []string `yaml:"recipients"`
}

// SOPSOptions represents SOPS-specific encryption options.
type SOPSOptions struct {
	EncryptedRegex    string `yaml:"encrypted_regex,omitempty"`
	EncryptedSuffix   string `yaml:"encrypted_suffix,omitempty"`
	UnencryptedRegex  string `yaml:"unencrypted_regex,omitempty"`
	UnencryptedSuffix string `yaml:"unencrypted_suffix,omitempty"`
}

// GetBackendSummary returns a human-readable summary of configured backends.
func (p *Profile) GetBackendSummary() string {
	if p.Age != nil && len(p.Age.Recipients) > 0 {
		return "age"
	}
	return "none"
}

// HasBackends returns true if the profile has at least one backend configured.
func (p *Profile) HasBackends() bool {
	return p.Age != nil && len(p.Age.Recipients) > 0
}
