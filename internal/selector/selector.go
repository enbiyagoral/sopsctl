// Package selector provides profile selection interfaces and implementations.
package selector

import (
	"github.com/enbiyagoral/sopsctl/internal/config"
)

// Selector is the interface for profile selection.
type Selector interface {
	// Select prompts the user to select a profile from the given list.
	Select(profiles []*config.Profile) (*config.Profile, error)
}
