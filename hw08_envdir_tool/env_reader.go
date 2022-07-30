package main

import (
	"bufio"
	"io/ioutil"
	"os"
	"strings"
)

type Environment map[string]EnvValue

// EnvValue helps to distinguish between empty files and files with the first empty line.
type EnvValue struct {
	Value      string
	NeedRemove bool
}

var replacer = strings.NewReplacer("0x00", "\n", "\x00", "\n", "=", "")

// ReadDir reads a specified directory and returns map of env variables.
// Variables represented as files where filename is name of variable, file first line is a value.
func ReadDir(dir string) (Environment, error) {
	filesInfo, err := ioutil.ReadDir(dir)
	if err != nil {
		return nil, err
	}

	env := Environment{}

	for _, fileInfo := range filesInfo {
		file, err := os.Open(dir + "/" + fileInfo.Name())
		if err != nil {
			return nil, err
		}

		scanner := bufio.NewScanner(file)

		if scanner.Scan() {
			value := replacer.Replace(scanner.Text())
			value = strings.TrimRight(value, "\t ")
			env[fileInfo.Name()] = EnvValue{Value: value}
		} else {
			env[fileInfo.Name()] = EnvValue{NeedRemove: true}
		}

		if err := scanner.Err(); err != nil {
			return nil, err
		}

		if err := file.Close(); err != nil {
			return nil, err
		}
	}

	return env, nil
}
