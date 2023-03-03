// Copyright 2021 The KCL Authors. All rights reserved.

// KCL playground.
package main

import (
	"flag"
	"fmt"
	"os/exec"
	"runtime"
	"strings"
	"time"

	"kusionstack.io/kcl-playground/pkg/play"
)

func main() {
	deploy_mode := flag.Bool("deploy", false, "Whether is deploy mode")
	flag.Parse()
	opt := play.Options{
		PlayMode:   true,
		AllowShare: false,
	}
	if *deploy_mode {
		play.Run(":80", &opt)
	} else {
		addr := "localhost:2023"
		fmt.Printf("listen at http://%s\n", addr)
		go func() {
			time.Sleep(time.Second * 2)
			openBrowser(addr)
		}()
		play.Run(addr, &opt)
	}
}

func openBrowser(url string) error {
	if !strings.HasPrefix(url, "http") {
		url = "http://" + url
	}
	switch runtime.GOOS {
	case "linux":
		return exec.Command("xdg-open", url).Start()
	case "windows":
		return exec.Command("rundll32", "url.dll,FileProtocolHandler", url).Start()
	case "darwin":
		return exec.Command("open", url).Start()
	default:
		return fmt.Errorf("unsupported platform")
	}
}
