package sast

import (
	"fmt"
	"os/exec"
	"bytes"
	"encoding/json"
)

// type SemgrepResults struct {
//     Errors []string
// 	Results []map[string]interface{}
// }

type SemgrepResults struct {
    Errors []string
	Results []byte
}

type SemgrepOutput struct {
	CheckID string
	Lines string
	Message string
}

func ParseOutput(rawResults []map[string]interface{}) string {
	var s2 []SemgrepOutput
	var s3 SemgrepOutput
	var tempmessage map[string]interface {}
	for key, result := range rawResults {
		fmt.Println("Reading value of key ", key)
		// s2[key].Check_id = result["check_id"]
		// s2[key].Lines = result["extra"]["lines"]
		// s2[key].Message = result["extra"]["message"]
		// fmt.Println(s2)
		tempmessage = result["extra"].(map[string]interface {})
		// fmt.Println(tempmessage["lines"])
		s3.CheckID = result["check_id"].(string)
		s3.Lines = tempmessage["lines"].(string)
		s3.Message = tempmessage["message"].(string)
		s2 = append(s2, s3)
	}
	s4, err := json.Marshal(s2)
	if err != nil {
		fmt.Println("Error while marshalling ", err)
	}
	fmt.Println(string(s4))
	return string(s4)
}

func SemgrepScan(config []string, dir string) ([]byte, error){
	fmt.Println(config)
	var configuration string
	// var output []string
	var output []byte
	// var sOutput SemgrepResults
	for _, cnf := range config{
		if cnf == ""{
			configuration = "--config=auto"
		} else {
			configuration = "--config=p/"+cnf
		}
		cmd := exec.Command("semgrep", configuration, "--json", dir)
		fmt.Println(cmd)
		var outb, errb bytes.Buffer
		cmd.Stdout = &outb
		cmd.Stderr = &errb
		err := cmd.Run()
		if err != nil{
			fmt.Println("Command errored out with error %v", err)
			return nil, err
		}
		// fmt.Println("out:", outb.String(), "\nerr:", errb.String())
		// output.WriteString(outb.String())
		// err = json.Unmarshal(outb.Bytes(), &sOutput)
		// if err != nil{
		// 	fmt.Println("Error unmarshaling json ", err)
		// }
		// output = append(output, ParseOutput(sOutput.Results))
		// fmt.Println(outb)
		// fmt.Println(outb.String())
		// output, err = json.Marshal(outb.Bytes())
		output = append(output, outb.Bytes()...)
		// fmt.Println(output)
		// if err != nil{
		// 	fmt.Println("Marshalling errored out with error %v", err)
		// 	return nil, err
		// }
	}
	// fmt.Println(output)

	return output, nil
}