This package is a layout template engine for [echo](https://github.com/labstack/echo).
### Installation
```
go get -u github.com/znbang/echo-layout
```

### Example
#### views/layouts/main.html
```html
<!DOCTYPE html>
<html lang="en">
<head>
    <meta charset="UTF-8">
    <title>{{block "title" .}}Layout{{end}}</title>
</head>
<body>
{{block "content" .}}
This is layout.html.
{{end}}
</body>
</html>
```
#### views/index.html
```html
<!DOCTYPE html>
<html lang="en">
<head>
    <meta charset="UTF-8">
    <title>{{define "title"}}Index{{end}}</title>
</head>
<body>
{{define "content"}}
This is index.html.
{{end}}
</body>
</html>
```

#### main.go
```go
package main

import (
	"embed"
	"net/http"

	"github.com/labstack/echo/v4"
	"github.com/znbang/echo-layout/html"
)

//go:embed views/*
var viewsFS embed.FS

func main() {
	e := echo.New()
	e.Renderer = html.NewFileSystem(viewsFS, "views", "layouts/main", ".html")
	e.GET("/", func(c echo.Context) error {
		return c.Render(http.StatusOK, "index", echo.Map{})
	})
	e.Logger.Fatal(e.Start(":9000"))
}
```

#### Generated HTML
```html
<!DOCTYPE html>
<html lang="en">
<head>
    <meta charset="UTF-8">
    <title>Index</title>
</head>
<body>

This is index.html.

</body>
</html>
```