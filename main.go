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

	args := os.Args[1:]
	copy(args, args)

	var scriptArg string

	if len(args) > 0 {
		scriptArg = args[0]
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
				args[0] = path.Join(scriptyDir, name)
				cmd := exec.Command("bash", args...)
				cmd.Stdout = os.Stdout
				cmd.Stderr = os.Stderr
				cmd.Run()
				break
			}
		}
		if !found {
			log.Fatal("argument not found in scripts: ", scriptArg)
		}
	}
}
