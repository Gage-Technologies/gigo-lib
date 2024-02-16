package utils

import (
	"crypto/rand"
	"encoding/hex"
	"errors"
	"fmt"
	"github.com/gage-technologies/gigo-lib/storage"
	"github.com/golang-jwt/jwt"
	"net/url"
	"strconv"
	"strings"
	"time"
)

const SkipIpValidation = "skip_ip_validation"

// NormalizeURLForHash
// Normalizes a URL for hash calculation
// Strips url of all query parameters and fragments
// Args:
//
//	uri    - string, url that will be normalized for hashing
//
// Returns:
//
//	out    - string, url that has been stripped of all of its query parameters and fragments
func NormalizeURLForHash(uri string) (string, error) {
	urlObj, err := url.Parse(uri)
	if err != nil {
		return "", err
	}

	return fmt.Sprintf("%s://%s%s", urlObj.Scheme, urlObj.Host, urlObj.Path), nil
}

// CreateInternalJWT Creates a JWT for internal_api API authentication
//
//			service   - string, the service that this JWT will be used for
//			ip        - string, the IP address of the machine that will be using this JWT
//	     hours     - int, amount of hours the JWT will be active for
//
// Returns:
//
//	out       - string, JWT containing the pertinent claims and signed by the global private key
func CreateInternalJWT(storageEngine storage.Storage, service string, ip string, hours int) (string, error) {
	// retrieve key
	buf, _, err := storageEngine.GetFile("keys/private.pem")
	if err != nil {
		return "", fmt.Errorf("failed to retrieve private key: %v", err)
	}

	// load private key from filesystem
	privateKey, _, err := LoadKeyFileRSA(buf)
	if err != nil {
		return "", err
	}

	// create claims
	claims := jwt.MapClaims{}
	// add ip to claims
	claims["ip"] = ip
	// add service to claims
	claims["service"] = service
	// add expiration to claims
	exp := time.Duration(int(time.Hour) * hours)
	claims["exp"] = time.Now().Add(exp).Unix()
	// sign claims with private key for validation at a later time
	signedContent := jwt.NewWithClaims(jwt.SigningMethodRS256, claims)
	token, err := signedContent.SignedString(privateKey)
	if err != nil {
		return "", err
	}

	return token, nil
}

// CreateExternalJWT Creates a JWT for http API authentication
//
//			userID    - string, the user ID that this JWT will be used for
//			ip        - string, the IP address of the machine that will be using this JWT
//	     hours     - int, amount of hours the JWT will be active for
//			minutes   - int, amount of minutes the JWT will be active for
//			payload   - map[string]interface{}, optional set of key, value pairs to be added to the JWT payload
//
// Returns:
//
//	out      - string, JWT containing the pertinent claims and signed by the global private key
func CreateExternalJWT(storageEngine storage.Storage, userID string, ip string, hours int, minutes int, payload map[string]interface{}) (string, error) {
	// retrieve key
	buf, _, err := storageEngine.GetFile("keys/private.pem")
	if err != nil {
		return "", fmt.Errorf("failed to retrieve private key: %v", err)
	}

	// load private key from filesystem
	privateKey, _, err := LoadKeyFileRSA(buf)
	if err != nil {
		return "", err
	}

	// create claims
	claims := jwt.MapClaims{}

	// add ip to claims
	claims["ip"] = ip

	// add user to claims
	claims["user"] = userID

	// add expiration to claims
	exp := time.Duration((int(time.Hour) * hours) + (int(time.Minute) * minutes))
	claims["exp"] = time.Now().Add(exp).Unix()

	// check if additional payload was passed
	if payload != nil {
		// add each value from payload to the claims of the JWT
		for k, v := range payload {
			claims[k] = v
		}
	}

	// sign claims with private key for validation at a later time
	signedContent := jwt.NewWithClaims(jwt.SigningMethodRS256, claims)
	token, err := signedContent.SignedString(privateKey)
	if err != nil {
		return "", err
	}

	return token, nil
}

// CreateExternalJWTNoIP Creates a JWT for http API authentication
//
//			userID    - string, the user ID that this JWT will be used for
//			ip        - string, the IP address of the machine that will be using this JWT
//	     hours     - int, amount of hours the JWT will be active for
//			minutes   - int, amount of minutes the JWT will be active for
//			payload   - map[string]interface{}, optional set of key, value pairs to be added to the JWT payload
//
// Returns:
//
//	out      - string, JWT containing the pertinent claims and signed by the global private key
func CreateExternalJWTNoIP(storageEngine storage.Storage, userID string, hours int, minutes int, payload map[string]interface{}) (string, error) {
	// retrieve key
	buf, _, err := storageEngine.GetFile("keys/private.pem")
	if err != nil {
		return "", fmt.Errorf("failed to retrieve private key: %v", err)
	}

	// load private key from filesystem
	privateKey, _, err := LoadKeyFileRSA(buf)
	if err != nil {
		return "", err
	}

	// create claims
	claims := jwt.MapClaims{}

	// add user to claims
	claims["user"] = userID

	// add expiration to claims
	exp := time.Duration((int(time.Hour) * hours) + (int(time.Minute) * minutes))
	claims["exp"] = time.Now().Add(exp).Unix()

	// check if additional payload was passed
	if payload != nil {
		// add each value from payload to the claims of the JWT
		for k, v := range payload {
			claims[k] = v
		}
	}

	// sign claims with private key for validation at a later time
	signedContent := jwt.NewWithClaims(jwt.SigningMethodRS256, claims)
	token, err := signedContent.SignedString(privateKey)
	if err != nil {
		return "", err
	}

	return token, nil
}

// ValidateInternalJWT Validates a JWT for internal_api API authentication
//
//	     tokenString   - string, JWT string that will be validated
//			service       - string, service that this JWT will be used for
//			ip            - string, IP address of the machine that will be using this JWT
//			validateIP    - bool, validate the IP stored inside the JWT
//
// Returns:
//
//	out           - bool, success status for the JWT validation
func ValidateInternalJWT(storageEngine storage.Storage, tokenString string, service string, ip string, validateIP bool) (bool, error) {
	// retrieve key
	buf, _, err := storageEngine.GetFile("keys/public.pem")
	if err != nil {
		return false, fmt.Errorf("failed to retrieve public key: %v", err)
	}

	// load private key from filesystem
	_, publicKey, err := LoadKeyFileRSA(buf)
	if err != nil {
		return false, err
	}

	// decode JWT
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		// check that signature type is the same as expected
		if _, ok := token.Method.(*jwt.SigningMethodRSA); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}

		// return private key
		return publicKey, nil
	})

	if err != nil {
		// check if error is the result of incorrect algorithm; this means tampering
		if strings.Contains(err.Error(), "unexpected signing method") {
			return false, err
		}

		jwtErr, ok := err.(*jwt.ValidationError)
		if ok {
			// check if error is the result of expiration or premature use
			if jwtErr.Errors&(jwt.ValidationErrorExpired|jwt.ValidationErrorNotValidYet) != 0 {
				return false, errors.New("invalid JWT time")
			}

			// check if JWT is re-signed or if signature is broken; this means tampering
			if jwtErr.Errors&jwt.ValidationErrorSignatureInvalid != 0 {
				return false, errors.New("invalid signature")
			}
		}

		// return generic failure
		return false, err
	}

	// extract claims
	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok || !token.Valid {
		return false, errors.New("failed to parse claims")
	}

	// check the validity of the service claim
	if val, ok := claims["service"]; ok {
		// only check service if specific service is passed
		if service != "all" {
			if val != service {
				return false, errors.New(fmt.Sprintf("incorrect service; expected %s; received %s", val, service))
			}
		} else {
			// fill string with all service names to make for easy multi point check
			allServices := "analyzer extractor scraper manager"

			// if all tokens are allowed validate that token is from service
			if !strings.Contains(allServices, val.(string)) {
				return false, errors.New(fmt.Sprintf("incorrect service; non working token passed"))
			}
		}
	} else {
		return false, errors.New(fmt.Sprintf("missing service"))
	}

	// check the validity of the IP if selected
	if validateIP {
		// check the validity of the ip claim
		if val, ok := claims["ip"]; ok {
			if val != ip {
				return false, errors.New(fmt.Sprintf("incorrect ip; expected %s; received %s", val, ip))
			}
		} else {
			return false, errors.New(fmt.Sprintf("missing ip"))
		}
	}

	// success
	return true, nil
}

// ValidateExternalJWT Validates a JWT for internal_api API authentication
//
//	     tokenString   - string, JWT string that will be validated
//			ip            - string, IP address of the machine that will be using this JWT
//			payload       - map[string]interface{}, optional set of key, value pairs to be checked for in the JWT payload
//
// Returns:
//
//	valid         - bool, success status for the JWT validation
//	userID        - string, id of the user stored in the JWT
//	payload       - map[string]interface, payload of the jwt
func ValidateExternalJWT(storageEngine storage.Storage, tokenString string, ip string, payload map[string]interface{}) (bool, int64, map[string]interface{}, error) {
	// retrieve key
	buf, _, err := storageEngine.GetFile("keys/public.pem")
	if err != nil {
		return false, 0, nil, fmt.Errorf("failed to retrieve public key: %v", err)
	}

	// load private key from filesystem
	_, publicKey, err := LoadKeyFileRSA(buf)
	if err != nil {
		return false, 0, nil, err
	}

	// decode JWT
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		// check that signature type is the same as expected
		if _, ok := token.Method.(*jwt.SigningMethodRSA); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}

		// return private key
		return publicKey, nil
	})

	if err != nil {
		// check if error is the result of incorrect algorithm; this means tampering
		if strings.Contains(err.Error(), "unexpected signing method") {
			return false, 0, nil, err
		}

		jwtErr, ok := err.(*jwt.ValidationError)
		if ok {
			// check if error is the result of expiration or premature use
			if jwtErr.Errors&(jwt.ValidationErrorExpired|jwt.ValidationErrorNotValidYet) != 0 {
				return false, 0, nil, errors.New("invalid JWT time")
			}

			// check if JWT is re-signed or if signature is broken; this means tampering
			if jwtErr.Errors&jwt.ValidationErrorSignatureInvalid != 0 {
				return false, 0, nil, errors.New("invalid signature")
			}
		}

		// return generic failure
		return false, 0, nil, err
	}

	// extract claims
	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok || !token.Valid {
		return false, 0, nil, errors.New("failed to parse claims")
	}

	// check the validity of the ip claim
	if val, ok := claims["ip"]; ok {
		if ip != SkipIpValidation && val != ip {
			return false, 0, nil, errors.New(fmt.Sprintf("incorrect ip; expected %s; received %s", val, ip))
		}
	} else {
		return false, 0, nil, errors.New("missing ip")
	}

	// retrieve the userID from the JWT
	userID, ok := claims["user"]
	if !ok {
		return false, 0, nil, errors.New("missing user id")
	}

	// check if optional payloads were passed
	if payload != nil {
		// check that each payload is in the JWT
		for k, v := range payload {
			// attempt to read key value from claims
			val, ok := claims[k]

			// return error if key did not exist
			if !ok {
				return false, 0, nil, errors.New(fmt.Sprintf("missing payload %s", k))
			}

			// return if payload value did not match claim value
			if val != v {
				return false, 0, nil, errors.New(fmt.Sprintf("incorrect payload %s", k))
			}
		}
	}

	uID, err := strconv.ParseInt(userID.(string), 10, 64)
	if err != nil {
		return false, 0, nil, errors.New(fmt.Sprintf("invalid user id: %s", userID))
	}

	// success
	return true, uID, claims, nil
}

func GenerateEmailToken() (string, error) {
	// hold generated token
	token := make([]byte, 16)

	// generate random token
	_, err := rand.Read(token)
	if err != nil {
		return "", err
	}

	// base64 encode
	return hex.EncodeToString(token), nil
}
