package main

import (
	"encoding/hex"
	"fmt"
	"git.cfr.re/dhcp.git/dhcpPacket"
	//	"git.cfr.re/dhcp.git/dhcpUDP"
	"net"
)

func main() {

	pack := dhcpPacket.NewDhcpPacket()
	pack.SetBroadcast(true)
	pack.SetOp(2)
	pack.SetXid(12985)
	pack.SetYiaddr("192.168.1.1")
	pack.SetChaddr("82:fb:4a:38:2b:46")

	pack.Options.Add(1, []byte{0xff, 0xff, 0xff, 0x00})
	pack.Options.Add(51, []byte{0xff, 0xff, 0xff, 0x00})
	pack.Options.Add(53, []byte{0x02})

	tmp := net.ParseIP("192.168.1.1")
	addr := tmp[12:]

	pack.Options.Add(3, addr)
	pack.Options.Add(6, addr)
	pack.Options.Add(54, addr)
	pack.Options.Add(255, nil)

	nPack := dhcpPacket.NewDhcpPacket()
	erre := dhcpPacket.ParseDhcpPacket(pack.ToBytes(), nPack)
	if erre != nil {
		fmt.Println(erre)
	}

	nPack.Options.Add(255, nil)

	// Test d'envoi de packet

	raddr := net.UDPAddr{IP: net.ParseIP("255.255.255.255"), Port: 68}

	conn, err := net.DialUDP("udp", nil, &raddr)
	if err != nil {
		fmt.Print(err)
	}

	n, err := conn.Write(pack.ToBytes())
	n, err = conn.Write(nPack.ToBytes())
	if err != nil {
		fmt.Println(err)
	}

	fmt.Println("Written : ", n)

	fmt.Println(nPack.GetXid())
	fmt.Println(hex.EncodeToString(pack.ToBytes()))

	//	handler := func(data []byte) {
	//		pkg := dhcpPacket.NewDhcpPacket()
	//		dhcpPacket.ParseDhcpPacket(data, pkg)
	//		fmt.Println(pkg.GetChaddr().String())
	//	}

	//dhcpUDP.InitListener(handler)

}
