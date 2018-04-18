package filescanner

import (
	"errors"
	"net/url"
	"net/http"
	"os"
	"io"
	"fmt"
	"strings"
	"bufio"
)

type DockerFileScanResult struct {
	Approved bool
	Messages []string
	Errors   []string
}

const WildCard = "*"

func FindFrom(commands []Command) ([]Command, error) {
	var com []Command
	var err error
	for _, v := range commands {
		if v.Cmd == "from" {
			com = append(com, v)
			break
		} else {
			err = errors.New("Cant find from")
		}
	}
	return com, err
}

func Scan(dockerFilePath string, whiteListPath string) (DockerFileScanResult, error) {
	whitelist, _ := whiteListResolver(whiteListPath)
	var result = DockerFileScanResult{
	}
	command, _ := ParseFile(dockerFilePath)
	froms, err := FindFrom(command)
	if err != nil {
		return result, err
	}

	var match = false
	for _, f := range froms {
		name := f.Value[0]
		match = whitelistMatch(whitelist, name)
		if match == true {
			result.Messages = append(result.Messages, fmt.Sprint("Base Image: ", name, " 	Approved"))
		}
	}
	if match == false {
		result.Errors = append(result.Errors, "Image Not Found in WhiteList")
	}
	return result, err
}

func whiteListResolver(whiteListPath string) (io.Reader, error) {
	var finalError error
	var reader io.Reader
	// check if url
	_, err := url.ParseRequestURI(whiteListPath)
	if err == nil {
		content, error := http.Get(whiteListPath)
		if error != nil {
			finalError = error
		}
		reader = content.Body

	} else {
		file, error := os.Open(whiteListPath)
		if error != nil {
			finalError = error
		}
		reader = file
	}
	return reader, finalError
}

//private checkers

func isDockerHubImage(image string) bool {
	return !(strings.Contains(image, "/"))
}

func whitelistMatch(list io.Reader, image string) bool {
	var eof error
	reader := bufio.NewReader(list)
	var result = false
	//iterate through the file line by line
	for eof == nil {
		line, er := reader.ReadBytes('\n')
		eof = er
		//get line lowercase and without the '\n'
		lineStr := strings.TrimRight(strings.ToLower(string(line)), "\n")
		//sanitize the incoming image
		imageLower := strings.Trim(strings.ToLower(image), " \t\r")

		if strings.Contains(lineStr, WildCard) {
			lineWithoutStar := strings.TrimRight(lineStr, WildCard)
			result = strings.Contains(imageLower, lineWithoutStar)
			break
		} else {
			result = strings.EqualFold(imageLower, lineStr)
		}
	}
	return result
}
