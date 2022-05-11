// Copyright 2021 The KCL Authors. All rights reserved.

package play

func Run(addr string) error {
	s := NewWebServer(nil)
	return s.Run(addr)
}

func RunPlayground(addr string) error {
	s := NewWebServer(&Option{
		PlayMode:   true,
		AllowShare: true,
	})
	return s.Run(addr)
}
