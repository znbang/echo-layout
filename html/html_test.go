package html

import (
	"bytes"
	"io/fs"
	"testing"
	"testing/fstest"
)

func getFileSystem() fs.FS {
	return fstest.MapFS{
		"views/layout/main.html": {
			Data: []byte("layout content: {{block \"index\" .}}main{{end}}"),
		},
		"views/index.html": {
			Data: []byte("{{define \"index\"}}index{{end}}"),
		},
		"views/partial.html": {
			Data: []byte("partial"),
		},
	}
}

func TestEngine_Render_index(t *testing.T) {
	var buf bytes.Buffer

	e := NewFileSystem(getFileSystem(), "views", "layout/main", ".html")
	err := e.Render(&buf, "index", nil, nil)
	if err != nil {
		t.Error(e)
	}

	expect := "layout content: index"
	result := buf.String()
	if expect != result {
		t.Errorf("unexpected result: %v, exptected: %v", result, expect)
	}
}

func TestEngine_Render_partial(t *testing.T) {
	var buf bytes.Buffer

	e := NewFileSystem(getFileSystem(), "views", "layout/main", ".html")
	err := e.Render(&buf, "partial.html", nil, nil)
	if err != nil {
		t.Error(e)
	}

	expect := "partial"
	result := buf.String()
	if expect != result {
		t.Errorf("unexpected result: %v, exptected: %v", result, expect)
	}
}
