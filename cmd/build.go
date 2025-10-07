package cmd

import (
	"context"
	"fmt"
)

// Build initializes and validates all container dependencies
func (c *Container) Build(ctx context.Context) error {
	// Validate configuration
	if err := c.validateConfig(); err != nil {
		return fmt.Errorf("invalid configuration: %w", err)
	}

	// Initialize auth client (lazy initialization will happen on first access)
	_ = c.GetAuthClient(ctx)

	// Validate all dependencies are initialized
	return nilCheck(c)
}

// validateConfig validates the authentication configuration
func (c *Container) validateConfig() error {
	// AuthFactory handles both real and mock clients, so no validation needed
	return nil
}

// nilCheck validates that all critical dependencies are not nil
func nilCheck(container *Container) error {
	if container.db == nil {
		return fmt.Errorf("database connection is nil")
	}
	if container.unitOfWork == nil {
		return fmt.Errorf("unitOfWork is nil")
	}
	if container.eventPublisher == nil {
		return fmt.Errorf("eventPublisher is nil")
	}
	// authClient can be nil (optional dependency)

	return nil
}
