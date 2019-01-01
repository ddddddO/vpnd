package lib

import (
	"bufio"
	"fmt"
	"log"
	"os/exec"
	"strings"
)

func TmpCommand() {
	log.Println("-TMP COMMAND START-")
	tcmd := exec.Command("bash", "-c", "touch /home/pi/Vpnclient/test.txt")
	if err := tcmd.Start(); err != nil {
		log.Fatalf("miss: %v", err)
	}
	tcmd.Wait()
	log.Println("-TMP COMMAND END-")
}

// TODO:1コマンド複数チェック・複数reコマンドjsonで設定できないか
// FIXME:json optionいらなくなるかも
// NOTE:exec.Command("bash", "-c", "コマンド")は、exec.Run×/exec.Start〇なものがあるっぽい。
func Command(c *Config) {
	for _, v := range c.Commands {
		log.Printf("-%s COMMAND START-\n", v.Command)

		var (
			cmd       = exec.Command("bash", "-c", v.Command)
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
				reCmd       = exec.Command("bash", "-c", v.ReCommandConfig.ReCommand)
				//reStdout, _ = cmd.StdoutPipe()
				//reStdin, _ = cmd.StdinPipe()
				//reStderr, _ = cmd.StderrPipe()
			)
			if err := reCmd.Start(); err != nil {
				log.Fatalf("failed to recommand '%s': %v\n", v.ReCommandConfig.ReCommand, err)
			}
			reCmd.Wait()
			log.Printf("-%s ReCOMMAND END-\n", v.ReCommandConfig.ReCommand)
		}
	}
}

func VPNCommand() {
	vpnClientCmd := exec.Command("/home/pi/Vpnclient/vpnclient/vpnclient", "start")
	if err := vpnClientCmd.Run(); err != nil {
		log.Fatalf("failed to command '%s': %v\n", "vpnclientcmd", err)
	}

	vpnCmd := exec.Command("/home/pi/Vpnclient/vpnclient/vpncmd")
	vpnStdout, _ := vpnCmd.StdoutPipe()
	vpnStdin, _ := vpnCmd.StdinPipe()
	if err := vpnCmd.Start(); err != nil {
		log.Fatalf("failed to command '%s': %v\n", "vpncmd", err)	
	}

	sessionFlg := false
	vpnSc := bufio.NewScanner(vpnStdout)
	log.Println("-VPN COMMAND START-")
	for vpnSc.Scan() {
		s := vpnSc.Text()
		if strings.Contains(s, "3. VPN Tools コマンドの使用 (証明書作成や通信速度測定)") {
			vpnStdin.Write([]byte("2\n"))
			continue
		}

		if strings.Contains(s, "何も入力せずに Enter を押すと、localhost (このコンピュータ) に接続します。") {
			vpnStdin.Write([]byte("\n"))
			continue
		}

		// AccountStatusGetコマンド実行
		// ここちょっと考えないといけない・・？
		if strings.Contains(s, `VPN Client "localhost" に接続しました。`) {
			vpnStdin.Write([]byte(fmt.Sprintf("AccountStatusGet %s\n", "MYIPSE")))
			continue
		}

		if strings.Contains(s, "セッション接続状態") {
			if strings.Contains(s, "接続完了 (セッション確立済み)") {
				sessionFlg = true
				continue
			}
		}

		if strings.Contains(s, "コマンドは正常に終了しました。") || strings.Contains(s, "指定された接続設定は接続されていません。") {
			vpnStdin.Write([]byte("QUIT\n"))
		}
	}
	vpnCmd.Wait()
	log.Println("-VPN COMMAND END-")

	// 再vpncmd実行で、AccountConnectコマンド実行する
	if !sessionFlg {
		vpnReCmd := exec.Command("/home/pi/Vpnclient/vpnclient/vpncmd")
		vpnReStdout, _ := vpnReCmd.StdoutPipe()
		vpnReStdin, _ := vpnReCmd.StdinPipe()	
		if err := vpnReCmd.Start(); err != nil {
			log.Fatalf("failed to recommand '%s': %v\n", "vpncmd", err)	
		}
	
		vpnReSc := bufio.NewScanner(vpnReStdout)
		log.Println("-VPN ReCOMMAND START-")
		for vpnReSc.Scan() {
			s := vpnReSc.Text()
			if strings.Contains(s, "3. VPN Tools コマンドの使用 (証明書作成や通信速度測定)") {
				vpnReStdin.Write([]byte("2\n"))
				continue
			}
	
			if strings.Contains(s, "何も入力せずに Enter を押すと、localhost (このコンピュータ) に接続します。") {
				vpnReStdin.Write([]byte("\n"))
				continue
			}
	
			// AccountConnectコマンド実行
			if strings.Contains(s, `VPN Client "localhost" に接続しました。`) {
				vpnReStdin.Write([]byte(fmt.Sprintf("AccountConnect %s\n", "MYIPSE")))
				sessionFlg = true
				continue
			}
	/*
			if strings.Contains(s, "セッション接続状態") {
				if strings.Contains(s, "接続完了 (セッション確立済み)") {
					sessionFlg = true
					continue
				}
			}
	*/
			if strings.Contains(s, "コマンドは正常に終了しました。") {
				vpnReStdin.Write([]byte("QUIT\n"))
			}
		}
		vpnReCmd.Wait()
		log.Println("-VPN ReCOMMAND END-")
		
		if !sessionFlg {
			log.Fatalln("そのとき考える")
		}
	}
}
