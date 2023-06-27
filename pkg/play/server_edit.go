// Copyright 2023 The KCL Authors. All rights reserved.

package play

import (
	_ "embed"
	"html/template"
	"net/http"

	"github.com/boltdb/bolt"
	"github.com/gin-gonic/gin"
	"kusionstack.io/kclvm-go"
)

func (p *WebServer) initEditHandler() {
	// play
	p.router.GET("/", func(c *gin.Context) {
		http.Redirect(c.Writer, c.Request, "/-/play/index.html", http.StatusSeeOther)
	})

	p.router.GET("/play.html", func(c *gin.Context) {
		http.Redirect(c.Writer, c.Request, "/-/play/index.html", http.StatusSeeOther)
	})
	p.router.GET("/play/index.html", func(c *gin.Context) {
		http.Redirect(c.Writer, c.Request, "/-/play/index.html", http.StatusSeeOther)
	})
	p.router.GET("/-/play/", func(c *gin.Context) {
		http.Redirect(c.Writer, c.Request, "/-/play/index.html", http.StatusSeeOther)
	})

	p.router.GET("/-/play/index.html", func(c *gin.Context) {
		p.EditHandler(c.Writer, c.Request)
	})
}
func (p *WebServer) EditHandler(w http.ResponseWriter, r *http.Request) {
	editTemplate.Execute(w, &editData{p.getSnippet(w, r), p.opts.AllowShare, string(kclvm.KclvmAbiVersion)})
}

func (p *WebServer) getSnippet(w http.ResponseWriter, r *http.Request) *Snippet {
	var snip *Snippet
	// Retrieve via the parameter id
	if r.Method == "GET" {
		p.db.View(func(tx *bolt.Tx) error {
			data := tx.Bucket(p.bucketSnippets).Get([]byte(r.FormValue("id")))
			if data != nil {
				snip = &Snippet{Body: data}
			}
			return nil
		})
	}
	if snip == nil {
		snip = &defaultSnippet
	}
	return snip
}

//go:embed _edit.tmpl.html
var editTemplateString string

var editTemplate = template.Must(template.New("playground/index.html").Parse(editTemplateString))

const defaultCode = `apiVersion = "apps/v1"
kind = "Deployment"
metadata = {
    name = "nginx"
    labels.app = "nginx"
}
spec = {
    replicas = 3
    selector.matchLabels = metadata.labels
    template.metadata.labels = metadata.labels
    template.spec.containers = [
        {
            name = metadata.name
            image = "${metadata.name}:1.14.2"
            ports = [{ containerPort = 80 }]
        }
    ]
}
`

var defaultSnippet = Snippet{Body: []byte(defaultCode)}
