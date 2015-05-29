package scripty

import (
	"flag"
	"fmt"
	"log"
	"os"
	"os/exec"
)

type scriptInfo struct {
	Name   string
	Suffix string
	Path   string
}

const (
	usagePrefix           = "Usage: scripty [options | <script_name>]\n\n"
	noScriptyDirError     = "No config file or scripty dir '%s' found"
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

func parseArgs() (scriptArg string, args []string) {
	flag.Parse()

	args = os.Args[1:]
	copy(args, args)

	if len(args) > 0 && !*listOnly && !*detailOnly {
		scriptArg = args[0]
	}
	return
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

func Main() {
	scriptArg, args := parseArgs()
	if scriptArg == "" && !*listOnly && !*detailOnly {
		fmt.Print(usagePrefix)
		flag.PrintDefaults()
		return
	}

	// if we can't read the current directory, fail
	cwd, err := os.Getwd()

	if err != nil {
		log.Fatal(cantReadDir, ": ", cwd)
	}

	config := getConfig(cwd)

	scriptInfos := getScriptInfos(config.ScriptyDir)

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
