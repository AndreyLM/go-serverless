package main

import (
	"encoding/json"
	"errors"
	"flag"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"os"
	"os/exec"
	"path"
	"strings"
)

// LProject - project
type LProject struct {
	Name   string
	Bucket string
	Role   string
	path   string
}

// NewLProject -
func NewLProject(fname string) (LProject, error) {
	project := LProject{}
	data, err := ioutil.ReadFile(fname)

	if err != nil {
		return project, err
	}

	if err = json.Unmarshal(data, &project); err != nil {
		return project, err
	}

	project.path = path.Dir(fname)

	if strings.HasPrefix(project.Role, "arn:") {
		return project, nil
	}

	rmp, err := RoleMap()
	if err != nil {
		return project, err
	}

	nRole, ok := rmp[project.Role]
	if !ok {
		return project, errors.New("cannot get role: " + project.Role)
	}
	project.Role = nRole

	return project, nil
}

// UploadLambda - uploads
func (lp LProject) UploadLambda(name string) error {
	fpath := path.Join(lp.path, name)

	os.Setenv("GOOS", "linux")
	os.Setenv("GOARCH", "amd64")

	fmt.Println("Building: " + fpath + ".go")

	if _, err := run("go", "build", "-o", fpath, fpath+".go"); err != nil {
		return err
	}

	fmt.Println("Zipping: " + fpath + ".zip")
	if _, err := run("zip", "-j", fpath+".zip", fpath); err != nil {
		return err
	}

	lamname := lp.Name + "_" + name

	fmt.Println("aws cp")
	upcmd := exec.Command("aws", "s3", "cp", fpath+".zip", "s3://"+lp.Bucket+"/"+lamname+".zip")
	upOut, err := upcmd.StdoutPipe()
	if err != nil {
		return err
	}

	fmt.Println("Starting Upload of   " + lamname)

	if err = upcmd.Start(); err != nil {
		return err
	}

	io.Copy(os.Stdout, upOut)
	if err = upcmd.Wait(); err != nil {
		return err
	}

	fl, err := NewFunctionList()
	if err != nil {
		return err
	}

	fmt.Println("Create/update function")
	var resp []byte

	if fl.HasFunction(lamname) {
		resp, err = UpdateLambdaFunction(lamname, lp.Bucket)
	} else {
		resp, err = CreateLambdaFunction(lamname, lp.Role, name, lp.Bucket)
	}

	if err != nil {
		return err
	}
	fmt.Println(string(resp))

	return nil
}

func main() {
	lname := flag.String("n", "", "Name of Lambda")
	confloc := flag.String("c", "project.json", "json file  name")
	flag.Parse()

	proj, err := NewLProject(*confloc)
	if err != nil {
		log.Fatal(err)
	}

	err = proj.UploadLambda(*lname)
	if err != nil {
		log.Fatal(err)
	}
}
