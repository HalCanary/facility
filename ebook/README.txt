package ebook // import "github.com/HalCanary/facility/ebook"

Copyright 2022 Hal Canary Use of this program is governed by the file LICENSE.

VARIABLES

var UnsupportedUrlError = errors.New("unsupported url")
    Returned by a EbookGeneratorFunction when the URL can not be handled.


FUNCTIONS

func ConvertToEbook(src, dst string, arguments ...string) error
func PopulateInfo(info *EbookInfo, doc *dom.Node)
func RegisterEbookGenerator(downloadFunction EbookGeneratorFunction)
    Register the given function.


TYPES

type Chapter struct {
	Title    string
	Url      string
	Content  *Node
	Modified time.Time
}

type EbookGeneratorFunction func(url string, doPopulate bool) (EbookInfo, error)

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
    Ebook content and metadata.

func DownloadEbook(url string, doPopulate bool) (EbookInfo, error)
    Return the result of the first registered download function that does not
    return UnsupportedUrlError. @param url - the URL of the title page of
    the book. @param doPopulate - if true, download and populate the entire
    EbookInfo data structure, not just its metadata.

func (info EbookInfo) CalculateLastModified() time.Time
    Return the time of most recently modified chapter.

func (info EbookInfo) Print(dst io.Writer)

func (info EbookInfo) Write(dst io.Writer) error
    Write the ebook as an Epub.

func (info EbookInfo) WriteHtml(dst io.Writer) error

type Node = dom.Node

func Cleanup(node *Node) *Node
    Clean up a HTML fragment.

func ResolveLinks(node *Node, ref *url.URL) *Node

