package main

import (
	"encoding/hex"
	"fmt"
	"git.cfr.re/dhcp.git/dhcpPacket"
	//	"git.cfr.re/dhcp.git/dhcpUDP"
	"github.com/mattn/go-sqlite3"
	"net"
)

var DB_DRIVER string
sql.Register(DB_DRIVER, &sqlite3.SQLiteDriver{})
database, err := sql.Open(DB_DRIVER, "mysqlite_3")

tx, err := database.Begin()
result, err := database.Exec(
 "CREATE TABLE IF NOT EXISTS IP_table ( id integer PRIMARY KEY, AddressIP varchar(255) NOT NULL, F",)

if err != nil {
 fmt.Println("Failed to create the database")
}

tx.Commit()

func main() {

	handler := func(data []byte) {
		pkg := dhcpPacket.NewDhcpPacket()
		dhcpPacket.ParseDhcpPacket(data, pkg)
		if 
	}

	dhcpUDP.InitListener(handler)

}
