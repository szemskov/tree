package main

import (
	"fmt"
	"io"
	"os"
	"sort"
	"strings"
)

const (
	separator      = string(os.PathSeparator)
	emptyIndent    = "	"
	commonIndent   = "│	"
	regularPrefix  = "├───"
	lastPrefix     = "└───"
	emptySizeLabel = "empty"
)

func printFile(out io.Writer, file os.FileInfo, prefix string) {
	pathNames := strings.Split(file.Name(), separator)

	size := emptySizeLabel
	if file.Size() > 0 {
		size = fmt.Sprintf("%db", file.Size())
	}

	_, err := fmt.Fprintf(out, "%s%s (%s)\n", prefix, pathNames[len(pathNames)-1], size)
	if err != nil {
		panic(err.Error())
	}

}

func printDir(out io.Writer, dir os.FileInfo, prefix string) {
	pathNames := strings.Split(dir.Name(), separator)

	_, err := fmt.Fprintf(out, "%s%s\n", prefix, pathNames[len(pathNames)-1])
	if err != nil {
		panic(err.Error())
	}
}

func getLastDirectoryIndex(files []os.FileInfo) int {
	lastIndex := -1
	for index, meta := range files {
		if meta.IsDir() {
			lastIndex = index
		}
	}

	return lastIndex
}

func printTreeLevel(out io.Writer, path string, printFiles bool, indent string) error {
	file, err := os.Open(path)
	if err != nil {
		return err
	}

	files, err := file.Readdir(0)
	if err != nil {
		return err
	}

	sort.Slice(files, func(i, j int) bool { return files[i].Name() < files[j].Name() })

	lastDirectoryIndex := getLastDirectoryIndex(files)

	for index, meta := range files {

		if !meta.IsDir() && !printFiles {
			continue
		}

		prefix := fmt.Sprintf("%s%s", indent, regularPrefix)
		levelIndent := fmt.Sprintf("%s%s", indent, commonIndent)
		if (!printFiles && index == lastDirectoryIndex) || index == len(files) -1  {
			prefix = fmt.Sprintf("%s%s", indent, lastPrefix)
			levelIndent = fmt.Sprintf("%s%s", indent, emptyIndent)
		}

		if !meta.IsDir() {
			printFile(out, meta, prefix)
			continue
		}

		printDir(out, meta, prefix)

		nextLevelPath := fmt.Sprintf("%s%s%s", path, separator, meta.Name())
		if err = printTreeLevel(out, nextLevelPath, printFiles, levelIndent); err != nil {
			return err
		}
	}

	return nil
}

func dirTree(out io.Writer, path string, printFiles bool) error {
	// print path
	fmt.Println(path)

	return printTreeLevel(out, path, printFiles, "")
}

func main() {
	out := os.Stdout
	if !(len(os.Args) == 2 || len(os.Args) == 3) {
		panic("usage go run main.go . [-f]")
	}
	path := os.Args[1]
	printFiles := len(os.Args) == 3 && os.Args[2] == "-f"
	err := dirTree(out, path, printFiles)
	if err != nil {
		panic(err.Error())
	}
}
