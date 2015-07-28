package main

import (
	"bytes"
	"io/ioutil"
	"log"
	"os"
	"os/exec"
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

func main() {
	cmd := exec.Command("git", "status", "-s")
	var out bytes.Buffer
	cmd.Stdout = &out
	err := cmd.Run()
	if err != nil {
		panic(err)
	}

	resultLines := strings.Split(out.String(), "\n")
	if resultLines == nil {
		return
	}

	re := regexp.MustCompile(`package\s(\w+)`)
	names := make(map[string]bool)
	for _, filenameStatus := range resultLines {
		if filenameStatus != "" {
			explodedStatus := strings.Split(strings.TrimSpace(filenameStatus), " ")
			if len(explodedStatus) > 1 {
				filename := explodedStatus[len(explodedStatus)-1]
				status := explodedStatus[0]
				if status == "??" {
					continue
				}

				content := readFile(filename)

				pkg := re.FindStringSubmatch(string(content))
				if len(pkg) > 1 {
					names[pkg[1]] = true
				}

			}
		}
	}

	if len(names) > 0 {
		prefix := "["
		i := 0
		for name, _ := range names {
			if i > 0 {
				prefix += "|"
			}
			prefix += name
			i++
		}
		prefix += "] "
		filename := os.Args[1]
		msg := readFile(filename)
		ioutil.WriteFile(filename, []byte(prefix+msg), 0x777)
	}
}
