package dom

// Copyright 2022 Hal Canary
// Use of this program is governed by the file LICENSE.

import (
	"bytes"
	"testing"

	"github.com/HalCanary/facility/expect"
)

func maketest() *Node {
	return Element("html", Attr{"lang": "en"},
		Elem("head",
			Element("meta", Attr{
				"http-equiv": "Content-Type", "content": "text/html; charset=utf-8"}),
			Element("meta", Attr{
				"name": "viewport", "content": "width=device-width, initial-scale=1.0"}),
			Elem("title", Text("TITLE")),
		),
		Elem("body", Element("p", Attr{"id": "foo"}, Text("hi"))),
	)
}

const expected = `<!DOCTYPE html>
<html lang="en"><head><meta content="text/html; charset=utf-8" http-equiv="Content-Type"/><meta content="width=device-width, initial-scale=1.0" name="viewport"/><title>TITLE</title></head><body><p id="foo">hi</p></body></html>
`

const expected2 = `<!DOCTYPE html>
<html lang="en"><head><meta content="text/html; charset=utf-8" http-equiv="Content-Type"><meta content="width=device-width, initial-scale=1.0" name="viewport"><title>TITLE</title></head><body><p id="foo">hi</p></body></html>
`

func TestDom(t *testing.T) {
	var b bytes.Buffer
	x := maketest()
	RenderHTML(x, &b)
	expect.Equal(t, expected, b.String())

	b = bytes.Buffer{}
	x = maketest()
	RenderHTMLExperimental(x, &b)
	expect.Equal(t, expected2, b.String())

	fnd := FindNodeByAttribute(x, "name", "viewport")
	expect.True(t, fnd != nil && fnd.Data == "meta")

	fnd = FindNodeById(x, "foo")
	expect.True(t, fnd != nil && fnd.Data == "p")

	fnd = FindNodeByTag(x, "title")
	expect.True(t, fnd != nil)

	fnd = FindNodeByTagAndAttrib(x, "p", "id", "foo")
	expect.True(t, fnd != nil && fnd.Data == "p")

	result := FindNodesByTagAndAttrib(x, "meta", "", "")
	expect.Equal(t, 2, len(result))
}
