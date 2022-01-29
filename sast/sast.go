package sast

func semgrepScan(config []string, dir string){
	fmt.Println(config)
	var configuration string
	var output bytes.Buffer
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
		}
		// fmt.Println("out:", outb.String(), "\nerr:", errb.String())
		output.WriteString(outb.String())
	}
	fmt.Println(output.String())
}