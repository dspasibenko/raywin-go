package cast

import (
	"testing"
)

func TestStringToByteArray(t *testing.T) {
	// Test case 1: Normal string conversion
	input := "hello"
	expected := []byte{'h', 'e', 'l', 'l', 'o'}

	result := StringToByteArray(input)

	if len(result) != len(expected) {
		t.Fatalf("expected length %d, got %d", len(expected), len(result))
	}

	for i := range expected {
		if result[i] != expected[i] {
			t.Errorf("expected byte %d to be %v, got %v", i, expected[i], result[i])
		}
	}

	// Test case 2: Empty string
	input = ""
	expected = []byte{}

	result = StringToByteArray(input)

	if len(result) != len(expected) {
		t.Fatalf("expected length %d, got %d", len(expected), len(result))
	}

	// Convert string to byte array
	input = ByteArrayToString([]byte{'h'})
	byteArray := StringToByteArray(input)

	// Ensure lengths match
	if len(byteArray) != 1 {
		t.Fatalf("expected length %d, got %d", len(input), len(byteArray))
	}

	// Modify the underlying byte slice to test if it affects the string
	byteArray[0] = 'H' // Changing 'h' to 'H'

	// Verify the string also reflects the change (zero-copy behavior)
	if input[0] != 'H' {
		t.Fatalf("expected input string to be affected, but it was not")
	}

	// Edge case: empty string
	empty := ""
	emptyBytes := StringToByteArray(empty)
	if len(emptyBytes) != 0 {
		t.Fatalf("expected empty slice for empty string, got length %d", len(emptyBytes))
	}
}
