package scrapper

import (
	"bytes"
	"fmt"
	"io"

	"github.com/andybalholm/cascadia"
	"golang.org/x/net/html"
)

type Document struct {
	Selection *Selection
	RootNode  *html.Node
}

func NewDocumentFromReader(r io.Reader) (*Document, error) {
	root, err := html.Parse(r)
	if err != nil {
		return nil, err
	}

	d := &Document{
		RootNode: root,
	}
	d.Selection = &Selection{
		Nodes:    []*html.Node{root},
		document: d,
	}

	return d, nil
}

func (d *Document) Find(selector string) (*Selection, error) {
	return d.Selection.Find(selector)
}

type Selection struct {
	Nodes    []*html.Node
	document *Document
	prevSel  *Selection
}

func (s *Selection) Text() string {
	buf := bytes.Buffer{}
	for _, n := range s.Nodes {
		writeNodeText(&buf, n)
	}

	return buf.String()
}

func writeNodeText(w io.Writer, n *html.Node) {
	if n.Type == html.TextNode {
		fmt.Fprint(w, n.Data)
	}

	if n.FirstChild == nil {
		return
	}

	for c := n.FirstChild; c != nil; c = c.NextSibling {
		writeNodeText(w, c)
	}
}

func (s *Selection) Find(selector string) (*Selection, error) {
	m, err := cascadia.Compile(selector)
	if err != nil {
		return nil, fmt.Errorf("can't compile selector %s : %w", selector, err)
	}

	nodes := findWithMatcher(s.Nodes, m)

	if len(nodes) == 0 {
		return nil, nil
	}

	return &Selection{
		Nodes:    nodes,
		document: s.document,
		prevSel:  s,
	}, nil
}

type Matcher interface {
	Match(*html.Node) bool
	MatchAll(*html.Node) []*html.Node
	Filter([]*html.Node) []*html.Node
}

func findWithMatcher(nodes []*html.Node, m Matcher) []*html.Node {
	set := map[*html.Node]struct{}{}

	for _, n := range nodes {
		children := getMatchingChild(n, m)
		for _, c := range children {
			set[c] = struct{}{}
		}
	}

	result := []*html.Node{}
	for n := range set {
		result = append(result, n)
	}

	return result
}

func getMatchingChild(n *html.Node, m Matcher) []*html.Node {
	result := []*html.Node{}
	for c := n.FirstChild; c != nil; c = c.NextSibling {
		if c.Type == html.ElementNode {
			result = append(result, m.MatchAll(c)...)
		}
	}
	return result
}
