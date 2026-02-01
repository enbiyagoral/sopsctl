// Package sops provides SOPS argument building and execution.
package sops

import (
	"github.com/enbiyagoral/sopsctl/internal/config"
)

// ArgsBuilder converts a profile configuration to SOPS CLI arguments.
type ArgsBuilder struct{}

// NewArgsBuilder creates a new ArgsBuilder.
func NewArgsBuilder() *ArgsBuilder {
	return &ArgsBuilder{}
}

// Build generates SOPS CLI arguments from a profile.
func (b *ArgsBuilder) Build(profile *config.Profile, command string, file string) []string {
	args := make([]string, 0, 16)

	// Age backend
	if profile.Age != nil {
		for _, recipient := range profile.Age.Recipients {
			args = append(args, "--age", recipient)
		}
	}

	// SOPS options
	if profile.SOPS.EncryptedRegex != "" {
		args = append(args, "--encrypted-regex", profile.SOPS.EncryptedRegex)
	}
	if profile.SOPS.EncryptedSuffix != "" {
		args = append(args, "--encrypted-suffix", profile.SOPS.EncryptedSuffix)
	}
	if profile.SOPS.UnencryptedRegex != "" {
		args = append(args, "--unencrypted-regex", profile.SOPS.UnencryptedRegex)
	}
	if profile.SOPS.UnencryptedSuffix != "" {
		args = append(args, "--unencrypted-suffix", profile.SOPS.UnencryptedSuffix)
	}

	// Command and file
	args = append(args, command, file)

	return args
}

// BuildDecrypt generates arguments for decrypt.
func (b *ArgsBuilder) BuildDecrypt(file string) []string {
	return []string{"decrypt", file}
}

// BuildEdit generates arguments for edit.
func (b *ArgsBuilder) BuildEdit(profile *config.Profile, file string) []string {
	if profile == nil {
		return []string{"edit", file}
	}
	return b.Build(profile, "edit", file)
}
