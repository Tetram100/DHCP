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

func (o *Options) Add(code uint32, value []byte) {

	var option Option

	if value != nil {
		option = Option{code, uint32(len(value)), value}
	} else {
		option = Option{code, 0, nil}
	}

	o.options = append(o.options, option)
}

func (o *Option) ToBytes() (bytePacket []byte) {
	code := intToBytes(o.code)
	bytePacket = append(bytePacket, code[0])
	if o.length != 0 {
		length := intToBytes(o.length)
		bytePacket = append(bytePacket, length[0])
	}
	bytePacket = append(bytePacket, o.value[:]...)

	return
}
