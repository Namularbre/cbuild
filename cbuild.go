package main

import (
	"bufio"
	"errors"
	"flag"
	"log"
	"os"
	"os/exec"
	"slices"
	"strings"
)

const BuildFile string = "cbuild.txt"
const CommentTag string = "#"
const DefaultExecutableName = "program"

func scanBuildFile(buildFile string, pwd string) ([]string, error) {
	file, err := os.OpenFile(buildFile, os.O_RDONLY, 755)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)

	var sourceFiles []string

	for scanner.Scan() {
		text := scanner.Text()

		// ignore commented lines and empty lines
		if strings.HasPrefix(text, CommentTag) || len(strings.Trim(text, " ")) == 0 {
			continue
		}

		// add the absolute path of the source file, so we're sure what file we are compiling
		sourceFile := pwd + "/" + text

		if slices.Contains(sourceFiles, sourceFile) {
			return nil, errors.New("there is two times the same source file :" + sourceFile)
		}

		sourceFiles = append(sourceFiles, sourceFile)
	}

	return sourceFiles, nil
}

func compile(sourcesFiles []string, outputFile string) error {
	// we put the sources files in the args of gcc, and then we add the executable file name
	var gccArgs []string

	// Prepare the arguments
	gccArgs = append(gccArgs, sourcesFiles...)
	gccArgs = append(gccArgs, "-o")
	gccArgs = append(gccArgs, outputFile)

	log.Printf("gcc arguments: %s\n", gccArgs)

	cmd := exec.Command("gcc", gccArgs...)

	out, err := cmd.CombinedOutput()
	if err != nil {
		log.Printf("error: %s\n", string(out))
		return err
	}

	log.Printf("GCC output: %s\n", string(out))

	return nil
}

func main() {
	executableName := flag.String("out", DefaultExecutableName, "Name of the output executable file")

	flag.Parse()

	pwd, err := os.Getwd()
	cbuildFilePath := pwd + "/" + BuildFile

	if err != nil {
		panic(err)
	}

	log.Printf("gathering file for %s\n", cbuildFilePath)

	sourceFiles, err := scanBuildFile(cbuildFilePath, pwd)
	if err != nil {
		log.Fatalf("error while gathering files : %s\n", err.Error())
	}

	log.Printf("building project in : %s, with files : %s\n", cbuildFilePath, sourceFiles)

	err = compile(sourceFiles, *executableName)
	if err != nil {
		log.Fatalf("error while compiling : %s\n", err.Error())
	}
}
