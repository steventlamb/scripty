package main

import (
	"bufio"
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"os/exec"
	"path"
	"strings"
)

type scriptInfo struct {
	Name   string
	Suffix string
	Path   string
}

const (
	usagePrefix           = "Usage: scripty [options | <script_name>]\n\n"
	noScriptyDirError     = "No scripty dir '%s' found"
	scriptyDirEnvVar      = "SCRIPTY_DIR"
	defaultScriptyDirName = "scripts"
	cantReadDir           = "can't read dir"
	argNotFound           = "argument not found in scripts"
	readUntilBodyLimit    = 20
)

var (
	suffixWhiteList = []string{".sh", ".py"}
	listOnly        = flag.Bool("l", false, "Print all available scripts (in machine readable format)")
	detailOnly      = flag.Bool("d", false, "Print all available scripts (with docstring, if available)")
)

func findScriptyDir(startPath string) string {
	scriptyDirName := os.Getenv(scriptyDirEnvVar)

	if scriptyDirName == "" {
		scriptyDirName = defaultScriptyDirName
	}

	// make sure we haven't recursed all the way up
	if path.Clean(startPath) == "/" {
		if *listOnly || *detailOnly {
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

	if len(args) > 0 && !*listOnly && !*detailOnly {
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

func isIgnoredHeaderLine(line string) bool {
	return line == "" || strings.HasPrefix(line, "#!")
}

func readFirstComment(path string) (string, error) {
	file, err := os.Open(path)
	defer file.Close()

	if err != nil {
		return "", err
	}

	scanner := bufio.NewScanner(file)

	scanner.Scan()
	currentLine := scanner.Text()

	// skip whitespace lines and the shebang
	for i := 0; i < readUntilBodyLimit && isIgnoredHeaderLine(currentLine); i++ {
		scanner.Scan()
		currentLine = scanner.Text()
	}

	// if you've reached something that isn't a comment, there are no docs
	// at the top, so return nothing
	if !strings.HasPrefix(currentLine, "#") {
		return "", nil
	}

	return strings.TrimRight(strings.TrimLeft(currentLine, "# "), " "), nil
}

func (info *scriptInfo) getDescription() string {
	description, err := readFirstComment(info.Path)

	if err != nil {
		log.Fatal(err)
	}
	return description
}

func getScriptInfos(nodePath string) []*scriptInfo {
	files, err := ioutil.ReadDir(nodePath)

	if err != nil {
		log.Fatal(cantReadDir, ": ", nodePath)
	}

	nodeInfos := make([]*scriptInfo, 0)
	for _, file := range files {
		if !file.Mode().IsRegular() {
			childPath := path.Join(nodePath, file.Name())
			nodeInfos = append(nodeInfos, getScriptInfos(childPath)...)
		} else {
			nodeInfos = append(nodeInfos, getScriptInfo(nodePath, file))
		}
	}
	return nodeInfos
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
	if scriptArg == "" && !*listOnly && !*detailOnly {
		fmt.Print(usagePrefix)
		flag.PrintDefaults()
		return
	}

	scriptyDir := getScriptyDir()

	scriptInfos := getScriptInfos(scriptyDir)

	if scriptArg == "" {
		var longName string
		for _, scriptInfo := range scriptInfos {
			if *detailOnly {
				fmt.Printf("%-25.25s %s\n", scriptInfo.Name, scriptInfo.getDescription())
				if len(scriptInfo.Name) > 25 {
					longName = scriptInfo.Name
				}
			} else {
				fmt.Println(scriptInfo.Name)
			}
		}
		if longName != "" {
			log.Fatalf("'%s' truncated for readability! Use 'scripty -l' instead.\n", longName)
		}
		return
	}

	for _, scriptInfo := range scriptInfos {
		if scriptArg == scriptInfo.Name ||
			scriptArg == (scriptInfo.Name+scriptInfo.Suffix) {
			args[0] = scriptInfo.Path
			runCommandInteractively(args)
			return
		}
	}

	log.Fatal(argNotFound, ": ", scriptArg)

}
