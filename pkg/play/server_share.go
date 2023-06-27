// Copyright 2023 The KCL Authors. All rights reserved.

package play

import (
	"bytes"
	"crypto/rand"
	"fmt"
	"io"
	"log"
	"net/http"

	"github.com/boltdb/bolt"
	"github.com/gin-gonic/gin"
)

func (p *WebServer) initShareHandler() {
	if !p.opts.AllowShare {
		return
	}

	p.router.POST("/-/play/share", func(c *gin.Context) {
		p.shareHandler(c.Writer, c.Request)
	})

	var err error
	p.db, err = bolt.Open(p.opts.DatabaseFile, 0600, nil)
	if err != nil {
		log.Fatal(err)
	}

	if err = p.createBucket(p.bucketSnippets); err != nil {
		log.Fatal(err)
	}

	if err = p.createBucket(p.bucketConfig); err != nil {
		log.Fatal(err)
	}

	if err = p.createBucket(p.bucketCache); err != nil {
		log.Fatal(err)
	}

	err = p.db.Update(func(tx *bolt.Tx) error {
		var err error
		b := tx.Bucket(p.bucketConfig)
		p.salt = b.Get([]byte("salt"))
		if p.salt == nil {
			p.salt = make([]byte, 30)
			if _, err = rand.Read(p.salt); err != nil {
				return err
			}
			b.Put([]byte("salt"), p.salt)
		}
		return nil
	})

	if err != nil {
		log.Fatal(err)
	}

}

func (p *WebServer) createBucket(name []byte) error {
	return p.db.Update(func(tx *bolt.Tx) error {
		_, err := tx.CreateBucketIfNotExists(name)
		return err
	})
}

func (p *WebServer) shareHandler(w http.ResponseWriter, r *http.Request) {
	if !p.opts.AllowShare || r.Method != "POST" {
		http.Error(w, "Forbidden", http.StatusForbidden)
		return
	}

	var body bytes.Buffer
	_, err := io.Copy(&body, io.LimitReader(r.Body, maxSnippetSize+1))
	r.Body.Close()
	if err != nil {
		log.Printf("Error reading body: %q", err)
		http.Error(w, "Server Error", http.StatusInternalServerError)
		return
	}
	if body.Len() > maxSnippetSize {
		http.Error(w, "Snippet is too large", http.StatusRequestEntityTooLarge)
		return
	}

	snip := &Snippet{Body: body.Bytes()}
	id := snip.Id(p.salt)
	key := []byte(id)

	err = p.db.Update(func(tx *bolt.Tx) error {
		b := tx.Bucket(p.bucketSnippets)
		return b.Put(key, snip.Body)
	})

	if err != nil {
		log.Printf("Error putting snippet: %q", err)
		http.Error(w, "Server Error", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Access-Control-Allow-Methods", "POST")
	w.Header().Set("Access-Control-Allow-Origin", p.opts.AllowOriginHeader)

	fmt.Fprint(w, id)
}
