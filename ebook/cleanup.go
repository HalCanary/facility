package ebook

// Copyright 2022 Hal Canary
// Use of this program is governed by the file LICENSE.

import (
	"net/url"
	"regexp"
	"strings"

	"github.com/HalCanary/facility/dom"
)

type Node = dom.Node

// Clean up a HTML fragment.
func Cleanup(node *Node) *Node {
	node = cleanupStyle(node)
	cleanupTables(node)
	cleanupCenter(node)
	cleanupDoubled(node)
	return node
}

var whiteSpaceOnly = regexp.MustCompile("^\\pZ*$")
var spaceOnly = regexp.MustCompile("^\\pZs*$")
var semiRegexp = regexp.MustCompile("\\s*;\\s*")

func styler(v string) string {
	var result []string
	for _, term := range semiRegexp.Split(v, -1) {
		switch term {
		case "background-attachment: initial",
			"background-clip: initial",
			"background-image: initial",
			"background-origin: initial",
			"background-position: initial",
			"background-repeat: initial",
			"background-size: initial",
			"break-before: page",
			"margin-bottom: 0in",
			"background: transparent",
			"font-family: Arial",
			"font-family: Segoe UI",
			"font-family: Segoe UI, sans-serif",
			"font-family: Segoe UI, serif",
			"font-style: normal",
			"font-variant: normal",
			"font-weight: normal",
			"margin-bottom: 0",
			"page-break-before: always",
			"text-decoration: none",
			"":
			// do nothing
		case "font-family: Courier New, monospace":
			result = append(result, "font-family:monospace")
		default:
			result = append(result, term)
		}
	}
	return strings.Join(result, ";")
}

func cleanupCenter(node *Node) {
	if node != nil && node.Type == dom.ElementNode {
		if node.Data == "center" {
			node.Data = "div"
			if i := getNodeAttributeIndex(node, "class"); i >= 0 {
				node.Attr[i].Val = node.Attr[i].Val + " mid"
			} else {
				dom.AddAttribute(node, "class", "mid")
			}
		}
		if node.Data == "big" {
			node.Data = "span"
			dom.AddAttribute(node, "style", "font-size:larger")
		}
		c := node.FirstChild
		for c != nil {
			next := c.NextSibling
			cleanupCenter(c)
			c = next
		}
	}
}

func cleanupDoubled(node *Node) {
	if node.Type == dom.ElementNode {
		data := node.Data
		for c := node.FirstChild; c != nil; {
			next := c.NextSibling
			cleanupDoubled(c)
			if data == "ul" && c.Type == dom.ElementNode && c.Data == data {
				dom.Remove(c)
				for c2 := c.FirstChild; c2 != nil; {
					n2 := c2.NextSibling
					c.RemoveChild(c2)
					node.InsertBefore(c2, next)
					c2 = n2
				}
			}
			c = next
		}
	}
}

func cleanupTables(node *Node) {
	if node != nil && node.Type == dom.ElementNode {
		if i := getNodeAttributeIndex(node, "border"); i >= 0 {
			v := node.Attr[i].Val
			if v != "1" && v != "" {
				if v == "none" {
					node.Attr[i].Val = ""
				} else {
					node.Attr[i].Val = "1"
				}
			}
		}

		c := node.FirstChild
		for c != nil {
			next := c.NextSibling
			cleanupTables(c)
			c = next
		}
		if node.FirstChild == nil {
			switch node.Data {
			case "tbody", "dd", "dl":
				dom.Remove(node)
			}
		}
	}
}

func resolve(oldUrl string, ref *url.URL) string {
	if u, _ := url.Parse(oldUrl); u != nil {
		return ref.ResolveReference(u).String()
	}
	return oldUrl
}

// Resolve all links in the tree, relative to `ref`.
func ResolveLinks(node *Node, ref *url.URL) *Node {
	if node != nil && node.Type == dom.ElementNode {
		if attr := getNodeAttribute(node, "href"); attr != nil {
			attr.Val = resolve(attr.Val, ref)
		}
		if attr := getNodeAttribute(node, "src"); attr != nil {
			attr.Val = resolve(attr.Val, ref)
		}
		for c := node.FirstChild; c != nil; c = c.NextSibling {
			ResolveLinks(c, ref)
		}
	}
	return node
}

func cleanupStyle(node *Node) *Node {
	if node != nil {
		switch node.Type {
		case dom.TextNode:
			if node.Data != "" && whiteSpaceOnly.MatchString(node.Data) {
				if !spaceOnly.MatchString(node.Data) {
					node.Data = "\n"
				}
			}
		case dom.ElementNode:
			if node.Data == "p" {
				if isWhitespaceOnly(node) {
					dom.Remove(node)
					return nil
				}
				if i := getNodeAttributeIndex(node, "align"); i >= 0 {
					switch node.Attr[i].Val {
					case "left":
						node.Attr = append(node.Attr[:i], node.Attr[i+1:]...)
					}
				}
			}
			if i := getNodeAttributeIndex(node, "style"); i >= 0 {
				v := styler(node.Attr[i].Val)
				if v == "" {
					node.Attr = append(node.Attr[:i], node.Attr[i+1:]...)
				} else {
					node.Attr[i].Val = v
				}
			}
			child := node.FirstChild
			for child != nil {
				next := child.NextSibling
				cleanupStyle(child)
				child = next
			}

			if node.Data == "span" && len(node.Attr) == 0 {
				if parent := node.Parent; parent != nil {
					nextSibling := node.NextSibling
					child := node.FirstChild
					for child != nil {
						next := child.NextSibling
						node.RemoveChild(child)
						parent.InsertBefore(child, nextSibling)
						child = next
					}
					parent.RemoveChild(node)
				}
			}
			if node.Data == "img" {
				if i := getNodeAttributeIndex(node, "src"); i >= 0 {
					if node.Attr[i].Val == "" {
						node.Attr[i].Val = "data:null;,"
					}
				} else {
					node.Attr = append(node.Attr, dom.Attribute{Key: "src", Val: "data:null;,"})
				}
			}
		}
	}
	return node
}

func countChildren(node *Node) int {
	var count int = 0
	if node != nil {
		for child := node.FirstChild; child != nil; child = child.NextSibling {
			count++
		}
	}
	return count
}

func getNodeAttributeIndex(node *Node, key string) int {
	if node != nil {
		for idx, attr := range node.Attr {
			if attr.Namespace == "" && attr.Key == key {
				return idx
			}
		}
	}
	return -1
}

func getNodeAttribute(node *Node, key string) *dom.Attribute {
	if node != nil {
		for idx, attr := range node.Attr {
			if attr.Namespace == "" && attr.Key == key {
				return &node.Attr[idx]
			}
		}
	}
	return nil
}

func isWhitespaceOnly(node *Node) bool {
	if node != nil {
		switch node.Type {
		case dom.TextNode:
			return whiteSpaceOnly.MatchString(node.Data)
		case dom.ElementNode:
			for child := node.FirstChild; child != nil; child = child.NextSibling {
				if !isWhitespaceOnly(child) {
					return false
				}
			}
		}
	}
	return true
}
