package zitimesh

import (
	"fmt"
	"github.com/openziti/sdk-golang/ziti"
	"github.com/openziti/sdk-golang/ziti/enroll"
	"strings"
)

func EnrollIdentity(token string) (*ziti.Config, error) {
	// parse the identity token
	tkn, _, err := enroll.ParseToken(strings.TrimSpace(token))
	if err != nil {
		return nil, fmt.Errorf("failed to parse token: %w", err)
	}

	fmt.Println("tkn: ", tkn)

	// enroll the identity into a configuration
	conf, err := enroll.Enroll(enroll.EnrollmentFlags{
		Token:  tkn,
		KeyAlg: "RSA",
	})
	if err != nil {
		return nil, fmt.Errorf("failed to enroll identity: %w", err)
	}

	return conf, nil
}
