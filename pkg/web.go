package pkg

import (
	"fmt"
	"html/template"
	"io"
	"net/http"
	"path/filepath"

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
	client, err := NewClient()
	if err != nil {
		return err
	}
	route, err := client.GetFirstTLSRoute()
	if err != nil {
		return err
	}

	return c.Render(http.StatusOK, "index.html", map[string]interface{}{
		"webhook_route_host": route.Spec.Host,
	})
}

func redirect(c echo.Context) error {
	code := c.QueryParam("code")
	return c.String(http.StatusOK, "OK: "+code)
	// ghclient := github.NewClient(nil)
	// ctx := context.Background()
	// ret, _, err := ghclient.Apps.CompleteAppManifest(ctx, code)
	// if err != nil {
	//	return err
	// }
	// ret.GetClientSecret()
	// return c.String(http.StatusOK, "OK")
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
