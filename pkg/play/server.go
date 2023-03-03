// Copyright 2021 The KCL Authors. All rights reserved.

package play

import (
	"embed"
	"io/fs"
	"net/http"
	"strings"

	"github.com/boltdb/bolt"
	"github.com/gin-gonic/gin"
)

const maxSnippetSize = 64 * 1024

type WebServer struct {
	opts *Options

	fs     fs.FS
	static *gin.Engine
	router *gin.Engine

	db             *bolt.DB
	bucketSnippets []byte
	bucketCache    []byte
	bucketConfig   []byte
	salt           []byte
}

//go:embed _static
var fs_static embed.FS

func getStaticFS() fs.FS {
	fs, err := fs.Sub(fs_static, "_static")
	if err != nil {
		panic(err)
	}
	return fs
}

func NewWebServer(opts *Options) *WebServer {
	if opts == nil {
		opts = &Options{}
	}

	if opts.AllowShare && opts.DatabaseFile == "" {
		opts.DatabaseFile = "kcl-play.db"
	}

	p := &WebServer{
		opts:   opts,
		fs:     getStaticFS(),
		static: gin.Default(),
		router: gin.Default(),

		bucketSnippets: []byte("snippets"),
		bucketCache:    []byte("cache"),
		bucketConfig:   []byte("config"),
		salt:           []byte{},
	}

	p.static.StaticFS("/static", http.FS(p.fs))

	p.initCompileHandler()
	p.initEditHandler()
	p.initFmtHandler()
	p.initShareHandler()

	return p
}

func (p *WebServer) Run(addr string) error {
	return http.ListenAndServe(addr,
		http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			if p.opts.PlayMode {
				if r.URL.Path == "/" || r.URL.Path == "/index.html" {
					p.router.ServeHTTP(w, r)
					return
				}
			}

			switch {
			case strings.HasPrefix(r.URL.Path, "/-/"):
				p.router.ServeHTTP(w, r)
			case r.URL.Path == "/":
				p.router.ServeHTTP(w, r)
			case r.URL.Path == "/play.html":
				p.router.ServeHTTP(w, r)
			default:
				p.static.ServeHTTP(w, r)
			}
		}),
	)
}
