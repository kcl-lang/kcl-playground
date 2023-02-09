// Copyright 2021 The KCL Authors. All rights reserved.

// KCL playground.
package main

import (
	"fmt"

	"kusionstack.io/kcl-playground/pkg/play"
)

func main() {
	addr := ":2022"
	fmt.Printf("listen at http://%s\n", addr)

	play.Run(addr, &play.Option{
		PlayMode: true,
	})
}
