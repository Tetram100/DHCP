package main

import (
	"database/sql"
	"fmt"
	"git.cfr.re/dhcp.git/dhcpPacket"
	"git.cfr.re/dhcp.git/dhcpUDP"
	_ "github.com/mattn/go-sqlite3"
	"log"
	"net"
)

var database *sql.DB

func initDB() {
	database, err := sql.Open("sqlite3", "mysqlite_3")
	if err != nil {
		log.Fatal(err)
	}
	defer database.Close()

	tx, err := database.Begin()
	_, err = database.Exec(
		"CREATE TABLE IF NOT EXISTS IP_table ( id integer PRIMARY KEY, AddressIP varchar(255) NOT NULL)")

	if err != nil {
		fmt.Println("Failed to create the database")
		log.Fatal(err)
	}

	tx.Commit()
}

func response(request *dhcpPacket.DhcpPacket) {

	packet_response := dhcpPacket.NewDhcpPacket()

	// Partie fixe dans le corps
	packet_response.SetOp(2)
	packet_response.SetCiaddr("0.0.0.0")
	packet_response.SetSiaddr("192.168.12.1") // TODO - Récuprer l'IP du serveur
	packet_response.SetGiaddr("0.0.0.0")

	// Partie fixe dans les options
	// TODO - Récupérer à partir de la conf
	packet_response.SetMessageType(2)
	packet_response.SetDhcpServer("192.168.12.1") // TODO - Récupérer l'IP du serveur
	packet_response.SetLeaseTime(43200)
	packet_response.SetSubnetMask("255.255.255.0")
	packet_response.SetDnsServer([]string{"8.8.8.8"})
	packet_response.SetRouter("192.168.12.1")

	// Réponse particulières à la demande du client
	packet_response.SetXid(request.GetXid())
	mac_request := request.GetChaddr()
	packet_response.SetChaddr(mac_request)

	/*
		rows1, err := database.QueryRow("SELECT * FROM IP_table WHERE MAC = ?", (mac_disover,))
		rows2, err := database.QueryRow("SELECT * FROM IP_table WHERE MAC = ?", ("",))
		rows2, err := database.QueryRow("SELECT * FROM IP_table WHERE ( julianday('now') - julianday(allocated_at) > 3600)")
		if rows1 != nil {
			packet_response.dhcpPacket.SetYiaddr(rows1[1])
			// mettre à jour allocated_at
		} else if rows2 !=nil {
			packet_response.dhcpPacket.SetYiaddr(rows2[1])
			// mettre à jour allocated_at
		} else if rows3 !=nil {
			packet_response.dhcpPacket.SetYiaddr(rows3[1])
			// mettre à jour allocated_at
		} else {
			fmt.Println("All IPs are used")
		}
	*/

	// DEBUG ---
	packet_response.SetYiaddr("192.168.12.2")
	packet_response.Options.Add(255, nil)

	// Envoi du packet

	raddr := net.UDPAddr{IP: net.ParseIP("255.255.255.255"), Port: 68}
	laddr, err := net.ResolveUDPAddr("udp", "192.168.12.1:6767")

	conn, err := net.DialUDP("udp", laddr, &raddr)
	if err != nil {
		fmt.Print(err)
	}

	n, err := conn.Write(packet_response.ToBytes())
	if err != nil {
		fmt.Println(err)
	}

	fmt.Println("Offer Send - Bytes written : ", n)
	conn.Close()
}

func ack(discover *dhcpPacket.DhcpPacket) {

	packet_response := dhcpPacket.NewDhcpPacket()

	// Partie fixe dans le corps
	packet_response.SetOp(2)
	packet_response.SetCiaddr("0.0.0.0")
	packet_response.SetSiaddr("192.168.12.1") // TODO - Récuprer l'IP du serveur
	packet_response.SetGiaddr("0.0.0.0")

	// Partie fixe dans les options
	// TODO - Récupérer à partir de la conf
	packet_response.SetMessageType(5)
	packet_response.SetDhcpServer("192.168.12.1") // TODO - Récupérer l'IP du serveur
	packet_response.SetLeaseTime(43200)
	packet_response.SetSubnetMask("255.255.255.0")
	packet_response.SetDnsServer([]string{"8.8.8.8"})
	packet_response.SetRouter("192.168.12.1")

	// Réponse particulières à la demande du client
	packet_response.SetXid(discover.GetXid())
	mac_discover := discover.GetChaddr()
	packet_response.SetChaddr(mac_discover)

	/*
		rows1, err := database.QueryRow("SELECT * FROM IP_table WHERE MAC = ?", (mac_disover,))
		rows2, err := database.QueryRow("SELECT * FROM IP_table WHERE MAC = ?", ("",))
		rows2, err := database.QueryRow("SELECT * FROM IP_table WHERE ( julianday('now') - julianday(allocated_at) > 3600)")
		if rows1 != nil {
			packet_response.dhcpPacket.SetYiaddr(rows1[1])
			// mettre à jour allocated_at
		} else if rows2 !=nil {
			packet_response.dhcpPacket.SetYiaddr(rows2[1])
			// mettre à jour allocated_at
		} else if rows3 !=nil {
			packet_response.dhcpPacket.SetYiaddr(rows3[1])
			// mettre à jour allocated_at
		} else {
			fmt.Println("All IPs are used")
		}
	*/

	// DEBUG ---
	packet_response.SetYiaddr("192.168.12.2")
	packet_response.Options.Add(255, nil)

	// Envoi du packet

	raddr := net.UDPAddr{IP: net.ParseIP("255.255.255.255"), Port: 68}
	laddr, err := net.ResolveUDPAddr("udp", "192.168.12.1:6767")

	conn, err := net.DialUDP("udp", laddr, &raddr)
	if err != nil {
		fmt.Print(err)
	}

	n, err := conn.Write(packet_response.ToBytes())
	if err != nil {
		fmt.Println(err)
	}

	fmt.Println("ACK Send - Bytes written : ", n)
	conn.Close()
}

func main() {

	initDB()

	handler := func(data []byte) {
		pkg := dhcpPacket.NewDhcpPacket()
		dhcpPacket.ParseDhcpPacket(data, pkg)
		if pkg.GetMessageType() == 1 {
			response(pkg)
		}

		if pkg.GetMessageType() == 3 {
			ack(pkg)
		}
	}

	dhcpUDP.InitListener(handler)

}
