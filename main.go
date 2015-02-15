package main

import (
	"flag"
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
	scriptyDir        = "scripts"
	chooseMsg         = "choose one of the following:"
	scriptRunner      = "bash"
	cantReadDir       = "can't read dir"
	argNotFound       = "argument not found in scripts"
	defaultSuffix     = ".sh"
)

var (
	listOnly = flag.Bool("l", false, "Print all possible scripts")
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

func parseArgs() (scriptArg string, args []string) {
	flag.Parse()

	args = os.Args[1:]
	copy(args, args)

	if len(args) > 0 && !*listOnly {
		scriptArg = args[0]
	}
	return
}

func getScriptyDir() string {
	cwd, err := os.Getwd()
	if err != nil {
		log.Fatal(cantReadDir, ": ", cwd)
	}

	return findScriptyDir(cwd)
}

func runCommandInteractively(args []string) {
	cmd := exec.Command(scriptRunner, args...)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	cmd.Run()
}

func main() {
	scriptArg, args := parseArgs()

	scriptyDir := getScriptyDir()

	files, err := ioutil.ReadDir(scriptyDir)

	if err != nil {
		log.Fatal(cantReadDir, ": ", scriptyDir)
	}

	if scriptArg == "" {
		if !*listOnly {
			fmt.Println(chooseMsg)
		}
		for _, file := range files {
			name := file.Name()
			if strings.HasSuffix(name, defaultSuffix) {
				name = strings.NewReplacer(defaultSuffix, "").Replace(name)
			}
			fmt.Println(name)
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
				runCommandInteractively(args)
				break
			}
		}
		if !found {
			log.Fatal(argNotFound, ": ", scriptArg)
		}
	}
}
