package oibot

import (
	"bytes"
	"encoding/binary"
	"fmt"
	"time"

	"github.com/tarm/serial"
)

type Roomba struct {
	port *serial.Port
	path string
	baud int
}

func MakeRoomba(path string, baud int) (*Roomba, error) {
	var r *Roomba
	if _, ok := codeForBaudRate[baud]; !ok {
		return nil, fmt.Errorf("oibot.MakeRoomba(): invalid baud rate: %d", baud)
	}
	if port, err := serial.OpenPort(&serial.Config{Name: path, Baud: baud}); nil != err {
		return nil, fmt.Errorf("oibot.MakeRoomba(): serial.OpenPort(%s, %d) failed: %s", path, baud, err)
	} else {
		r = &Roomba{port: port, path: path, baud: baud}
	}
	return r, nil
}

func (r *Roomba) Pack(data ...interface{}) ([]byte, error) {
	buf := new(bytes.Buffer)
	for _, bin := range data {
		if err := binary.Write(buf, binary.BigEndian, bin); nil != err {
			return nil, fmt.Errorf("oibot.pack(): binary.Write() failed: ", err)
		}
	}
	return buf.Bytes(), nil
}

func (r *Roomba) WriteCode(code OpCode) error {
	if n, err := r.port.Write([]byte{byte(code)}); n != 1 || nil != err {
		return fmt.Errorf("(*Roomba).WriteCode(): port.Write(%d) failed: %s", code, err)
	}
	return nil
}

func (r *Roomba) Write(code OpCode, buf ...interface{}) error {
	var (
		bin []byte
		err error
	)
	if err = r.WriteCode(code); nil != err {
		return err
	}
	if bin, err = r.Pack(buf...); nil != err {
		return err
	}
	if n, err := r.port.Write(bin); n != len(bin) || nil != err {
		return fmt.Errorf("(*Roomba).Write(): port.Write(<[DATA]>) failed: %s", err)
	}
	return nil
}

func (r *Roomba) Read(buf []byte) (n int, err error) {
	return r.port.Read(buf)
}

// =====================================================================================================================

func (r *Roomba) Start() error {
	return r.WriteCode(opcStart)
}

func (r *Roomba) Reset() error {
	return r.WriteCode(opcReset)
}

func (r *Roomba) Stop() error {
	return r.WriteCode(opcStop)
}

func (r *Roomba) Baud(baud int) error {
	if code, ok := codeForBaudRate[baud]; !ok {
		return fmt.Errorf("(*Roomba).Baud(): invalid baud rate: %d", baud)
	} else {
		err := r.Write(opcBaud, code)
		time.Sleep(100 * time.Millisecond)
		return err
	}
}

func (r *Roomba) Control() error {
	return r.WriteCode(opcControl)
}

func (r *Roomba) Safe() error {
	return r.WriteCode(opcSafe)
}

func (r *Roomba) Full() error {
	return r.WriteCode(opcFull)
}

func (r *Roomba) Power() error {
	return r.WriteCode(opcPower)
}

func (r *Roomba) Clean() error {
	return r.WriteCode(opcClean)
}

func (r *Roomba) MaxClean() error {
	return r.WriteCode(opcMaxClean)
}

func (r *Roomba) Spot() error {
	return r.WriteCode(opcSpot)
}

func (r *Roomba) SeekDock() error {
	return r.WriteCode(opcForceSeekingDock)
}

func (r *Roomba) Drive(velocity int16, radius int16) error {
	if velocity < MinDriveVelocityMMPS || velocity > MaxDriveVelocityMMPS {
		return fmt.Errorf("(*Roomba).Drive(): invalid velocity: %d", velocity)
	}
	if radius < MinDriveRadiusMM || radius > MaxDriveRadiusMM {
		return fmt.Errorf("(*Roomba).Drive(): invalid radius: %d", radius)
	}
	return r.Write(opcDrive, velocity, radius)
}

func (r *Roomba) DriveStop() error {
	return r.Drive(0, 0)
}

func (r *Roomba) DriveWheels(rightVelocity int16, leftVelocity int16) error {
	if rightVelocity < MinDriveVelocityMMPS || rightVelocity > MaxDriveVelocityMMPS {
		return fmt.Errorf("(*Roomba).Drive(): invalid right wheel velocity: %d", rightVelocity)
	}
	if leftVelocity < MinDriveVelocityMMPS || leftVelocity > MaxDriveVelocityMMPS {
		return fmt.Errorf("(*Roomba).Drive(): invalid left wheel velocity: %d", leftVelocity)
	}
	return r.Write(opcDriveWheels, rightVelocity, leftVelocity)
}
