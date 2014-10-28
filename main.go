package main

import (
	"database/sql"
	"fmt"
	"git.cfr.re/dhcp.git/dhcpPacket"
	"git.cfr.re/dhcp.git/dhcpUDP"
	"github.com/mattn/go-sqlite3"
	"log"
	"net"
)

var database *sql.DB
var allocation_time int = 3600 // must be in second
var Our_Network string = "192.168.1.0/24" // With CIDR notation
var IP_server string = "192.168.1.1"

func create_ip(Network_cidr string) (list []string) {
	list := []string{}
	ip, ipnet, err := net.ParseCIDR(Network_cidr)
	if err != nil {
		log.Fatal(err)
	}
	for ip := ip.Mask(ipnet.Mask); ipnet.Contains(ip); inc(ip) {
		list = Append(list, ip)
	}
}

func inc(ip net.IP) {
	for j := len(ip) - 1; j >= 0; j-- {
		ip[j]++
		if ip[j] > 0 {
			break
		}
	}
}

func initDB() {
	database, err := sql.Open("sqlite3", "mysqlite_3")
	if err != nil {
		log.Fatal(err)
	}
	defer database.Close()

	tx, err := database.Begin()
	_, err = database.Exec(
		"CREATE TABLE IF NOT EXISTS IP_table ( id integer PRIMARY KEY, AddressIP varchar(255) NOT NULL, MAC varchar(255), release_date TIMESTAMP DEFAULT NOW())")

	IPs := create_ip(Our_Network)
	for ip := IPs.Front(); ip != nil; ip = ip.Next() {
		if ip != IP_server {
			_, err = database.Exec(
				"INSERT INTO IP_table (id, AddressIP, MAC, release_date) VALUES (?, ?, ?, ?)", nil, ip, "", nil)
		}
	}

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
	packet_response.SetSiaddr(IP_server)
	packet_response.SetGiaddr("0.0.0.0")

	// Partie fixe dans les options
	// TODO - Récupérer à partir de la conf
	packet_response.SetMessageType(2)
	packet_response.SetDhcpServer(IP_server)
	packet_response.SetLeaseTime(43200)
	packet_response.SetSubnetMask("255.255.255.0")
	packet_response.SetDnsServer([]string{"8.8.8.8"})
	packet_response.SetRouter("192.168.12.1")

	// Réponse particulières à la demande du client
	packet_response.SetXid(request.GetXid())
	mac_request := request.GetChaddr()
	packet_response.SetChaddr(mac_request)

	tx, err := database.Begin()
	rows1, err := database.QueryRow("SELECT * FROM IP_table WHERE MAC = ?", mac_request)
	rows2, err := database.QueryRow("SELECT * FROM IP_table WHERE MAC = ?", "")
	rows2, err := database.QueryRow("SELECT * FROM IP_table WHERE ( NOW() > release_date )")
	// We allocate the IP for 3 minutes to give the time to receive a dhcp request
	if rows1 != nil {
		packet_response.SetYiaddr(rows1[1])
		_, err = database.Exec(
			"UPDATE IP_table SET release_date=(NOW()+180) WHERE id = ?", rows1[0])
	} else if rows2 !=nil {
		packet_response.SetYiaddr(rows2[1])
		_, err = database.Exec(
			"UPDATE IP_table SET release_date=(NOW()+180), MAC = ? WHERE id = ?", (mac_request, rows2[0]))
	} else if rows3 !=nil {
		packet_response.SetYiaddr(rows3[1])
		_, err = database.Exec(
			"UPDATE IP_table SET release_date=(NOW()+180), MAC = ? WHERE id = ?", (mac_request, rows3[0]))
	} else {
		fmt.Println("All IPs are used")
	}
	if err != nil {
		fmt.Println("Failed to check the database")
		log.Fatal(err)
	}
	tx.Commit()
	packet_response.Options.Add(255, nil)

	// DEBUG ---
	// packet_response.SetYiaddr("192.168.12.2")
	// packet_response.Options.Add(255, nil)

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
	packet_response.SetSiaddr(IP_server)
	packet_response.SetGiaddr("0.0.0.0")

	// Partie fixe dans les options
	// TODO - Récupérer à partir de la conf
	packet_response.SetMessageType(5)
	packet_response.SetDhcpServer(IP_server)
	packet_response.SetLeaseTime(43200)
	packet_response.SetSubnetMask("255.255.255.0")
	packet_response.SetDnsServer([]string{"8.8.8.8"})
	packet_response.SetRouter("192.168.12.1")

	// Réponse particulières à la demande du client
	packet_response.SetXid(discover.GetXid())
	mac_discover := discover.GetChaddr()
	packet_response.SetChaddr(mac_discover)

	tx, err := database.Begin()
	rows1, err := database.QueryRow("SELECT * FROM IP_table WHERE MAC = ?", mac_request)
	// We allocate the IP for allocation_time
	if rows1 != nil {
		packet_response.SetYiaddr(rows1[1])
		_, err = database.Exec(
			"UPDATE IP_table SET release_date=(NOW() + ?) WHERE id = ?", (allocation_time, rows1[0]))
	} else {
		fmt.Println("He has wait too long and needs a dhcp discover")
	}
	if err != nil {
		fmt.Println("Failed to check the database")
		log.Fatal(err)
	}
	tx.Commit()
	packet_response.Options.Add(255, nil)

	// DEBUG ---
	// packet_response.SetYiaddr("192.168.12.2")
	// packet_response.Options.Add(255, nil)

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
