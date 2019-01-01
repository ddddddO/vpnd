package main

import (
	"flag"

	"github.com/ddddddO/vpnd/lib"
)

var (
	config string
)

/*
  備忘:~/Vpnclient/vpn_memo/memo.txt に簡易接続手順メモしてある
	   ~/Vpnclient/vpn_memo/内にvpn設定後の状態メモ
       ~/Vpnclient/vpn_memo/aft_reboot内にvpn設定後からOS再起動後の状態メモ
*/
func main() {
	flag.StringVar(&config, "config", "./config.json", "config file path")
	flag.Parse()

	c := lib.NewConfig()
	c.Unmarshal(config)

	lib.VPNCommand()
	lib.TmpCommand()
	lib.Command(c)	
}
