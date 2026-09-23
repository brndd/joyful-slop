package gremlin

import (
	"encoding/xml"
	"fmt"
	"io"
	"strings"
)

type node struct {
	Name     string
	Attr     map[string]string
	Text     string
	Children []*node
}

func parseXML(r io.Reader) (*node, error) {
	dec := xml.NewDecoder(r)
	var root *node
	stack := []*node{}
	for {
		tok, err := dec.Token()
		if err == io.EOF {
			break
		}
		if err != nil {
			return nil, fmt.Errorf("parse XML: %w", err)
		}
		switch t := tok.(type) {
		case xml.StartElement:
			n := &node{Name: t.Name.Local, Attr: map[string]string{}}
			for _, a := range t.Attr {
				n.Attr[a.Name.Local] = a.Value
			}
			if len(stack) == 0 {
				if root != nil {
					return nil, fmt.Errorf("multiple XML roots")
				}
				root = n
			} else {
				stack[len(stack)-1].Children = append(stack[len(stack)-1].Children, n)
			}
			stack = append(stack, n)
		case xml.CharData:
			if len(stack) > 0 {
				stack[len(stack)-1].Text += string(t)
			}
		case xml.EndElement:
			if len(stack) == 0 || stack[len(stack)-1].Name != t.Name.Local {
				return nil, fmt.Errorf("unbalanced XML element %q", t.Name.Local)
			}
			stack[len(stack)-1].Text = strings.TrimSpace(stack[len(stack)-1].Text)
			stack = stack[:len(stack)-1]
		}
	}
	if root == nil || len(stack) != 0 {
		return nil, fmt.Errorf("incomplete XML document")
	}
	return root, nil
}

func (n *node) children(name string) []*node {
	var out []*node
	for _, child := range n.Children {
		if child.Name == name {
			out = append(out, child)
		}
	}
	return out
}

func (n *node) child(name string) (*node, error) {
	children := n.children(name)
	if len(children) != 1 {
		return nil, fmt.Errorf("%s: expected one <%s>, found %d", n.Name, name, len(children))
	}
	return children[0], nil
}

func (n *node) value(name string) (string, error) {
	child, err := n.child(name)
	if err != nil {
		return "", err
	}
	if child.Text == "" {
		return "", fmt.Errorf("%s/%s is empty", n.Name, name)
	}
	return child.Text, nil
}

func properties(n *node) (map[string][]string, error) {
	out := map[string][]string{}
	for _, p := range n.children("property") {
		name, err := p.value("name")
		if err != nil {
			return nil, err
		}
		value, err := p.value("value")
		if err != nil {
			return nil, err
		}
		out[name] = append(out[name], value)
	}
	return out, nil
}

func oneProperty(props map[string][]string, name string) (string, error) {
	values := props[name]
	if len(values) != 1 {
		return "", fmt.Errorf("expected one property %q, found %d", name, len(values))
	}
	return values[0], nil
}
