// Copyright 2022 Hal Canary
// Use of this program is governed by the file LICENSE.
//
// This is a lightweight wrapper around x/net/html, providing some extra
// functionality.
package dom

import (
	"bytes"
	"encoding/xml"
	"io"
	"regexp"
	"sort"
	"strings"

	"golang.org/x/net/html"
)

type (
	Node      = html.Node
	Attribute = html.Attribute
	Attr      = map[string]string
)

const (
	ElementNode = html.ElementNode
	TextNode    = html.TextNode
)

var whitespaceRegexp = regexp.MustCompile("\\s+")

// Wrapper for html.Parse.
func Parse(source io.Reader) (*Node, error) {
	return html.Parse(source)
}

// Wrapper for html.ParseFragment.
func ParseFragment(source io.Reader, context *Node) ([]*Node, error) {
	return html.ParseFragment(source, context)
}

var dashesRegexp = regexp.MustCompile("---*")

// Return a HTML comment with the given data.
func Comment(data string) *Node {
	if data == "" {
		return nil
	}
	data = dashesRegexp.ReplaceAllStringFunc(data, func(s string) string {
		return string(bytes.Repeat([]byte{'~'}, len(s)))
	})
	return &Node{Type: html.CommentNode, Data: data}
}

// Return a HTML node with the given text.
func Text(data string) *Node {
	return &Node{Type: html.TextNode, Data: data}
}

// Return an element with given attributes and children.
func Element(tag string, attributes Attr, children ...*Node) *Node {
	node := &Node{Type: html.ElementNode, Data: tag}
	keys := make([]string, 0, len(attributes))
	for k := range attributes {
		keys = append(keys, k)
	}
	sort.Strings(keys)
	for _, k := range keys {
		node.Attr = append(node.Attr, makeAttribute(k, attributes[k]))
	}
	return Append(node, children...)
}

// Return a Node with the given raw html.
func RawHtml(data string) *Node {
	return &Node{Type: html.RawNode, Data: data}
}

// Append children to the node, returning the node.
func Append(node *Node, children ...*Node) *Node {
	if node != nil && (node.Type == html.ElementNode || node.Type == html.DocumentNode) {
		for _, c := range children {
			if c != nil {
				node.AppendChild(c)
			}
		}
	}
	return node
}

func makeAttribute(k, v string) html.Attribute {
	if ns, key, found := strings.Cut(k, ":"); found {
		return html.Attribute{Namespace: ns, Key: key, Val: v}
	} else {
		return html.Attribute{Key: k, Val: v}
	}
}

// Add an attribute to an element node, returning the node.
func AddAttribute(node *Node, k, v string) *Node {
	if node != nil && node.Type == html.ElementNode {
		node.Attr = append(node.Attr, makeAttribute(k, v))
	}
	return node
}

// Return an element with the given children.
func Elem(tag string, children ...*Node) *Node {
	return Element(tag, nil, children...)
}

// Generates HTML5 doc.
func RenderHTML(root *Node, w io.Writer) error {
	d := Node{Type: html.DocumentNode}
	Append(&d, &Node{Type: html.DoctypeNode, Data: "html"}, Text("\n"), root)
	e := html.Render(w, &d)
	w.Write([]byte{'\n'})
	return e
}

// Generates HTML5 doc.
func RenderHTMLExperimental(root *Node, w io.Writer) error {
	d := Node{Type: html.DocumentNode}
	Append(&d, &Node{Type: html.DoctypeNode, Data: "html"}, Text("\n"), root)
	cw := checkedWriter{Writer: w}
	renderXHTML(&cw, &d, false)
	cw.Write([]byte{'\n'})
	return cw.Error
}

// Generates XHTML1 doc.
func RenderXHTMLDoc(root *Node, w io.Writer) error {
	if root == nil || w == nil {
		return nil
	}
	io.WriteString(w, xml.Header)
	cw := checkedWriter{Writer: w}
	renderXHTML(&cw, root, true)
	cw.Write([]byte{'\n'})
	return cw.Error
}

type checkedWriter struct {
	io.Writer
	Error error
}

func (w *checkedWriter) Write(b []byte) {
	if w.Error == nil {
		_, w.Error = w.Writer.Write(b)
	}
}

func (w *checkedWriter) WriteString(s string) {
	if w.Error == nil {
		_, w.Error = io.WriteString(w.Writer, s)
	}
}

var xhtmlattribs = map[string]struct{}{
	"alt":        struct{}{},
	"border":     struct{}{},
	"class":      struct{}{},
	"content":    struct{}{},
	"dir":        struct{}{},
	"href":       struct{}{},
	"http-equiv": struct{}{},
	"id":         struct{}{},
	"lang":       struct{}{},
	"name":       struct{}{},
	"src":        struct{}{},
	"style":      struct{}{},
	"title":      struct{}{},
	"type":       struct{}{},
	"xmlns":      struct{}{},
}

var htmlVoidElements = map[string]struct{}{
	"area":   struct{}{},
	"base":   struct{}{},
	"br":     struct{}{},
	"col":    struct{}{},
	"embed":  struct{}{},
	"hr":     struct{}{},
	"img":    struct{}{},
	"input":  struct{}{},
	"link":   struct{}{},
	"meta":   struct{}{},
	"source": struct{}{},
	"track":  struct{}{},
	"wbr":    struct{}{},
}

func renderXHTML(w *checkedWriter, node *Node, xhtml bool) {
	switch node.Type {
	case html.DoctypeNode:
		if node.Data == "html" {
			w.WriteString("<!DOCTYPE html>")
		}
	case html.DocumentNode:
		for c := node.FirstChild; c != nil; c = c.NextSibling {
			if w.Error == nil {
				renderXHTML(w, c, xhtml)
			}
		}
	case html.ElementNode:
		w.Write([]byte{'<'})
		w.WriteString(node.Data)
		for _, attr := range node.Attr {
			ok := !xhtml || attr.Namespace != ""
			if !ok {
				_, ok = xhtmlattribs[attr.Key]
			}
			if ok {
				w.Write([]byte{' '})
				if attr.Namespace != "" {
					w.WriteString(attr.Namespace)
					w.Write([]byte{':'})
				}
				w.WriteString(attr.Key)
				w.Write([]byte{'=', '"'})
				w.WriteString(html.EscapeString(attr.Val))
				w.Write([]byte{'"'})
			}
		}
		if node.FirstChild == nil {
			if xhtml {
				w.Write([]byte{'/', '>'})
			} else {
				_, isVoidElement := htmlVoidElements[node.Data]
				if isVoidElement {
					w.Write([]byte{'>'})
				} else {
					w.Write([]byte{'/', '>'})
				}
			}
		} else {
			w.Write([]byte{'>'})
			for c := node.FirstChild; c != nil; c = c.NextSibling {
				if w.Error == nil {
					if (node.Data == "script" || node.Data == "style") && c.Type == html.TextNode {
						w.WriteString(c.Data)
					} else {
						renderXHTML(w, c, xhtml)
					}
				}
			}
			w.Write([]byte{'<', '/'})
			w.WriteString(node.Data)
			w.Write([]byte{'>'})
		}
	case html.TextNode:
		w.WriteString(html.EscapeString(node.Data))
	case html.CommentNode:
		w.Write([]byte{'<', '!', '-', '-'})
		w.WriteString(node.Data)
		w.Write([]byte{'-', '-', '>'})
	}
}

// Find the matching attributes, ignoring namespace.
func GetAttribute(node *Node, key string) string {
	if node != nil {
		for _, attr := range node.Attr {
			if attr.Key == key {
				return attr.Val
			}
		}
	}
	return ""
}

// Extract and combine all Text Nodes under given node.
func ExtractText(root *Node) string {
	var result strings.Builder
	if root == nil {
		return result.String()
	}
	node := root
	for {
		if node.Type == html.TextNode {
			result.WriteString(whitespaceRegexp.ReplaceAllString(node.Data, " "))
		} else if node.Type == html.ElementNode {
			switch node.Data {
			case "br":
				result.WriteString("\n")
			case "hr":
				result.WriteString("\n* * *\n")
			case "p":
				result.WriteString("\n\n")
			case "img":
				result.WriteString(GetAttribute(node, "alt"))
			}
		}
		if node.FirstChild != nil {
			node = node.FirstChild
			continue
		}
		for {
			if node == root || node == nil {
				return result.String()
			}
			if node.NextSibling != nil {
				node = node.NextSibling
				break
			}
			node = node.Parent
		}
	}
}

// Return the approximate number of bytes that `ExtractText` would return.
func TextBytes(root *Node) int {
	var result int = 0
	if root == nil {
		return result
	}
	node := root
	for {
		if node.Type == html.TextNode {
			result += len(root.Data)
		}
		if node.FirstChild != nil {
			node = node.FirstChild
			continue
		}
		for {
			if node == root || node == nil {
				return result
			}
			if node.NextSibling != nil {
				node = node.NextSibling
				break
			}
			node = node.Parent
		}
	}
}

// Remove a node from its parent.
func Remove(node *Node) *Node {
	if node != nil && node.Parent != nil {
		node.Parent.RemoveChild(node)
	}
	return node
}
