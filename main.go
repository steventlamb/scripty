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
	usage = ("Usage: scripty [-l] [<script_name>]\n\n" +
		"Run 'scripty -l' to see possible scripts\n")
	noScriptyDirError     = "No scripty dir '%s' found"
	scriptyDirEnvVar      = "SCRIPTY_DIR"
	defaultScriptyDirName = "scripts"
	cantReadDir           = "can't read dir"
	argNotFound           = "argument not found in scripts"
)

var (
	suffixWhiteList = []string{".sh", ".py"}
	listOnly        = flag.Bool("l", false, "Print all possible scripts")
)

func findScriptyDir(startPath string) string {
	scriptyDirName := os.Getenv(scriptyDirEnvVar)

	if scriptyDirName == "" {
		scriptyDirName = defaultScriptyDirName
	}

	// make sure we haven't recursed all the way up
	if path.Clean(startPath) == "/" {
		log.Fatal(fmt.Sprintf(noScriptyDirError+"\n", scriptyDirName))
	}

	files, _ := ioutil.ReadDir(startPath)
	for _, file := range files {
		if file.Name() == scriptyDirName {
			return path.Join(startPath, scriptyDirName)
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
	cmd := exec.Command(args[0], args[1:]...)
	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	err := cmd.Run()

	if err != nil {
		log.Fatal("SCRIPTY ERROR: ", err)
	}
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
			fmt.Print(usage)
		} else {
			for _, file := range files {
				name := file.Name()
				for _, suffix := range suffixWhiteList {
					if strings.HasSuffix(name, suffix) {
						name = strings.NewReplacer(suffix, "").Replace(name)
						break
					}
				}
				fmt.Println(name)
			}
		}
	} else {
		var foundScript bool
		for _, file := range files {
			// if the filename ends in .sh, then we can
			// optionally omit it from the scriptArg for
			// convenience. Otherwise, match exactly.
			name := file.Name()

			choices := stringSet{scriptArg: true}
			for _, suffix := range suffixWhiteList {
				if strings.HasSuffix(name, suffix) {
					choices[scriptArg+suffix] = true
					break
				}
			}

			_, foundScript = choices[name]
			if foundScript {
				args[0] = path.Join(scriptyDir, name)
				runCommandInteractively(args)
				break
			}
		}
		if !foundScript {
			log.Fatal(argNotFound, ": ", scriptArg)
		}
	}
}
