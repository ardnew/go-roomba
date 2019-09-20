package oibot

import (
	"bytes"
	"encoding/binary"
	"fmt"
	"log"
	"time"

	"github.com/tarm/serial"
)

type Roomba struct {
	port     *serial.Port
	infoLog  *log.Logger
	errorLog *log.Logger
	path     string
	baud     int
}

func MakeRoomba(infoLog *log.Logger, errorLog *log.Logger, path string, baud int) *Roomba {
	var r *Roomba
	if _, ok := codeForBaudRate[baud]; !ok {
		errorLog.Panic(fmt.Errorf("invalid baud rate: %d", baud))
	}
	if port, err := serial.OpenPort(&serial.Config{Name: path, Baud: baud}); nil != err {
		errorLog.Panic(fmt.Errorf("failed to open serial port: %s (%d): %s", path, baud, err))
	} else {
		r = &Roomba{port: port, infoLog: infoLog, errorLog: errorLog, path: path, baud: baud}
	}
	return r
}

func (r *Roomba) Pack(data ...interface{}) []byte {
	buf := new(bytes.Buffer)
	for _, bin := range data {
		if err := binary.Write(buf, binary.BigEndian, bin); nil != err {
			r.errorLog.Panic(fmt.Errorf("failed to pack binary data: ", err))
		}
	}
	return buf.Bytes()
}

func (r *Roomba) WriteCode(code OpCode) {
	if n, err := r.port.Write([]byte{byte(code)}); n != 1 || nil != err {
		r.errorLog.Panic(fmt.Errorf("failed to write opcode (%d) to serial port: %s", code, err))
	}
	time.Sleep(SerialTransferDelayMS)
}

func (r *Roomba) Write(code OpCode, buf ...interface{}) int {
	r.WriteCode(code)
	bin := r.Pack(buf...)
	count := 0
	if n, err := r.port.Write(bin); n != len(bin) || nil != err {
		r.errorLog.Panic(fmt.Errorf("failed to write opcode (%d) data to serial port: %s", code, err))
	} else {
		count = n
	}
	time.Sleep(SerialTransferDelayMS)
	return count
}

func (r *Roomba) Read(buf []byte) int {
	count := 0
	if n, err := r.port.Read(buf); nil != err {
		r.errorLog.Panic(fmt.Errorf("failed to read from serial port: %s", err))
	} else {
		count = n
	}
	time.Sleep(SerialTransferDelayMS)
	return count
}

func (r *Roomba) Sensor(packet *SensorPacket) []byte {
	r.Write(opcQuery, packet.id)
	data := make([]byte, packet.size)
	current, remaining := 0, int(packet.size)
	for current < remaining {
		frame := data[current:]
		remaining -= current
		current = r.Read(frame)
	}
	return data
}

func (r *Roomba) sensorListID(packet ...*SensorPacket) []byte {
	id := make([]byte, len(packet))
	for i, d := range packet {
		id[i] = d.id
	}
	return id
}

func (r *Roomba) SensorList(packet ...*SensorPacket) [][]byte {
	numPackets := byte(len(packet))
	if numPackets > 0 {
		queryList := []byte{numPackets}
		queryList = append(queryList, r.sensorListID(packet...)...)
		r.Write(opcQueryList, queryList)
		data := make([][]byte, numPackets)
		for i, p := range packet {
			data[i] = make([]byte, p.size)
			current, remaining := 0, int(p.size)
			for current < remaining {
				frame := data[i][current:]
				remaining -= current
				current = r.Read(frame)
			}
		}
		return data
	}
	return nil
}

// =====================================================================================================================

func (r *Roomba) Start() {
	r.Write(opcStart)
}

func (r *Roomba) Passive() { // alias for start command
	r.Write(opcStart)
}

func (r *Roomba) Reset() {
	r.Write(opcReset)
}

func (r *Roomba) Stop() {
	r.Write(opcStop)
}

func (r *Roomba) Baud(baud int) {
	if code, ok := codeForBaudRate[baud]; !ok {
		r.errorLog.Panic(fmt.Errorf("will not change to invalid baud rate: %d", baud))
	} else {
		_ = r.Write(opcBaud, code)
		time.Sleep(100 * time.Millisecond)
	}
}

func (r *Roomba) Control() {
	r.Write(opcControl)
}

func (r *Roomba) Safe() {
	r.Write(opcSafe)
}

func (r *Roomba) Full() {
	r.Write(opcFull)
}

func (r *Roomba) Power() {
	r.Write(opcPower)
}

func (r *Roomba) Clean() {
	r.Write(opcClean)
}

func (r *Roomba) MaxClean() {
	r.Write(opcMaxClean)
}

func (r *Roomba) Spot() {
	r.Write(opcSpot)
}

func (r *Roomba) SeekDock() {
	r.Write(opcForceSeekingDock)
}

func (r *Roomba) Drive(velocity int16, radius int16) {
	if velocity < MinDriveVelocityMMPS || velocity > MaxDriveVelocityMMPS {
		r.errorLog.Panic(fmt.Errorf("invalid drive velocity: %d", velocity))
	}
	if radius < MinDriveRadiusMM || radius > MaxDriveRadiusMM {
		r.errorLog.Panic(fmt.Errorf("invalid drive radius: %d", radius))
	}
	_ = r.Write(opcDrive, velocity, radius)
}

func (r *Roomba) DriveStop() {
	r.Drive(0, 0)
}

func (r *Roomba) DriveWheels(rightVelocity int16, leftVelocity int16) {
	if rightVelocity < MinDriveVelocityMMPS || rightVelocity > MaxDriveVelocityMMPS {
		r.errorLog.Panic(fmt.Errorf("invalid right wheel velocity: %d", rightVelocity))
	}
	if leftVelocity < MinDriveVelocityMMPS || leftVelocity > MaxDriveVelocityMMPS {
		r.errorLog.Panic(fmt.Errorf("invalid left wheel velocity: %d", leftVelocity))
	}
	_ = r.Write(opcDriveWheels, rightVelocity, leftVelocity)
}

func (r *Roomba) Mode() OpenInterfaceMode {
	data := r.Sensor(spcOpenInterfaceMode)
	if len(data) > 0 {
		return OpenInterfaceMode(data[0])
	}
	return oimOff
}

func (r *Roomba) Battery() []byte {
	data := make([]byte, 6)
	for i, d := range r.SensorList(
		spcChargingState, spcVoltage, spcCurrent, spcBatteryCharge, spcBatteryCapacity, spcChargerAvailable) {
		data[i] = d[0]
	}
	return data
}
