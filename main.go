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

const (
	noScriptyDirError = "No scripty dir found"
	scriptyDir = "scripts"
	chooseMsg = "choose one of the following:"
	scriptRunner = "bash"
	cantReadDir = "can't read dir"
	argNotFound = "argument not found in scripts"
	defaultSuffix = ".sh"
)

func findScriptyDir(startPath string) string {
	// make sure we haven't recursed all the way up
	if path.Clean(startPath) == "/" {
		log.Fatal(noScriptyDirError)
	}

	files, _ := ioutil.ReadDir(startPath)
	for _, file := range files {
		if file.Name() == scriptyDir {
			return path.Join(startPath, scriptyDir)
		}
	}
	return findScriptyDir(path.Join(startPath, ".."))

}

func main() {
	cwd, err := os.Getwd()
	if err != nil {
		log.Fatal(cantReadDir, ": ", cwd)
	}

	args := os.Args[1:]
	copy(args, args)

	var scriptArg string

	if len(args) > 0 {
		scriptArg = args[0]
	}

	scriptyDir := findScriptyDir(cwd)
	files, err := ioutil.ReadDir(scriptyDir)

	if err != nil {
		log.Fatal(cantReadDir, ": ", scriptyDir)
	}

	if scriptArg == "" {
		fmt.Println(chooseMsg)
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
			if strings.HasSuffix(name, defaultSuffix) {
				choices = stringSet{scriptArg: true, scriptArg + defaultSuffix: true}
			} else {
				choices = stringSet{scriptArg: true}
			}
			_, found = choices[name]
			if found {
				args[0] = path.Join(scriptyDir, name)
				cmd := exec.Command(scriptRunner, args...)
				cmd.Stdout = os.Stdout
				cmd.Stderr = os.Stderr
				cmd.Run()
				break
			}
		}
		if !found {
			log.Fatal(argNotFound, ": ", scriptArg)
		}
	}
}
