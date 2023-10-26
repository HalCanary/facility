// Copyright 2022 Hal Canary
// Use of this program is governed by the file LICENSE.
package ebook

import (
	"encoding/base64"
	"encoding/xml"
	"fmt"
	"io"
	"log"
	"math/rand"
	"net/http"
	"net/url"
	"strings"
	"time"

	"github.com/HalCanary/facility/dom"
	"github.com/HalCanary/facility/zipper"
)

// One Chapter of an Ebook.
type Chapter struct {
	Title    string
	Url      string
	Content  *Node
	Modified time.Time
}

// Ebook content and metadata.
type EbookInfo struct {
	Authors  string
	Comments string
	Title    string
	Source   string
	Language string
	Chapters []Chapter
	Modified time.Time
	Cover    []byte
}

const bookStyle = `
div p{text-indent:2em;margin-top:0;margin-bottom:0}
div p:first-child{text-indent:0;}
table, th, td { border:2px solid #808080; padding:3px; }
table { border-collapse:collapse; margin:3px; }
ol.flat {list-style-type:none;}
ol.flat li {list-style:none; display:inline;}
ol.flat li::after {content:"]";}
ol.flat li::before {content:"[";}
div.mid {margin: 0 auto;}
div.mid p {text-indent:0;}
div.center {margin-left:auto;margin-right:auto;}
`

const conatainer_xml = xml.Header + `<container version="1.0" xmlns="urn:oasis:names:tc:opendocument:xmlns:container">
<rootfiles>
<rootfile full-path="book/content.opf" media-type="application/oebps-package+xml"/>
</rootfiles>
</container>
`

// Return the time of most recently modified chapter.
func (info EbookInfo) CalculateLastModified() time.Time {
	var result time.Time = info.Modified
	for _, ch := range info.Chapters {
		if !ch.Modified.IsZero() && ch.Modified.After(result) {
			result = ch.Modified
		}
	}
	return result
}

func meta(name, content string) *Node {
	return dom.Element("meta", dom.Attr{"name": name, "content": content})
}

func head(title, style, comment string) *Node {
	return dom.Elem("head",
		dom.Element("meta", dom.Attr{
			"http-equiv": "Content-Type", "content": "text/html; charset=utf-8"}),
		dom.Comment(comment),
		meta("viewport", "width=device-width, initial-scale=1.0"),
		dom.Elem("title", dom.Text(title)),
		dom.Elem("style", dom.Text(style)),
	)
}

func nl() *Node {
	return dom.Text("\n")
}

func dataUrl(src []byte) string {
	return fmt.Sprintf("data:%s;base64,%s",
		http.DetectContentType(src),
		base64.StdEncoding.EncodeToString(src))
}

func (info *EbookInfo) Cleanup() {
	for i, chapter := range info.Chapters {
		info.Chapters[i].Content = Cleanup(chapter.Content)
		if chUrl, _ := url.Parse(chapter.Url); chUrl != nil {
			info.Chapters[i].Content = ResolveLinks(info.Chapters[i].Content, chUrl)
		}
	}
}

// Write the ebook as a single HTML file.
func (info EbookInfo) WriteHtml(dst io.Writer) error {
	body := dom.Elem("body", nl())
	if len(info.Cover) > 0 {
		dom.Append(body,
			dom.Element("div", dom.Attr{"style": "text-align:center"},
				nl(),
				imgElem(dataUrl(info.Cover), "[COVER]"),
				nl()),
			nl(),
			dom.Elem("hr"),
			nl(),
		)
	}
	dom.Append(body,
		dom.Elem("div", nl(),
			dom.Elem("h1", dom.Text(info.Title)), nl(),
			dom.Elem("div", dom.Text("Author: "+info.Authors)), nl(),
			dom.Elem("div",
				dom.Text("Source: "),
				dom.Element("a", dom.Attr{"href": info.Source}, dom.Text(info.Source)),
			), nl(),
		), nl(),
		dom.Elem("hr"), nl(),
	)
	for i, chapter := range info.Chapters {
		attr := dom.Attr{"class": "chapter", "id": fmt.Sprintf("ch%03d", i)}
		div := dom.Elem("div",
			dom.Element("h2", attr, dom.Text(chapter.Title)), nl())
		if chapter.Url != "" {
			dom.Append(div, dom.Comment(fmt.Sprintf("\n%s\n", chapter.Url)), nl())
		}
		if !chapter.Modified.IsZero() {
			dom.Append(div, dom.Elem("div", dom.Elem("em", dom.Text(chapter.Modified.Format("2006-01-02")))), nl())
		}
		dom.Append(div, nl(), dom.Elem("hr"), nl(), chapter.Content, nl(), dom.Elem("hr"), nl())
		if i+1 == len(info.Chapters) {
			dom.Append(div, dom.Elem("div", link(info.Source, info.Source)), nl(), dom.Elem("hr"), nl())
		}
		dom.Append(body, div, dom.Text("\n\n"))
	}
	description := info.Source
	if len(info.Comments) > 0 {
		description = description + "\n\n" + info.Comments
	}
	htmlNode := dom.Element("html", dom.Attr{"lang": info.Language}, nl(),
		dom.Elem("head", nl(),
			dom.Element("meta", dom.Attr{"charset": "utf-8"}), nl(),
			meta("viewport", "width=device-width, initial-scale=1.0"), nl(),
			dom.Elem("title", dom.Text(info.Title)), nl(),
			meta("DC.title", info.Title), nl(),
			meta("DC.creator.aut", info.Authors), nl(),
			meta("DC.description", description), nl(),
			meta("DC.source", info.Source), nl(),
			meta("DC.language", info.Language), nl(),
			meta("DC.date.modified", info.Modified.Format("2006-01-02")), nl(),
			dom.Elem("style", dom.Text(bookStyle)), nl(),
		),
		nl(), body, nl(),
	)
	err := dom.RenderHTML(htmlNode, dst)
	for _, chapter := range info.Chapters {
		dom.Remove(chapter.Content)
	}
	return err
}

// Print information about the book.
func (info EbookInfo) Print(dst io.Writer) {
	fmt.Fprintf(dst, "Authors:  %q\n", info.Authors)
	fmt.Fprintf(dst, "Comments: %q\n", info.Comments)
	fmt.Fprintf(dst, "Title:    %q\n", info.Title)
	fmt.Fprintf(dst, "Source:   %q\n", info.Source)
	fmt.Fprintf(dst, "Language: %q\n", info.Language)
	fmt.Fprintf(dst, "Cover:    %d bytes\n", len(info.Cover))
	fmt.Fprintf(dst, "Modified: %s\n", info.Modified.Format(time.RFC3339))
	fmt.Fprintf(dst, "Chapters: %d\n", len(info.Chapters))
	for _, ch := range info.Chapters {
		fmt.Fprintf(dst, "* Title: %q\n", ch.Title)
		fmt.Fprintf(dst, "  Url:   %q\n", ch.Url)
		fmt.Fprintf(dst, "  Text:  %d bytes\n", dom.TextBytes(ch.Content))
		fmt.Fprintf(dst, "  Mod:   %s\n", ch.Modified.Format(time.RFC3339))
	}
}

// Write the ebook as an Epub.
func (info EbookInfo) Write(dst io.Writer) error {
	var (
		uid   string = randomUUID()
		cover []byte
	)
	if len(info.Cover) > 0 {
		var err error
		cover, err = saveJpegWithScale(info.Cover, 400, 600)
		if err != nil {
			log.Printf("Cover error: %v", err)
			cover = nil
		}
	}

	zw := zipper.Make(dst)
	defer zw.Close()

	modTime := info.Modified
	if !modTime.IsZero() {
		modTime = modTime.UTC()
	}
	if w := zw.CreateStore("mimetype", time.Time{}); w != nil {
		_, zw.Error = w.Write([]byte("application/epub+zip"))
	}
	if w := zw.CreateDeflate("META-INF/container.xml", modTime); w != nil {
		_, zw.Error = w.Write([]byte(conatainer_xml))
	}
	if w := zw.CreateDeflate("book/"+"toc.ncx", modTime); w != nil {
		zw.Error = makeNCX(info, uid, w)
	}
	if w := zw.CreateDeflate("book/"+"content.opf", modTime); w != nil {
		zw.Error = makePackage(info, uid, w, len(cover) > 0)
	}
	if w := zw.CreateDeflate("book/"+"frontmatter.xhtml", modTime); w != nil {
		zw.Error = writeFrontmatter(info, w, len(cover) > 0)
	}
	if w := zw.CreateDeflate("book/"+"toc.xhtml", modTime); w != nil {
		zw.Error = writeToc(info, w)
	}
	if len(cover) > 0 {
		if w := zw.CreateStore("book/"+"cover.jpg", modTime); w != nil {
			_, zw.Error = w.Write(cover)
		}
	}
	for i, chapter := range info.Chapters {
		if w := zw.CreateDeflate(fmt.Sprintf("book/"+"%04d.xhtml", i), chapter.Modified); w != nil {
			var churl string
			if i+1 == len(info.Chapters) {
				churl = chapter.Url
			}
			zw.Error = writeChapter(chapter, churl, info.Language, w)
		}
	}
	return zw.Error
}

func writeFrontmatter(info EbookInfo, dst io.Writer, cover bool) error {
	description := dom.Elem("div")
	for _, p := range strings.Split(info.Comments, "\n\n") {
		pnode := dom.Elem("p")
		for i, c := range strings.Split(p, "\n\n") {
			if i > 0 {
				dom.Append(pnode, dom.Elem("br"))
			}
			dom.Append(pnode, dom.Text(c))
		}
		dom.Append(description, pnode)
	}
	var img *dom.Node
	if cover {
		img = dom.Element("img", dom.Attr{"src": "cover.jpg", "alt": "[COVER]"})
	}
	htmlNode := dom.Element("html", dom.Attr{"xmlns": "http://www.w3.org/1999/xhtml", "xml:lang": info.Language},
		head(info.Title, bookStyle, ""),
		dom.Elem("body",
			dom.Elem("h1", dom.Text(info.Title)),
			img,
			dom.Elem("div", dom.Text(info.Authors)),
			dom.Elem("div", dom.Text(info.Source)),
			dom.Elem("div", dom.Elem("em", dom.Text(info.Modified.Format("2006-01-02")))),
			description,
		),
	)
	return dom.RenderXHTMLDoc(htmlNode, dst)
}

func writeChapter(chapter Chapter, url, lang string, dst io.Writer) error {
	body := dom.Elem("body")
	if chapter.Url != "" {
		dom.Append(body, dom.Comment(fmt.Sprintf("\n%s\n", chapter.Url)))
	}
	dom.Append(body, dom.Element("h2", dom.Attr{"class": "chapter"}, dom.Text(chapter.Title)))
	if !chapter.Modified.IsZero() {
		dom.Append(body, dom.Elem("p", dom.Elem("em", dom.Text(chapter.Modified.Format("2006-01-02")))))
	}
	dom.Append(body, dom.Elem("hr"), chapter.Content, dom.Elem("hr"))
	if url != "" {
		dom.Append(body, dom.Elem("div", link(url, url)), dom.Elem("hr"))
	}
	htmlNode := dom.Element("html",
		dom.Attr{"xmlns": "http://www.w3.org/1999/xhtml", "xml:lang": lang},
		head(chapter.Title, bookStyle, ""),
		body,
	)
	return dom.RenderXHTMLDoc(htmlNode, dst)
}

func writeToc(info EbookInfo, dst io.Writer) error {
	links := dom.Element("ol", dom.Attr{"class": "flat"})
	for i, ch := range info.Chapters {
		label := fmt.Sprintf("%d. %s", i+1, ch.Title)
		dom.Append(links, dom.Elem("li", link(fmt.Sprintf("%04d.xhtml", i), label)))
	}
	htmlNode := dom.Element("html",
		dom.Attr{
			"xmlns":      "http://www.w3.org/1999/xhtml",
			"xml:lang":   info.Language,
			"xmlns:epub": "http://www.idpf.org/2007/ops",
		},
		head(info.Title, bookStyle, ""),
		dom.Elem("body",
			dom.Element("nav",
				dom.Attr{"epub:type": "toc", "style": "display:none;"},
				dom.Elem("h2", dom.Text("Contents")),
				links),
		),
	)
	return dom.RenderXHTMLDoc(htmlNode, dst)
}

func link(url, text string) *Node {
	if url == "" {
		return nil
	}
	return dom.Element("a", dom.Attr{"href": url}, dom.Text(text))
}

func imgElem(url, alt string) *Node {
	if url == "" {
		return nil
	}
	return dom.Element("img", dom.Attr{"src": url, "alt": alt})
}

func randomUUID() string {
	var v [16]byte
	rand.Read(v[:])
	return fmt.Sprintf("%x-%x-%x-%x-%x", v[0:4], v[4:6], v[6:8], v[8:10], v[10:16])
}
