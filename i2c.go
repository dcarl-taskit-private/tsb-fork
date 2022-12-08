package tsb

import (
	"fmt"
	"time"
)

// I2C represents a connection to I2C-device.
type I2C struct {
	addr   uint8
	jack   byte
	server Server
}

// NewI2C opens a connection for I2C-device.
func NewI2c(addr uint8, jack byte, server Server) (*I2C, error) {
	CheckJack(jack)
	server.I2cSetAdr(jack, JackModeReg)
	server.I2cWrite(jack, []byte{JackI2c})
	server.I2cSetAdr(jack, addr)
	i2c := &I2C{server: server, addr: addr, jack: jack}
	return i2c, nil
}

// GetAddr return device occupied address in the bus.
func (i2c *I2C) GetAddr() uint8 {
	return i2c.addr
}

func (i2c *I2C) write(buf []byte) (int, error) {
	return i2c.server.I2cWrite(i2c.jack, buf)
}

// WriteBytes send bytes to the remote I2C-device.
func (i2c *I2C) WriteBytes(buf []byte) (int, error) {
	return i2c.write(buf)
}

func (i2c *I2C) read(buf []byte) (int, error) {
	return i2c.server.I2cRead(i2c.jack, buf)
}

// ReadBytes read bytes from I2C-device.
// Number of bytes read correspond to buf parameter length.
func (i2c *I2C) ReadBytes(buf []byte) (int, error) {
	n, err := i2c.read(buf)
	if err != nil {
		return n, err
	}
	return n, nil
}

// Close I2C-connection.
func (i2c *I2C) Close() error {
	return nil
}

// ReadRegBytes read count of n byte's sequence from I2C-device
// starting from reg address.
// SMBus (System Management Bus) protocol over I2C.
func (i2c *I2C) ReadRegBytes(reg byte, n int) ([]byte, int, error) {
	_, err := i2c.WriteBytes([]byte{reg})
	if err != nil {
		return nil, 0, err
	}
	buf := make([]byte, n)
	c, err := i2c.ReadBytes(buf)
	if err != nil {
		return nil, 0, err
	}
	return buf, c, nil
}

// ReadRegU8 reads byte from I2C-device register specified in reg.
// SMBus (System Management Bus) protocol over I2C.
func (i2c *I2C) ReadRegU8(reg byte) (byte, error) {
	_, err := i2c.WriteBytes([]byte{reg})
	if err != nil {
		return 0, err
	}
	buf := make([]byte, 1)
	_, err = i2c.ReadBytes(buf)
	if err != nil {
		return 0, err
	}
	//Debugf("Read U8 %d from reg 0x%0X", buf[0], reg)
	return buf[0], nil
}

// WriteRegU8 writes byte to I2C-device register specified in reg.
// SMBus (System Management Bus) protocol over I2C.
func (i2c *I2C) WriteRegU8(reg byte, value byte) error {
	buf := []byte{reg, value}
	_, err := i2c.WriteBytes(buf)
	if err != nil {
		return err
	}
	return nil
}

// ReadRegU16BE reads unsigned big endian word (16 bits)
// from I2C-device starting from address specified in reg.
// SMBus (System Management Bus) protocol over I2C.
func (i2c *I2C) ReadRegU16BE(reg byte) (uint16, error) {
	_, err := i2c.WriteBytes([]byte{reg})
	if err != nil {
		return 0, err
	}
	buf := make([]byte, 2)
	_, err = i2c.ReadBytes(buf)
	if err != nil {
		return 0, err
	}
	w := uint16(buf[0])<<8 + uint16(buf[1])
	return w, nil
}

// ReadRegU16LE reads unsigned little endian word (16 bits)
// from I2C-device starting from address specified in reg.
// SMBus (System Management Bus) protocol over I2C.
func (i2c *I2C) ReadRegU16LE(reg byte) (uint16, error) {
	w, err := i2c.ReadRegU16BE(reg)
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
func (i2c *I2C) ReadRegS16BE(reg byte) (int16, error) {
	_, err := i2c.WriteBytes([]byte{reg})
	if err != nil {
		return 0, err
	}
	buf := make([]byte, 2)
	_, err = i2c.ReadBytes(buf)
	if err != nil {
		return 0, err
	}
	w := int16(buf[0])<<8 + int16(buf[1])
	return w, nil
}

// ReadRegS16LE reads signed little endian word (16 bits)
// from I2C-device starting from address specified in reg.
// SMBus (System Management Bus) protocol over I2C.
func (i2c *I2C) ReadRegS16LE(reg byte) (int16, error) {
	w, err := i2c.ReadRegS16BE(reg)
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
func (i2c *I2C) WriteRegU16BE(reg byte, value uint16) error {
	buf := []byte{reg, byte((value & 0xFF00) >> 8), byte(value & 0xFF)}
	_, err := i2c.WriteBytes(buf)
	if err != nil {
		return err
	}
	return nil
}

// WriteRegU16LE writes unsigned little endian word (16 bits)
// value to I2C-device starting from address specified in reg.
// SMBus (System Management Bus) protocol over I2C.
func (i2c *I2C) WriteRegU16LE(reg byte, value uint16) error {
	w := (value*0xFF00)>>8 + value<<8
	return i2c.WriteRegU16BE(reg, w)
}

// WriteRegS16BE writes signed big endian word (16 bits)
// value to I2C-device starting from address specified in reg.
// SMBus (System Management Bus) protocol over I2C.
func (i2c *I2C) WriteRegS16BE(reg byte, value int16) error {
	buf := []byte{reg, byte((uint16(value) & 0xFF00) >> 8), byte(value & 0xFF)}
	_, err := i2c.WriteBytes(buf)
	if err != nil {
		return err
	}
	return nil
}

// WriteRegS16LE writes signed little endian word (16 bits)
// value to I2C-device starting from address specified in reg.
// SMBus (System Management Bus) protocol over I2C.
func (i2c *I2C) WriteRegS16LE(reg byte, value int16) error {
	w := int16((uint16(value)*0xFF00)>>8) + value<<8
	return i2c.WriteRegS16BE(reg, w)
}

func (s *Server) I2cSetAdr(jack byte, adr byte) (err error) {
	w := make([]byte, 2)
	w[0] = 0x80
	w[1] = adr
	td := TsbData{Ch: []byte{byte(jack)}, Typ: []byte{TypI2c}, Payload: w}
	s.tdPutCh <- td
	select {
	case a := <-s.jack[jack].ReadChan[TypI2c]:
		if a == 1 {
			return nil
		} else {
			return fmt.Errorf("wrong response: %x", a)
		}
	case <-time.After(1 * time.Second):
		return fmt.Errorf("timeout")
	}
}

func (s *Server) I2cRead(jack byte, b []byte) (n int, err error) {
	w := make([]byte, 1)
	n = len(b)
	if n > 127 {
		return 0, fmt.Errorf("only 127 bytes to write are allowed")
	}
	w[0] = byte(n)
	td := TsbData{Ch: []byte{byte(jack)}, Typ: []byte{TypI2c}, Payload: w}
	s.tdPutCh <- td
	for i := 0; i < n; i++ {
		select {
		case b[i] = <-s.jack[jack].ReadChan[TypI2c]:
		case <-time.After(1 * time.Second):
			return 0, fmt.Errorf("timeout")
		}
	}
	return n, nil
}

func (s *Server) I2cWrite(jack byte, b []byte) (n int, err error) {
	n = len(b)
	if n > 127 {
		return 0, fmt.Errorf("only 127 bytes to write are allowed")
	}
	w := make([]byte, n+1)
	w[0] = byte(n + 128)
	for i := 0; i < n; i++ {
		w[i+1] = b[i]
	}
	td := TsbData{Ch: []byte{byte(jack)}, Typ: []byte{TypI2c}, Payload: w[:n+1]}
	s.tdPutCh <- td
	select {
	case a := <-s.jack[jack].ReadChan[TypI2c]:
		if a == byte(n) {
			return 0, nil
		} else {
			return 0, fmt.Errorf("wrong response: %x", a)
		}
	case <-time.After(1 * time.Second):
		return 0, fmt.Errorf("timeout")
	}
}
