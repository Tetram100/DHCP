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

func main() {

	initDB()

	handler := func(data []byte) {
		pkg := dhcpPacket.NewDhcpPacket()
		dhcpPacket.ParseDhcpPacket(data, pkg)
	}

	dhcpUDP.InitListener(handler)

}
