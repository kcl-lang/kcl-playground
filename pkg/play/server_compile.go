// Copyright 2023 The KCL Authors. All rights reserved.

package play

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"path/filepath"

	"github.com/acarl005/stripansi"
	"github.com/gin-gonic/gin"

	"kcl-lang.io/kcl-go"
)

func (p *WebServer) initCompileHandler() {
	p.router.POST("/-/play/compile", func(c *gin.Context) {
		p.compileHandler(c.Writer, c.Request)
	})
}

func (p *WebServer) compileHandler(w http.ResponseWriter, r *http.Request) {
	var req Request
	version := r.PostFormValue("version")
	if version == "2" {
		req.Body = r.PostFormValue("body")
	} else {
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			http.Error(w, fmt.Sprintf("error decoding request: %v", err), http.StatusBadRequest)
			return
		}
	}
	resp, err := p.compileAndRun(&req)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	if err := json.NewEncoder(w).Encode(resp); err != nil {
		http.Error(w, fmt.Sprintf("error encoding response: %v", err), http.StatusInternalServerError)
		return
	}
}

func (p *WebServer) compileAndRun(req *Request) (*Response, error) {
	tmpDir, err := ioutil.TempDir("", "sandbox")
	if err != nil {
		return nil, fmt.Errorf("error creating temp directory: %v", err)
	}
	defer os.RemoveAll(tmpDir)
	// Write kcl code in the temp file.
	kFile := filepath.Join(tmpDir, "prog.k")
	err = os.WriteFile(kFile, []byte(req.Body), 0666)
	if err != nil {
		return nil, err
	}

	result, err := kclvm.Run(kFile)
	if err != nil {
		resp := &Response{Errors: stripansi.Strip(err.Error())}
		return resp, nil
	}

	resp := &Response{
		Events: []Event{
			{
				Message: result.GetRawYamlResult(),
				Kind:    "stdout",
			},
		},
	}

	return resp, nil
}
