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

type scriptInfo struct {
	Name   string
	Suffix string
	Path   string
}

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
		if *listOnly {
			os.Exit(1)
		} else {
			log.Fatal(fmt.Sprintf(noScriptyDirError+"\n", scriptyDirName))
		}
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

func getScriptInfos(scriptyDir string) []*scriptInfo {
	files, err := ioutil.ReadDir(scriptyDir)

	if err != nil {
		log.Fatal(cantReadDir, ": ", scriptyDir)
	}

	scriptInfos := make([]*scriptInfo, len(files))
	for i, file := range files {
		if file.Mode().IsDir() {
			// TODO
		}
		scriptInfos[i] = getScriptInfo(scriptyDir, file)
	}
	return scriptInfos
}

func getScriptInfo(scriptyDir string, file os.FileInfo) *scriptInfo {
	scriptInfo := &scriptInfo{}
	name := file.Name()
	scriptInfo.Path = path.Join(scriptyDir, name)
	for _, suffix := range suffixWhiteList {
		if strings.HasSuffix(name, suffix) {
			name = strings.NewReplacer(suffix, "").Replace(name)
			scriptInfo.Suffix = suffix
			break
		}
	}
	scriptInfo.Name = name
	return scriptInfo
}

func main() {
	scriptArg, args := parseArgs()

	if scriptArg == "" && !*listOnly {
		fmt.Print(usage)
		return
	}

	scriptyDir := getScriptyDir()

	scriptInfos := getScriptInfos(scriptyDir)

	if scriptArg == "" {
		for _, scriptInfo := range scriptInfos {
			fmt.Println(scriptInfo.Name)
		}
		return
	}

	for _, scriptInfo := range scriptInfos {
		if scriptArg == scriptInfo.Name ||
			scriptArg == (scriptInfo.Name + scriptInfo.Suffix) {
			args[0] = scriptInfo.Path
			runCommandInteractively(args)
			return
		}
	}

	log.Fatal(argNotFound, ": ", scriptArg)

}
