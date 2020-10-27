package ipmi

type BootDevice uint8

// In same order as the spec
const (
	BootDeviceNoOverride         BootDevice = 0b0000
	BootDeviceForcePXE           BootDevice = 0b0001
	BootDeviceForceHDD           BootDevice = 0b0010
	BootDeviceForceHDDSafe       BootDevice = 0b0011
	BootDeviceForceDiagnostics   BootDevice = 0b0100
	BootDeviceForceDVD           BootDevice = 0b0101
	BootDeviceForceBIOS          BootDevice = 0b0110
	BootDeviceForceRemoteFloppy  BootDevice = 0b0111
	BootDeviceForceRemotePrimary BootDevice = 0b1001
	BootDeviceForceRemoteDVD     BootDevice = 0b1000
	BootDeviceForceRemoteHDD     BootDevice = 0b1011
	BootDeviceForceFloppy        BootDevice = 0b1111
)

type FirmwareVerbosity uint8

const (
	FirmwareVerbosityDefault FirmwareVerbosity = 0b00
	FirmwareVerbosityQuiet   FirmwareVerbosity = 0b01
	FirmwareVerbosityVerbose FirmwareVerbosity = 0b10
)

type ConsoleRedirectionControl uint8

const (
	ConsoleRedirectionControlDefault  ConsoleRedirectionControl = 0b00
	ConsoleRedirectionControlSuppress ConsoleRedirectionControl = 0b01
	ConsoleRedirectionControlEnabled  ConsoleRedirectionControl = 0b10
)

type BIOSMuxControlOverride uint8

const (
	BIOSMuxControlOverrideRecommended BIOSMuxControlOverride = 0b000
	BIOSMuxControlOverrideBMC         BIOSMuxControlOverride = 0b001
	BIOSMuxControlOverrideSystem      BIOSMuxControlOverride = 0b010
)

type BootFlags struct {
	// Byte 1
	Valid      bool // The bit should be set to indicate that valid flag data is present
	Persistent bool // If the options should be persistent or apply to next boot only
	UEFI       bool // Use UEFI or BIOS

	// Byte 2
	ClearCMOS          bool
	LockKeyboard       bool
	BootDevice         BootDevice
	ScreenBlank        bool
	LockOutResetButton bool

	// Byte 3
	LockOutPowerButton        bool
	FirmwareVerbosity         FirmwareVerbosity
	ForceProgressEventTraps   bool
	UserPasswordBypass        bool
	LockOutSleepButton        bool
	ConsoleRedirectionControl ConsoleRedirectionControl

	// Byte 4
	BIOSSharedModeOverride bool
	BIOSMuxControlOverride BIOSMuxControlOverride

	// Byte 5
	DeviceInstanceSelector uint8
}

type BootOptionsParameterData struct {
	// Parameter0: Not implemented
	// Parameter1: Not implemented
	// Parameter2: Not implemented
	// Parameter3: Not implemented
	// Parameter4: Not implemented
	// Parameter5: Boot Flags
	BootFlags BootFlags

	// Parameter6: Not implemented
	// Parameter7: Not implemented
}
