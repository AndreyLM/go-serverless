package main

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"log"
	"os/exec"
)

func run(prog string, args ...string) ([]byte, error) {
	cmd := exec.Command(prog, args...)
	
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

// RoleList - role lise
type RoleList struct {
	Roles []Role
}

// Role - role
type Role struct {
	RoleName string
	Arn      string
}

// RoleMap - get aws roles
func RoleMap() (map[string]string, error) {
	res := make(map[string]string)
	data, err := run("aws", "iam", "list-roles")
	if err != nil {
		return res, err
	}

	var rList RoleList
	if err = json.Unmarshal(data, &rList); err != nil {
		return res, err
	}

	for _, v := range rList.Roles {
		res[v.RoleName] = v.Arn
	}

	return res, nil
}

func main() {
	rm, err := RoleMap()
	if err != nil {
		log.Fatal(err)
	}
	for k, v := range rm {
		fmt.Println(k + "--" + v)
	}
}
