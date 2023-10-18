package session

import (
	"bytes"
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"encoding/base64"
	"encoding/binary"
	"fmt"
	"golang.org/x/crypto/argon2"
)

// GenerateServicePassword
//
//	Generates a 128 byte random secret using crypto/rand
//	and returns a base64 encoded string
//
//	Returns
//	  string - base64 encoded secret
func GenerateServicePassword() (string, error) {
	// create 128 byte buffer for random secret
	b := make([]byte, 128)

	// read random bytes
	_, err := rand.Read(b)
	if err != nil {
		return "", fmt.Errorf("failed to create random secret: %v", err)
	}

	return base64.RawStdEncoding.EncodeToString(b), nil
}

// EncryptServicePassword
//
//	Encrypts the service password using AES256-GCM with
//	the passed encryption key
//
//	Args
//	  password (string) - string to be encrypted
//	  key ([]byte) - encryption key
//
//	Returns
//	  string - base64 encoded encrypted string
func EncryptServicePassword(password string, key []byte) (string, error) {
	// create slice to hold salt for cipher
	salt := make([]byte, aes.BlockSize)

	// generate a salt for cipher
	if _, err := rand.Read(salt); err != nil {
		return "", fmt.Errorf("failed to generate salt for file encryption cipher: %v", err)
	}

	// derive the cipher key from the passed key and the salt created
	cipherKey := argon2.IDKey(key, salt, 4, 60*1024, 4, 32)

	// create a new cipher block
	block, err := aes.NewCipher(cipherKey)
	if err != nil {
		return "", fmt.Errorf("failed to create cipher block: %v", err)
	}

	// create buffer to write encrypted password too
	buffer := new(bytes.Buffer)

	// write salt as the first bytes into the writer
	_, err = buffer.Write(salt)
	if err != nil {
		return "", fmt.Errorf("failed to write salt to file: %v", err)
	}

	// convert password to bytes
	payload := []byte(password)

	// write 8 bytes of content size as the second chunk to writer
	err = binary.Write(buffer, binary.LittleEndian, uint64(len(payload)))
	if err != nil {
		return "", fmt.Errorf("failed to write paylod size: %v", err)
	}

	// create a new GCM encrypter
	encrypter, err := cipher.NewGCM(block)
	if err != nil {
		return "", fmt.Errorf("failed to create gcm encrypter: %v", err)
	}

	// create slice to hold nonce for cipher
	nonce := make([]byte, encrypter.NonceSize())

	// generate a new nonce this block in the cipher
	if _, err := rand.Read(nonce); err != nil {
		return "", fmt.Errorf("failed to generate nonce for encryption cipher block: %v", err)
	}

	// write nonce to buffer
	_, err = buffer.Write(nonce)
	if err != nil {
		return "", fmt.Errorf("failed to write nonce to file: %v", err)
	}

	// pad buffer to a block size multiple if it is not
	if len(payload)%aes.BlockSize != 0 {
		// calculate the size to pad
		padSize := aes.BlockSize - (len(payload) % aes.BlockSize)

		// create slice to hold padding
		padding := make([]byte, padSize)

		// load random values into padding
		if _, err := rand.Read(padding); err != nil {
			return "", fmt.Errorf("failed to generate buffer to pad block during encryption: %v", err)
		}

		// append padding to the payload
		payload = append(payload, padding...)
	}

	// encrypt buffer and write to cipher buffer
	cipherBuf := encrypter.Seal(nil, nonce, payload, nil)

	// write cipher buffer to writer
	_, err = buffer.Write(cipherBuf)
	if err != nil {
		return "", fmt.Errorf("failed to write cipher buffer to writer: %v", err)
	}

	return base64.RawStdEncoding.EncodeToString(buffer.Bytes()), nil
}

// DecryptServicePassword
//
//		Decrypts the service password using AES256-GCM with
//		the passed encryption key
//
//		Args
//	      - cipherText (string) - base64 encoded encrypted payload
//		  - key ([]byte) - encryption key
//
//		Returns
//		  string - decrypted string
func DecryptServicePassword(cipherText string, key []byte) (string, error) {
	// decode cipher text from base64
	cipherBytes, err := base64.RawStdEncoding.DecodeString(cipherText)
	if err != nil {
		return "", fmt.Errorf("failed to decode cipher text: %v", err)
	}

	// create buffer from cipher text
	buffer := bytes.NewBuffer(cipherBytes)

	// create slice to hold salt for cipher
	salt := make([]byte, aes.BlockSize)

	// read the first BLOCK_SIZE bytes to retrieve the key salt
	_, err = buffer.Read(salt)
	if err != nil {
		return "", fmt.Errorf("failed to read salt from reader: %v", err)
	}

	// create slice to hold little endian encoded payload size
	payloadSizeBytes := make([]byte, 8)

	// read the second 8 bytes to retrieve the payload size
	_, err = buffer.Read(payloadSizeBytes)
	if err != nil {
		return "", fmt.Errorf("failed to payload size bytes from reader: %v", err)
	}

	// load payload from little endian
	payloadSize := int(binary.LittleEndian.Uint64(payloadSizeBytes))

	// derive the cipher key from the passed key and the salt created
	cipherKey := argon2.IDKey(key, salt, 4, 60*1024, 4, 32)

	// create a new cipher block
	block, err := aes.NewCipher(cipherKey)
	if err != nil {
		return "", fmt.Errorf("failed to create cipher block: %v", err)
	}

	// create a new GCM decrypter
	decrypter, err := cipher.NewGCM(block)
	if err != nil {
		return "", fmt.Errorf("failed to create gcm decrypter: %v", err)
	}

	// create slice to hold nonce
	nonce := make([]byte, decrypter.NonceSize())

	// read nonce from buffer
	_, err = buffer.Read(nonce)
	if err != nil {
		return "", fmt.Errorf("failed to read nonce from reader: %v", err)
	}

	// decrypt buffer and write to plain buffer
	plainBuf, err := decrypter.Open(nil, nonce, buffer.Bytes(), nil)
	if err != nil {
		return "", fmt.Errorf("failed to decrypt block: %v", err)
	}

	// ensure that plain text buffer is at least payload size
	if len(plainBuf) < payloadSize {
		return "", fmt.Errorf("plain text buffer is shorter than expected %d != %d", len(plainBuf), payloadSize)
	}

	return string(plainBuf[:payloadSize]), nil
}
