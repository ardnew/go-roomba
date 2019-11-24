package oibot

import (
	"bytes"
	"encoding/binary"
	"fmt"
	"log"
	"time"

	"github.com/tarm/serial"
)

type OIBot struct {
	port     *serial.Port
	infoLog  *log.Logger
	errorLog *log.Logger
	path     string
	baud     int
	timeout  time.Duration
}

func MakeOIBot(infoLog *log.Logger, errorLog *log.Logger, init bool, path string, baud int, rtime time.Duration) *OIBot {
	var o *OIBot
	if _, ok := codeForBaudRate[baud]; !ok {
		errorLog.Panic(fmt.Errorf("invalid baud rate: %d", baud))
	}
	config := &serial.Config{Name: path, Baud: baud}
	if rtime > NeverReadTimeoutMS {
		config.ReadTimeout = rtime
	}
	if port, err := serial.OpenPort(config); nil != err {
		errorLog.Panic(fmt.Errorf("failed to open serial port: %s (%d): %s", path, baud, err))
	} else {
		o = &OIBot{port: port, infoLog: infoLog, errorLog: errorLog, path: path, baud: baud, timeout: rtime}
		if init {
			o.Baud(baud)
		}
		o.Passive()
	}
	return o
}

func (o *OIBot) Flush() {
	if err := o.port.Flush(); nil != err {
		o.errorLog.Panic(fmt.Errorf("failed to flush serial port: %s", err))
	}
}

func (o *OIBot) Close() {
	if err := o.port.Close(); nil != err {
		o.errorLog.Panic(fmt.Errorf("failed to close serial port: %s", err))
	}
}

func (o *OIBot) Pack(data ...interface{}) []byte {
	buf := new(bytes.Buffer)
	for _, bin := range data {
		if err := binary.Write(buf, binary.BigEndian, bin); nil != err {
			o.errorLog.Panic(fmt.Errorf("failed to pack binary data: ", err))
		}
	}
	return buf.Bytes()
}

func (o *OIBot) WriteCode(code OpCode) {
	if n, err := o.port.Write([]byte{byte(code)}); n != 1 || nil != err {
		o.errorLog.Panic(fmt.Errorf("failed to write opcode (%d) to serial port: %s", code, err))
	}
	time.Sleep(SerialTransferDelayMS)
}

func (o *OIBot) Write(code OpCode, buf ...interface{}) int {
	o.WriteCode(code)
	bin := o.Pack(buf...)
	o.infoLog.Printf("%+v", bin)
	count := 0
	if n, err := o.port.Write(bin); n != len(bin) || nil != err {
		o.errorLog.Panic(fmt.Errorf("failed to write opcode (%d) data to serial port: %s", code, err))
	} else {
		count = n
	}
	time.Sleep(SerialTransferDelayMS)
	return count
}

func (o *OIBot) Read(buf []byte) int {
	count := 0
	if n, err := o.port.Read(buf); nil != err {
		o.errorLog.Panic(fmt.Errorf("failed to read from serial port: %+v", err))
	} else {
		count = n
	}
	return count
}

func (o *OIBot) Sensor(packet *SensorPacket) []byte {
	defer func() {
		if o.timeout > NeverReadTimeoutMS {
			recover()
		}
	}()
	o.Write(opcQuery, packet.id)
	data := make([]byte, packet.size)
	current, remaining := 0, int(packet.size)
	for current < remaining {
		frame := data[current:]
		remaining -= current
		current = o.Read(frame)
	}
	return data
}

func (o *OIBot) sensorListID(packet ...*SensorPacket) []byte {
	id := make([]byte, len(packet))
	for i, d := range packet {
		id[i] = d.id
	}
	return id
}

func (o *OIBot) SensorList(packet ...*SensorPacket) [][]byte {
	defer func() {
		if o.timeout > NeverReadTimeoutMS {
			recover()
		}
	}()
	numPackets := byte(len(packet))
	if numPackets > 0 {
		queryList := []byte{numPackets}
		queryList = append(queryList, o.sensorListID(packet...)...)
		o.Write(opcQueryList, queryList)
		data := make([][]byte, numPackets)
		for i, p := range packet {
			data[i] = make([]byte, p.size)
			current, remaining := 0, int(p.size)
			for current < remaining {
				frame := data[i][current:]
				remaining -= current
				current = o.Read(frame)
			}
		}
		return data
	}
	return nil
}

// =============================================================================

func (o *OIBot) Start() {
	o.Write(opcStart)
}

func (o *OIBot) Passive() { // alias for start command
	o.Write(opcStart)
}

func (o *OIBot) Reset() {
	o.Write(opcReset)
}

func (o *OIBot) Stop() {
	o.Write(opcStop)
}

func (o *OIBot) Baud(baud int) {
	if code, ok := codeForBaudRate[baud]; !ok {
		o.errorLog.Panic(fmt.Errorf("will not change to invalid baud rate: %d", baud))
	} else {
		_ = o.Write(opcBaud, code)
		time.Sleep(100 * time.Millisecond)
	}
}

func (o *OIBot) Control() {
	o.Write(opcControl)
}

func (o *OIBot) Safe() {
	o.Write(opcSafe)
}

func (o *OIBot) Full() {
	o.Write(opcFull)
}

func (o *OIBot) Power() {
	o.Write(opcPower)
}

func (o *OIBot) Clean() {
	o.Write(opcClean)
}

func (o *OIBot) MaxClean() {
	o.Write(opcMaxClean)
}

func (o *OIBot) Spot() {
	o.Write(opcSpot)
}

func (o *OIBot) SeekDock() {
	o.Write(opcForceSeekingDock)
}

func (o *OIBot) Drive(velocity int16, radius int16) {
	if velocity < MinDriveVelocityMMPS || velocity > MaxDriveVelocityMMPS {
		o.errorLog.Panic(fmt.Errorf("invalid drive velocity: %d", velocity))
	}
	if StraightDriveRadiusMM != radius {
		if radius < MinDriveRadiusMM || radius > MaxDriveRadiusMM {
			o.errorLog.Panic(fmt.Errorf("invalid drive radius: %d", radius))
		}
	}
	_ = o.Write(opcDrive, velocity, radius)
}

func (o *OIBot) DriveStop() {
	o.Drive(0, 0)
}

func (o *OIBot) DriveWheels(rightVelocity int16, leftVelocity int16) {
	if rightVelocity < MinDriveVelocityMMPS || rightVelocity > MaxDriveVelocityMMPS {
		o.errorLog.Panic(fmt.Errorf("invalid right wheel velocity: %d", rightVelocity))
	}
	if leftVelocity < MinDriveVelocityMMPS || leftVelocity > MaxDriveVelocityMMPS {
		o.errorLog.Panic(fmt.Errorf("invalid left wheel velocity: %d", leftVelocity))
	}
	_ = o.Write(opcDriveWheels, rightVelocity, leftVelocity)
}

func (o *OIBot) Mode() OpenInterfaceMode {
	data := o.Sensor(spcOpenInterfaceMode)
	if len(data) > 0 {
		return OpenInterfaceMode(data[0])
	}
	return OIMOff
}

// =============================================================================

var (

	// full general info status message
	infoPacket = []*SensorPacket{
		// OI mode
		spcOpenInterfaceMode,
		// battery/charger
		spcChargingState, spcVoltage, spcCurrent, spcBatteryCharge,
		spcBatteryCapacity, spcChargerAvailable,
	}

	// battery-only status message
	batteryPacket = infoPacket[1:]
)

type BatteryStatus struct {
	ChargingState      byte
	VoltagemV          uint16
	CurrentmA          int16
	BatteryChargemAh   uint16
	BatteryCapacitymAh uint16
	ChargerAvailable   byte
}

func batteryStatus(data [][]byte) *BatteryStatus {
	return &BatteryStatus{
		ChargingState:      data[0][0],
		VoltagemV:          uint16((uint16(data[1][0]) << 8) | uint16(data[1][1])),
		CurrentmA:          int16((uint16(data[2][0]) << 8) | uint16(data[2][1])),
		BatteryChargemAh:   uint16((uint16(data[3][0]) << 8) | uint16(data[3][1])),
		BatteryCapacitymAh: uint16((uint16(data[4][0]) << 8) | uint16(data[4][1])),
		ChargerAvailable:   data[5][0],
	}
}

func (o *OIBot) Battery() (*BatteryStatus, bool) {
	data := o.SensorList(batteryPacket...)
	if nil != data && len(data) == len(batteryPacket) {
		return batteryStatus(data), true
	}
	return nil, false
}

type InfoStatus struct {
	Mode    OpenInterfaceMode
	Battery *BatteryStatus
}

func (o *OIBot) Info() (*InfoStatus, bool) {
	data := o.SensorList(infoPacket...)
	if nil != data && len(data) == len(infoPacket) {
		return &InfoStatus{
			Mode:    OpenInterfaceMode(data[0][0]),
			Battery: batteryStatus(data[1:]),
		}, true
	}
	return nil, false
}
