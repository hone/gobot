package r2q5

import (
	"bytes"
	"fmt"
	"log"
	"sync"
	"time"

	"gobot.io/x/gobot"
	"gobot.io/x/gobot/platforms/ble"
)

type Driver struct {
	name              string
	connection        gobot.Connection
	seq               uint8
	mtx               sync.Mutex
	collisionResponse []uint8
	packetChannel     chan *packet
	gobot.Eventer
}

// Service: 1800
//   Characteristic: 2a00 VHandle(0x03)
//   Characteristic: 2a01 VHandle(0x05)
//   Characteristic: 2a04 VHandle(0x07)
//
// Service: 1801
//   Characteristic: 2a05 VHandle(0x0A)
//
// Service: 00020001574f4f2053706865726f2121
//   Characteristic: 00020003574f4f2053706865726f2121 VHandle(0x0E)
//   Characteristic: 00020002574f4f2053706865726f2121 VHandle(0x10)
//   Characteristic: 00020004574f4f2053706865726f2121 VHandle(0x13)
//   Characteristic: 00020005574f4f2053706865726f2121 VHandle(0x15)
//
// Service: 180f
//   Characteristic: 2a19 VHandle(0x18)
//
// Service: 00010001574f4f2053706865726f2121
//   Characteristic: 00010002574f4f2053706865726f2121 VHandle(0x1C)
//   Characteristic: 00010003574f4f2053706865726f2121 VHandle(0x1F)

const (
	// BLE characteristic IDs
	wakeCharacteristic     = ""
	antiDosCharacteristic  = "00020005574f4f2053706865726f2121"
	commandsCharacteristic = "00010002574f4f2053706865726f2121"
	responseCharacteristic = "00020004574f4f2053706865726f2121"

	// Error event
	Error = "error"
)

type packet struct {
	header   []uint8
	body     []uint8
	checksum uint8
	footer   []uint8
}

func NewDriver(a ble.BLEConnector) *Driver {
	n := &Driver{
		name:          gobot.DefaultName("R2-Q5"),
		connection:    a,
		Eventer:       gobot.NewEventer(),
		packetChannel: make(chan *packet, 1024),
	}
	return n
}

// Connection returns the connection to this R2-Q5
func (b *Driver) Connection() gobot.Connection { return b.connection }

// Name returns the name for the Driver
func (b *Driver) Name() string { return b.name }

// SetName sets the Name for the Driver
func (b *Driver) SetName(n string) { b.name = n }

// adaptor returns BLE adaptor
func (b *Driver) adaptor() ble.BLEConnector {
	return b.Connection().(ble.BLEConnector)
}

// Start tells driver to get ready to do work
func (b *Driver) Start() (err error) {
	b.Init()

	// send commands
	go func() {
		for {
			packet := <-b.packetChannel
			err := b.write(packet)
			if err != nil {
				b.Publish(b.Event(Error), err)
			}
		}
	}()

	go func() {
		for {
			b.adaptor().ReadCharacteristic(responseCharacteristic)
			time.Sleep(100 * time.Millisecond)
		}
	}()

	return
}

// Halt stops R2-Q5 driver (void)
func (b *Driver) Halt() (err error) {
	b.Sleep()
	time.Sleep(750 * time.Microsecond)
	return
}

// Init is used to initialize the R2-Q5
func (b *Driver) Init() (err error) {
	b.AntiDOSOff()
	b.Wake()

	return
}

func (b *Driver) AntiDOSOff() (err error) {
	log.Print("Use the Force!")
	str := "usetheforce...band"
	buf := &bytes.Buffer{}
	buf.WriteString(str)

	err = b.adaptor().WriteCharacteristic(antiDosCharacteristic, buf.Bytes())
	if err != nil {
		fmt.Println("AntiDOSOff error:", err)
		return err
	}

	return
}

func (b *Driver) Wake() (err error) {
	log.Print("Wake!")
	b.packetChannel <- b.craftPacket(0x13, 0x0d, []uint8{})
	return
}

func (b *Driver) Macro(num uint8) {
	b.packetChannel <- b.craftPacket(0x17, 0x05, []uint8{0x00, num})
}

func (b *Driver) Sleep() {
	log.Print("Sleep...")
	b.packetChannel <- b.craftPacket(0x13, 0x01, []uint8{})
}

func (b *Driver) write(packet *packet) (err error) {
	buf := append(packet.header, packet.body...)
	buf = append(buf, packet.checksum)
	buf = append(buf, packet.footer...)
	err = b.adaptor().WriteCharacteristic(commandsCharacteristic, buf)
	if err != nil {
		fmt.Println("send command error:", err)
		return err
	}

	b.mtx.Lock()
	defer b.mtx.Unlock()
	b.seq++
	return
}

func (b *Driver) craftPacket(did byte, cid byte, body []uint8) *packet {
	b.mtx.Lock()
	defer b.mtx.Unlock()
	packet := new(packet)
	packet.header = []uint8{0x8d, 0x0a, did, cid, b.seq}
	packet.body = body
	packet.checksum = b.calculateChecksum(packet)
	packet.footer = []uint8{0xd8}
	return packet
}

func (b *Driver) calculateChecksum(packet *packet) uint8 {
	buf := append(packet.header, packet.body...)
	return calculateChecksum(buf[1:])
}

func calculateChecksum(buf []byte) byte {
	var calculatedChecksum uint16
	for i := range buf {
		calculatedChecksum += uint16(buf[i])
	}
	return uint8(^(calculatedChecksum % 256))
}
