// Copyright 2021 The KCL Authors. All rights reserved.

package play

func Run(addr string, opt *Option) error {
	s := NewWebServer(opt)
	return s.Run(addr)
}
