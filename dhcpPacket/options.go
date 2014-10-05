package dhcpPacket

type Options struct {
	options []Option
}

type Option struct {
	code   uint32
	length uint32
	value  []byte
}

func (o *Options) ToBytes() (bytePacket []byte) {
	for _, option := range o.options {
		tmp := option.ToBytes()
		bytePacket = append(bytePacket, tmp[:]...)
	}

	return
}

func (o *Options) Add(code uint32, value []byte, length uint32) {
	option := Option{code, uint32(len(value)), value}
	o.options = append(o.options, option)
}

func (o *Option) ToBytes() (bytePacket []byte) {
	code := intToBytes(o.code)
	length := intToBytes(o.length)
	bytePacket = append(bytePacket, code[0], length[0])
	bytePacket = append(bytePacket, o.value[:]...)

	return
}

/*
Options à implémenter
- 53 --> DHCP Type
- 50
- 55
- 1
- 3
- 51
- 54
- 6





*/
