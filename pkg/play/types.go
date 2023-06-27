// Copyright 2023 The KCL Authors. All rights reserved.

package play

import (
	"crypto/sha1"
	"fmt"
	"time"
)

type Options struct {
	PlayMode          bool
	AllowShare        bool
	AllowOriginHeader string
	DisableCache      bool
	CompileURL        string
	DatabaseFile      string
}

type Snippet struct {
	Body []byte
}

type Request struct {
	Body string
}

type Response struct {
	Errors string
	Events []Event
}

type Event struct {
	Message string
	Kind    string        // "stdout" or "stderr"
	Delay   time.Duration // time to wait before printing Message
}

type editData struct {
	Snippet *Snippet
	Share   bool
	Version string
}

type fmtResponse struct {
	Body  string
	Error string
}

func (s *Snippet) Id(salt []byte) string {
	h := sha1.New()
	h.Write(salt)
	h.Write(s.Body)
	sum := h.Sum(nil)
	return fmt.Sprintf("%x", sum)
}
