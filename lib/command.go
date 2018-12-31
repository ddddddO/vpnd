package lib

import (
	"bufio"
	"fmt"
	"log"
	"os/exec"
	"strings"
)

// TODO:1コマンド複数チェック・複数reコマンドjsonで設定できないか
// FIXME:json optionいらなくなるかも
func Command(c *Config) {
	for _, v := range c.Commands {
		log.Printf("-%s COMMAND START-\n", v.Command)

		var (
			cmd       = exec.Command("sh", "-c", v.Command, v.Option)
			stdout, _ = cmd.StdoutPipe()
			//stdin, _ = cmd.StdinPipe()
			//stderr, _ = cmd.StderrPipe()
		)
		if err := cmd.Start(); err != nil {
			log.Fatalf("failed to command '%s': %v\n", v.Command, err)
		}
	
		var (
			sc       = bufio.NewScanner(stdout)
			commandFlg = false
		)
		for sc.Scan() {
			s := sc.Text()
			if strings.Contains(s, v.Check) {
				commandFlg = true
				break
			}
		}
		cmd.Wait()
		log.Printf("-%s COMMAND END-\n", v.Command)

		if !commandFlg {
			log.Printf("-%s ReCOMMAND START-\n", v.ReCommandConfig.ReCommand)

			var (
				reCmd       = exec.Command("sh", "-c", v.ReCommandConfig.ReCommand, v.ReCommandConfig.Option)
				//reStdout, _ = cmd.StdoutPipe()
				//reStdin, _ = cmd.StdinPipe()
				//reStderr, _ = cmd.StderrPipe()
			)
			if err := reCmd.Run(); err != nil {
				log.Fatalf("failed to recommand '%s': %v\n", v.ReCommandConfig.ReCommand, err)
			}
			log.Printf("-%s ReCOMMAND END-\n", v.ReCommandConfig.ReCommand)
		}
	}
}

func VPNCommand() {
	vpnCmd := exec.Command("./vpnclient/vpncmd")
	vpnStdout, _ := vpnCmd.StdoutPipe()
	vpnStdin, _ := vpnCmd.StdinPipe()
	if err := vpnCmd.Start(); err != nil {
		log.Fatalf("failed to command '%s': %v\n", "vpncmd", err)
	}

	vpnSc := bufio.NewScanner(vpnStdout)
	log.Println("-VPN COMMAND START-")
	for vpnSc.Scan() {
		s := vpnSc.Text()
		fmt.Println(s)
		if strings.Contains(s, "3. VPN Tools コマンドの使用 (証明書作成や通信速度測定)") {
			vpnStdin.Write([]byte("2\n"))
			continue
		}

		if strings.Contains(s, "何も入力せずに Enter を押すと、localhost (このコンピュータ) に接続します。") {
			vpnStdin.Write([]byte("\n"))
			continue
		}

		// AccountStatusGetコマンド実行
		if strings.Contains(s, `VPN Client "localhost" に接続しました。`) {
			vpnStdin.Write([]byte(fmt.Sprintf("AccountStatusGet %s\n", "MYIPSE")))
			continue
		}

		if strings.Contains(s, "セッション接続状態") {
			if strings.Contains(s, "接続完了 (セッション確立済み)") {
				continue
			}
		}

		if strings.Contains(s, "コマンドは正常に終了しました。") {
			vpnStdin.Write([]byte("QUIT\n"))
		}
	}
	vpnCmd.Wait()
	log.Println("-VPN COMMAND END-")
}