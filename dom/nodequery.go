package dom

// Copyright 2022 Hal Canary
// Use of this program is governed by the file LICENSE.

// Return the first matching node
func FindNodeByTag(node *Node, tag string) *Node {
	return FindNodeByTagAndAttrib(node, tag, "", "")
}

// Return the first matching node with `id` attribude set to id.
func FindNodeById(node *Node, id string) *Node {
	return FindNodeByTagAndAttrib(node, "", "id", id)
}

// Return the first matching node with key attribude set to value.
func FindNodeByAttribute(node *Node, key, value string) *Node {
	return FindNodeByTagAndAttrib(node, "", key, value)
}

// Return all matching nodes.  If tag is "", match key="value".
// If key is "", match on tag.
func FindNodesByTagAndAttrib(root *Node, tag, key, value string) []*Node {
	return findNodesByTagAndAttrib(root, tag, key, value, false)
}

// Return the first matching node.  If tag is "", match key="value".
// If key is "", match on tag.
func FindNodeByTagAndAttrib(root *Node, tag, key, value string) *Node {
	result := findNodesByTagAndAttrib(root, tag, key, value, true)
	if len(result) == 0 {
		return nil
	}
	return result[0]
}

func findNodesByTagAndAttrib(root *Node, tag, key, value string, stopEarly bool) []*Node {
	var result []*Node
	// I unrolled a recursive function to use no recursion, or heap allocation,
	// only loops.  This is allowed by the existance of Parent pointer.
	if root == nil {
		return result
	}
	node := root
	for {
		if node.Type == ElementNode && (tag == "" || node.Data == tag) {
			if key == "" {
				result = append(result, node)
				if stopEarly {
					return result
				}
			} else {
				for _, attr := range node.Attr {
					if attr.Key == key && attr.Val == value {
						result = append(result, node)
						if stopEarly {
							return result
						}
					}
				}
			}
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
