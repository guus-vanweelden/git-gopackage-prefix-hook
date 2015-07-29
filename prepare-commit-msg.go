package main

import (
	"bytes"
	"io/ioutil"
	"log"
	"os"
	"os/exec"
	"path/filepath"
	"regexp"
	"strings"
)

func readFile(filename string) string {
	contents, err := ioutil.ReadFile(filename)
	if err != nil {
		log.Printf("Error during reading file: %s", err)
		return ""
	}
	return string(contents)
}

func runCommand(command string, values ...string) string {
	cmd := exec.Command(command, values...)
	var out bytes.Buffer
	cmd.Stdout = &out
	err := cmd.Run()
	if err != nil {
		panic(err)
	}
	return out.String()
}

func main() {
	gitDir := runCommand("git", "rev-parse", "--git-dir")
	gitDiff := runCommand("git", "diff", "--name-only", "--cached")

	absolutePath := filepath.Dir(gitDir)

	resultLines := strings.Split(gitDiff, "\n")
	if resultLines == nil {
		return
	}

	re := regexp.MustCompile(`package\s(\w+)`)
	names := make(map[string]bool)
	var prefix string
	for _, filename := range resultLines {
		if filename != "" {
			content := readFile(absolutePath + "/" + filename)

			pkg := re.FindStringSubmatch(string(content))

			if len(pkg) < 2 {
				continue
			}
			if _, found := names[pkg[1]]; !found {
				if prefix != "" {
					prefix += "|"
				}
				prefix += pkg[1]
				names[pkg[1]] = true
			}

		}
	}

	if prefix != "" {
		prefix = "[" + prefix + "] "
		filename := os.Args[1]
		msg := readFile(filename)
		ioutil.WriteFile(filename, []byte(prefix+msg), 0x777)
	}
}
