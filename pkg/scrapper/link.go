package scrapper

import (
	"fmt"
)

func (c *Client) Click(selector, desc string) error {
	selection, err := c.Find(selector)
	if err != nil {
		return fmt.Errorf("searching link for %s got: %w", desc, err)
	}
	if selection == nil {
		return fmt.Errorf("for %s link %s not found", desc, selector)
	}
	if len(selection.Nodes) > 1 {
		return fmt.Errorf("for %s too many links found got %d want 1", desc, len(selection.Nodes))
	}

	href := getAttribute(selection.Nodes[0], "href")
	if href == "" {
		return fmt.Errorf("for %s to link found but href was empty", desc)
	}

	return c.Get(href)
}
