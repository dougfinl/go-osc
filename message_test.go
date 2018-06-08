package osc

import (
	"bytes"
	"testing"
)

func TestNewEmptyMessage(t *testing.T) {
	msg := NewEmptyMessage()

	if msg.Address != "/" {
		t.Errorf("Address pattern is \"%s\", expected \"/\"", msg.Address)
	}

	if msg.Arguments != nil {
		t.Error("Message does not have 0 arguments")
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

func TestTypeTagString(t *testing.T) {
	// A message with no arguments should produce an empty type tag string
	msg1 := NewEmptyMessage()
	expected1 := ","
	result1, err1 := msg1.TypeTagString()

	if err1 != nil {
		t.Error(err1)
	} else if result1 != expected1 {
		t.Errorf("Type tag string is \"%v\", expected \"%v\"", result1, expected1)
	}

	// Test every supported argument type
	msg2 := NewEmptyMessage()
	msg2.AddArgument(nil) // N
	msg2.AddArgument(int32(10))
	msg2.AddArgument(float32(12.5))
	msg2.AddArgument("test")
	msg2.AddArgument([]byte{'a', 'b'})
	msg2.AddArgument(true)
	msg2.AddArgument(false)
	msg2.AddArgument(int64(9e10))
	msg2.AddArgument(float64(10.1))
	expected2 := ",NifsbTFhd"
	result2, err2 := msg2.TypeTagString()

	if err2 != nil {
		t.Error(err2)
	} else if result2 != expected2 {
		t.Errorf("Type tag string is \"%v\", expected \"%v\"", result2, expected2)
	}
}

func TestBytes(t *testing.T) {
	msg1 := NewEmptyMessage()
	expected1 := []byte{'/', '\x00', '\x00', '\x00', ',', '\x00', '\x00', '\x00'}
	result1, err1 := msg1.MarshalBinary()

	if err1 != nil {
		t.Error(err1)
	} else if !bytes.Equal(result1, expected1) {
		t.Errorf("Got %v, expected %v", result1, expected1)
	}

	msg2 := NewMessage("/oscillator/4/frequency")
	msg2.AddArgument(float32(440))
	expected2 := []byte{'/', 'o', 's', 'c', 'i', 'l', 'l', 'a', 't', 'o', 'r', '/', '4', '/', 'f', 'r', 'e', 'q', 'u', 'e', 'n', 'c', 'y', '\x00', ',', 'f', '\x00', '\x00', '\x43', '\xdc', '\x00', '\x00'}
	result2, err2 := msg2.MarshalBinary()

	if err2 != nil {
		t.Error(err2)
	} else if !bytes.Equal(result2, expected2) {
		t.Errorf("Got %v, expected %v", result2, expected2)
	}

	msg3 := NewMessage("/foo")
	msg3.AddArgument(int32(1000))
	msg3.AddArgument(int32(-1))
	msg3.AddArgument("hello")
	msg3.AddArgument(float32(1.234))
	msg3.AddArgument(float32(5.678))
	expected3 := []byte{'/', 'f', 'o', 'o', '\x00', '\x00', '\x00', '\x00', ',', 'i', 'i', 's', 'f', 'f', '\x00', '\x00', '\x00', '\x00', '\x03', '\xe8', '\xff', '\xff', '\xff', '\xff', '\x68', '\x65', '\x6c', '\x6c', '\x6f', '\x00', '\x00', '\x00', '\x3f', '\x9d', '\xf3', '\xb6', '\x40', '\xb5', '\xb2', '\x2d'}
	result3, err3 := msg3.MarshalBinary()

	if err3 != nil {
		t.Error(err3)
	} else if !bytes.Equal(result3, expected3) {
		t.Errorf("Got %v, expected %v", result3, expected3)
	}

	msg4 := NewMessage("/bytes")
	msg4.AddArgument([]byte{'a', 'b', 'c', 'd', 'e'})
	expected4 := []byte{'/', 'b', 'y', 't', 'e', 's', '\x00', '\x00', ',', 'b', '\x00', '\x00', '\x00', '\x00', '\x00', '\x05', 'a', 'b', 'c', 'd', 'e', '\x00', '\x00', '\x00'}
	result4, err4 := msg4.MarshalBinary()

	if err4 != nil {
		t.Error(err4)
	} else if !bytes.Equal(result4, expected4) {
		t.Errorf("Got %v, expected %v", result4, expected4)
	}
}

func TestUnmarshalBinary(t *testing.T) {
	data1 := []byte{'/', '\x00', '\x00', '\x00', ',', '\x00', '\x00', '\x00'}
	expected1 := NewMessage("/")
	var result1 Message
	err1 := result1.UnmarshalBinary(data1)

	if err1 != nil {
		t.Error(err1)
	} else if !result1.Equals(&expected1) {
		t.Errorf("Got %v, expected %v", result1, expected1)
	}

	data2 := []byte{'/', 'o', 's', 'c', 'i', 'l', 'l', 'a', 't', 'o', 'r', '/', '4', '/', 'f', 'r', 'e', 'q', 'u', 'e', 'n', 'c', 'y', '\x00', ',', 'f', '\x00', '\x00', '\x43', '\xdc', '\x00', '\x00'}
	expected2 := NewMessage("/oscillator/4/frequency")
	expected2.AddArgument(float32(440))
	var result2 Message
	err2 := result2.UnmarshalBinary(data2)

	if err2 != nil {
		t.Error(err2)
	} else if !result2.Equals(&expected2) {
		t.Errorf("Got %v, expected %v", result2, expected2)
	}

	data3 := []byte{'/', 'f', 'o', 'o', '\x00', '\x00', '\x00', '\x00', ',', 'i', 'i', 's', 'f', 'f', '\x00', '\x00', '\x00', '\x00', '\x03', '\xe8', '\xff', '\xff', '\xff', '\xff', '\x68', '\x65', '\x6c', '\x6c', '\x6f', '\x00', '\x00', '\x00', '\x3f', '\x9d', '\xf3', '\xb6', '\x40', '\xb5', '\xb2', '\x2d'}
	expected3 := NewMessage("/foo")
	expected3.AddArgument(int32(1000))
	expected3.AddArgument(int32(-1))
	expected3.AddArgument("hello")
	expected3.AddArgument(float32(1.234))
	expected3.AddArgument(float32(5.678))
	var result3 Message
	err3 := result3.UnmarshalBinary(data3)

	if err3 != nil {
		t.Error(err3)
	} else if !result3.Equals(&expected3) {
		t.Errorf("Got %v, expected %v", result3, expected3)
	}

	data4 := []byte{'/', 'b', 'y', 't', 'e', 's', '\x00', '\x00', ',', 'b', '\x00', '\x00', '\x00', '\x00', '\x00', '\x05', 'a', 'b', 'c', 'd', 'e', '\x00', '\x00', '\x00'}
	expected4 := NewMessage("/bytes")
	expected4.AddArgument([]byte{'a', 'b', 'c', 'd', 'e'})
	var result4 Message
	err4 := result4.UnmarshalBinary(data4)

	if err4 != nil {
		t.Error(err4)
	} else if !result4.Equals(&expected4) {
		t.Errorf("Got %v, expected %v", result4, expected4)
	}
}
