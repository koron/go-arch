package arch

import (
	"debug/pe"
	"errors"
	"fmt"
	"os"
	"strings"
)

// CPU is type of CPU architecture.
type CPU int

const (
	// X86 means Intel x86 (32 bit).
	X86 CPU = iota + 1

	// AMD64 means AMD/Intel 64 bit.
	AMD64
)

// ParseCPU parses string as CPU.
func ParseCPU(s string) CPU {
	cpu, err := parseString(s)
	if err != nil {
		return 0
	}
	return cpu
}

func (cpu CPU) String() string {
	switch cpu {
	case X86:
		return "X86"
	case AMD64:
		return "AMD64"
	default:
		return "(UNKNOWN)"
	}
}

// ErrorUnknownArch is returned when failed to deetect architecture.
var ErrorUnknownArch = errors.New("unknown architecture")

// OS returns architecture of operating system.
func OS() (CPU, error) {
	if v, ok := os.LookupEnv("PROCESSOR_ARCHITEW6432"); ok {
		return parseString(v)
	}
	v, ok := os.LookupEnv("PROCESSOR_ARCHITECTURE")
	if !ok {
		return 0, ErrorUnknownArch
	}
	return parseString(v)
}

func parseString(v string) (CPU, error) {
	switch strings.ToUpper(v) {
	case "X86":
		return X86, nil
	case "AMD64":
		return AMD64, nil
	default:
		return 0, ErrorUnknownArch
	}
}

// Exe returns architecture of execute file.
func Exe(name string) (CPU, error) {
	f, err := pe.Open(name)
	if err != nil {
		if os.IsNotExist(err) {
			return OS()
		}
		return 0, fmt.Errorf("failed to pe.Open %#q: %s", name, err)
	}
	defer f.Close()

	switch f.FileHeader.Machine {
	case 0x014c:
		return X86, nil
	case 0x8664:
		return AMD64, nil
	}
	return 0, ErrorUnknownArch
}
