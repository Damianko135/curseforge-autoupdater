package api

import (
	"fmt"
)

// Legacy function to maintain backward compatibility
func (c *Client) CheckIfExists(id int) (bool, error) {
	return c.CheckIfModExists(id)
}

// Legacy ModInfo struct for backward compatibility
type LegacyModInfo struct {
	ID   int    `json:"id"`
	Name string `json:"name"`
}

// GetLegacyModInfo retrieves basic mod information for backward compatibility
func (c *Client) GetLegacyModInfo(id int) (*LegacyModInfo, error) {
	modInfo, err := c.GetMod(id)
	if err != nil {
		return nil, fmt.Errorf("failed to get mod info: %w", err)
	}

	return &LegacyModInfo{
		ID:   modInfo.ID,
		Name: modInfo.Name,
	}, nil
}
