package scrapper

import (
	"fmt"

	"golang.org/x/net/html"
)

type Form struct {
	Method    string
	Action    string
	Attribute map[string]string
}

func GetForm(doc *Document, name string) (Form, error) {
	form := Form{Attribute: map[string]string{}}
	targetform := fmt.Sprintf(`form[name="%s"]`, name)

	selection, err := doc.Find(targetform)
	if err != nil {
		return form, fmt.Errorf("while searching for form %w", err)
	}
	if selection == nil {
		return form, fmt.Errorf("form %s not found", name)
	}

	if len(selection.Nodes) != 1 {
		return form, fmt.Errorf("cannot find form with selector '%s' found %d result (expect 1)", targetform, len(selection.Nodes))
	}

	nodeForm := selection.Nodes[0]
	form.Method = getAttribute(nodeForm, "method")

	if form.Method == "" {
		return form, fmt.Errorf("cannot process form '%s' no method found", targetform)
	}

	form.Action = getAttribute(nodeForm, "action")

	subItem := []string{"input", "select", "button", "textarea"}

	for _, item := range subItem {
		formItem := fmt.Sprintf(`form[name="%s"] %s`, name, item)
		selectedItem, err := doc.Find(formItem)

		if err != nil {
			return form, fmt.Errorf("while searching for the form item %s of %s : %w", item, name, err)
		}
		if selectedItem == nil {
			continue
		}
		for _, nodeItem := range selectedItem.Nodes {
			name := getAttribute(nodeItem, "name")
			form.Attribute[name] = item
		}
	}

	return form, nil
}

func getAttribute(node *html.Node, name string) string {
	for _, att := range node.Attr {
		if att.Key == name {
			return att.Val
		}
	}

	return ""
}
