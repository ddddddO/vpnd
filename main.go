package main

import (
	"bufio"
	"encoding/json"
	"io/ioutil"
	"log"
	"os/exec"
	"strings"
)

type Config struct {
	Command string `json:"command"`
	Target  string `json:"target"` //TODO:実行環境の設定に変更
	Ip      string `json:"ip"`
	Netmask string `json:"netmask"`
}

func main() {
	raw, err := ioutil.ReadFile("./config.json")
	if err != nil {
		log.Fatalf("failed to read file %s: %v\n", "config.json", err)
	}

	var c Config
	err = json.Unmarshal(raw, &c)
	if err != nil {
		log.Fatal("failed to json unmarshal: %v\n", err)
	}

	var (
		cmd       = exec.Command(c.Command, c.Target)
		stdout, _ = cmd.StdoutPipe()
		//stdin, _ = cmd.StdinPipe()
		//stderr, _ = cmd.StderrPipe()
	)
	if err := cmd.Start(); err != nil {
		log.Fatalf("failed to command '%s': %v\n", c.Command, err)
	}

	var (
		sc       = bufio.NewScanner(stdout)
		ifconFlg = false
	)
	for sc.Scan() {
		s := sc.Text()
		if strings.Contains(s, c.Ip) && strings.Contains(s, c.Netmask) {
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
