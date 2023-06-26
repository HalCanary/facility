package ebook

import (
	"strings"

	"github.com/HalCanary/facility/dom"
)

func getAttribute(dst *string, root *dom.Node, tag, key, value, attribute string) {
	if *dst == "" {
		*dst = dom.GetAttribute(dom.FindNodeByTagAndAttrib(root, tag, key, value), attribute)
	}
}

func getTextByTagAndAttrib(dst *string, root *dom.Node, tag, key, value string) {
	if *dst == "" {
		*dst = strings.TrimSpace(dom.ExtractText(dom.FindNodeByTagAndAttrib(root, tag, key, value)))
	}
}

// Populate `info` based on common patterns.
func PopulateInfo(info *EbookInfo, doc *dom.Node) {
	if info == nil || doc == nil {
		return
	}
	getAttribute(&info.Title, doc, "meta", "name", "twitter:title", "content")
	getTextByTagAndAttrib(&info.Title, doc, "h1", "", "")
	getAttribute(&info.Authors, doc, "meta", "property", "books:author", "content")
	getAttribute(&info.Authors, doc, "meta", "name", "twitter:creator", "content")
	getTextByTagAndAttrib(&info.Authors, doc, "a", "rel", "author")
	getTextByTagAndAttrib(&info.Comments, doc, "div", "property", "description")
	getTextByTagAndAttrib(&info.Comments, doc, "div", "class", "description")
	getAttribute(&info.Comments, doc, "meta", "property", "og:description", "content")
	getAttribute(&info.Language, doc, "html", "", "", "lang")
}
