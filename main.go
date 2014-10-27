package main

import (
	"database/sql"
	"fmt"
	"git.cfr.re/dhcp.git/dhcpPacket"
	"git.cfr.re/dhcp.git/dhcpUDP"
	_ "github.com/mattn/go-sqlite3"
	"log"
)

var database *sql.DB
ip_server := ""

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

func response(discover dhcpPacket) {
	packet_response := dhcpPacket.NewDhcpPacket()
	packet_response.dhcpPacket.SetOp(2)
	packet_response.dhcpPacket.SetXid(discover.dhcpPacket.GetXid())
	packet_response.dhcpPacket.SetCiaddr("0.0.0.0")
	mac_disover = discover.dhcpPacket.GetChaddr()
	packet_response.dhcpPacket.SetChaddr(mac_disover)

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

	packet_response.dhcpPacket.SetSiaddr(ip_server)

}

func main() {

	initDB()

	handler := func(data []byte) {
		pkg := dhcpPacket.NewDhcpPacket()
		dhcpPacket.ParseDhcpPacket(data, pkg)
		if option == {
			response(pkg)
		}
	}

	dhcpUDP.InitListener(handler)

}
