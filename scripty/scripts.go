package scripty

import (
	"bufio"
	"io/ioutil"
	"log"
	"os"
	"path"
	"strings"
)

func isIgnoredHeaderLine(line string) bool {
	return line == "" || strings.HasPrefix(line, "#!")
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
