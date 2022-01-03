package encryption

import (
	"regexp"
	"testing"
)

func TestEncodeDecode(t *testing.T) {
	key := "thisCoolKeyis§//w34#adsad"
	keyLength := len(key)
	input := "A very secure message that no one can read!"

	SetEncryptionKey(key)
	// Just making sure the string is not changed outside in some weird way
	if len(key) != keyLength {
		t.Logf("Length of key changed from %d to %d\n", keyLength, len(key))
		t.Fail()
	}

	encString, err := EncryptConfigString(input)
	if err != nil {
		t.Log("There should be no error when encrypting as string but got:")
		t.Log(err)
		t.FailNow()
	}
	if matched, err := regexp.MatchString(`ENC\(.*\)`, encString); err != nil || !matched {
		t.Logf("Encoded string '%s' does not match pattern ENC(.*)\n", encString)
		t.Fail()
	}
	encByte, err := Encrypt([]byte(input))
	if err != nil {
		t.Log("There should be no error when encrypting but got:")
		t.Log(err)
		t.Fail()
	}

	decString, err := DecryptConfigString(encString)
	if err != nil {
		t.Log("There should be no error when decrypting from string but got:")
		t.Log(err)
		t.Fail()
	}
	if input != decString {
		t.Logf("Given input '%s' does not match decoded value '%s'\n", input, decString)
		t.Fail()
	}
	decString, err = DecryptConfigString(input)
	if err != nil {
		t.Log("There should be no error when decrypting unencrypted string but got:")
		t.Log(err)
		t.Fail()
	}
	if input != decString {
		t.Logf("Given input '%s' does not match decoded value '%s'\n", input, decString)
		t.Fail()
	}

	decByte, err := Decrypt(encByte)
	if err != nil {
		t.Logf("There should be no error when decrypting but got: %v\n", err)
		t.Fail()
	}
	if input != string(decByte) {
		t.Logf("Given input '%s' does not match decoded value '%s'\n", input, string(decByte))
		t.Fail()
	}
}
