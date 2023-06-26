package dom // import "github.com/HalCanary/facility/dom"

Copyright 2022 Hal Canary Use of this program is governed by the file LICENSE.

This is a lightweight wrapper around x/net/html, providing some extra
functionality.

CONSTANTS

const (
	ElementNode = html.ElementNode
	TextNode    = html.TextNode
)

FUNCTIONS

func ExtractText(root *Node) string
    Extract and combine all Text Nodes under given node.

func GetAttribute(node *Node, key string) string
    Find the matching attributes, ignoring namespace.

func RenderHTML(root *Node, w io.Writer) error
    Generates HTML5 doc.

func RenderHTMLExperimental(root *Node, w io.Writer) error
    Generates HTML5 doc.

func RenderXHTMLDoc(root *Node, w io.Writer) error
    Generates XHTML1 doc.

func TextBytes(root *Node) int

TYPES

type Attr = map[string]string

type Attribute = html.Attribute

type Node = html.Node

func AddAttribute(node *Node, k, v string) *Node
    Add an attribute to an element node, returning the node.

func Append(node *Node, children ...*Node) *Node
    Append children to the node, returning the node.

func Comment(data string) *Node
    Return a HTML comment with the given data.

func Elem(tag string, children ...*Node) *Node
    Return an element with the given children.

func Element(tag string, attributes Attr, children ...*Node) *Node
    Return an element with given attributes and children.

func FindNodeByAttribute(node *Node, key, value string) *Node
    Return the first matching node with key attribude set to value.

func FindNodeById(node *Node, id string) *Node
    Return the first matching node with `id` attribude set to id.

func FindNodeByTag(node *Node, tag string) *Node
    Return the first matching node

func FindNodeByTagAndAttrib(root *Node, tag, key, value string) *Node
    Return the first matching node. If tag is "", match key="value". If key is
    "", match on tag.

func FindNodesByTagAndAttrib(root *Node, tag, key, value string) []*Node
    Return all matching nodes. If tag is "", match key="value". If key is "",
    match on tag.

func Parse(source io.Reader) (*Node, error)
    Wrapper for html.Parse.

func ParseFragment(source io.Reader, context *Node) ([]*Node, error)
    Wrapper for html.ParseFragment.

func RawHtml(data string) *Node
    Return a Node with the given raw html.

func Remove(node *Node) *Node
    Remove a node from its parent.

func Text(data string) *Node
    Return a HTML node with the given text.

