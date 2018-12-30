package main

import (
	"fmt"
	"log"
	"os"
	"os/exec"
)

func main() {
	//cmd := exec.Command("ls", "-al")
	c := "ifconfig"
	cmd := exec.Command(c)
	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout
	if err := cmd.Run(); err != nil {
		log.Fatalf("failed to command '%s': %v\n", c, err)
	}

	var w []byte
	_, err := cmd.Stdout.Write(w)
	if err != nil {
		log.Fatalln("failed to stdout write '%s': %v\n", c, err)
	}

	fmt.Print(string(w))

}
