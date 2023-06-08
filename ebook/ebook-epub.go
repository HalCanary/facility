package ebook

// Copyright 2022 Hal Canary
// Use of this program is governed by the file LICENSE.

import (
	"bytes"
	"encoding/xml"
	"fmt"
	"io"
)

func makePackage(info EbookInfo, uuid string, dst io.Writer, cover bool) error {
	manifestItems := []xmlItem{
		xmlItem{Id: "frontmatter", Href: "frontmatter.xhtml", MediaType: "application/xhtml+xml"},
		xmlItem{Id: "toc", Href: "toc.xhtml", MediaType: "application/xhtml+xml",
			Attributes: []xml.Attr{xml.Attr{Name: xml.Name{Local: "properties"}, Value: "nav"}}},
		xmlItem{Id: "ncx", Href: "toc.ncx", MediaType: "application/x-dtbncx+xml"},
	}
	itemrefs := []xmlItemref{
		xmlItemref{Idref: "frontmatter"},
		xmlItemref{Idref: "toc"},
	}
	if cover {
		manifestItems = append(manifestItems, xmlItem{Id: "cover", Href: "cover.jpg", MediaType: "image/jpeg",
			Attributes: []xml.Attr{xml.Attr{Name: xml.Name{Local: "properties"}, Value: "cover-image"}}})
	}
	for i, _ := range info.Chapters {
		fn := fmt.Sprintf("%04d", i)
		id := "ch" + fn
		manifestItems = append(manifestItems, xmlItem{Id: id, Href: fn + ".xhtml", MediaType: "application/xhtml+xml"})
		itemrefs = append(itemrefs, xmlItemref{Idref: id})
	}
	modified := info.Modified.UTC().Format("2006-01-02T15:04:05Z")
	description := fmt.Sprintf("%s\n\nSOURCE: %s\nCHAPTERS: %d\n", info.Comments, info.Source, len(info.Chapters))
	p := xmlPackage{
		Xmlns:            "http://www.idpf.org/2007/opf",
		XmlnsOpf:         "http://www.idpf.org/2007/opf",
		Version:          "3.0",
		UniqueIdentifier: "BookID",
		Metadata: xmlMetaData{
			XmlnsDC: "http://purl.org/dc/elements/1.1/",
			Properties: []xmlMetaProperty{
				xmlMetaProperty{Property: "dcterms:modified", Value: modified},
			},
			MetaItems: []xmlMetaItems{
				xmlMetaItems{
					XMLName:    xml.Name{Local: "dc:identifier"},
					Value:      uuid,
					Attributes: []xml.Attr{xml.Attr{Name: xml.Name{Local: "id"}, Value: "BookID"}},
				},
				xmlMetaItems{XMLName: xml.Name{Local: "dc:title"}, Value: info.Title},
				xmlMetaItems{XMLName: xml.Name{Local: "dc:language"}, Value: info.Language},
				xmlMetaItems{XMLName: xml.Name{Local: "dc:creator"}, Value: info.Authors},
				xmlMetaItems{XMLName: xml.Name{Local: "dc:description"}, Value: description},
				xmlMetaItems{XMLName: xml.Name{Local: "dc:source"}, Value: info.Source},
				xmlMetaItems{XMLName: xml.Name{Local: "dc:date"}, Value: modified},
			},
		},
		ManifestItems: manifestItems,
		Spine: xmlSpine{
			Toc:      "ncx",
			Itemrefs: itemrefs,
		},
		GuideRefs: []xmlGuideReference{
			xmlGuideReference{Title: "Cover page", Type: "cover", Href: "frontmatter.xhtml"},
			xmlGuideReference{Title: "Table of content", Type: "toc", Href: "toc.xhtml"},
		},
	}
	encoded, err := xml.MarshalIndent(&p, "", " ")
	if err != nil {
		return err
	}
	encoded = bytes.ReplaceAll(encoded, []byte("></item>"), []byte("/>"))
	encoded = bytes.ReplaceAll(encoded, []byte("></itemref>"), []byte("/>"))
	encoded = bytes.ReplaceAll(encoded, []byte("></reference>"), []byte("/>"))
	_, err = dst.Write(encoded)
	return err
}

type xmlPackage struct {
	XMLName          xml.Name            `xml:"package"`
	Xmlns            string              `xml:"xmlns,attr"`
	XmlnsOpf         string              `xml:"xmlns:opf,attr"`
	Version          string              `xml:"version,attr"`
	UniqueIdentifier string              `xml:"unique-identifier,attr"`
	Metadata         xmlMetaData         `xml:"metadata"`
	ManifestItems    []xmlItem           `xml:"manifest>item"`
	Spine            xmlSpine            `xml:"spine"`
	GuideRefs        []xmlGuideReference `xml:"guide>reference"`
}

type xmlMetaData struct {
	XmlnsDC    string            `xml:"xmlns:dc,attr"`
	Properties []xmlMetaProperty `xml:"meta"`
	MetaItems  []xmlMetaItems    `xml:",any"`
}

type xmlMetaItems struct {
	XMLName    xml.Name
	Value      string     `xml:",chardata"`
	Attributes []xml.Attr `xml:",attr,any"`
}

type xmlMetaProperty struct {
	Property string `xml:"property,attr"`
	Value    string `xml:",chardata"`
}

type xmlItem struct {
	Id         string     `xml:"id,attr"`
	Href       string     `xml:"href,attr"`
	MediaType  string     `xml:"media-type,attr"`
	Attributes []xml.Attr `xml:",attr,any"`
}

type xmlSpine struct {
	Toc      string       `xml:"toc,attr"`
	Itemrefs []xmlItemref `xml:"itemref"`
}

type xmlItemref struct {
	Idref string `xml:"idref,attr"`
}

type xmlGuideReference struct {
	Title string `xml:"title,attr"`
	Type  string `xml:"type,attr"`
	Href  string `xml:"href,attr"`
}

func makeNCX(info EbookInfo, uid string, dst io.Writer) error {
	nav := []navPointXml{
		navPointXml{Class: "chapter", Id: "frontmatter", PlayOrder: 0, Label: "Front Matter", Content: contentXml{Src: "frontmatter.xhtml"}},
	}
	for i, ch := range info.Chapters {
		fn := fmt.Sprintf("%04d", i)
		id := "ch" + fn
		label := fmt.Sprintf("%d. %s", i+1, ch.Title)
		nav = append(nav, navPointXml{Class: "chapter", Id: id, PlayOrder: i + 1, Label: label, Content: contentXml{Src: "" + fn + ".xhtml"}})
	}
	ncx := ncxXml{
		Xmlns:   "http://www.daisy.org/z3986/2005/ncx/",
		Version: "2005-1",
		Lang:    "en",
		Metas: []metaNcxXml{
			metaNcxXml{Name: "dtb:uid", Content: uid},
			metaNcxXml{Name: "dtb:depth", Content: "1"},
			metaNcxXml{Name: "dtb:totalPageCount", Content: "0"},
			metaNcxXml{Name: "dtb:maxPageNumber", Content: "0"},
		},
		Title:     info.Title,
		Author:    info.Authors,
		NavPoints: nav,
	}
	encoded, err := xml.MarshalIndent(&ncx, "", " ")
	if err != nil {
		return err
	}
	encoded = bytes.ReplaceAll(encoded, []byte("></meta>"), []byte("/>"))
	encoded = bytes.ReplaceAll(encoded, []byte("></content>"), []byte("/>"))
	_, err = dst.Write(encoded)
	return err
}

type ncxXml struct {
	XMLName   xml.Name      `xml:"ncx"`
	Xmlns     string        `xml:"xmlns,attr"`
	Version   string        `xml:"version,attr"`
	Lang      string        `xml:"xml:lang,attr"`
	Metas     []metaNcxXml  `xml:"head>meta"`
	Title     string        `xml:"docTitle>text"`
	Author    string        `xml:"docAuthor>text"`
	NavPoints []navPointXml `xml:"navMap>navPoint"`
}
type metaNcxXml struct {
	Name    string `xml:"name,attr"`
	Content string `xml:"content,attr"`
}
type navPointXml struct {
	Class     string     `xml:"class,attr"`
	Id        string     `xml:"id,attr"`
	PlayOrder int        `xml:"playOrder,attr"`
	Label     string     `xml:"navLabel>text"`
	Content   contentXml `xml:"content"`
}
type contentXml struct {
	Src string `xml:"src,attr"`
}
