package main

import (
	"github.com/sirupsen/logrus"
	"io"
	"os"
	"os/exec"
)

func setLogLevel() {
	if os.Getenv("LOG_LEVEL") == "DEBUG" {
		logrus.SetLevel(logrus.DebugLevel)
	}
}

func startReturnOutput(command *exec.Cmd) (io.ReadCloser, io.ReadCloser, error) {
	readerStdout, err := command.StdoutPipe()
	if err != nil {
		return nil, nil, err
	}

	readerStderr, err := command.StderrPipe()
	if err != nil {
		return nil, nil, err
	}

	if err := command.Start(); err != nil {
		readerStdout.Close()
		readerStderr.Close()
		return nil, nil, err
	}

	return readerStdout, readerStderr, nil
}

