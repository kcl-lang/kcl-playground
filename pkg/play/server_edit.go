// Copyright 2021 The KCL Authors. All rights reserved.

package play

import (
	_ "embed"
	"html/template"
	"net/http"

	"github.com/gin-gonic/gin"
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
	snip := &Snippet{Body: []byte(edit_helloPlayground)}
	edit_Template.Execute(w, &editData{snip, p.opt.AllowShare})
}

//go:embed _edit.tmpl.html
var edit_tmpl string

var edit_Template = template.Must(template.New("playground/index.html").Parse(edit_tmpl))

const edit_helloPlayground = `import kcl_plugin.hello

name = "kcl"
age = 1
two = hello.add(1, 1)

schema Person:
    name: str = "kcl"
    age: int = 1

x0 = Person{}
x1 = Person{age=101}
`
