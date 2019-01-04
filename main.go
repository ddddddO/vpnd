package main

import (
	"flag"
	"time"

	"github.com/ddddddO/vpnd/lib"
)

var (
	config   string
	interval time.Duration
)

/*
  備忘:~/Vpnclient/vpn_memo/memo.txt に簡易接続手順メモしてある
	   ~/Vpnclient/vpn_memo/内にvpn設定後の状態メモ
       ~/Vpnclient/vpn_memo/aft_reboot内にvpn設定後からOS再起動後の状態メモ
*/
func main() {
	flag.StringVar(&config, "config", "./config.json", "config file path")
	flag.DurationVar(&interval, "interval", time.Minute*5, "proc interval")
	flag.Parse()

	t := time.NewTicker(interval)
	c := lib.NewConfig()
	c.Unmarshal(config)

	//lib.TmpCommand() // debug用
	lib.VPNCommand()
	lib.Command(c)

	for {
		select {
		case <-t.C:
			lib.VPNCommand()
			lib.Command(c)
		}
	}
}
