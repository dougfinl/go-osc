package osc

import (
	"bytes"
	"testing"
)

func TestEncodeBundle(t *testing.T) {
	// A new bundle should only encode the header and immediate time tag
	test1 := NewBundle()
	expected1 := []byte{'#', 'b', 'u', 'n', 'd', 'l', 'e', '\x00', '\x00', '\x00', '\x00', '\x00', '\x00', '\x00', '\x00', '\x01'}
	result1, err1 := test1.MarshalBinary()

	if err1 != nil {
		t.Error(err1)
	} else if !bytes.Equal(result1, expected1) {
		t.Errorf("Got %v, but expected %v", result1, expected1)
	}

	// Bundle with an empty child bundle
	test2 := NewBundle()
	test2.AddPacket(NewBundle())
	expected2 := []byte{'#', 'b', 'u', 'n', 'd', 'l', 'e', '\x00', '\x00', '\x00', '\x00', '\x00', '\x00', '\x00', '\x00', '\x01', '\x00', '\x00', '\x00', '\x10', '#', 'b', 'u', 'n', 'd', 'l', 'e', '\x00', '\x00', '\x00', '\x00', '\x00', '\x00', '\x00', '\x00', '\x01'}
	result2, err2 := test2.MarshalBinary()

	if err2 != nil {
		t.Error(err2)
	} else if !bytes.Equal(result2, expected2) {
		t.Errorf("Got %v, but expected %v", result2, expected2)
	}

	// Bundle with two empty child bundles
	test3 := NewBundle()
	test3.AddPacket(NewBundle())
	test3.AddPacket(NewBundle())
	expected3 := []byte{'#', 'b', 'u', 'n', 'd', 'l', 'e', '\x00', '\x00', '\x00', '\x00', '\x00', '\x00', '\x00', '\x00', '\x01', '\x00', '\x00', '\x00', '\x10', '#', 'b', 'u', 'n', 'd', 'l', 'e', '\x00', '\x00', '\x00', '\x00', '\x00', '\x00', '\x00', '\x00', '\x01', '\x00', '\x00', '\x00', '\x10', '#', 'b', 'u', 'n', 'd', 'l', 'e', '\x00', '\x00', '\x00', '\x00', '\x00', '\x00', '\x00', '\x00', '\x01'}
	result3, err3 := test3.MarshalBinary()

	if err3 != nil {
		t.Error(err3)
	} else if !bytes.Equal(result3, expected3) {
		t.Errorf("Got %v, but expected %v", result3, expected3)
	}

	// Bundle with a child message and a child bundle
	test4 := NewBundle()
	msg4 := NewMessage("/foo")
	msg4.AddArgument([]byte{'a', 'r', 'g'})
	test4.AddPacket(msg4)
	test4.AddPacket(NewBundle())
	expected4 := []byte{'#', 'b', 'u', 'n', 'd', 'l', 'e', '\x00', '\x00', '\x00', '\x00', '\x00', '\x00', '\x00', '\x00', '\x01', '\x00', '\x00', '\x00', '\x14', '/', 'f', 'o', 'o', '\x00', '\x00', '\x00', '\x00', ',', 'b', '\x00', '\x00', '\x00', '\x00', '\x00', '\x03', 'a', 'r', 'g', '\x00', '\x00', '\x00', '\x00', '\x10', '#', 'b', 'u', 'n', 'd', 'l', 'e', '\x00', '\x00', '\x00', '\x00', '\x00', '\x00', '\x00', '\x00', '\x01'}
	result4, err4 := test4.MarshalBinary()

	if err4 != nil {
		t.Error(err4)
	} else if !bytes.Equal(result4, expected4) {
		t.Errorf("Got %v, but expected %v", result4, expected4)
	}
}

func TestDecodeBundle(t *testing.T) {
	// Empty bundle
	test1 := []byte{'#', 'b', 'u', 'n', 'd', 'l', 'e', '\x00', '\x00', '\x00', '\x00', '\x00', '\x00', '\x00', '\x00', '\x01'}
	expected1 := NewBundle()
	result1, err1 := NewBundleFromData(test1)

	if err1 != nil {
		t.Error(err1)
	} else if !result1.Equals(expected1) {
		t.Errorf("Got %v, but expected %v", result1, expected1)
	}

	// Bundle with an empty child bundle
	test2 := []byte{'#', 'b', 'u', 'n', 'd', 'l', 'e', '\x00', '\x00', '\x00', '\x00', '\x00', '\x00', '\x00', '\x00', '\x01', '\x00', '\x00', '\x00', '\x10', '#', 'b', 'u', 'n', 'd', 'l', 'e', '\x00', '\x00', '\x00', '\x00', '\x00', '\x00', '\x00', '\x00', '\x01'}
	expected2 := NewBundle()
	expected2.AddPacket(NewBundle())
	result2, err2 := NewBundleFromData(test2)

	if err2 != nil {
		t.Error(err2)
	} else if !result2.Equals(expected2) {
		t.Errorf("Got %v, but expected %v", result2, expected2)
	}

	// Bundle with two empty child bundles
	test3 := []byte{'#', 'b', 'u', 'n', 'd', 'l', 'e', '\x00', '\x00', '\x00', '\x00', '\x00', '\x00', '\x00', '\x00', '\x01', '\x00', '\x00', '\x00', '\x10', '#', 'b', 'u', 'n', 'd', 'l', 'e', '\x00', '\x00', '\x00', '\x00', '\x00', '\x00', '\x00', '\x00', '\x01', '\x00', '\x00', '\x00', '\x10', '#', 'b', 'u', 'n', 'd', 'l', 'e', '\x00', '\x00', '\x00', '\x00', '\x00', '\x00', '\x00', '\x00', '\x01'}
	expected3 := NewBundle()
	expected3.AddPacket(NewBundle())
	expected3.AddPacket(NewBundle())
	result3, err3 := NewBundleFromData(test3)

	if err3 != nil {
		t.Error(err3)
	} else if !result3.Equals(expected3) {
		t.Errorf("Got %v, but expected %v", result3, expected3)
	}

	// Bundle with a child message and a child bundle
	test4 := []byte{'#', 'b', 'u', 'n', 'd', 'l', 'e', '\x00', '\x00', '\x00', '\x00', '\x00', '\x00', '\x00', '\x00', '\x01', '\x00', '\x00', '\x00', '\x14', '/', 'f', 'o', 'o', '\x00', '\x00', '\x00', '\x00', ',', 'b', '\x00', '\x00', '\x00', '\x00', '\x00', '\x03', 'a', 'r', 'g', '\x00', '\x00', '\x00', '\x00', '\x10', '#', 'b', 'u', 'n', 'd', 'l', 'e', '\x00', '\x00', '\x00', '\x00', '\x00', '\x00', '\x00', '\x00', '\x01'}
	expected4 := NewBundle()
	msg4 := NewMessage("/foo")
	msg4.AddArgument([]byte{'a', 'r', 'g'})
	expected4.AddPacket(msg4)
	expected4.AddPacket(NewBundle())
	result4, err4 := NewBundleFromData(test4)

	if err4 != nil {
		t.Error(err4)
	} else if !result4.Equals(expected4) {
		t.Errorf("Got %v, but expected %v", result4, expected4)
	}
}
