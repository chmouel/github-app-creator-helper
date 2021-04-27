package main

import (
	"context"
	"fmt"
	"html/template"
	"io"
	"log"
	"net/http"
	"path/filepath"

	"github.com/google/go-github/v35/github"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

// TemplateRenderer is a custom html/template renderer for Echo framework
type TemplateRenderer struct {
	templates *template.Template
}

func (t *TemplateRenderer) Render(w io.Writer, name string, data interface{}, c echo.Context) error {
	return t.templates.ExecuteTemplate(w, name, data)
}

func index(c echo.Context) error {
	return c.Render(http.StatusOK, "index.html", map[string]interface{}{})
}

func redirect_first(c echo.Context) error {
	fp, err := c.FormParams()
	if err != nil {
		return err
	}
	fmt.Printf("%+v\n", fp)
	return c.String(http.StatusOK, "OK")
}

func redirect(c echo.Context) error {
	code := c.QueryParam("code")
	ghclient := github.NewClient(nil)
	ctx := context.Background()
	ret, _, err := ghclient.Apps.CompleteAppManifest(ctx, code)
	if err != nil {
		return err
	}
	fmt.Printf("%+v\n", ret)
	ret.GetClientSecret()
	return c.String(http.StatusOK, "OK")
}

func Server() (err error) {
	e := echo.New()
	e.Renderer = &TemplateRenderer{
		templates: template.Must(template.ParseGlob(filepath.Join("html", "*.html"))),
	}
	e.Debug = true
	e.Use(middleware.LoggerWithConfig(middleware.LoggerConfig{}))
	e.GET("/redirect", redirect)
	e.GET("/", index)

	return (e.Start(
		fmt.Sprintf("localhost:3000")))
}

func main() {
	err := Server()
	if err != nil {
		log.Fatal(err)
	}
}
