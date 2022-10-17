package tsb

func (s *Server) PortInit(jack byte) (err error) {
	CheckJack(jack)
	s.jack[jack].ReadChan[TypPort] = make(chan byte, 1024)
	return nil
}

func (s *Server) PortPutc(jack byte, c byte) (err error) {
	td := TsbData{Ch: []byte{byte(jack)}, Typ: []byte{TypPort}, Payload: []byte{c}}
	s.tdPutCh <- td
	/*
		fmt.Printf("td.Ch: %d, Typ: %d, Payload: %x\n", td.Ch[0], td.Typ[0], td.Payload)
		encoded := Encode(td)
		fmt.Printf("encode: %x\n", encoded)
		cobs := CobsEncode(encoded)
		fmt.Printf("cobs: %x\n", cobs)
		packet := CobsDecode(cobs)
		fmt.Printf("packet: %x\n", packet)
		td2, _ := Decode(packet)
		fmt.Printf("td2.Ch: %d, Typ: %d, Payload: %x\n", td2.Ch[0], td2.Typ[0], td2.Payload)
	*/
	return nil
}
