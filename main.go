package main

import (
	"github.com/robmerrell/comandante"
	"flag"
	"fmt"
	"os"
	"strings"
	"github.com/austbot/lineage/lib/filescanner"
)

var whiteListPath string
var dockerFilePath string

func main() {
	bin := comandante.New("lineage", "Validate Docker Image Ancestry")
	bin.IncludeHelp()
	scanDockerFile := comandante.NewCommand("scan-file", "Scan a Dockerfile", ScanDockerFileCtrl)
	scanDockerFile.FlagInit = func(set *flag.FlagSet) {
		set.StringVar(&whiteListPath, "whitelist", "whitelist.txt", "A file path or url.")
		set.StringVar(&dockerFilePath, "dockerfile", "Dockerfile", "A file path to a Dockerfile")
	}
	bin.RegisterCommand(scanDockerFile)

	if err := bin.Run(); err != nil {
		fmt.Fprintln(os.Stderr, err)
	}
}

func ScanDockerFileCtrl() error {
	print("Input: ", dockerFilePath, " ", whiteListPath, "\n")
	result, _ := filescanner.Scan(dockerFilePath, whiteListPath)
	if len(result.Errors) > 0 {
		fmt.Println("Result:", strings.Join(result.Errors, " "))
		os.Exit(1)
	}
	if len(result.Messages) > 0 {
		fmt.Println("Result:", strings.Join(result.Messages, " "))
		os.Exit(0)
	}
	return nil
}
