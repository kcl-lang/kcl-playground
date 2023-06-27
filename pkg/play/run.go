// Copyright 2023 The KCL Authors. All rights reserved.

package play

func Run(addr string, opts *Options) error {
	s := NewWebServer(opts)
	return s.Run(addr)
}
