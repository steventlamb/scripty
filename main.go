package main

import (
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"os/exec"
	"path"
	"strings"
)

type stringSet map[string]bool

func findScriptyDir(startPath string) string {
	// make sure we haven't recursed all the way up
	if path.Clean(startPath) == "/" {
		log.Fatal("No scripty dir found")
	}

	files, _ := ioutil.ReadDir(startPath)
	for _, file := range files {
		// TODO: make this configurable with a .scripty file
		// instead.
		if file.Name() == "scripts" {
			return path.Join(startPath, "scripts")
		}
	}
	return findScriptyDir(path.Join(startPath, ".."))

}

func main() {
	cwd, err := os.Getwd()
	if err != nil {
		log.Fatal("can't read cwd %s", cwd)
	}

	var scriptArg string

	switch len(os.Args) {
	case 1:
		scriptArg = ""
	case 2:
		scriptArg = os.Args[1]
	default:
		log.Fatal("too many args")
	}

	scriptyDir := findScriptyDir(cwd)
	files, _ := ioutil.ReadDir(scriptyDir)

	if scriptArg == "" {
		fmt.Println("choose one of the following:")
		for _, file := range files {
			fmt.Println(file.Name())
		}
	} else {
		var found bool
		for _, file := range files {
			// if the filename ends in .sh, then we can
			// optionally omit it from the scriptArg for
			// convenience. Otherwise, match exactly.
			name := file.Name()
			var choices stringSet
			if strings.HasSuffix(name, ".sh") {
				choices = stringSet{scriptArg: true, scriptArg + ".sh": true}
			} else {
				choices = stringSet{scriptArg: true}
			}
			_, found = choices[name]
			if found {
				execPath := path.Join(scriptyDir, name)
				cmd := exec.Command("bash", execPath)
				cmd.Stdout = os.Stdout
				cmd.Stderr = os.Stderr
				cmd.Run()
			}
		}
		if !found {
			log.Fatal("argument not found in scripts")
		}
	}
}
