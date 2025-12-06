package scrapper

import (
	"fmt"
	"regexp"
	"strings"
)

var (
	singleSpacePattern = regexp.MustCompile(`\s+`)
)

func trim(in string) string {
	return strings.TrimSpace(singleSpacePattern.ReplaceAllString(in, " "))
}

func (c *Client) TestNodeDontExist(previousError error, selector, desc string) error {
	if previousError != nil {
		return previousError
	}
	selection, err := c.Find(selector)
	if err != nil {
		return fmt.Errorf("searching for %s got: %w", desc, err)
	}
	if selection != nil {
		return fmt.Errorf("%s was found with value %s expected to be not found", desc, selection.Text())
	}
	return nil
}

func (c *Client) TestValue(previousError error, selector, desc, expected string) error {
	if previousError != nil {
		return previousError
	}
	selection, err := c.Find(selector)
	if err != nil {
		return fmt.Errorf("searching for %s got: %w", desc, err)
	}
	if selection == nil {
		return fmt.Errorf("%s not found", desc)
	}
	result := trim(selection.Text())
	if !strings.EqualFold(result, expected) {
		return fmt.Errorf("%s found with value '%s' but expected '%s'", desc, result, expected)
	}
	return nil
}

func (c *Client) TestValueRegex(previousError error, selector, desc, expected string) error {
	if previousError != nil {
		return previousError
	}
	selection, err := c.Find(selector)
	if err != nil {
		return fmt.Errorf("searching for %s got: %w", desc, err)
	}
	if selection == nil {
		return fmt.Errorf("%s not found", desc)
	}
	result := trim(selection.Text())
	matched, err := regexp.MatchString(expected, result)
	if err != nil || !matched {
		return fmt.Errorf("%s found with value '%s' but expected '%s'", desc, result, expected)
	}
	return nil
}

func (c *Client) TestValueRegexIsNot(previousError error, selector, desc, expected string) error {
	if previousError != nil {
		return previousError
	}
	selection, err := c.Find(selector)
	if err != nil {
		return fmt.Errorf("searching for %s got: %w", desc, err)
	}
	if selection == nil {
		return fmt.Errorf("%s not found", desc)
	}
	result := trim(selection.Text())
	matched, err := regexp.MatchString(expected, result)
	if err != nil || matched {
		return fmt.Errorf("%s found with value '%s'", desc, result)
	}
	return nil
}

func (c *Client) TestLength(previousError error, selector, desc string, expectedLength int) error {
	if previousError != nil {
		return previousError
	}
	selection, err := c.Find(selector)
	if err != nil {
		return fmt.Errorf("searching for %s got: %w", desc, err)
	}
	if selection == nil {
		return fmt.Errorf("%s not found", desc)
	}
	if len(selection.Nodes) != expectedLength {
		return fmt.Errorf("expected node {%s} length to be '%d' but got '%d'", selector, expectedLength, len(selection.Nodes))
	}
	return nil
}

func (c *Client) TestAttributeValue(previousError error, selector, desc, attrName, value string) error {
	if previousError != nil {
		return previousError
	}
	selection, err := c.Find(selector)
	if err != nil {
		return fmt.Errorf("searching for %s got: %w", desc, err)
	}
	if selection == nil {
		return fmt.Errorf("%s not found", desc)
	}
	v := getAttribute(selection.Nodes[0], attrName)
	if v != value {
		return fmt.Errorf("expected selector {%s} attribute '%s' to be '%s' but got '%s'", selector, attrName, value, v)
	}
	return nil
}

func (c *Client) GetAttribute(previousError error, selector, desc, attrName string) (string, error) {
	if previousError != nil {
		return "", previousError
	}
	selection, err := c.Find(selector)
	if err != nil {
		return "", fmt.Errorf("searching for %s got: %w", desc, err)
	}
	if selection == nil {
		return "", fmt.Errorf("%s not found", desc)
	}
	v := getAttribute(selection.Nodes[0], attrName)

	return v, nil
}

func (c *Client) TestAttributeValueRegex(previousError error, selector, desc, attrName, regex string) error {
	if previousError != nil {
		return previousError
	}
	selection, err := c.Find(selector)
	if err != nil {
		return fmt.Errorf("searching for %s got: %w", desc, err)
	}
	if selection == nil {
		return fmt.Errorf("%s not found", desc)
	}
	v := getAttribute(selection.Nodes[0], attrName)
	matched, err := regexp.MatchString(regex, v)
	if err != nil || !matched {
		return fmt.Errorf("%s found with value '%s' but expected '%s'", desc, v, regex)
	}
	return nil
}

func (c *Client) TestAttributeValueIsNot(previousError error, selector, desc, attrName, value string) error {
	if previousError != nil {
		return previousError
	}
	selection, err := c.Find(selector)
	if err != nil {
		return fmt.Errorf("searching for %s got: %w", desc, err)
	}
	if selection == nil {
		return fmt.Errorf("%s not found", desc)
	}
	v := getAttribute(selection.Nodes[0], attrName)
	if v == value {
		return fmt.Errorf("expected selector {%s} attribute '%s' not to be '%s' but got '%s'", selector, attrName, value, v)
	}
	return nil
}

func (c *Client) TestTable(previousError error, selector, desc string, expected [][]string) error {
	if previousError != nil {
		return previousError
	}
	table, err := c.Find(selector)
	if err != nil {
		return fmt.Errorf("searching for %s got: %w", desc, err)
	}
	if table == nil {
		return fmt.Errorf("%s not found", desc)
	}

	tableBody, err := table.Find("tbody")
	if err != nil {
		return fmt.Errorf("while searching for tbody in table of %s : %w", desc, err)
	}
	if tableBody == nil {
		return fmt.Errorf("table %s no tbody", desc)
	}

	err = c.checkTableContent(tableBody, expected)
	if err != nil {
		return fmt.Errorf("while checking table content of %s : %w", desc, err)
	}

	return nil
}

func (c *Client) TestTableRegex(previousError error, selector, desc string, expected [][]string) error {
	if previousError != nil {
		return previousError
	}
	table, err := c.Find(selector)
	if err != nil {
		return fmt.Errorf("searching for %s got: %w", desc, err)
	}
	if table == nil {
		return fmt.Errorf("%s not found", desc)
	}

	tableBody, err := table.Find("tbody")
	if err != nil {
		return fmt.Errorf("while searching for tbody in table of %s : %w", desc, err)
	}
	if tableBody == nil {
		return fmt.Errorf("table %s no tbody", desc)
	}

	err = c.checkTableContentRegex(tableBody, expected)
	if err != nil {
		return fmt.Errorf("while checking table content of %s : %w", desc, err)
	}

	return nil
}

func (c *Client) TestNodeDoesNotExist(previousError error, selector, desc string) error {
	if previousError != nil {
		return previousError
	}
	selection, err := c.Find(selector)
	if err != nil {
		return fmt.Errorf("searching for %s got: %w", desc, err)
	}
	if selection != nil && selection.Nodes != nil {
		return fmt.Errorf("node %s should not exist", selector)
	}

	return nil
}

func (c *Client) checkTableContent(s *Selection, expected [][]string) error {
	for i, line := range expected {
		for j, cell := range line {
			cellSelector := fmt.Sprintf("tr:nth-child(%d)>td:nth-child(%d)", i+1, j+1)
			result, err := s.Find(cellSelector)
			if err != nil {
				return fmt.Errorf("searching for cells got: %w", err)
			}
			if result == nil {
				return fmt.Errorf("cell %d,%d not found", j, i)
			}
			resultValue := trim(result.Text())
			if resultValue != cell {
				return fmt.Errorf("in cell %d,%d : found with value '%s' but expected '%s'", j, i, resultValue, cell)
			}
		}
	}
	return nil
}

func (c *Client) checkTableContentRegex(s *Selection, expected [][]string) error {
	for i, line := range expected {
		for j, cell := range line {
			cellSelector := fmt.Sprintf("tr:nth-child(%d)>td:nth-child(%d)", i+1, j+1)
			result, err := s.Find(cellSelector)
			if err != nil {
				return fmt.Errorf("searching for cells got: %w", err)
			}
			if result == nil {
				return fmt.Errorf("cell %d,%d not found", j, i)
			}
			resultValue := trim(result.Text())
			matched, err := regexp.MatchString(cell, resultValue)
			if err != nil || !matched {
				return fmt.Errorf("in cell %d,%d : found with value '%s' but expected '%s'", j, i, resultValue, cell)
			}
		}
	}
	return nil
}

func (c *Client) TestFieldValue(previousError error, selector, desc, value string) error {
	if previousError != nil {
		return previousError
	}
	selection, err := c.Find(selector)
	if err != nil {
		return fmt.Errorf("searching for %s got: %w", desc, err)
	}
	if selection == nil {
		return fmt.Errorf("%s not found", desc)
	}
	if selection.Nodes == nil {
		return fmt.Errorf("%s not found", desc)
	}
	v := getAttribute(selection.Nodes[0], "value")
	if v != value {
		return fmt.Errorf("expected selector {%s} value to be '%s' but got '%s'", selector, value, v)
	}

	return nil
}

func (c *Client) TestURL(previousError error, url string) error {
	if previousError != nil {
		return previousError
	}

	if c.currentURL != url {
		return fmt.Errorf("current Url: %s, want: %s", c.currentURL, url)
	}
	return nil
}

func (c *Client) TestURLStartWith(previousError error, url string) error {
	if previousError != nil {
		return previousError
	}

	if len(url) > len(c.currentURL) {
		return fmt.Errorf("current Url: %s, want to start with: %s", c.currentURL, url)
	}

	if c.currentURL[:len(url)] != url {
		return fmt.Errorf("current Url: %s, want to start with: %s", c.currentURL, url)
	}
	return nil
}
