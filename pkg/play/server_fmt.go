// Copyright 2021 The KCL Authors. All rights reserved.

package play

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"

	"kusionstack.io/kclvm-go/pkg/tools/format"
)

func (p *WebServer) initFmtHandler() {
	p.router.POST("/-/play/fmt", func(c *gin.Context) {
		p.fmtHandler(c.Writer, c.Request)
	})
}

func (p *WebServer) fmtHandler(w http.ResponseWriter, r *http.Request) {

	resp, err := p.fmtCode([]byte(r.FormValue("body")))
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	if err := json.NewEncoder(w).Encode(resp); err != nil {
		http.Error(w, fmt.Sprintf("error encoding response: %v", err), http.StatusInternalServerError)
		return
	}
}

func (p *WebServer) fmtCode(code []byte) (*fmtResponse, error) {
	output, err := format.FormatCode(code)
	if err != nil {
		resp := &fmtResponse{
			Error: err.Error(),
		}
		return resp, nil
	}

	resp := &fmtResponse{
		Body: string(output),
	}

	return resp, nil
}
