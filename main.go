package main

import (
	"database/sql"
	"fmt"
	"git.cfr.re/dhcp.git/dhcpPacket"
	"git.cfr.re/dhcp.git/dhcpUDP"
	_ "github.com/mattn/go-sqlite3"
	"log"
	"net"
	"strconv"
	"time"
)

var database *sql.DB
var allocation_time int = 3600             // must be in second
var Our_Network string = "192.168.12.0/24" // With CIDR notation
var IP_server string = "192.168.12.1"
var IP_DNS string = "8.8.8.8"

type sqlRow struct {
	Id           int
	IP           string
	MAC          string
	Release_date time.Time
}

func create_ip(Network_cidr string) (list []string) {
	ip, ipnet, err := net.ParseCIDR(Network_cidr)
	if err != nil {
		log.Fatal(err)
	}
	for ip := ip.Mask(ipnet.Mask); ipnet.Contains(ip); inc(ip) {
		list = append(list, ip.String())
	}

	return
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

	_, err := database.Exec("DROP TABLE IF EXISTS IP_table")

	_, err := database.Exec(
		"CREATE TABLE IP_table (id integer PRIMARY KEY, AddressIP varchar(255) NOT NULL, MAC varchar(255), release_date TIMESTAMP DEFAULT CURRENT_TIMESTAMP);")
	if err != nil {
		log.Fatal(err)
	}

	IPs := create_ip(Our_Network)
	// We delete the network and broadcast address
	IPs = append(IPs[:0], IPs[1:]...)
	IPs = append(IPs[:(len(IPs) - 1)])
	for _, ip := range IPs {
		if ip != IP_server {
			_, err = database.Exec(
				"INSERT INTO IP_table (id, AddressIP, MAC) VALUES (?, ?, ?)", nil, ip, "")
		}
	}

	if err != nil {
		fmt.Println("Failed to create the database")
		log.Fatal(err)
	}

}

func response(request *dhcpPacket.DhcpPacket) {

	var err error

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
	packet_response.SetLeaseTime(allocation_time)
	packet_response.SetSubnetMask("255.255.255.0")
	packet_response.SetDnsServer([]string{IP_DNS})
	packet_response.SetRouter(IP_server)

	// Réponse particulières à la demande du client
	packet_response.SetXid(request.GetXid())
	mac_request := request.GetChaddr()
	packet_response.SetChaddr(mac_request)

	var row1 sqlRow

	stmt, err := database.Prepare("SELECT * FROM IP_table WHERE MAC = ?")
	if err != nil {
		log.Fatal(err)
	}

	err1 := stmt.QueryRow(mac_request.String()).Scan(&row1.Id, &row1.IP, &row1.MAC, &row1.Release_date)
	if err1 != nil {
		fmt.Println("Erreur lors de la requête 1")
		fmt.Println(err1)
	}

	stmt, err = database.Prepare("SELECT * FROM IP_table WHERE MAC = ?")
	if err != nil {
		log.Fatal(err)
	}

	var row2 sqlRow
	err2 := stmt.QueryRow("").Scan(&row2.Id, &row2.IP, &row2.MAC, &row2.Release_date)
	if err2 != nil {
		fmt.Println("Erreur lors de la requête 2")
		fmt.Println(err2)
	}

	stmt, err = database.Prepare("SELECT * FROM IP_table WHERE ( CURRENT_TIMESTAMP > release_date )")
	if err != nil {
		log.Fatal(err)
	}

	var row3 sqlRow
	err3 := stmt.QueryRow().Scan(&row3.Id, &row3.IP, &row3.MAC, &row3.Release_date)
	if err3 != nil {
		fmt.Println("Erreur lors de la requête 3")
		fmt.Println(err3)
	}

	// We allocate the IP for 3 minutes to give the time to receive a dhcp request
	if err1 == nil {
		packet_response.SetYiaddr(row1.IP)
		_, err = database.Exec(
			"UPDATE IP_table SET release_date=datetime(CURRENT_TIMESTAMP, '+3 minutes') WHERE id = ?", row1.Id)
	} else if err2 == nil {
		packet_response.SetYiaddr(row2.IP)
		_, err = database.Exec(
			"UPDATE IP_table SET release_date=datetime(CURRENT_TIMESTAMP, '+3 minutes'), MAC = ? WHERE id = ?", mac_request.String(), row2.Id)
	} else if err3 == nil {
		packet_response.SetYiaddr(row3.IP)
		_, err = database.Exec(
			"UPDATE IP_table SET release_date=datetime(CURRENT_TIMESTAMP, '+3 minutes'), MAC = ? WHERE id = ?", mac_request.String(), row3.Id)
	} else {
		fmt.Println("All IPs are used")
		packet_response.MessageType(6)
		packet_response.SetYiaddr("0.0.0.0")
		packet_response.SetSiaddr("0.0.0.0")
	}
	if err != nil {
		fmt.Println("Failed to check the database")
		log.Fatal(err)
	}
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

	if packet_response.GetMessageType() == 5 {
		fmt.Println("OFFER Sent - Bytes written : ", n)
	} else {
		fmt.Println("NACK Sent - Bytes written : ", n)
	}
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
	packet_response.SetLeaseTime(allocation_time)
	packet_response.SetSubnetMask("255.255.255.0")
	packet_response.SetDnsServer([]string{IP_DNS})
	packet_response.SetRouter(IP_server)

	// Réponse particulières à la demande du client
	packet_response.SetXid(discover.GetXid())
	mac_request := discover.GetChaddr()
	packet_response.SetChaddr(mac_request)

	var row sqlRow

	stmt, err := database.Prepare("SELECT * FROM IP_table WHERE MAC = ?")
	if err != nil {
		fmt.Println(err)
	}

	err = stmt.QueryRow(mac_request.String()).Scan(&row.Id, &row.IP, &row.MAC, &row.Release_date)
	if err != nil {

		if err == sql.ErrNoRows {
			packet_response.SetMessageType(6)
			packet_response.SetYiaddr("0.0.0.0")
			packet_response.SetSiaddr("0.0.0.0")
		} else {
			fmt.Println("Failed to check the database")
			fmt.Println(err)
		}

	} else {
		// We allocate the IP for allocation_time
		packet_response.SetYiaddr(row.IP)

		timeModifier := "+" + strconv.Itoa(allocation_time) + " seconds"

		_, err = database.Exec(
			"UPDATE IP_table SET release_date=datetime(CURRENT_TIMESTAMP, ?) WHERE id = ?", timeModifier, 1)
		if err != nil {
			fmt.Println(err)
		}
	}

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

	if packet_response.GetMessageType() == 5 {
		fmt.Println("ACK Sent - Bytes written : ", n)
	} else {
		fmt.Println("NACK Sent - Bytes written : ", n)
	}
	conn.Close()
}

func release(discover *dhcpPacket.DhcpPacket) {
	mac_request := discover.GetChaddr()

	var row sqlRow
	stmt, err := database.Prepare("SELECT * FROM IP_table WHERE MAC = ?")
	if err != nil {
		fmt.Println(err)
	}

	err = stmt.QueryRow(mac_request.String()).Scan(&row.Id, &row.IP, &row.MAC, &row.Release_date)
	if err != nil {
		fmt.Println("Illegal DHCPRELEASE")
		return
	} else {
		_, err = database.Exec(
			"UPDATE IP_table SET release_date=CURRENT_TIMESTAMP WHERE id = ?", 1)
		if err != nil {
			fmt.Println(err)
		}
	}

}

func main() {

	var err error

	database, err = sql.Open("sqlite3", "mysqlite_3")
	if err != nil {
		log.Fatal(err)
	}
	defer database.Close()

	initDB()

	handler := func(data []byte) {
		pkg := dhcpPacket.NewDhcpPacket()
		dhcpPacket.ParseDhcpPacket(data, pkg)
		if pkg.GetMessageType() == 1 {
			response(pkg)
		} else if pkg.GetMessageType() == 3 {
			ack(pkg)
		} else if pkg.GetMessageType() == 7 {
			release(pkg)
		}
	}

	dhcpUDP.InitListener(handler)

}
