package main

import (
	"flag"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"strings"
)

var Version = "development"

func main() {

	// Declare variables to hold the values from the command line
	var builder strings.Builder
	var path string
	var exclude []string
	var excludeList string
	var replaceCharacter string

	// declare the command line flags
	flag.StringVar(&path, "p", "envfile", "Specify the path to the envfile")
	flag.StringVar(&excludeList, "e", "", "Command delimited list of variables to exclude from the file")
	flag.StringVar(&replaceCharacter, "r", " ", "Character to replace newlines with")

	// parse the flags
	flag.Parse()

	// create the string builder
	builder = strings.Builder{}

	// turn the list of excludes into a slice
	exclude = strings.Split(excludeList, ",")

	// Determine the path to the envfile
	// If it is not absolute, then prepend the current dir onto it
	isAbsolute := filepath.IsAbs(path)
	if !isAbsolute {

		// get the current directory
		cwd, err := os.Getwd()

		if err != nil {
			log.Fatalf("Unable to get current directory: %s\n", err.Error())
		}

		path = filepath.Join(cwd, path)
	}

	// get all the environment variables and iterate around them
	for _, env := range os.Environ() {

		// split the environment variable using = as the delimiter
		// this is so that newlines can be surpressed
		parts := strings.Split(env, "=")

		// if the name is in the exlude list then move onto the next iteration
		name := parts[0]
		if sliceContains(exclude, name) {
			continue
		}

		// replace the newline character with a space
		value := strings.Replace(parts[1], "\n", replaceCharacter, -1)

		// Add the key and the value to the string builder
		builder.WriteString(fmt.Sprintf("%s=%s\n", name, value))
	}

	// get the string from the builder
	output := builder.String()

	// Write the output to the file
	if err := os.WriteFile(path, []byte(output), 0666); err != nil {
		log.Fatalf("Error writing out file: %s\n", err.Error())
	}
}

// sliceContains performs a case insensitive match to see if the slice
// contains the specified value
func sliceContains(slice []string, value string) bool {
	var result bool

	for _, x := range slice {
		if strings.EqualFold(x, value) {
			result = true
			break
		}
	}

	return result
}
