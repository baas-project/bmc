package ipmi

import (
	"fmt"
	"github.com/google/gopacket"
	"github.com/google/gopacket/layers"
)

// Section 28.12
type SetSystemBootOptionsReq struct {
	layers.BaseLayer

	// valid flag
	ParameterValid bool

	// boot option parameter selector
	ParameterSelector uint8

	// BootOptionsParameterData contains all (supported) parameters.
	BootOptionsParameterData
}

func (s SetSystemBootOptionsReq) LayerType() gopacket.LayerType {
	return LayerTypeSetSystemBootOptionsReq
}

func (s SetSystemBootOptionsReq) SerializeTo(b gopacket.SerializeBuffer, _ gopacket.SerializeOptions) error {
	switch s.ParameterSelector {
	case 5:
		bytes, err := b.PrependBytes(6)
		if err != nil {
			return err
		}

		bytes[0] |= s.ParameterSelector & 0b01111111
		if s.ParameterValid {
			bytes[0] |= 1<<7
		}

		// Byte 1
		if s.BootFlags.Valid {
			bytes[1] |= 1<<7
		}
		if s.BootFlags.Persistent {
			bytes[1] |= 1<<6
		}
		if s.BootFlags.UEFI {
			bytes[1] |= 1<<5
		}

		// Byte 2
		if s.BootFlags.ClearCMOS {
			bytes[2] |= 1<<7
		}
		bytes[2] |= (uint8(s.BootFlags.BootDevice) << 2) & 0b00111100
		if s.BootFlags.ScreenBlank {
			bytes[2] |= 1<<1
		}
		if s.BootFlags.LockOutResetButton {
			bytes[2] |= 1<<0
		}

		// Byte 3
		if s.BootFlags.LockOutPowerButton {
			bytes[3] |= 1<<7
		}
		bytes[3] |= (uint8(s.BootFlags.FirmwareVerbosity) << 5) & 0b01100000
		if s.BootFlags.ForceProgressEventTraps {
			bytes[3] |= 1<<4
		}
		if s.BootFlags.UserPasswordBypass {
			bytes[3] |= 1<<3
		}
		if s.BootFlags.LockOutSleepButton {
			bytes[3] |= 1<<2
		}
		bytes[3] |= uint8(s.BootFlags.FirmwareVerbosity) & 0b00000011

		// Byte 4
		if s.BootFlags.BIOSSharedModeOverride {
			bytes[4] |= 1<<3
		}
		bytes[4] |= uint8(s.BootFlags.BIOSMuxControlOverride) & 0b00000011

		// Byte 5
		bytes[5] |= s.BootFlags.DeviceInstanceSelector & 0b00011111

	default:
		return fmt.Errorf("unsupported parameter type %v", s.ParameterSelector)
	}

	return nil
}

type SetSystemBootOptionsCmd struct {
	Req SetSystemBootOptionsReq
}

func (c *SetSystemBootOptionsCmd) Name() string {
	return "Set System Boot Options"
}

func (*SetSystemBootOptionsCmd) Operation() *Operation {
	return &OperationSetSystemBootOptionsRsp
}

func (c *SetSystemBootOptionsCmd) Request() gopacket.SerializableLayer {
	return &c.Req
}

func (*SetSystemBootOptionsCmd) Response() gopacket.DecodingLayer {
	return nil
}
