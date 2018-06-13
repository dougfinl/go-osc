package osc

import (
	"bytes"
	"testing"
	"time"
)

func TestEncodeTimeTag(t *testing.T) {
	// A time tag with Immediate=true should return a special value.
	test1 := NewImmediateTimeTag()
	expected1 := []byte{'\x00', '\x00', '\x00', '\x00', '\x00', '\x00', '\x00', '\x01'}
	result1 := encodeTimeTag(test1)

	if !bytes.Equal(result1, expected1) {
		t.Errorf("Encoded tag Immediate value incorrect. Got %v, expected %v", result1, expected1)
	}

	// Jan 1, 2018 should be 3723753600 (0xDDF3F880) seconds since OSC epoch.
	test2 := NewTimeTag(time.Date(2018, 1, 1, 0, 0, 0, 0, time.UTC))
	expected2 := []byte{'\xDD', '\xF3', '\xF8', '\x80', '\x00', '\x00', '\x00', '\x00'}
	result2 := encodeTimeTag(test2)

	if !bytes.Equal(result2, expected2) {
		t.Errorf("New vale is %v, expected %v", result2, expected2)
	}

	// Add 0.5s to previous test to test encoding of fractional time
	test3 := NewTimeTag(time.Date(2018, 1, 1, 0, 0, 0, 0, time.UTC))
	test3.time = test3.time.Add(500 * time.Millisecond)
	expected3 := []byte{'\xDD', '\xF3', '\xF8', '\x80', '\x1D', '\xCD', '\x65', '\x00'}
	result3 := encodeTimeTag(test3)

	if !bytes.Equal(result3, expected3) {
		t.Errorf("New vale is %v, expected %v", result3, expected3)
	}
}

func TestDecodeTimeTag(t *testing.T) {
	// Should result in a time tag with Immediate=true
	test1 := []byte{'\x00', '\x00', '\x00', '\x00', '\x00', '\x00', '\x00', '\x01'}
	expected1 := NewImmediateTimeTag()
	result1, err1 := decodeTimeTag(bytes.NewBuffer(test1))

	if err1 != nil {
		t.Error(err1)
	} else if result1 != expected1 {
		t.Errorf("New value is %v, expected %v", result1, expected1)
	}

	// 0xDDF3F880 seconds since OSC epoch should return Jan 1, 2018 UTC.
	test2 := []byte{'\xDD', '\xF3', '\xF8', '\x80', '\x00', '\x00', '\x00', '\x00'}
	expected2 := NewTimeTag(time.Date(2018, 1, 1, 0, 0, 0, 0, time.UTC))
	result2, err2 := decodeTimeTag(bytes.NewBuffer(test2))

	if err2 != nil {
		t.Error(err2)
	} else if result2 != expected2 {
		t.Errorf("New value is %v, expected %v", result2, expected2)
	}

	// Same as previous test but with 0.5s added
	test3 := []byte{'\xDD', '\xF3', '\xF8', '\x80', '\x1D', '\xCD', '\x65', '\x00'}
	expected3 := NewTimeTag(time.Date(2018, 1, 1, 0, 0, 0, 0, time.UTC))
	expected3.time = expected3.time.Add(500 * time.Millisecond)
	result3, err3 := decodeTimeTag(bytes.NewBuffer(test3))

	if err3 != nil {
		t.Error(err3)
	} else if result3 != expected3 {
		t.Errorf("New value is %v, expected %v", result3, expected3)
	}
}

func TestPadTo32Bits(t *testing.T) {
	// 0-byte slice should not change size
	test1 := []byte{}
	expected1 := []byte{}
	result1 := padTo32Bits(test1)

	if !bytes.Equal(result1, expected1) {
		t.Errorf("New value is %v, expected %v", result1, expected1)
	}

	// Single-byte slice should become 4 bytes (32 bits)
	test2 := []byte{'/'}
	expected2 := []byte{'/', '\x00', '\x00', '\x00'}
	result2 := padTo32Bits(test2)

	if !bytes.Equal(result2, expected2) {
		t.Errorf("New value is %v, expected %v", result2, expected2)
	}

	// Random test
	test3 := []byte("/oscillator/4/frequency")
	expected3 := []byte{'/', 'o', 's', 'c', 'i', 'l', 'l', 'a', 't', 'o', 'r', '/', '4', '/', 'f', 'r', 'e', 'q', 'u', 'e', 'n', 'c', 'y', '\x00'}
	result3 := padTo32Bits(test3)

	if !bytes.Equal(result3, expected3) {
		t.Errorf("New value if %v, expected %v", result3, expected3)
	}
}
