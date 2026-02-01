package selector

import (
	"fmt"

	"github.com/ktr0731/go-fuzzyfinder"

	"github.com/enbiyagoral/sopsctl/internal/config"
)

// FZFSelector implements Selector using go-fuzzyfinder.
type FZFSelector struct{}

// NewFZFSelector creates a new FZF-based selector.
func NewFZFSelector() *FZFSelector {
	return &FZFSelector{}
}

// Select prompts the user to select a profile using fuzzy finder.
func (s *FZFSelector) Select(profiles []*config.Profile) (*config.Profile, error) {
	if len(profiles) == 0 {
		return nil, fmt.Errorf("no profiles available")
	}

	idx, err := fuzzyfinder.Find(
		profiles,
		func(i int) string {
			p := profiles[i]
			return fmt.Sprintf("%s  %s  [%s]",
				p.Name,
				p.Description,
				p.GetBackendSummary(),
			)
		},
		fuzzyfinder.WithPromptString("Select profile > "),
		fuzzyfinder.WithHeader("Use ↑/↓ to navigate, Enter to select, Esc to cancel"),
	)

	if err != nil {
		if err == fuzzyfinder.ErrAbort {
			return nil, fmt.Errorf("selection cancelled")
		}
		return nil, fmt.Errorf("selection failed: %w", err)
	}

	return profiles[idx], nil
}
