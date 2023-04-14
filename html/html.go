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
	fileSystem fs.FS
	directory  string
	extension  string
	layout     string
	mutex      sync.RWMutex
	templates  map[string]*template.Template
}

func (e *Engine) Render(w io.Writer, name string, data interface{}, c echo.Context) error {
	useLayout := true
	if filepath.Ext(name) == e.extension {
		useLayout = false
		name = strings.TrimSuffix(name, e.extension)
	}

	e.mutex.Lock()
	tpl, ok := e.templates[name]
	e.mutex.Unlock()

	if !ok {
		tplList := make([]string, 0)
		if useLayout && e.layout != "" {
			tplList = append(tplList, e.layout)
		}
		tplList = append(tplList, name)

		tpl = template.New(name)
		for _, v := range tplList {
			var temp *template.Template
			if v == name {
				temp = tpl
			} else {
				temp = tpl.New(v)
			}
			content, err := e.load(v)
			if err != nil {
				return err
			}
			_, err = temp.Parse(content)
			if err != nil {
				return fmt.Errorf("parse %v failed: %v", v, e)
			}
		}

		e.mutex.Lock()
		e.templates[name] = tpl
		e.mutex.Unlock()
	}

	if useLayout && e.layout != "" {
		return tpl.ExecuteTemplate(w, e.layout, data)
	} else {
		return tpl.ExecuteTemplate(w, name, data)
	}
}

func (e *Engine) load(name string) (string, error) {
	file, err := e.fileSystem.Open(name + e.extension)
	if err != nil {
		return "", fmt.Errorf("open %v failed: %v", name, err)
	}
	content, err := io.ReadAll(file)
	if err != nil {
		return "", fmt.Errorf("read %v failed: %v", name, err)
	}
	return string(content), err
}
