package arch

import (
	"fmt"
	"os"
	"reflect"
	"testing"
)

func assertEquals(t *testing.T, actual, expected interface{}, format string, a ...interface{}) {
	if !reflect.DeepEqual(actual, expected) {
		msg := fmt.Sprintf(format, a...)
		t.Errorf("not equal: %s\nactual=%+v\nexpected=%+v", msg, actual, expected)
	}
}

func TestParseCPU(t *testing.T) {
	f := func(s string, cpu CPU) {
		rcpu := ParseCPU(s)
		assertEquals(t, rcpu, cpu, "ParseCPU(%q)", s)
	}
	f("x86", X86)
	f("X86", X86)
	f("amd64", AMD64)
	f("AMD64", AMD64)
	f("Amd64", AMD64)
	f("aMd64", AMD64)
	f("amD64", AMD64)
	f("AMd64", AMD64)
	f("AmD64", AMD64)
	f("aMD64", AMD64)
	f("", 0)
	f("foo", 0)
	f("bar", 0)
}

func TestOS(t *testing.T) {
	f := func(s1, s2 string, cpu CPU, err error) {
		if s1 == "" {
			os.Unsetenv("PROCESSOR_ARCHITEW6432")
		} else {
			os.Setenv("PROCESSOR_ARCHITEW6432", s1)
		}
		if s2 == "" {
			os.Unsetenv("PROCESSOR_ARCHITECTURE")
		} else {
			os.Setenv("PROCESSOR_ARCHITECTURE", s2)
		}
		rcpu, rerr := OS()
		assertEquals(t, rcpu, cpu, "OS() returns unexpected CPU for %q and %q", s1, s2)
		assertEquals(t, rerr, err, "OS() returns unexpected error")
	}
	f("amd64", "x86", AMD64, nil)
	f("x86", "amd64", X86, nil)
	f("amd64", "", AMD64, nil)
	f("x86", "", X86, nil)
	f("", "amd64", AMD64, nil)
	f("", "x86", X86, nil)
	f("foo", "amd64", 0, ErrorUnknownArch)
	f("foo", "x86", 0, ErrorUnknownArch)
	f("", "bar", 0, ErrorUnknownArch)
	f("", "bar", 0, ErrorUnknownArch)
	f("amd64", "bar", AMD64, nil)
	f("x86", "bar", X86, nil)
}

func TestExe(t *testing.T) {
	f := func(p string, cpu CPU, err error) {
		rcpu, rerr := Exe(p)
		assertEquals(t, rcpu, cpu, "Exe(%#q) returns unexpected CPU", p)
		assertEquals(t, rerr, err, "Exe(%#q) returns unexpected error", p)
	}
	f("test_data/hello-vc10-32.ex_", X86, nil)
	f("test_data/hello-vc10-64.ex_", AMD64, nil)
}
