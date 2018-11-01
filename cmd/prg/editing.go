package main

import (
	"errors"
	"fmt"
	"io/ioutil"
	"os"
	"os/exec"
	"path"
)

func tempFile(reason string, initialContent string) (string, error) {
	if reason != "" {
		reason = "." + reason
	}

	prefix := fmt.Sprintf(".progress%v.", reason)
	tempFile, err := ioutil.TempFile("", prefix)

	if err != nil {
		return "", errors.New("failed to create a temporary file for editing")
	}

	defer tempFile.Close()

	stat, err := tempFile.Stat()

	if err != nil {
		return "", errors.New("failed to stat temporary file")
	}

	absolutePath := path.Join(os.TempDir(), stat.Name())

	// Attempt to clean up file after it's used
	if _, err := tempFile.WriteString(initialContent); err != nil {
		return absolutePath, errors.New("failed to write initial contents to temporary file")
	}

	return absolutePath, nil
}

func edit(content string) (string, error) {
	editor := os.Getenv("EDITOR")

	if editor == "" {
		return "", errors.New("the EDITOR environment variable must be set in order to launch an editor")
	}

	absolutePath, err := tempFile("editor_session", content)

	if _, err := os.Stat(absolutePath); !os.IsNotExist(err) {
		defer os.Remove(absolutePath)
	}

	command := exec.Command(editor, absolutePath)

	command.Stderr = os.Stderr
	command.Stdin = os.Stdin
	command.Stdout = os.Stdout

	outputBytes, err := command.Output()
	output := fmt.Sprintf("%v", outputBytes)
	return output, err
}
