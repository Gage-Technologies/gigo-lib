package network

import (
	"errors"
	"fmt"
	"github.com/buger/jsonparser"
	"io/ioutil"
	"net"
	"net/http"
	"strings"
)

var cidrs []*net.IPNet

// initialize constant variables
func init() {
	maxCidrBlocks := []string{
		"127.0.0.1/8",    // localhost
		"10.0.0.0/8",     // 24-bit block
		"172.16.0.0/12",  // 20-bit block
		"169.254.0.0/26", // link local address
		"192.168.0.0/24", // 16-bit block
		"::1/128",        // localhost IPv6
		"fc00::/7",       // unique local address IPv6
		"fe80::/10",      // link local address IPv6
	}

	// create array to hold CIDR blocks
	cidrs = make([]*net.IPNet, len(maxCidrBlocks))
	// parse strings into CIDR objects
	for i, maxCidrBlock := range maxCidrBlocks {
		// parse address string into CIDR block
		_, cidr, _ := net.ParseCIDR(maxCidrBlock)
		// add CIDR block into array
		cidrs[i] = cidr
	}
}

// Determines if an IP address is private or not
// Checks if the address in within a private CIDR block
//
//	address   - string, IP address that will be checked for being private
//
// Returns:
//
//	out       - bool, whether the passed IP address is private
func privateAddress(address string) (bool, error) {
	// parse ip address into ip object
	ip := net.ParseIP(address)

	// exit if address was not a valid
	if ip == nil {
		return false, errors.New("invalid ip address")
	}

	// check if ip is classified as private
	if ip.IsLoopback() || ip.IsLinkLocalUnicast() || ip.IsLinkLocalMulticast() || ip.IsPrivate() {
		return true, nil
	}

	// check if ip address is within any private CIDR blocks
	for i := range cidrs {
		// return true if ip is within private block
		if cidrs[i].Contains(ip) {
			return true, nil
		}
	}

	// exit with false if ip was not in ant private blocks
	return false, nil
}

// GetRemoteAddr Extracts the remote IP address from a request object
// If no IP address ForwardFor header is found then the Host header will be used
//
//	r     - *http.Request, request object that will be used to retrieve remote IP address
//
// Returns:
//
//	out   - string, remote IP address for the passed request object
func GetRemoteAddr(r *http.Request, publicIp string) string {
	// retrieve header that could contain the remote IP address
	xRealIP := r.Header.Get("X-Real-Ip")
	xForwardedFor := r.Header.Values("X-Forwarded-For")
	xOriginalForwardedFor := r.Header.Values("X-Original-Forwarded-For")

	// return Host header if neither forward header are filled
	if len(xRealIP) == 0 && len(xForwardedFor) == 0 && len(xOriginalForwardedFor) == 0 {
		// retrieve Host header
		remoteIP := r.RemoteAddr

		// remove port from address
		if strings.Contains(remoteIP, ":") {
			remoteIP, _, _ = net.SplitHostPort(r.RemoteAddr)
		}

		// return for private ip
		if p, _ := privateAddress(remoteIP); p {
			return publicIp
		}

		// return ip address
		return remoteIP
	}

	// parse X-Original-Forwarded-For header if it exists
	if len(xOriginalForwardedFor) > 0 {
		for i, h := range xOriginalForwardedFor {
			// split header into array of addresses
			addresses := strings.Split(h, ",")

			// scraper for the first public address
			for j, address := range addresses {
				// clip any whitespace from the address string
				address = strings.TrimSpace(address)

				// check if the address is private
				private, err := privateAddress(address)
				// return address if it is public
				if !private && err == nil {
					return address
				}

				// return last address in list if all were private
				if i == len(xOriginalForwardedFor)-1 && j == len(addresses)-1 {
					// return for private ip
					if p, _ := privateAddress(address); p {
						return publicIp
					}
					return address
				}
			}
		}
	}

	// parse X-Forward-For header if it exists
	if len(xForwardedFor) > 0 {
		for i, h := range xForwardedFor {
			// split header into array of addresses
			addresses := strings.Split(h, ",")

			// scraper for the first public address
			for j, address := range addresses {
				// clip any whitespace from the address string
				address = strings.TrimSpace(address)

				// check if the address is private
				private, err := privateAddress(address)
				// return address if it is public
				if !private && err == nil {
					return address
				}

				// return last address in list if all were private
				if i == len(xForwardedFor)-1 && j == len(addresses)-1 {
					// return for private ip
					if p, _ := privateAddress(address); p {
						return publicIp
					}
					return address
				}
			}
		}
	}

	// return for private or empty ip
	if p, _ := privateAddress(xRealIP); p || len(xRealIP) == 0 {
		return publicIp
	}

	// return X-Real-Ip if X-Forward-For was empty
	return xRealIP
}

// GetRequestIP Extracts the remote IP address from a request object
// If no IP address ForwardFor header is found then the Host header will be used
//
//	r     - *http.Request, request object that will be used to retrieve remote IP address
//
// Returns:
//
//	out   - string, remote IP address for the passed request object
func GetRequestIP(r *http.Request) string {
	// retrieve header that could contain the remote IP address
	xRealIP := r.Header.Get("X-Real-Ip")
	xForwardedFor := r.Header.Values("X-Forwarded-For")
	xOriginalForwardedFor := r.Header.Values("X-Original-Forwarded-For")

	// return Host header if neither forward header are filled
	if len(xRealIP) == 0 && len(xForwardedFor) == 0 && len(xOriginalForwardedFor) == 0 {
		// retrieve Host header
		remoteIP := r.RemoteAddr

		// remove port from address
		if strings.Contains(remoteIP, ":") {
			remoteIP, _, _ = net.SplitHostPort(r.RemoteAddr)
		}

		// return ip address
		return remoteIP
	}

	// parse X-Forward-For header if it exists
	if len(xOriginalForwardedFor) > 0 {
		for i, h := range xOriginalForwardedFor {
			// split header into array of addresses
			addresses := strings.Split(h, ",")

			// scraper for the first public address
			for j, address := range addresses {
				// clip any whitespace from the address string
				address = strings.TrimSpace(address)

				// check if the address is private
				private, err := privateAddress(address)
				// return address if it is public
				if !private && err == nil {
					return address
				}

				// return last address in list if all were private
				if i == len(xOriginalForwardedFor)-1 && j == len(addresses)-1 {
					return address
				}
			}
		}
	}

	// parse X-Forward-For header if it exists
	if len(xForwardedFor) > 0 {
		for i, h := range xForwardedFor {
			// split header into array of addresses
			addresses := strings.Split(h, ",")

			// scraper for the first public address
			for j, address := range addresses {
				// clip any whitespace from the address string
				address = strings.TrimSpace(address)

				// check if the address is private
				private, err := privateAddress(address)
				// return address if it is public
				if !private && err == nil {
					return address
				}

				// return last address in list if all were private
				if i == len(xForwardedFor)-1 && j == len(addresses)-1 {
					return address
				}
			}
		}
	}

	// return X-Real-Ip if X-Forward-For was empty
	return xRealIP
}

// GetSelfPublicIP Retrieves the current servers remote IP address
// Returns:
//
//	out   - string, remote IP address for this server
func GetSelfPublicIP() (string, error) {
	// execute api request for public ip
	res, err := http.Get("https://api.ipify.org?format=json")
	if err != nil {
		return "", fmt.Errorf("failed execute public ip api request: %v", err)
	}

	// attempt to read body of api response
	b, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return "", fmt.Errorf("failed execute public ip api request: %v", err)
	}

	// retrieve the ip from the api json response
	data, _, _, err := jsonparser.Get(b, "ip")
	if err != nil {
		return "", fmt.Errorf("failed execute public ip api request: %v", err)
	}

	// ensure response is not empty
	if len(data) == 0 {
		return "", fmt.Errorf("failed execute public ip api request: empty ip field returned")
	}

	// return ip
	return string(data), nil
}
