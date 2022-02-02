package sast

import (
	"log"
	"os/exec"
	"bytes"
	"errors"
)

// Checks if Semgrep is installed on OS
func SemgrepCheck() bool {
	cmd := exec.Command("/bin/sh", "-c", "command -v semgrep")
	if err := cmd.Run(); err != nil {
			return false
	}
	return true
}


// Scans given directory with given config using Semgrep.
// Semgrep needs to be present as an executable in PATH environment variable  
func SemgrepScan(config []string, dir string) ([]byte, error){
	var configuration string
	var output []byte

	if !SemgrepCheck(){
		err := errors.New("Semgrep command doesn't exist. Please install semgrep. Refer https://semgrep.dev/")
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
			log.Printf("Semgrep command errored out with error %v", err)
			return nil, err
		}
		output = append(output, outBuffer.Bytes()...)
	}
	return output, nil
}