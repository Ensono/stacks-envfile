package main

import (
	"flag"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"regexp"
	"strings"
)

var Version = "development"

func main() {

	// Declare variables to hold the values from the command line
	var builder strings.Builder
	var debug bool
	var exclude []string
	var excludeList string
	var include []string
	var includeList string
	var path string
	var replaceCharacter string

	// declare the command line flags
	flag.BoolVar(&debug, "d", false, "Output the values to the screen. CAUTION: this may show sensitive values")
	flag.StringVar(&excludeList, "e", "", "Comma delimited list of variables to exclude from the file")
	flag.StringVar(&includeList, "i", "", "Comma delimited list of variables to include in the envfile")
	flag.StringVar(&path, "p", "envfile", "Specify the path to the envfile")
	flag.StringVar(&replaceCharacter, "r", " ", "Character to replace newlines with")

	// parse the flags
	flag.Parse()

	// if an exclusion list and inclusion list have been specified throw an error
	if includeList != "" && excludeList != "" {
		log.Fatalf("Include and exclude lists are mutally exclusive")
	}

	// create the string builder
	builder = strings.Builder{}

	// turn the lists into a slice
	exclude = strings.Split(excludeList, ",")
	include = strings.Split(includeList, ",")

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

		// Get the name of the variable
		name := parts[0]

		// Determine if the variable should be excluded
		shouldExclude := sliceContains(exclude, name)
		shouldInclude := sliceContains(include, name)
		if shouldExclude && !shouldInclude {
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

	if debug {
		fmt.Println(output)
	}
}

// sliceContains performs a case insensitive match to see if the slice
// contains the specified value
func sliceContains(slice []string, value string) bool {
	var result bool

	for _, x := range slice {

		// create regular expression pattern to test against
		// this allows multiple variables to be added or excluded
		pattern := fmt.Sprintf(`(?i)\b%s\b`, x)
		re := regexp.MustCompile(pattern)

		// match the value against the re
		result = re.MatchString(value)

		// if strings.EqualFold(x, value) {
		// 	result = true
		// 	break
		// }
	}

	return result
}
