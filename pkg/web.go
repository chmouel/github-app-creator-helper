package pkg

import (
	"context"
	"fmt"
	"html/template"
	"io"
	"net/http"
	"path/filepath"

	"github.com/google/go-github/v35/github"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

const secretName = "github-app-secret"

// TemplateRenderer is a custom html/template renderer for Echo framework
type TemplateRenderer struct {
	templates *template.Template
}

func (t *TemplateRenderer) Render(w io.Writer, name string, data interface{}, c echo.Context) error {
	return t.templates.ExecuteTemplate(w, name, data)
}

func index(c echo.Context) error {
	var warning string

	client, err := NewClient()
	if err != nil {
		return err
	}
	route, err := client.GetFirstTLSRoute()
	if err != nil {
		return err
	}

	if sec, _ := client.kubeClient.CoreV1().Secrets(client.Namespace).Get(secretName, metav1.GetOptions{}); sec.Name != "" {
		warning = fmt.Sprintf("We already have a secret called '%s' in namespace '%s', you need to delete it before you click on Submit or it would not work.\n", secretName, client.Namespace)
	}

	return c.Render(http.StatusOK, "index.html", map[string]interface{}{
		"webhook_route_host": route.Spec.Host,
		"warning":            warning,
	})
}

func redirect(c echo.Context) error {
	code := c.QueryParam("code")
	ghclient := github.NewClient(nil)
	ctx := context.Background()
	manifest, _, err := ghclient.Apps.CompleteAppManifest(ctx, code)
	if err != nil {
		return err
	}
	client, err := NewClient()
	if err != nil {
		return err
	}
	client.kubeClient.CoreV1().Secrets(client.Namespace).Create(&corev1.Secret{
		ObjectMeta: metav1.ObjectMeta{
			Name: secretName,
		},
		Data: map[string][]byte{
			"application_id": []byte(fmt.Sprintf("%d", manifest.GetID())),
			"private.key":    []byte(manifest.GetPEM()),
			"webhook.secret": []byte(manifest.GetWebhookSecret()),
		},
	})
	return c.Redirect(301, manifest.GetHTMLURL())
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
