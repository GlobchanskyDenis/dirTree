package main

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"os"
	"strconv"
)

func findSize(file string, printFiles bool) string {
	fileData, err := ioutil.ReadFile(file)
	if printFiles && err == nil {
		if len(fileData) == 0 {
			return " (empty)\n"
		} else {
			return " (" + strconv.Itoa(len(fileData)) + "b)\n"
		}
	}
	return "\n"
}

func findFilePrefix(isNotLastLine int) string {
	if isNotLastLine != 0 {
		return "├───"
	}
	return "└───"
}

func newPrefix(oldPrefix []string, isNotLastLine int) []string {
	if isNotLastLine == 0 {
		return append(oldPrefix, "\t")
	} else {
		return append(oldPrefix, "│\t")
	}
}

func addPrefixToBuffer(out *bytes.Buffer, prefix []string) {
	for _, pref := range prefix {
		out.Write([]byte(pref))
	}
}

func isDirectory(filename string) bool {
	_, err := ioutil.ReadDir(filename)
	if err != nil {
		return false
	}
	return true
}

func findDirTree(out *bytes.Buffer, higherDir string, prefix []string, printFiles bool) {

	var filePrefix string
	var dirLen int
	var i int
	var size string
	var err error
	var files []os.FileInfo
	var toPrint []string

	files, err = ioutil.ReadDir(higherDir)
	if err != nil {
		return
	}

	// fmt.Println("\n", files, "\n")

	// toPrint = make([]string, 100)
	for _, file := range files {
		if printFiles || isDirectory(higherDir+"/"+file.Name()) {
			toPrint = append(toPrint, file.Name())
			// fmt.Println(higherDir+"/"+file.Name())
		}
	}

	// fmt.Println("\t", toPrint)
	
	dirLen = len(toPrint)
	i = 0
	for i < dirLen {
		filePrefix = findFilePrefix(dirLen - 1 - i)
		addPrefixToBuffer(out, prefix)

		//	find Postfix (size + EOL)
		size = findSize(higherDir+"/"+toPrint[i],
			printFiles)
		//	add line in buffer
		out.Write([]byte(filePrefix + toPrint[i] + size))

		//	recursion
		findDirTree(out, higherDir+"/"+toPrint[i],
			newPrefix(prefix, dirLen-1-i),
			printFiles)
		i++
	}
	return
}

func dirTree(out *bytes.Buffer, higherDir string, printFiles bool) error {
	var prefix []string = make([]string, 1)

	findDirTree(out, higherDir, prefix, printFiles)
	return nil
}

func main() {
	var out *bytes.Buffer = new(bytes.Buffer)
	var Args []string
	var printFiles bool

	Args = os.Args
	if len(Args) == 2 && Args[1] == "-f" {
		printFiles = true
	} else if len(Args) != 1 {
		fmt.Println("bad flag")
		return
	}
	dirTree(out, ".", printFiles)
	fmt.Print(out)
}
