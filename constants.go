package oibot

import "time"

// =====================================================================================================================
type OpCode byte

// -------------------- ----- --------------------- --------------------- --------------------- ---------------------
//  OpCode               ID    Input Byte 1          Input Byte 2          Input Byte 3          Input Byte 4
// -------------------- ----- --------------------- --------------------- --------------------- ---------------------
//  Reset                  7
//  Start                128
//  Baud                 129   BaudCode
//  Control              130
//  Safe                 131
//  Full                 132
//  Power                133
//  Spot                 134
//  Clean                135
//  Max Clean            136
//  Drive                137   VelocityHigh          VelocityLow           RadiusHigh            RadiusLow
//  Motors               138   MotorsState
//  LEDs                 139   LEDsState             PowerColor            PowerIntensity
//  Song                 140   SongNum               SongLength
//  Play                 141   SongNum
//  Query                142   Packet
//  Force Seeking Dock   143
//  PWM Motors           144   MainBrushPWM          SideBrushPWM          VacuumPWM
//  Drive Wheels         145   RightVelocityHigh     RightVelocityLow      LeftVelocityHigh      LeftVelocityLow
//  Drive PWM            146   RightPWMHigh          RightPWMLow           LeftPWMHigh           LeftPWMLow
//  Stream               148   NumPackets
//  Query List           149   NumPackets
//  Do Stream            150   StreamState
//  Scheduling LEDs      162   Weekdays              SchedulingLEDsState
//  Digit LEDs Raw       163   Digit3                Digit2                Digit1                Digit0
//  Digit LEDs ASCII     164   Digit3                Digit2                Digit1                Digit0
//  Buttons              165   Buttons
//  Schedule             167   Days                  SunHour               SunMin                MonHour       (etc.)
//  Set Day/Time         168   Day                   Hour                  Minute
//  Stop                 173
// -------------------- ----- --------------------- --------------------- --------------------- ---------------------
const (
	opcReset            OpCode = 7
	opcStart            OpCode = 128
	opcBaud             OpCode = 129
	opcControl          OpCode = 130
	opcSafe             OpCode = 131
	opcFull             OpCode = 132
	opcPower            OpCode = 133
	opcSpot             OpCode = 134
	opcClean            OpCode = 135
	opcMaxClean         OpCode = 136
	opcDrive            OpCode = 137
	opcMotors           OpCode = 138
	opcLEDs             OpCode = 139
	opcSong             OpCode = 140
	opcPlay             OpCode = 141
	opcQuery            OpCode = 142
	opcForceSeekingDock OpCode = 143
	opcPWMMotors        OpCode = 144
	opcDriveWheels      OpCode = 145
	opcDrivePWM         OpCode = 146
	opcStream           OpCode = 148
	opcQueryList        OpCode = 149
	opcDoStream         OpCode = 150
	opcSchedulingLEDs   OpCode = 162
	opcDigitLEDsRaw     OpCode = 163
	opcDigitLEDsASCII   OpCode = 164
	opcButtons          OpCode = 165
	opcSchedule         OpCode = 167
	opcSetDayTime       OpCode = 168
	opcStop             OpCode = 173
)

const (
	MaxDriveVelocityMMPS   int16 = 500
	MinDriveVelocityMMPS   int16 = -500
	MaxDriveRadiusMM       int16 = 2000
	MinDriveRadiusMM       int16 = -2000
	StraightDriveRadiusMM  int16 = 0x7FFF
	DriveWheelSeparationMM int16 = 298
)

// =====================================================================================================================
type SensorPacket struct {
	id   byte
	size byte
}

// -------------------------- ---- ------ ------------------------------- ---------------- -------
//  Sensor Packet              ID   Size   Group Membership                Range            Units
// -------------------------- ---- ------ ------------------------------- ---------------- -------
//  Bumps Wheeldrops            7   1      0 1         6 100                    0 - 15
//  Wall                        8   1      0 1         6 100                    0 - 1
//  Cliff Left                  9   1      0 1         6 100                    0 - 1
//  Cliff Front Left           10   1      0 1         6 100                    0 - 1
//  Cliff Front Right          11   1      0 1         6 100                    0 - 1
//  Cliff Right                12   1      0 1         6 100                    0 - 1
//  Virtual Wall               13   1      0 1         6 100                    0 - 1
//  Overcurrents               14   1      0 1         6 100                    0 - 29
//  Dirt Detect                15   1      0 1         6 100                    0 - 255
//  Unused 1                   16   1      0 1         6 100                    0 - 255
//  IR OpCode                  17   1      0   2       6 100                    0 - 255
//  Buttons                    18   1      0   2       6 100                    0 - 255
//  Distance                   19   2      0   2       6 100               -32768 - 32767   mm
//  Angle                      20   2      0   2       6 100               -32768 - 32767   degrees
//  Charging State             21   1      0     3     6 100                    0 - 6
//  Voltage                    22   2      0     3     6 100                    0 - 65535   mV
//  Current                    23   2      0     3     6 100               -32768 - 32767   mA
//  Temperature                24   1      0     3     6 100                 -128 - 127     deg C
//  Battery Charge             25   2      0     3     6 100                    0 - 65535   mAh
//  Battery Capacity           26   2      0     3     6 100                    0 - 65535   mAh
//  Wall Signal                27   2              4   6 100                    0 - 1023
//  Cliff Left Signal          28   2              4   6 100                    0 - 4095
//  Cliff Front Left Signal    29   2              4   6 100                    0 - 4095
//  Cliff Front Right Signal   30   2              4   6 100                    0 - 4095
//  Cliff Right Signal         31   2              4   6 100                    0 - 4095
//  Unused 2                   32   1              4   6 100                    0 - 255
//  Unused 3                   33   2              4   6 100                    0 - 65535
//  Charger Available          34   1              4   6 100                    0 - 3
//  Open Interface Mode        35   1                5 6 100                    0 - 3
//  Song Number                36   1                5 6 100                    0 - 4
//  Song Playing?              37   1                5 6 100                    0 - 1
//  Oi Stream Num Packets      38   1                5 6 100                    0 - 108
//  Velocity                   39   2                5 6 100                 -500 - 500     mm/s
//  Radius                     40   2                5 6 100               -32768 - 32767   mm
//  Velocity Right             41   2                5 6 100                 -500 - 500     mm/s
//  Velocity Left              42   2                5 6 100                 -500 - 500     mm/s
//  Encoder Counts Left        43   2                    100 101                0 - 65535
//  Encoder Counts Right       44   2                    100 101                0 - 65535
//  Light Bumper               45   1                    100 101                0 - 127
//  Light Bump Left            46   2                    100 101 106            0 - 4095
//  Light Bump Front Left      47   2                    100 101 106            0 - 4095
//  Light Bump Center Left     48   2                    100 101 106            0 - 4095
//  Light Bump Center Right    49   2                    100 101 106            0 - 4095
//  Light Bump Front Right     50   2                    100 101 106            0 - 4095
//  Light Bump Right           51   2                    100 101 106            0 - 4095
//  IR OpCode Left             52   1                    100 101                0 - 255
//  IR OpCode Right            53   1                    100 101                0 - 255
//  Left Motor Current         54   2                    100 101     107   -32768 - 32767   mA
//  Right Motor Current        55   2                    100 101     107   -32768 - 32767   mA
//  Main Brush Current         56   2                    100 101     107   -32768 - 32767   mA
//  Side Brush Current         57   2                    100 101     107   -32768 - 32767   mA
//  Stasis                     58   1                    100 101     107        0 - 3
// -------------------------- ---- ------ ------------------------------- ---------------- -------
var (
	spcBumpsWheeldrops       = &SensorPacket{id: 7, size: 1}
	spcWall                  = &SensorPacket{id: 8, size: 1}
	spcCliffLeft             = &SensorPacket{id: 9, size: 1}
	spcCliffFrontLeft        = &SensorPacket{id: 10, size: 1}
	spcCliffFrontRight       = &SensorPacket{id: 11, size: 1}
	spcCliffRight            = &SensorPacket{id: 12, size: 1}
	spcVirtualWall           = &SensorPacket{id: 13, size: 1}
	spcOvercurrents          = &SensorPacket{id: 14, size: 1}
	spcDirtDetect            = &SensorPacket{id: 15, size: 1}
	spcUnused1               = &SensorPacket{id: 16, size: 1}
	spcIROpCode              = &SensorPacket{id: 17, size: 1}
	spcButtons               = &SensorPacket{id: 18, size: 1}
	spcDistance              = &SensorPacket{id: 19, size: 2}
	spcAngle                 = &SensorPacket{id: 20, size: 2}
	spcChargingState         = &SensorPacket{id: 21, size: 1}
	spcVoltage               = &SensorPacket{id: 22, size: 2}
	spcCurrent               = &SensorPacket{id: 23, size: 2}
	spcTemperature           = &SensorPacket{id: 24, size: 1}
	spcBatteryCharge         = &SensorPacket{id: 25, size: 2}
	spcBatteryCapacity       = &SensorPacket{id: 26, size: 2}
	spcWallSignal            = &SensorPacket{id: 27, size: 2}
	spcCliffLeftSignal       = &SensorPacket{id: 28, size: 2}
	spcCliffFrontLeftSignal  = &SensorPacket{id: 29, size: 2}
	spcCliffFrontRightSignal = &SensorPacket{id: 30, size: 2}
	spcCliffRightSignal      = &SensorPacket{id: 31, size: 2}
	spcUnused2               = &SensorPacket{id: 32, size: 1}
	spcUnused3               = &SensorPacket{id: 33, size: 2}
	spcChargerAvailable      = &SensorPacket{id: 34, size: 1}
	spcOpenInterfaceMode     = &SensorPacket{id: 35, size: 1}
	spcSongNumber            = &SensorPacket{id: 36, size: 1}
	spcSongPlaying           = &SensorPacket{id: 37, size: 1}
	spcOIStreamNumPackets    = &SensorPacket{id: 38, size: 1}
	spcVelocity              = &SensorPacket{id: 39, size: 2}
	spcRadius                = &SensorPacket{id: 40, size: 2}
	spcVelocityRight         = &SensorPacket{id: 41, size: 2}
	spcVelocityLeft          = &SensorPacket{id: 42, size: 2}
	spcEncoderCountsLeft     = &SensorPacket{id: 43, size: 2}
	spcEncoderCountsRight    = &SensorPacket{id: 44, size: 2}
	spcLightBumper           = &SensorPacket{id: 45, size: 1}
	spcLightBumpLeft         = &SensorPacket{id: 46, size: 2}
	spcLightBumpFrontLeft    = &SensorPacket{id: 47, size: 2}
	spcLightBumpCenterLeft   = &SensorPacket{id: 48, size: 2}
	spcLightBumpCenterRight  = &SensorPacket{id: 49, size: 2}
	spcLightBumpFrontRight   = &SensorPacket{id: 50, size: 2}
	spcLightBumpRight        = &SensorPacket{id: 51, size: 2}
	spcIROpCodeLeft          = &SensorPacket{id: 52, size: 1}
	spcIROpCodeRight         = &SensorPacket{id: 53, size: 1}
	spcLeftMotorCurrent      = &SensorPacket{id: 54, size: 2}
	spcRightMotorCurrent     = &SensorPacket{id: 55, size: 2}
	spcMainBrushCurrent      = &SensorPacket{id: 56, size: 2}
	spcSideBrushCurrent      = &SensorPacket{id: 57, size: 2}
	spcStasis                = &SensorPacket{id: 58, size: 1}
)

// =====================================================================================================================
type SensorGroup struct {
	id     byte
	size   byte
	member []*SensorPacket
}

//  Group Packet ID   ~~  Packet Size   ##  Contains Packets
// --------------------------  ----- ------ ------------------------------- ---------------- -------
//  Sensor Group                ID    Size   Members
// --------------------------  ----- ------ ------------------------------- ---------------- -------
//  Status                      0     26      7 - 26
//  Obstacle                    1     10      7 - 16
//  Dock                        2     6      17 - 20
//  Battery                     3     10     21 - 26
//  Signal                      4     14     27 - 34
//  ModeData                    5     12     35 - 42
//  Sensor                      6     52      7 - 42
//  All                         100   80      7 - 58
//  Drive                       101   28     43 - 58
//  Proximity                   106   12     46 - 51
//  Actuator                    107   9      54 - 58
// --------------------------  ----- ------ ------------------------------- ---------------- -------
var (
	sgpStatus    = &SensorGroup{id: 0, size: 26, member: []*SensorPacket{spcBumpsWheeldrops, spcWall, spcCliffLeft, spcCliffFrontLeft, spcCliffFrontRight, spcCliffRight, spcVirtualWall, spcOvercurrents, spcDirtDetect, spcUnused1, spcIROpCode, spcButtons, spcDistance, spcAngle, spcChargingState, spcVoltage, spcCurrent, spcTemperature, spcBatteryCharge, spcBatteryCapacity}}
	sgpObstacle  = &SensorGroup{id: 1, size: 10, member: []*SensorPacket{spcBumpsWheeldrops, spcWall, spcCliffLeft, spcCliffFrontLeft, spcCliffFrontRight, spcCliffRight, spcVirtualWall, spcOvercurrents, spcDirtDetect, spcUnused1}}
	sgpDock      = &SensorGroup{id: 2, size: 6, member: []*SensorPacket{spcIROpCode, spcButtons, spcDistance, spcAngle}}
	sgpBattery   = &SensorGroup{id: 3, size: 10, member: []*SensorPacket{spcChargingState, spcVoltage, spcCurrent, spcTemperature, spcBatteryCharge, spcBatteryCapacity}}
	sgpSignal    = &SensorGroup{id: 4, size: 14, member: []*SensorPacket{spcWallSignal, spcCliffLeftSignal, spcCliffFrontLeftSignal, spcCliffFrontRightSignal, spcCliffRightSignal, spcUnused2, spcUnused3, spcChargerAvailable}}
	sgpModeData  = &SensorGroup{id: 5, size: 12, member: []*SensorPacket{spcOpenInterfaceMode, spcSongNumber, spcSongPlaying, spcOIStreamNumPackets, spcVelocity, spcRadius, spcVelocityRight, spcVelocityLeft}}
	sgpSensor    = &SensorGroup{id: 6, size: 52, member: []*SensorPacket{spcBumpsWheeldrops, spcWall, spcCliffLeft, spcCliffFrontLeft, spcCliffFrontRight, spcCliffRight, spcVirtualWall, spcOvercurrents, spcDirtDetect, spcUnused1, spcIROpCode, spcButtons, spcDistance, spcAngle, spcChargingState, spcVoltage, spcCurrent, spcTemperature, spcBatteryCharge, spcBatteryCapacity, spcWallSignal, spcCliffLeftSignal, spcCliffFrontLeftSignal, spcCliffFrontRightSignal, spcCliffRightSignal, spcUnused2, spcUnused3, spcChargerAvailable, spcOpenInterfaceMode, spcSongNumber, spcSongPlaying, spcOIStreamNumPackets, spcVelocity, spcRadius, spcVelocityRight, spcVelocityLeft}}
	sgpAll       = &SensorGroup{id: 100, size: 80, member: []*SensorPacket{spcBumpsWheeldrops, spcWall, spcCliffLeft, spcCliffFrontLeft, spcCliffFrontRight, spcCliffRight, spcVirtualWall, spcOvercurrents, spcDirtDetect, spcUnused1, spcIROpCode, spcButtons, spcDistance, spcAngle, spcChargingState, spcVoltage, spcCurrent, spcTemperature, spcBatteryCharge, spcBatteryCapacity, spcWallSignal, spcCliffLeftSignal, spcCliffFrontLeftSignal, spcCliffFrontRightSignal, spcCliffRightSignal, spcUnused2, spcUnused3, spcChargerAvailable, spcOpenInterfaceMode, spcSongNumber, spcSongPlaying, spcOIStreamNumPackets, spcVelocity, spcRadius, spcVelocityRight, spcVelocityLeft, spcEncoderCountsLeft, spcEncoderCountsRight, spcLightBumper, spcLightBumpLeft, spcLightBumpFrontLeft, spcLightBumpCenterLeft, spcLightBumpCenterRight, spcLightBumpFrontRight, spcLightBumpRight, spcIROpCodeLeft, spcIROpCodeRight, spcLeftMotorCurrent, spcRightMotorCurrent, spcMainBrushCurrent, spcSideBrushCurrent, spcStasis}}
	sgpDrive     = &SensorGroup{id: 101, size: 28, member: []*SensorPacket{spcEncoderCountsLeft, spcEncoderCountsRight, spcLightBumper, spcLightBumpLeft, spcLightBumpFrontLeft, spcLightBumpCenterLeft, spcLightBumpCenterRight, spcLightBumpFrontRight, spcLightBumpRight, spcIROpCodeLeft, spcIROpCodeRight, spcLeftMotorCurrent, spcRightMotorCurrent, spcMainBrushCurrent, spcSideBrushCurrent, spcStasis}}
	sgpProximity = &SensorGroup{id: 106, size: 12, member: []*SensorPacket{spcLightBumpLeft, spcLightBumpFrontLeft, spcLightBumpCenterLeft, spcLightBumpCenterRight, spcLightBumpFrontRight, spcLightBumpRight}}
	sgpActuator  = &SensorGroup{id: 107, size: 9, member: []*SensorPacket{spcLeftMotorCurrent, spcRightMotorCurrent, spcMainBrushCurrent, spcSideBrushCurrent, spcStasis}}
)

// =====================================================================================================================
type BaudRateCode byte

const (
	brc300BPS    BaudRateCode = 0
	brc600BPS    BaudRateCode = 1
	brc1200BPS   BaudRateCode = 2
	brc2400BPS   BaudRateCode = 3
	brc4800BPS   BaudRateCode = 4
	brc9600BPS   BaudRateCode = 5
	brc14400BPS  BaudRateCode = 6
	brc19200BPS  BaudRateCode = 7
	brc28800BPS  BaudRateCode = 8
	brc38400BPS  BaudRateCode = 9
	brc57600BPS  BaudRateCode = 10
	brc115200BPS BaudRateCode = 11
)

var (
	codeForBaudRate = map[int]BaudRateCode{
		300:    brc300BPS,
		600:    brc600BPS,
		1200:   brc1200BPS,
		2400:   brc2400BPS,
		4800:   brc4800BPS,
		9600:   brc9600BPS,
		14400:  brc14400BPS,
		19200:  brc19200BPS,
		28800:  brc28800BPS,
		38400:  brc38400BPS,
		57600:  brc57600BPS,
		115200: brc115200BPS,
	}
)

const (
	DefaultBaudRateBPS    int           = 115200
	SerialTransferDelayMS time.Duration = 50 * time.Millisecond
	SensorUpdateDelayMS   time.Duration = 15 * time.Millisecond // from OI spec
	DefaultReadTimeoutMS  time.Duration = 500 * time.Millisecond
	NeverReadTimeoutMS    time.Duration = 0
)

// =====================================================================================================================
type ChargingStateCode byte

const (
	cstNotCharging            ChargingStateCode = 0
	cstReconditioningCharging ChargingStateCode = 1
	cstFullCharging           ChargingStateCode = 2
	cstTrickleCharging        ChargingStateCode = 3
	cstWaiting                ChargingStateCode = 4
	cstChargingFaultCondition ChargingStateCode = 5
)

// =====================================================================================================================
type OpenInterfaceMode byte

const (
	BootupTimeMS time.Duration = 5000 * time.Millisecond
)

const (
	OIMOff     OpenInterfaceMode = 0
	OIMPassive OpenInterfaceMode = 1
	OIMSafe    OpenInterfaceMode = 2
	OIMFull    OpenInterfaceMode = 3
)

var (
	oiModeStr = [...]string{"OFF", "PASV", "SAFE", "FULL"}
)

func OIModeStr(mode OpenInterfaceMode) (string, bool) {
	if mode >= OIMOff && mode <= OIMFull {
		return oiModeStr[mode], true
	}
	return "", false
}

// =====================================================================================================================
type Direction uint

type angleBounds struct {
	min, max int16
}

type angleRune struct {
	dom []angleBounds
	gfx []rune
	dir Direction
}

const (
	AngleRuneWeightMax = byte(4)
	AngleRuneUnknown   = '⮔'
)

const (
	DirStop Direction = iota
	DirLeft
	DirLeftFwd
	DirFwd
	DirRightFwd
	DirRight
	DirRightAft
	DirAft
	DirLeftAft
)

func AngleRune(angle int16, weight byte) (rune, Direction) {
	// all angle vars are in units degrees. weight corresponds to the column of
	// runes in the table below, valid range is 0..4
	var (
		arrow = []angleRune{{
			dom: []angleBounds{{min: 202, max: 247}, {min: -157, max: -112}},
			gfx: []rune{'🡧', '🡯', '🡷', '🡿', '🢇'},
			dir: DirLeftAft,
		}, {
			dom: []angleBounds{{min: 112, max: 157}, {min: -247, max: -202}},
			gfx: []rune{'🡦', '🡮', '🡶', '🡾', '🢆'},
			dir: DirRightAft,
		}, {
			dom: []angleBounds{{min: 22, max: 67}, {min: -337, max: -292}},
			gfx: []rune{'🡥', '🡭', '🡵', '🡽', '🢅'},
			dir: DirRightFwd,
		}, {
			dom: []angleBounds{{min: 292, max: 337}, {min: -67, max: -22}},
			gfx: []rune{'🡤', '🡬', '🡴', '🡼', '🢄'},
			dir: DirLeftFwd,
		}, {
			dom: []angleBounds{{min: 157, max: 202}, {min: -202, max: -157}},
			gfx: []rune{'🡣', '🡫', '🡳', '🡻', '🢃'},
			dir: DirAft,
		}, {
			dom: []angleBounds{{min: -22, max: 22}, {min: 337, max: 382}, {min: -382, max: -337}},
			gfx: []rune{'🡡', '🡩', '🡱', '🡹', '🢁'},
			dir: DirFwd,
		}, {
			dom: []angleBounds{{min: 67, max: 112}, {min: -292, max: -247}},
			gfx: []rune{'🡢', '🡪', '🡲', '🡺', '🢂'},
			dir: DirRight,
		}, {
			dom: []angleBounds{{min: 247, max: 292}, {min: -112, max: -67}},
			gfx: []rune{'🡠', '🡨', '🡰', '🡸', '🢀'},
			dir: DirLeft,
		}}
	)

	if weight >= 0 && weight <= AngleRuneWeightMax {
		norm := angle % 360
		for _, a := range arrow {
			for _, d := range a.dom {
				if norm >= d.min && norm < d.max {
					return a.gfx[weight], a.dir
				}
			}
		}
	}
	return AngleRuneUnknown, DirStop
}
