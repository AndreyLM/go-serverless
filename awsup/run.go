package main

import (
	"log"
	"bytes"
	"errors"
	"io"
	"os/exec"
)

func run(prog string, args ...string) ([]byte, error) {
	cmd := exec.Command(prog, args...)
	log.Println(cmd.Process)
	outPipe, err := cmd.StdoutPipe()
	if err != nil {
		return []byte{}, err
	}

	errPipe, err := cmd.StderrPipe()
	if err != nil {
		return []byte{}, err
	}

	err = cmd.Start()
	if err != nil {
		return []byte{}, err
	}

	var outBuf bytes.Buffer
	var errBuf bytes.Buffer

	ch := make(chan bool)
	go func() {
		io.Copy(&outBuf, outPipe)
		ch <- true
	}()
	io.Copy(&errBuf, errPipe)
	<-ch

	err = cmd.Wait()
	if err != nil {
		return []byte{}, err
	}

	if len(errBuf.Bytes()) != 0 {
		return outBuf.Bytes(), errors.New(errBuf.String())
	}

	return outBuf.Bytes(), nil
}
