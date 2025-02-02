package main

import (
	"bytes"
	"io/ioutil"
	"path/filepath"
	"strings"
	"testing"

	"github.com/mmarkdown/mmark/v2/lang"
	"github.com/mmarkdown/mmark/v2/mparser"
	"github.com/mmarkdown/mmark/v2/render/xml"

	"github.com/gomarkdown/markdown"
	"github.com/gomarkdown/markdown/parser"
)

func TestMmarkXML(t *testing.T) {
	dir := "testdata"
	testFiles, err := ioutil.ReadDir(dir)
	if err != nil {
		t.Fatalf("could not read %s: %q", dir, err)
	}
	for _, f := range testFiles {
		if f.IsDir() {
			continue
		}

		if filepath.Ext(f.Name()) != ".md" {
			continue
		}

		base := f.Name()[:len(f.Name())-3]
		opts := xml.RendererOptions{
			Flags:    xml.CommonFlags | xml.XMLFragment,
			Comments: [][]byte{[]byte("//"), []byte("#")},
		}
		// if the file name has a prefix ending in a underscore that prefix is taken is the language
		// for this particular file and used.
		us := strings.Index(f.Name(), "_")
		l := "en"
		if us >= 0 {
			l = f.Name()[:us]
		}
		opts.Language = lang.New(l)

		renderer := xml.NewRenderer(opts)

		doTest(t, dir, base, renderer)
	}
}

func doTest(t *testing.T, dir, basename string, renderer markdown.Renderer) {
	filename := filepath.Join(dir, basename+".md")
	input, err := ioutil.ReadFile(filename)
	if err != nil {
		t.Errorf("couldn't open '%s', error: %v\n", filename, err)
		return
	}

	filename = filepath.Join(dir, basename+".xml")
	expected, err := ioutil.ReadFile(filename)
	if err != nil {
		t.Errorf("couldn't open '%s', error: %v\n", filename, err)
	}
	expected = bytes.TrimSpace(expected)

	p := parser.NewWithExtensions(mparser.Extensions)

	init := mparser.NewInitial(filename)
	p.Opts = parser.Options{
		ParserHook:    mparser.TitleHook,
		ReadIncludeFn: init.ReadInclude,
	}

	actual := markdown.ToHTML(input, p, renderer)
	actual = bytes.TrimSpace(actual)
	if bytes.Compare(actual, expected) != 0 {
		t.Errorf("\n    [%#v]\nExpected[%s]\nActual  [%s]",
			basename+".md", expected, actual)
	}
}
