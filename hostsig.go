// Copyright (c) 2016, Gareth Watts
// All rights reserved.

package main

import (
	"bytes"
	"encoding/base64"
	"errors"
	"fmt"
	"strings"

	"golang.org/x/crypto/ssh"
)

type hostSig struct {
	sigType string
	data    []byte
}

func hostSigFromString(s string) (sig *hostSig, err error) {
	parts := strings.SplitN(s, " ", 2)
	if len(parts) != 2 {
		return nil, errors.New("Invalid signature")
	}
	data, err := base64.StdEncoding.DecodeString(parts[1])
	if err != nil {
		return nil, fmt.Errorf("Invalid signature: %v", err)
	}
	return &hostSig{sigType: parts[0], data: data}, nil
}

func hostSigFromKey(key ssh.PublicKey) *hostSig {
	return &hostSig{
		sigType: key.Type(),
		data:    key.Marshal(),
	}
}

func (s hostSig) String() string {
	return s.sigType + " " + base64.StdEncoding.EncodeToString(s.data)
}

func (s hostSig) isEqual(other *hostSig) bool {
	return s.sigType == other.sigType && bytes.Equal(s.data, other.data)
}
