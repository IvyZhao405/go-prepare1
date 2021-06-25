package main

import (
	"fmt"
	"os/exec"
)

func main() {
	var(
		cmd *exec.Cmd
		output []byte
		err error
	)
	//create Cmd
	cmd = exec.Command("/bin/bash", "-c", "sleep 5;ls -l")

	//execute process
	if output, err = cmd.CombinedOutput();err != nil {
		fmt.Println(err)
		return
	}
	//print command output
	fmt.Println(string(output))
}
