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

	"kcl-lang.io/kcl-playground/pkg/play"
)

func main() {
	deploy_mode := flag.Bool("deploy", false, "Whether is deploy mode")
	flag.Parse()
	opts := play.Options{
		PlayMode:   true,
		AllowShare: true,
	}
	if *deploy_mode {
		play.Run(":80", &opts)
	} else {
		addr := "localhost:2023"
		fmt.Printf("listen at http://%s\n", addr)
		go func() {
			time.Sleep(time.Second * 2)
			openBrowser(addr)
		}()
		play.Run(addr, &opts)
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
