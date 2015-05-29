package scripty

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"path"
)

type config struct {
	ScriptyDir    string
	ExtraCommands map[string][]string
}

const (
	configFileName     = ".scripty.json"
	cantReadConfigJson = "Can't read scriptyDir from json"
)

func readConfig(startPath string) *config {
	jdata, err := ioutil.ReadFile(path.Join(startPath, configFileName))
	var options config

	if err != nil {
		log.Fatal(err)
	}

	err = json.Unmarshal(jdata, &options)

	// die on basic jason errors
	if err != nil {
		log.Fatal(err)
	}

	// this signifies that the key was not read properly from
	// the json file, it needs to be cleaned up
	if options.ScriptyDir == "" {
		log.Fatal(cantReadConfigJson)
	}

	// paths are relative in configs, but need to be absolute
	// to be consumed.
	options.ScriptyDir = path.Join(startPath, options.ScriptyDir)
	return &options
}

// this is either going to find a config file or a scripty dir
// if it finds a scripty dir, it'll make a config object with it
// and return that. If it finds a config, it'll read it into a
// config object and return that.
func getConfig(startPath string) *config {
	// make sure we haven't recursed all the way up
	if path.Clean(startPath) == "/" {
		// these commands expect no output on failure, so
		// exit quietly
		if *listOnly || *detailOnly {
			os.Exit(1)
		} else {
			log.Fatal(fmt.Sprintf(noScriptyDirError+"\n", defaultScriptyDirName))
		}
	}

	files, _ := ioutil.ReadDir(startPath)

	for _, file := range files {
		if file.Name() == configFileName {
			return readConfig(startPath)
		}
	}

	for _, file := range files {
		if file.Name() == defaultScriptyDirName {
			scriptyDirName := path.Join(startPath, defaultScriptyDirName)
			return &config{ScriptyDir: scriptyDirName}
		}
	}

	return getConfig(path.Join(startPath, ".."))
}
