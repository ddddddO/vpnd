package main

import (
	"bufio"
	"log"
	"os/exec"
	"strings"
)

func main() {
	var (
		command = "ifconfig"
		target = "wifi0" // TODO:実行環境の設定に変更
		cmd = exec.Command(command, target)
		stdout, _ = cmd.StdoutPipe()
		//stdin, _ = cmd.StdinPipe()
		//stderr, _ = cmd.StderrPipe()
	)
	if err := cmd.Start(); err != nil {
		log.Fatalf("failed to command '%s': %v\n", command, err)
	}

	var (
		sc = bufio.NewScanner(stdout)
		ip = "192.168.179.4"
		netmask = "255.255.255.0"
		ifconFlg = false 
	)
	for sc.Scan() {
		s := sc.Text()
		if strings.Contains(s, ip) && strings.Contains(s, netmask) {
			log.Println("in!!!")
			ifconFlg = true
			break
		}
	}
	cmd.Wait()

	if !ifconFlg {
		// 仮想LANカードにグローバルIPを設定する処理
		log.Println("to set g ip")
	}
}
