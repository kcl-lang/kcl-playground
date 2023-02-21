// Copyright 2021 The KCL Authors. All rights reserved.

package play

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"

	"github.com/gin-gonic/gin"

	"kusionstack.io/kclvm-go"
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

	result, err := kclvm.Run("prog.k", kclvm.WithCode(req.Body))
	if err != nil {
		resp := &Response{Errors: err.Error()}
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
