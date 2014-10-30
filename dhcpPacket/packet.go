package dhcpPacket

import (
	"bytes"
	"encoding/binary"
	"errors"
	"fmt"
	"net"
)

type DhcpPacket struct {
	op      byte
	htype   byte
	hlen    byte
	hops    byte
	xid     [4]byte
	secs    [2]byte
	flags   Flags
	ciaddr  [4]byte
	yiaddr  [4]byte
	siaddr  [4]byte
	giaddr  [4]byte
	chaddr  [16]byte
	sname   [64]byte
	file    [128]byte
	Options Options
}

type Flags struct {
	broadcast bool
}

// Returns the byte array forming the packet
func (d DhcpPacket) ToBytes() (bytePacket []byte) {
	bytePacket = append(bytePacket, d.op, d.htype, d.hlen, d.hops)
	bytePacket = append(bytePacket, d.xid[:]...)
	bytePacket = append(bytePacket, d.secs[:]...)

	flagsBytes := d.flags.ToBytes()
	bytePacket = append(bytePacket, flagsBytes[:]...)

	bytePacket = append(bytePacket, d.ciaddr[:]...)
	bytePacket = append(bytePacket, d.yiaddr[:]...)
	bytePacket = append(bytePacket, d.siaddr[:]...)
	bytePacket = append(bytePacket, d.giaddr[:]...)
	bytePacket = append(bytePacket, d.chaddr[:]...)
	bytePacket = append(bytePacket, d.sname[:]...)
	bytePacket = append(bytePacket, d.file[:]...)

	magic := [...]byte{0x63, 0x82, 0x53, 0x63}
	bytePacket = append(bytePacket, magic[:]...)

	optionsBytes := d.Options.ToBytes()
	bytePacket = append(bytePacket, optionsBytes[:]...)

	return
}

func NewDhcpPacket() (d *DhcpPacket) {

	d = new(DhcpPacket)
	d.htype = 0x01
	d.hlen = 0x06
	d.hops = 0x00

	return d

}

func (d *DhcpPacket) SetOp(value int) {
	if value == 1 {
		d.op = 0x01
	} else {
		d.op = 0x02
	}
}

func (d *DhcpPacket) GetOp() (value int) {
	if d.op == 0x01 {
		return 1
	} else {
		return 2
	}
}

func (d *DhcpPacket) GetXid() (value uint32) {
	value = binary.LittleEndian.Uint32(d.xid[:])
	return
}

func (d *DhcpPacket) SetXid(value uint32) {
	copy(d.xid[:], intToBytes(value))
}

func (d *DhcpPacket) SetCiaddr(value string) {
	tmp := net.ParseIP(value)
	copy(d.ciaddr[:], tmp[12:])
}

func (d *DhcpPacket) SetYiaddr(value string) {
	tmp := net.ParseIP(value)
	copy(d.yiaddr[:], tmp[12:])
}
func (d *DhcpPacket) SetSiaddr(value string) {
	tmp := net.ParseIP(value)
	copy(d.siaddr[:], tmp[12:])
}
func (d *DhcpPacket) SetGiaddr(value string) {
	tmp := net.ParseIP(value)
	copy(d.giaddr[:], tmp[12:])
}

func (d *DhcpPacket) SetChaddr(value net.HardwareAddr) {
	copy(d.chaddr[:], value)
}

func (d *DhcpPacket) GetChaddr() (address net.HardwareAddr) {
	address = d.chaddr[:]
	return
}
func (d *DhcpPacket) GetCiaddr() (address net.IP) {
	address = net.IPv4(d.ciaddr[0], d.ciaddr[1], d.ciaddr[2], d.ciaddr[3])
	return
}

func (d *DhcpPacket) SetBroadcast(value bool) {
	d.flags.SetBroadcast(value)
}

func (f *Flags) ToBytes() (bytePacket [2]byte) {

	if f.broadcast {
		bytePacket = [2]byte{0x80, 0x00}
		return
	}

	return [2]byte{0x00, 0x00}
}

func (f *Flags) SetBroadcast(value bool) {
	f.broadcast = value
}

func ParseDhcpPacket(b []byte, o *DhcpPacket) (err error) {

	if len(b) < 237 {
		return errors.New("Mauvais packet")
	}

	magic := [...]byte{0x63, 0x82, 0x53, 0x63}
	parsedMagic := [4]byte{}
	copy(parsedMagic[:], b[236:240])

	if parsedMagic != magic {
		return errors.New("Mauvais Magic Cookie")
	}

	o.op = b[0]
	o.htype = b[1]
	o.hlen = b[2]
	o.hops = b[3]
	copy(o.xid[:], b[4:8])
	copy(o.secs[:], b[8:10])
	o.flags = parseFlags(b[10:12])
	copy(o.ciaddr[:], b[12:16])
	copy(o.yiaddr[:], b[16:20])
	copy(o.siaddr[:], b[20:24])
	copy(o.giaddr[:], b[24:28])
	copy(o.chaddr[:], b[28:44])
	copy(o.sname[:], b[44:108])
	copy(o.file[:], b[108:236])

	// We skip the magic Cookie
	o.Options = parseOptions(b[240:])

	return

}

func (d *DhcpPacket) SetMessageType(value int) {

	valueRaw := intToBytes(uint32(value))
	tmp := valueRaw[0]
	d.Options.Add(53, []byte{tmp})
}

func (d *DhcpPacket) SetSubnetMask(value string) {
	tmp := net.ParseIP(value)
	d.Options.Add(1, tmp[12:])
}

func (d *DhcpPacket) SetRouter(value string) {
	tmp := net.ParseIP(value)
	d.Options.Add(3, tmp[12:])
}

func (d *DhcpPacket) SetDhcpServer(value string) {
	tmp := net.ParseIP(value)
	d.Options.Add(54, tmp[12:])
}

func (d *DhcpPacket) SetDnsServer(value []string) {

	var raw []byte

	for _, addr := range value {
		tmp := net.ParseIP(addr)
		raw = append(raw, []byte(tmp[12:])...)
	}

	d.Options.Add(6, raw)
}

func (d *DhcpPacket) SetLeaseTime(value int) {
	tmp := intToBytesBG(uint32(value))
	d.Options.Add(51, tmp)
}

func (d *DhcpPacket) GetMessageType() (msgType int) {
	value, err := d.Options.GetOption(53)
	if err != nil {
		fmt.Println(err)
	}

	msgType = byteToInt(value)

	return
}

func (d *DhcpPacket) GetRequestedIp() (ip net.IP) {
	value, err := d.Options.GetOption(50)
	if err != nil {
		fmt.Println(err)
	}

	ip = net.IP(value)
	return
}

func (d *DhcpPacket) GetHostName() (hostname string) {
	value, err := d.Options.GetOption(12)
	if err != nil {
		fmt.Println(err)
	}

	hostname = string(value)
	return
}

func (d *DhcpPacket) GetMaximumPacketSize() (size int) {
	value, err := d.Options.GetOption(57)
	if err != nil {
		fmt.Println(err)
	}

	size = int(binary.BigEndian.Uint16(value))
	return
}

func (d *DhcpPacket) GetParameterRequestList() (list []int) {
	value, err := d.Options.GetOption(55)
	if err != nil {
		fmt.Println(err)
	}

	buf := bytes.NewBuffer(value)

	for i := 0; i < len(value); i++ {
		codeRaw := make([]byte, 1)

		n, err := buf.Read(codeRaw)
		if err != nil {
			fmt.Println(err)
		}
		if n == 0 {
			fmt.Println("Unexpected End")
		}
		list = append(list, byteToInt(codeRaw))
	}

	return
}

func parseFlags(f []byte) (o Flags) {
	if f[0] == byte(0x80) {
		o.SetBroadcast(true)
	}

	return
}

// Helpers

func intToBytes(value uint32) (output []byte) {

	buf := new(bytes.Buffer)
	err := binary.Write(buf, binary.LittleEndian, value)
	if err != nil {
		fmt.Println("binary.Write failed:", err)
	}

	return buf.Bytes()
}

func intToBytesBG(value uint32) (output []byte) {

	buf := new(bytes.Buffer)
	err := binary.Write(buf, binary.BigEndian, value)
	if err != nil {
		fmt.Println("binary.Write failed:", err)
	}

	return buf.Bytes()
}

func byteToInt(value []byte) (output int) {
	value = append(value, byte(0x00))
	return int(binary.LittleEndian.Uint16(value))
}
