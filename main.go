package main

import (
	"github.com/ddddddO/vpnd/lib"
)
/*
  備忘:~/Vpnclient/vpn_memo/memo.txt に簡易接続手順メモしてある
*/
func main() {
	c := lib.NewConfig()
	c.Unmarshal()

	lib.Command(c)
	lib.VPNCommand() // TODO:ここから	
}
