// Copyright 2021 The KCL Authors. All rights reserved.

package play

import (
	"fmt"
	"os/exec"
	"strings"
)

var _ = fmt.Sprint

type KclvmRuntime interface {
	Fmt(workDir, kFile string) (output []byte, err error)
	Run(workDir string, kclFiles ...string) (output []byte, err error)
}

func DefaultKclvmRuntime() KclvmRuntime {
	return new(cmdKclvmRuntime)
}

type cmdKclvmRuntime struct{}

func (p *cmdKclvmRuntime) Fmt(workDir, kFile string) (output []byte, err error) {
	cmd := exec.Command("kcl", "--fmt", kFile, "--fmt-output")
	cmd.Dir = workDir

	output, err = cmd.CombinedOutput()
	if err != nil {
		if s := string(output); s != "" {
			s = strings.Replace(s, workDir, "", -1)
			return nil, fmt.Errorf("%s", s)
		}
		return nil, err
	}
	return
}

func (p *cmdKclvmRuntime) Run(workDir string, kclFiles ...string) (output []byte, err error) {
	cmd := exec.Command("kcl", kclFiles...)
	cmd.Dir = workDir

	output, err = cmd.CombinedOutput()
	if err != nil {
		if s := string(output); s != "" {
			s = strings.Replace(s, workDir, "", -1)
			return nil, fmt.Errorf("%s", s)
		}
		return nil, err
	}
	return
}
