package session

import "testing"

func TestEncryptDecryptServicePassword(t *testing.T) {
	pass, err := GenerateServicePassword()
	if err != nil {
		t.Fatalf("\n%s failed\n    Error: %v", t.Name(), err)
	}

	cipherText, err := EncryptServicePassword(pass, []byte("test-password"))
	if err != nil {
		t.Fatalf("\n%s failed\n    Error: %v", t.Name(), err)
	}

	plainText, err := DecryptServicePassword(cipherText, []byte("test-password"))
	if err != nil {
		t.Fatalf("\n%s failed\n    Error: %v", t.Name(), err)
	}

	if plainText != pass {
		t.Fatalf("\n%s failed\n    Error: incorrect plain text %q != %q", t.Name(), plainText, pass)
	}

	t.Logf("\n%s succeeded", t.Name())
}
