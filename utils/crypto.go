package utils

import (
	"crypto/rsa"
	"crypto/x509"
	"encoding/hex"
	"encoding/pem"
	"errors"
	"golang.org/x/crypto/bcrypt"
	"golang.org/x/crypto/sha3"
	"io"
	"io/ioutil"
	"os"
	"strings"
)

// LoadKeyFileRSA Loads RSA key from .pem file
//
//	path    - string, filepath to the .pem file
//
// Returns:
//
//	block   - *rsa.PrivateKey, RSA private key loaded from file
//	rest    - *rsa.PublicKey, RSA public key loaded from file
func LoadKeyFileRSA(buf io.ReadCloser) (*rsa.PrivateKey, *rsa.PublicKey, error) {
	// read raw file content
	r, err := io.ReadAll(buf)
	if err != nil {
		return nil, nil, err
	}
	// decode content into block
	block, _ := pem.Decode(r)

	// create variables to store results in
	var ok bool
	var privateKey *rsa.PrivateKey
	var publicKey *rsa.PublicKey

	if strings.Contains(block.Type, "PRIVATE") {
		// decode private key directly

		// try PKCS1 format first
		privateKey, err = x509.ParsePKCS1PrivateKey(block.Bytes)
		if err != nil && !strings.Contains(err.Error(), "ParsePKCS8PrivateKey") {
			return nil, nil, err
		}

		// try PKCS8 format
		if err != nil {
			pkey, err := x509.ParsePKCS8PrivateKey(block.Bytes)
			if err != nil {
				return nil, nil, err
			}
			privateKey = pkey.(*rsa.PrivateKey)
		}
	} else {
		// decode public key
		pKey, err := x509.ParsePKIXPublicKey(block.Bytes)
		if err != nil {
			return nil, nil, err
		}
		// attempt to cast interface{} to public key
		publicKey, ok = pKey.(*rsa.PublicKey)
		if !ok {
			return nil, nil, errors.New("failed to cast interface to public key")
		}
	}

	return privateKey, publicKey, nil
}

// LoadToken Loads authentication token from the specified path
// Args:
//
//	path   - string, the file path of the authentication token
//
// Returns:
//
//	out         - string, JWT token loaded from the passed file
func LoadToken(path string) (string, error) {
	// read raw file content
	tokenBytes, err := ioutil.ReadFile(path)
	if err != nil {
		return "", err
	}

	// load into string
	token := string(tokenBytes)

	return token, nil
}

// HashPassword Hashes a password with bcrypt
//
//	password   - string, the password that will be hashed
//
// Returns:
//
//	out        - string, a string representation of the hashed password
func HashPassword(password string) (string, error) {
	hash, err := bcrypt.GenerateFromPassword([]byte(password), 8)
	if err != nil {
		return "", err
	}
	return string(hash), nil
}

// CheckPassword Checks a hashed password against a string
//
//	password   - string, the password that will be compared
//	hash       - string, the hashed password that will be compared
//
// Returns:
//
//	out        - bool, whether the passed password matches the hash
func CheckPassword(password string, hash string) (bool, error) {
	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
	if err != nil {
		// return false with nil error is password does not match
		if err.Error() == "crypto/bcrypt: hashedPassword is not the hash of the given password" {
			return false, nil
		}
		return false, err
	}
	return true, nil
}

// HashFile Hashes the contents of a file with SHA3-256
//
//	fp   - string, the file path for the file that will be hashes
//
// Returns:
//
//	out  - string, hex encoded SHA3-256 hash made of passed file's contents
func HashFile(fp string) (string, error) {
	// open file
	file, err := os.Open(fp)
	if err != nil {
		return "", err
	}
	// delay file close
	defer file.Close()

	// create SHA256 hasher
	hasher := sha3.New256()

	// copy file data to hasher
	_, err = io.Copy(hasher, file)
	if err != nil {
		return "", err
	}

	// create output buffer
	buff := make([]byte, 32)
	// sum hash slices into buffer
	hasher.Sum(buff[:0])

	// hex encode hash and return
	return hex.EncodeToString(buff), nil
}

// HashData Hashes the passed data with SHA3-256
//
//	data   - []byte, the data that will be hashed
//
// Returns:
//
//	out    - string, hex encoded SHA3-256 hash made of the passed data
func HashData(data []byte) (string, error) {
	// create SHA256 hasher
	hasher := sha3.New256()

	// add data to hasher
	hasher.Write(data)

	// create output buffer
	buff := make([]byte, 32)
	// sum hash slices into buffer
	hasher.Sum(buff[:0])

	// hex encode hash and return
	return hex.EncodeToString(buff), nil
}
