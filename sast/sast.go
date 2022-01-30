package sast

import (
	"log"
	"os/exec"
	"bytes"
	"errors"
)

func SemgrepCheck() bool {
	cmd := exec.Command("/bin/sh", "-c", "command -v semgrep")
	if err := cmd.Run(); err != nil {
			return false
	}
	return true
}

func SemgrepScan(config []string, dir string) ([]byte, error){
	var configuration string
	var output []byte

	if !SemgrepCheck(){
		err := errors.New("Semgrep command doesn't exist. Please install semgrep.")
		log.Print(err)
		return nil, err
	}

	for _, cnf := range config{
		var outBuffer, errBuffer bytes.Buffer
		if cnf == ""{
			configuration = "--config=auto"
		} else {
			configuration = "--config=p/"+cnf
		}
		cmd := exec.Command("semgrep", configuration, "--json", dir)
		cmd.Stdout = &outBuffer
		cmd.Stderr = &errBuffer
		err := cmd.Run()
		if err != nil{
			log.Printf("Semgrep command errored out with error ", err)
			return nil, err
		}
		output = append(output, outBuffer.Bytes()...)
	}
	return output, nil
}