package lib

import (
	"bufio"
	"fmt"
	"log"
	"os/exec"
	"strings"
)

/* debug用
func TmpCommand() {
	log.Println("-TMP COMMAND START-")
	tcmd := exec.Command("bash", "-c", "ifconfig vpn_myipse-lan >> /home/pi/Vpnclient/ifconfig.txt")
	if err := tcmd.Start(); err != nil {
		log.Fatalf("miss: %v", err)
	}
	tcmd.Wait()
	log.Println("-TMP COMMAND END-")
}
*/

// TODO:1コマンド複数チェック・複数reコマンドjsonで設定できないか
func Command(c *Config) {
	for _, v := range c.Commands {
		log.Printf("-%s COMMAND START-\n", v.Command)

		var (
			cmd       = exec.Command("bash", "-c", v.Command)
			stdout, _ = cmd.StdoutPipe()
		)
		if err := cmd.Start(); err != nil {
			log.Fatalf("failed to command '%s': %v\n", v.Command, err)
		}

		var (
			sc         = bufio.NewScanner(stdout)
			shouldReCommand = true
		)
		for sc.Scan() {
			s := sc.Text()
			if strings.Contains(s, v.Check) || v.Check == "" {
				shouldReCommand = false
				break
			}
		}
		cmd.Wait()
		log.Printf("-%s COMMAND END-\n", v.Command)

		if shouldReCommand {
			log.Printf("-%s ReCOMMAND START-\n", v.ReCommandConfig.ReCommand)

			reCmd := exec.Command("bash", "-c", v.ReCommandConfig.ReCommand)			
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

	shouldReConnect := false
	vpnSc := bufio.NewScanner(vpnStdout)
	log.Println("-VPN COMMAND START-")
	for vpnSc.Scan() {
		s := vpnSc.Text()
		switch {
		case strings.Contains(s, "3. VPN Tools コマンドの使用 (証明書作成や通信速度測定)"):
			vpnStdin.Write([]byte("2\n"))
			continue

		case strings.Contains(s, "何も入力せずに Enter を押すと、localhost (このコンピュータ) に接続します。"):
			vpnStdin.Write([]byte("\n"))
			continue

		case strings.Contains(s, `VPN Client "localhost" に接続しました。`):
			vpnStdin.Write([]byte(fmt.Sprintf("AccountStatusGet %s\n", "MYIPSE")))
			continue

		case strings.Contains(s, "指定された接続設定は接続されていません。"):
			vpnStdin.Write([]byte(fmt.Sprintf("AccountConnect %s\n", "MYIPSE")))
			continue
		case strings.Contains(s, "指定された接続設定は現在接続中です。"):
			vpnStdin.Write([]byte(fmt.Sprintf("AccountDisconnect %s\n", "MYIPSE")))
			shouldReConnect = true
			continue
		case strings.Contains(s, "セッション接続状態") && strings.Contains(s, "接続完了 (セッション確立済み)"):
			vpnStdin.Write([]byte("QUIT\n"))
			break

		case strings.Contains(s, "コマンドは正常に終了しました。"):
			if shouldReConnect {
				vpnStdin.Write([]byte(fmt.Sprintf("AccountConnect %s\n", "MYIPSE")))
				shouldReConnect = false
				continue
			}
			vpnStdin.Write([]byte("QUIT\n"))
			break
		}
	}
	vpnCmd.Wait()
	log.Println("-VPN COMMAND END-")
}
