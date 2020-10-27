package ipmi

import (
	"fmt"
	"github.com/google/gopacket"
	"github.com/google/gopacket/layers"
)

// Section 28.13
type GetSystemBootOptionsReq struct {
	layers.BaseLayer

	// boot option parameter selector
	ParameterSelector uint8

	// SetSelector Selects a particular block or set of parameters under the given parameter selector.
	SetSelector uint8

	// BlockSelector Selects aparticular block within a setof parameters
	BlockSelector uint8
}

func (g GetSystemBootOptionsReq) LayerType() gopacket.LayerType {
	return LayerTypeGetSystemBootOptionsReq
}

func (g GetSystemBootOptionsReq) SerializeTo(b gopacket.SerializeBuffer, _ gopacket.SerializeOptions) error {
	bytes, err := b.PrependBytes(3)
	if err != nil {
		return err
	}

	bytes[0] = g.ParameterSelector
	bytes[1] = g.SetSelector
	bytes[2] = g.BlockSelector

	return nil
}

type GetSystemBootOptionsRsp struct {
	layers.BaseLayer

	ParameterNotSupported bool

	ParameterVersion uint8

	ParameterValid bool

	ParameterSelector uint8

	BootOptionsParameterData
}

func (*GetSystemBootOptionsRsp) LayerType() gopacket.LayerType {
	return LayerTypeGetSystemBootOptionsRsp
}

func (g *GetSystemBootOptionsRsp) CanDecode() gopacket.LayerClass {
	return g.LayerType()
}

func (*GetSystemBootOptionsRsp) NextLayerType() gopacket.LayerType {
	return gopacket.LayerTypePayload
}

func (g *GetSystemBootOptionsRsp) DecodeFromBytes(data []byte, df gopacket.DecodeFeedback) error {
	if len(data) < 3 {
		df.SetTruncated()
		return fmt.Errorf("GetBootDevRsp must be at least 3 bytes; got %v", len(data))
	}

	// FIXME
	// data: [1 5 0 0 0 0 0]

	// Header
	g.ParameterVersion = data[0]
	g.ParameterValid = !(data[1]&0b10000000 != 0)
	g.ParameterSelector = data[1] & 0b01111111

	// Parameter
	switch g.ParameterSelector {
	case 5:
		// Data 1
		g.BootFlags.Valid = data[2]&(1<<7) != 0
		g.BootFlags.Persistent = data[2]&(1<<6) != 0
		g.BootFlags.UEFI = data[2]&(1<<5) != 0

		// Data 2
		g.BootFlags.ClearCMOS = data[3]&(1<<7) != 0
		g.BootFlags.BootDevice = BootDevice((data[3] & 0b00111100) >> 2)
		g.BootFlags.ScreenBlank = data[3]&(1<<1) != 0
		g.BootFlags.LockOutResetButton = data[3]&(1<<0) != 0

		// Data 3
		g.BootFlags.LockOutPowerButton = data[4]&(1<<7) != 0
		g.BootFlags.FirmwareVerbosity = FirmwareVerbosity((data[4] & 0b01100000) >> 5)
		g.BootFlags.ForceProgressEventTraps = data[4]&(1<<4) != 0
		g.BootFlags.UserPasswordBypass = data[4]&(1<<3) != 0
		g.BootFlags.LockOutSleepButton = data[4]&(1<<2) != 0
		g.BootFlags.ConsoleRedirectionControl = ConsoleRedirectionControl(data[4] & 0b00000011)

		// Data 4
		g.BootFlags.BIOSSharedModeOverride = data[5]&(1<<3) != 0
		g.BootFlags.BIOSMuxControlOverride = BIOSMuxControlOverride(data[5] & 0b00000011)

		// Data 5
		g.BootFlags.DeviceInstanceSelector = data[6] & 0b00011111

	default:
		return fmt.Errorf("unsupported parameter type %v", g.ParameterSelector)
	}

	g.BaseLayer.Contents = data[:3]
	g.BaseLayer.Payload = data[3:]

	return nil
}

type GetSystemBootOptionsCmd struct {
	Req GetSystemBootOptionsReq
	Rsp GetSystemBootOptionsRsp
}

func (c *GetSystemBootOptionsCmd) Name() string {
	return "Get System Boot Options"
}

func (*GetSystemBootOptionsCmd) Operation() *Operation {
	return &OperationGetSystemBootOptionsReq
}

func (c *GetSystemBootOptionsCmd) Request() gopacket.SerializableLayer {
	return &c.Req
}

func (c *GetSystemBootOptionsCmd) Response() gopacket.DecodingLayer {
	return &c.Rsp
}


