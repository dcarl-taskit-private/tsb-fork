package tsb

// I2C represents a connection to I2C-device.
type I2C struct {
	addr uint8
	jack uint8
}

// NewI2C opens a connection for I2C-device.
func NewI2C(addr uint8, jack uint8) (*I2C, error) {
	v := &I2C{jack: jack, addr: addr}
	return v, nil
}

// GetJack return bus line, where I2C-device is allocated.
func (v *I2C) GetJack() uint8 {
	return v.jack
}

// GetAddr return device occupied address in the bus.
func (v *I2C) GetAddr() uint8 {
	return v.addr
}

func (v *I2C) write(buf []byte) (int, error) {
	return 0, nil
}

// WriteBytes send bytes to the remote I2C-device.
func (v *I2C) WriteBytes(buf []byte) (int, error) {
	//	lg.Debugf("Write %d hex bytes: [%+v]", len(buf), hex.EncodeToString(buf))
	return v.write(buf)
}

func (v *I2C) read(buf []byte) (int, error) {
	return 0, nil
}

// ReadBytes read bytes from I2C-device.
// Number of bytes read correspond to buf parameter length.
func (v *I2C) ReadBytes(buf []byte) (int, error) {
	n, err := v.read(buf)
	if err != nil {
		return n, err
	}
	//lg.Debugf("Read %d hex bytes: [%+v]", len(buf), hex.EncodeToString(buf))
	return n, nil
}

// Close I2C-connection.
func (v *I2C) Close() error {
	return nil
}

// ReadRegBytes read count of n byte's sequence from I2C-device
// starting from reg address.
// SMBus (System Management Bus) protocol over I2C.
func (v *I2C) ReadRegBytes(reg byte, n int) ([]byte, int, error) {
	//lg.Debugf("Read %d bytes starting from reg 0x%0X...", n, reg)
	_, err := v.WriteBytes([]byte{reg})
	if err != nil {
		return nil, 0, err
	}
	buf := make([]byte, n)
	c, err := v.ReadBytes(buf)
	if err != nil {
		return nil, 0, err
	}
	return buf, c, nil

}

// ReadRegU8 reads byte from I2C-device register specified in reg.
// SMBus (System Management Bus) protocol over I2C.
func (v *I2C) ReadRegU8(reg byte) (byte, error) {
	_, err := v.WriteBytes([]byte{reg})
	if err != nil {
		return 0, err
	}
	buf := make([]byte, 1)
	_, err = v.ReadBytes(buf)
	if err != nil {
		return 0, err
	}
	//lg.Debugf("Read U8 %d from reg 0x%0X", buf[0], reg)
	return buf[0], nil
}

// WriteRegU8 writes byte to I2C-device register specified in reg.
// SMBus (System Management Bus) protocol over I2C.
func (v *I2C) WriteRegU8(reg byte, value byte) error {
	buf := []byte{reg, value}
	_, err := v.WriteBytes(buf)
	if err != nil {
		return err
	}
	//lg.Debugf("Write U8 %d to reg 0x%0X", value, reg)
	return nil
}

// ReadRegU16BE reads unsigned big endian word (16 bits)
// from I2C-device starting from address specified in reg.
// SMBus (System Management Bus) protocol over I2C.
func (v *I2C) ReadRegU16BE(reg byte) (uint16, error) {
	_, err := v.WriteBytes([]byte{reg})
	if err != nil {
		return 0, err
	}
	buf := make([]byte, 2)
	_, err = v.ReadBytes(buf)
	if err != nil {
		return 0, err
	}
	w := uint16(buf[0])<<8 + uint16(buf[1])
	//lg.Debugf("Read U16 %d from reg 0x%0X", w, reg)
	return w, nil
}

// ReadRegU16LE reads unsigned little endian word (16 bits)
// from I2C-device starting from address specified in reg.
// SMBus (System Management Bus) protocol over I2C.
func (v *I2C) ReadRegU16LE(reg byte) (uint16, error) {
	w, err := v.ReadRegU16BE(reg)
	if err != nil {
		return 0, err
	}
	// exchange bytes
	w = (w&0xFF)<<8 + w>>8
	return w, nil
}

// ReadRegS16BE reads signed big endian word (16 bits)
// from I2C-device starting from address specified in reg.
// SMBus (System Management Bus) protocol over I2C.
func (v *I2C) ReadRegS16BE(reg byte) (int16, error) {
	_, err := v.WriteBytes([]byte{reg})
	if err != nil {
		return 0, err
	}
	buf := make([]byte, 2)
	_, err = v.ReadBytes(buf)
	if err != nil {
		return 0, err
	}
	w := int16(buf[0])<<8 + int16(buf[1])
	//lg.Debugf("Read S16 %d from reg 0x%0X", w, reg)
	return w, nil
}

// ReadRegS16LE reads signed little endian word (16 bits)
// from I2C-device starting from address specified in reg.
// SMBus (System Management Bus) protocol over I2C.
func (v *I2C) ReadRegS16LE(reg byte) (int16, error) {
	w, err := v.ReadRegS16BE(reg)
	if err != nil {
		return 0, err
	}
	// exchange bytes
	w = (w&0xFF)<<8 + w>>8
	return w, nil

}

// WriteRegU16BE writes unsigned big endian word (16 bits)
// value to I2C-device starting from address specified in reg.
// SMBus (System Management Bus) protocol over I2C.
func (v *I2C) WriteRegU16BE(reg byte, value uint16) error {
	buf := []byte{reg, byte((value & 0xFF00) >> 8), byte(value & 0xFF)}
	_, err := v.WriteBytes(buf)
	if err != nil {
		return err
	}
	//lg.Debugf("Write U16 %d to reg 0x%0X", value, reg)
	return nil
}

// WriteRegU16LE writes unsigned little endian word (16 bits)
// value to I2C-device starting from address specified in reg.
// SMBus (System Management Bus) protocol over I2C.
func (v *I2C) WriteRegU16LE(reg byte, value uint16) error {
	w := (value*0xFF00)>>8 + value<<8
	return v.WriteRegU16BE(reg, w)
}

// WriteRegS16BE writes signed big endian word (16 bits)
// value to I2C-device starting from address specified in reg.
// SMBus (System Management Bus) protocol over I2C.
func (v *I2C) WriteRegS16BE(reg byte, value int16) error {
	buf := []byte{reg, byte((uint16(value) & 0xFF00) >> 8), byte(value & 0xFF)}
	_, err := v.WriteBytes(buf)
	if err != nil {
		return err
	}
	//lg.Debugf("Write S16 %d to reg 0x%0X", value, reg)
	return nil
}

// WriteRegS16LE writes signed little endian word (16 bits)
// value to I2C-device starting from address specified in reg.
// SMBus (System Management Bus) protocol over I2C.
func (v *I2C) WriteRegS16LE(reg byte, value int16) error {
	w := int16((uint16(value)*0xFF00)>>8) + value<<8
	return v.WriteRegS16BE(reg, w)
}
