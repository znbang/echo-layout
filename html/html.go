package html

import (
	"fmt"
	"html/template"
	"io"
	"io/fs"
	"os"
	"path/filepath"
	"strings"
	"sync"

	"github.com/labstack/echo/v4"
)

func New(directory, layout, extension string) *Engine {
	return &Engine{
		fileSystem: os.DirFS(directory),
		directory:  directory,
		extension:  extension,
		layout:     layout,
		templates:  make(map[string]*template.Template),
	}
}

func NewFileSystem(fileSystem fs.FS, directory, layout, extension string) *Engine {
	subFS, err := fs.Sub(fileSystem, directory)
	if err != nil {
		panic(err)
	}
	return &Engine{
		fileSystem: subFS,
		directory:  directory,
		extension:  extension,
		layout:     layout,
		templates:  make(map[string]*template.Template),
	}
}

type Engine struct {
	loaded     bool
	fileSystem fs.FS
	directory  string
	extension  string
	layout     string
	mutex      sync.RWMutex
	templates  map[string]*template.Template
}

func (e *Engine) load() error {
	if e.loaded {
		return nil
	}

	e.mutex.Lock()
	defer e.mutex.Unlock()

	var layoutPath string
	if e.layout != "" {
		layoutPath = e.layout + e.extension
	}

	walkFunc := func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			return err
		}

		if d.IsDir() {
			return nil
		}

		if ext := filepath.Ext(path); ext != e.extension {
			return nil
		}

		if layoutPath != "" && strings.HasSuffix(path, layoutPath) {
			return nil
		}

		var tmpl *template.Template
		if layoutPath != "" {
			tmpl, err = template.ParseFS(e.fileSystem, layoutPath, path)
		} else {
			tmpl, err = template.ParseFS(e.fileSystem, path)
		}

		if err == nil {
			name := strings.TrimSuffix(path, e.extension)
			e.templates[name] = tmpl
		}

		return err
	}

	e.loaded = true

	return fs.WalkDir(e.fileSystem, ".", walkFunc)
}

func (e *Engine) Render(w io.Writer, name string, data interface{}, c echo.Context) error {
	if !e.loaded {
		if err := e.load(); err != nil {
			return err
		}
	}

	if tmpl := e.templates[name]; tmpl == nil {
		return fmt.Errorf("template not found: %v", name)
	} else {
		return tmpl.Execute(w, data)
	}
}
