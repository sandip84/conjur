/*
MIT License

Copyright (c) 2017 Aleksandr Fedotov

Permission is hereby granted, free of charge, to any person obtaining a copy
of this software and associated documentation files (the "Software"), to deal
in the Software without restriction, including without limitation the rights
to use, copy, modify, merge, publish, distribute, sublicense, and/or sell
copies of the Software, and to permit persons to whom the Software is
furnished to do so, subject to the following conditions:

The above copyright notice and this permission notice shall be included in all
copies or substantial portions of the Software.

THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN THE
SOFTWARE.
*/

package main

import (
	"bytes"
	"encoding/binary"
	"testing"

	"github.com/conjurinc/secretless-broker/internal/app/secretless/handlers/mysql/protocol"
	"github.com/stretchr/testify/assert"
)

func TestUnpackOkResponse(t *testing.T) {

	type UnpackOkResponseAssert struct {
		Packet   []byte
		HasError bool
		Error    error
		protocol.OkResponse
	}

	testData := []*UnpackOkResponseAssert{
		{
			[]byte{
				0x30, 0x00, 0x00, 0x01, 0x00, 0x01, 0x00, 0x22, 0x00, 0x00, 0x00, 0x28, 0x52, 0x6f, 0x77, 0x73,
				0x20, 0x6d, 0x61, 0x74, 0x63, 0x68, 0x65, 0x64, 0x3a, 0x20, 0x31, 0x20, 0x20, 0x43, 0x68, 0x61,
				0x6e, 0x67, 0x65, 0x64, 0x3a, 0x20, 0x31, 0x20, 0x20, 0x57, 0x61, 0x72, 0x6e, 0x69, 0x6e, 0x67,
				0x73, 0x3a, 0x20, 0x30,
			},
			false,
			nil,
			protocol.OkResponse{
				PacketType:   0x00,
				AffectedRows: uint64(1),
				LastInsertID: uint64(0),
				StatusFlags:  uint16(34),
				Warnings:     uint16(0)},
		},
		{
			[]byte{0x07, 0x00, 0x00, 0x02, 0x00, 0x00, 0x00, 0x02, 0x00, 0x00, 0x00},
			false,
			nil,
			protocol.OkResponse{
				PacketType:   0x00,
				AffectedRows: uint64(0),
				LastInsertID: uint64(0),
				StatusFlags:  uint16(2),
				Warnings:     uint16(0)},
		},
		{
			[]byte{0x07, 0x00, 0x00, 0x01, 0x00, 0x01, 0x02, 0x02, 0x00, 0x00, 0x00},
			false,
			nil,
			protocol.OkResponse{
				PacketType:   0x00,
				AffectedRows: uint64(1),
				LastInsertID: uint64(2),
				StatusFlags:  uint16(2),
				Warnings:     uint16(0)},
		},
	}

	for _, asserted := range testData {
		decoded, err := protocol.UnpackOkResponse(asserted.Packet)

		assert.Nil(t, err)

		if err == nil {
			assert.Equal(t, asserted.OkResponse.PacketType, decoded.PacketType)
			assert.Equal(t, asserted.OkResponse.AffectedRows, decoded.AffectedRows)
			assert.Equal(t, asserted.OkResponse.LastInsertID, decoded.LastInsertID)
			assert.Equal(t, asserted.OkResponse.StatusFlags, decoded.StatusFlags)
			assert.Equal(t, asserted.OkResponse.Warnings, decoded.Warnings)
		}
	}
}

func TestUnpackHandshakeV10(t *testing.T) {

	type UnpackHandshakeV10Assert struct {
		Packet   []byte
		HasError bool
		Error    error
		protocol.HandshakeV10
		CapabilitiesMap map[uint32]bool
	}

	testData := []*UnpackHandshakeV10Assert{
		{
			[]byte{
				0x4a, 0x00, 0x00, 0x00, 0x0a, 0x35, 0x2e, 0x35, 0x2e, 0x35, 0x36, 0x00, 0x5e, 0x06, 0x00, 0x00,
				0x48, 0x6a, 0x5b, 0x6a, 0x24, 0x71, 0x30, 0x3a, 0x00, 0xff, 0xf7, 0x08, 0x02, 0x00, 0x0f, 0x80,
				0x15, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x6f, 0x43, 0x40, 0x56, 0x6e,
				0x4b, 0x68, 0x4a, 0x79, 0x46, 0x30, 0x5a, 0x00, 0x6d, 0x79, 0x73, 0x71, 0x6c, 0x5f, 0x6e, 0x61,
				0x74, 0x69, 0x76, 0x65, 0x5f, 0x70, 0x61, 0x73, 0x73, 0x77, 0x6f, 0x72, 0x64, 0x00,
			},
			false,
			nil,
			protocol.HandshakeV10{
				ProtocolVersion:    byte(10),
				ServerVersion:      "5.5.56",
				ConnectionID:       uint32(1630),
				AuthPlugin:         "mysql_native_password",
				ServerCapabilities: binary.LittleEndian.Uint32([]byte{255, 247, 15, 128}),
				Salt: []byte{0x48, 0x6a, 0x5b, 0x6a, 0x24, 0x71, 0x30, 0x3a, 0x6f, 0x43, 0x40, 0x56, 0x6e, 0x4b,
					0x68, 0x4a, 0x79, 0x46, 0x30, 0x5a},
			},
			map[uint32]bool{
				protocol.ClientLongPassword: true, protocol.ClientFoundRows: true, protocol.ClientLongFlag: true,
				protocol.ClientConnectWithDB: true, protocol.ClientNoSchema: true, protocol.ClientCompress: true, protocol.ClientODBC: true,
				protocol.ClientLocalFiles: true, protocol.ClientIgnoreSpace: true, protocol.ClientProtocol41: true, protocol.ClientInteractive: true,
				protocol.ClientSSL: false, protocol.ClientIgnoreSIGPIPE: true, protocol.ClientTransactions: true, protocol.ClientMultiStatements: true,
				protocol.ClientMultiResults: true, protocol.ClientPSMultiResults: true, protocol.ClientPluginAuth: true, protocol.ClientConnectAttrs: false,
				protocol.ClientPluginAuthLenEncClientData: false, protocol.ClientCanHandleExpiredPasswords: false,
				protocol.ClientSessionTrack: false, protocol.ClientDeprecateEOF: false},
		},
		{
			[]byte{
				0x4a, 0x00, 0x00, 0x00, 0x0a, 0x35, 0x2e, 0x37, 0x2e, 0x31, 0x38, 0x00, 0x0f, 0x00, 0x00, 0x00,
				0x15, 0x12, 0x4b, 0x1f, 0x70, 0x2b, 0x33, 0x55, 0x00, 0xff, 0xff, 0x08, 0x02, 0x00, 0xff, 0xc1,
				0x15, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x01, 0x30, 0x0d, 0x0a, 0x28,
				0x06, 0x4a, 0x12, 0x5e, 0x45, 0x18, 0x05, 0x00, 0x6d, 0x79, 0x73, 0x71, 0x6c, 0x5f, 0x6e, 0x61,
				0x74, 0x69, 0x76, 0x65, 0x5f, 0x70, 0x61, 0x73, 0x73, 0x77, 0x6f, 0x72, 0x64, 0x00,
			},
			false,
			nil,
			protocol.HandshakeV10{
				ProtocolVersion:    byte(10),
				ServerVersion:      "5.7.18",
				ConnectionID:       uint32(15),
				AuthPlugin:         "mysql_native_password",
				ServerCapabilities: binary.LittleEndian.Uint32([]byte{255, 255, 255, 193}),
				Salt: []byte{0x15, 0x12, 0x4b, 0x1f, 0x70, 0x2b, 0x33, 0x55, 0x01, 0x30, 0x0d,
					0x0a, 0x28, 0x06, 0x4a, 0x12, 0x5e, 0x45, 0x18, 0x05},
			},
			map[uint32]bool{
				protocol.ClientLongPassword: true, protocol.ClientFoundRows: true, protocol.ClientLongFlag: true,
				protocol.ClientConnectWithDB: true, protocol.ClientNoSchema: true, protocol.ClientCompress: true, protocol.ClientODBC: true,
				protocol.ClientLocalFiles: true, protocol.ClientIgnoreSpace: true, protocol.ClientProtocol41: true, protocol.ClientInteractive: true,
				protocol.ClientSSL: true, protocol.ClientIgnoreSIGPIPE: true, protocol.ClientTransactions: true, protocol.ClientMultiStatements: true,
				protocol.ClientMultiResults: true, protocol.ClientPSMultiResults: true, protocol.ClientPluginAuth: true, protocol.ClientConnectAttrs: true,
				protocol.ClientPluginAuthLenEncClientData: true, protocol.ClientCanHandleExpiredPasswords: true,
				protocol.ClientSessionTrack: true, protocol.ClientDeprecateEOF: true},
		},
	}

	for _, asserted := range testData {
		decoded, err := protocol.UnpackHandshakeV10(asserted.Packet)

		if err != nil {
			assert.Equal(t, asserted.Error, err)
		} else {
			assert.Equal(t, asserted.HandshakeV10.ProtocolVersion, decoded.ProtocolVersion)
			assert.Equal(t, asserted.HandshakeV10.ServerVersion, decoded.ServerVersion)
			assert.Equal(t, asserted.HandshakeV10.ConnectionID, decoded.ConnectionID)
			assert.Equal(t, asserted.HandshakeV10.AuthPlugin, decoded.AuthPlugin)
			assert.Equal(t, asserted.HandshakeV10.Salt, decoded.Salt)
			assert.Equal(t, asserted.HandshakeV10.ServerCapabilities, decoded.ServerCapabilities)

			for flag, isSet := range asserted.CapabilitiesMap {
				if isSet {
					assert.True(t, decoded.ServerCapabilities&flag > 0)
					if decoded.ServerCapabilities&flag == 0 {
						println(flag)
					}
				} else {
					assert.True(t, decoded.ServerCapabilities&flag == 0)
				}
			}
		}
	}
}

func TestUnpackHandshakeResponse41(t *testing.T) {
	expected := protocol.HandshakeResponse41{
		Header:          []byte{0xaa, 0x0, 0x0, 0x1},
		CapabilityFlags: uint32(33464965),
		MaxPacketSize:   uint32(1073741824),
		ClientCharset:   uint8(8),
		Username:        "roger",
		AuthLength:      int64(20),
		AuthResponse: []byte{0xc0, 0xb, 0xbc, 0xb6, 0x6, 0xf5,
			0x4f, 0x4e, 0xf4, 0x1b, 0x87, 0xc0, 0xb8, 0x89, 0xae,
			0xc4, 0x49, 0x7c, 0x46, 0xf3},
		Database: "",
		PacketTail: []byte{0x58, 0x3, 0x5f, 0x6f, 0x73, 0xa, 0x6d, 0x61,
			0x63, 0x6f, 0x73, 0x31, 0x30, 0x2e, 0x31, 0x32, 0xc, 0x5f, 0x63,
			0x6c, 0x69, 0x65, 0x6e, 0x74, 0x5f, 0x6e, 0x61, 0x6d, 0x65, 0x8,
			0x6c, 0x69, 0x62, 0x6d, 0x79, 0x73, 0x71, 0x6c, 0x4, 0x5f, 0x70,
			0x69, 0x64, 0x5, 0x36, 0x36, 0x34, 0x37, 0x39, 0xf, 0x5f, 0x63,
			0x6c, 0x69, 0x65, 0x6e, 0x74, 0x5f, 0x76, 0x65, 0x72, 0x73, 0x69,
			0x6f, 0x6e, 0x6, 0x35, 0x2e, 0x37, 0x2e, 0x32, 0x30, 0x9, 0x5f,
			0x70, 0x6c, 0x61, 0x74, 0x66, 0x6f, 0x72, 0x6d, 0x6, 0x78, 0x38,
			0x36, 0x5f, 0x36, 0x34},
	}
	input := []byte{0xaa, 0x0, 0x0, 0x1, 0x85, 0xa2, 0xfe, 0x1, 0x0,
		0x0, 0x0, 0x40, 0x8, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0,
		0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0,
		0x0, 0x0, 0x0, 0x72, 0x6f, 0x67, 0x65, 0x72, 0x0, 0x14, 0xc0,
		0xb, 0xbc, 0xb6, 0x6, 0xf5, 0x4f, 0x4e, 0xf4, 0x1b, 0x87, 0xc0,
		0xb8, 0x89, 0xae, 0xc4, 0x49, 0x7c, 0x46, 0xf3, 0x6d, 0x79, 0x73,
		0x71, 0x6c, 0x5f, 0x6e, 0x61, 0x74, 0x69, 0x76, 0x65, 0x5f, 0x70,
		0x61, 0x73, 0x73, 0x77, 0x6f, 0x72, 0x64, 0x0, 0x58, 0x3, 0x5f,
		0x6f, 0x73, 0xa, 0x6d, 0x61, 0x63, 0x6f, 0x73, 0x31, 0x30, 0x2e,
		0x31, 0x32, 0xc, 0x5f, 0x63, 0x6c, 0x69, 0x65, 0x6e, 0x74, 0x5f,
		0x6e, 0x61, 0x6d, 0x65, 0x8, 0x6c, 0x69, 0x62, 0x6d, 0x79, 0x73,
		0x71, 0x6c, 0x4, 0x5f, 0x70, 0x69, 0x64, 0x5, 0x36, 0x36, 0x34,
		0x37, 0x39, 0xf, 0x5f, 0x63, 0x6c, 0x69, 0x65, 0x6e, 0x74, 0x5f,
		0x76, 0x65, 0x72, 0x73, 0x69, 0x6f, 0x6e, 0x6, 0x35, 0x2e, 0x37,
		0x2e, 0x32, 0x30, 0x9, 0x5f, 0x70, 0x6c, 0x61, 0x74, 0x66, 0x6f,
		0x72, 0x6d, 0x6, 0x78, 0x38, 0x36, 0x5f, 0x36, 0x34}

	output, err := protocol.UnpackHandshakeResponse41(input)

	assert.Equal(t, expected, *output)
	assert.Equal(t, nil, err)
}

func TestInjectCredentials(t *testing.T) {
	username := "testuser"
	password := "testpass"
	salt := []byte{0x2f, 0x50, 0x25, 0x34, 0x78, 0x17, 0x1, 0x44, 0x1d,
		0xc, 0x61, 0x4f, 0x5c, 0x69, 0x65, 0x6f, 0x25, 0x66, 0x7c, 0x64}
	expectedAuth := []byte{0xf, 0xf8, 0xe1, 0xa3, 0xe7, 0xe3, 0x5f, 0xd2,
		0xb1, 0x69, 0x8c, 0x39, 0x5b, 0xfa, 0x99, 0x4f, 0x53, 0xdd, 0xe5,
		0x35}
	expectedHeader := []byte{0xa4, 0x0, 0x0, 0x1}

	// test with handshake response that already has auth set to another value
	handshake := protocol.HandshakeResponse41{
		AuthLength: int64(20),
		AuthResponse: []byte{0xc0, 0xb, 0xbc, 0xb6, 0x6, 0xf5, 0x4f, 0x4e,
			0xf4, 0x1b, 0x87, 0xc0, 0xb8, 0x89, 0xae, 0xc4, 0x49, 0x7c, 0x46, 0xf3},
		Username: "madeupusername",
		Header:   []byte{0xaa, 0x0, 0x0, 0x1},
	}

	err := protocol.InjectCredentials(&handshake, salt, username, password)

	assert.Equal(t, username, handshake.Username)
	assert.Equal(t, int64(20), handshake.AuthLength)
	assert.Equal(t, expectedAuth, handshake.AuthResponse)
	assert.Equal(t, expectedHeader, handshake.Header)
	assert.Equal(t, nil, err)

	// test with handshake response with empty auth
	expectedHeader = []byte{0xb8, 0x0, 0x0, 0x1}
	handshake = protocol.HandshakeResponse41{
		AuthLength:   0,
		AuthResponse: []byte{},
		Username:     "madeupusername",
		Header:       []byte{0xaa, 0x0, 0x0, 0x1},
	}

	err = protocol.InjectCredentials(&handshake, salt, username, password)

	assert.Equal(t, username, handshake.Username)
	assert.Equal(t, int64(20), handshake.AuthLength)
	assert.Equal(t, expectedAuth, handshake.AuthResponse)
	assert.Equal(t, expectedHeader, handshake.Header)
	assert.Equal(t, nil, err)
}

func TestPackHandshakeResponse41(t *testing.T) {
	input := protocol.HandshakeResponse41{
		Header:          []byte{0xaa, 0x0, 0x0, 0x1},
		CapabilityFlags: uint32(33464965),
		MaxPacketSize:   uint32(1073741824),
		ClientCharset:   uint8(8),
		Username:        "roger",
		AuthLength:      int64(20),
		AuthResponse: []byte{0xc0, 0xb, 0xbc, 0xb6, 0x6, 0xf5,
			0x4f, 0x4e, 0xf4, 0x1b, 0x87, 0xc0, 0xb8, 0x89, 0xae,
			0xc4, 0x49, 0x7c, 0x46, 0xf3},
		Database: "",
		PacketTail: []byte{0x58, 0x3, 0x5f, 0x6f, 0x73, 0xa, 0x6d, 0x61,
			0x63, 0x6f, 0x73, 0x31, 0x30, 0x2e, 0x31, 0x32, 0xc, 0x5f, 0x63,
			0x6c, 0x69, 0x65, 0x6e, 0x74, 0x5f, 0x6e, 0x61, 0x6d, 0x65, 0x8,
			0x6c, 0x69, 0x62, 0x6d, 0x79, 0x73, 0x71, 0x6c, 0x4, 0x5f, 0x70,
			0x69, 0x64, 0x5, 0x36, 0x36, 0x34, 0x37, 0x39, 0xf, 0x5f, 0x63,
			0x6c, 0x69, 0x65, 0x6e, 0x74, 0x5f, 0x76, 0x65, 0x72, 0x73, 0x69,
			0x6f, 0x6e, 0x6, 0x35, 0x2e, 0x37, 0x2e, 0x32, 0x30, 0x9, 0x5f,
			0x70, 0x6c, 0x61, 0x74, 0x66, 0x6f, 0x72, 0x6d, 0x6, 0x78, 0x38,
			0x36, 0x5f, 0x36, 0x34},
	}
	expected := []byte{0xaa, 0x0, 0x0, 0x1, 0x85, 0xa2, 0xfe, 0x1, 0x0,
		0x0, 0x0, 0x40, 0x8, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0,
		0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0,
		0x0, 0x0, 0x0, 0x72, 0x6f, 0x67, 0x65, 0x72, 0x0, 0x14, 0xc0,
		0xb, 0xbc, 0xb6, 0x6, 0xf5, 0x4f, 0x4e, 0xf4, 0x1b, 0x87, 0xc0,
		0xb8, 0x89, 0xae, 0xc4, 0x49, 0x7c, 0x46, 0xf3, 0x6d, 0x79, 0x73,
		0x71, 0x6c, 0x5f, 0x6e, 0x61, 0x74, 0x69, 0x76, 0x65, 0x5f, 0x70,
		0x61, 0x73, 0x73, 0x77, 0x6f, 0x72, 0x64, 0x0, 0x58, 0x3, 0x5f,
		0x6f, 0x73, 0xa, 0x6d, 0x61, 0x63, 0x6f, 0x73, 0x31, 0x30, 0x2e,
		0x31, 0x32, 0xc, 0x5f, 0x63, 0x6c, 0x69, 0x65, 0x6e, 0x74, 0x5f,
		0x6e, 0x61, 0x6d, 0x65, 0x8, 0x6c, 0x69, 0x62, 0x6d, 0x79, 0x73,
		0x71, 0x6c, 0x4, 0x5f, 0x70, 0x69, 0x64, 0x5, 0x36, 0x36, 0x34,
		0x37, 0x39, 0xf, 0x5f, 0x63, 0x6c, 0x69, 0x65, 0x6e, 0x74, 0x5f,
		0x76, 0x65, 0x72, 0x73, 0x69, 0x6f, 0x6e, 0x6, 0x35, 0x2e, 0x37,
		0x2e, 0x32, 0x30, 0x9, 0x5f, 0x70, 0x6c, 0x61, 0x74, 0x66, 0x6f,
		0x72, 0x6d, 0x6, 0x78, 0x38, 0x36, 0x5f, 0x36, 0x34}

	output, err := protocol.PackHandshakeResponse41(&input)

	assert.Equal(t, expected, output)
	assert.Equal(t, nil, err)
}

func TestGetLenEncodedIntegerSize(t *testing.T) {
	inputArray := []byte{0xfc, 0xfd, 0xfe, 0xfb}
	expectedArray := []byte{2, 3, 8, 1}

	for k, v := range inputArray {
		output := protocol.GetLenEncodedIntegerSize(v)

		assert.Equal(t, expectedArray[k], output)
	}
}

func TestReadLenEncodedInteger(t *testing.T) {
	expected := uint64(251)
	input := bytes.NewReader([]byte{0xfc, 0xfb, 0x00})

	output, err := protocol.ReadLenEncodedInteger(input)

	assert.Equal(t, expected, output)
	assert.Equal(t, nil, err)
}

func TestReadLenEncodedString(t *testing.T) {
	expected := "ABCDEFGHIKLMONPQRSTYW"
	packet := bytes.NewReader([]byte{
		0x15, 0x41, 0x42, 0x43, 0x44, 0x45, 0x46, 0x47, 0x48, 0x49, 0x4b, 0x4c, 0x4d, 0x4f, 0x4e, 0x50,
		0x51, 0x52, 0x53, 0x54, 0x59, 0x57})

	decoded, length, err := protocol.ReadLenEncodedString(packet)

	assert.Equal(t, expected, decoded)
	assert.Equal(t, len(expected), int(length))
	assert.Equal(t, nil, err)
}

func TestReadEOFLengthString(t *testing.T) {
	expected := "SET sql_mode='STRICT_TRANS_TABLES'"
	encoded := []byte{
		0x53, 0x45, 0x54, 0x20, 0x73, 0x71, 0x6c, 0x5f, 0x6d, 0x6f, 0x64, 0x65, 0x3d, 0x27, 0x53, 0x54, 0x52,
		0x49, 0x43, 0x54, 0x5f, 0x54, 0x52, 0x41, 0x4e, 0x53, 0x5f, 0x54, 0x41, 0x42, 0x4c, 0x45, 0x53, 0x27,
	}

	decoded := protocol.ReadEOFLengthString(encoded)

	assert.Equal(t, expected, decoded)
}

func TestReadNullTerminatedString(t *testing.T) {
	x := bytes.NewReader([]byte{0x35, 0x2e, 0x37, 0x2e, 0x31, 0x38, 0x00})
	assert.Equal(t, "5.7.18", protocol.ReadNullTerminatedString(x))
}

func TestReadNullTerminatedBytes(t *testing.T) {
	input := bytes.NewReader([]byte{0x1d, 0xc, 0x61, 0x4f, 0x5c, 0x69,
		0x65, 0x6f, 0x25, 0x66, 0x7c, 0x64, 0x0, 0x6d, 0x79, 0x73, 0x71,
		0x6c, 0x5f, 0x6e, 0x61, 0x74, 0x69, 0x76, 0x65, 0x5f, 0x70, 0x61,
		0x73, 0x73, 0x77, 0x6f, 0x72, 0x64, 0x0})
	expected := []byte{0x1d, 0xc, 0x61, 0x4f, 0x5c, 0x69, 0x65, 0x6f,
		0x25, 0x66, 0x7c, 0x64}

	output := protocol.ReadNullTerminatedBytes(input)

	assert.Equal(t, expected, output)
}

func TestGetPacketHeader(t *testing.T) {
	input := bytes.NewReader([]byte{0x4a, 0x0, 0x0, 0x0, 0xa, 0x35, 0x2e,
		0x37, 0x2e, 0x32, 0x31, 0x0, 0x38, 0x9, 0x0, 0x0, 0x2f, 0x50, 0x25,
		0x34, 0x78, 0x17, 0x1, 0x44, 0x0, 0xff, 0xff, 0x8, 0x2, 0x0,
		0xff, 0xc1, 0x15, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0,
		0x0, 0x1d, 0xc, 0x61, 0x4f, 0x5c, 0x69, 0x65, 0x6f, 0x25, 0x66,
		0x7c, 0x64, 0x0, 0x6d, 0x79, 0x73, 0x71, 0x6c, 0x5f, 0x6e, 0x61,
		0x74, 0x69, 0x76, 0x65, 0x5f, 0x70, 0x61, 0x73, 0x73, 0x77, 0x6f,
		0x72, 0x64, 0x0})
	expected := []byte{0x4a, 0x0, 0x0, 0x0}

	output, err := protocol.GetPacketHeader(input)

	assert.Equal(t, expected, output)
	assert.Equal(t, nil, err)
}

func TestCheckPacketLength(t *testing.T) {
	inputLength := 4
	inputPacket := []byte{0xf, 0xf8, 0xe1, 0xa3}

	err := protocol.CheckPacketLength(inputLength, inputPacket)

	assert.Equal(t, nil, err)
}

func TestNativePassword(t *testing.T) {
	expected := []byte{0xf, 0xf8, 0xe1, 0xa3, 0xe7, 0xe3, 0x5f, 0xd2, 0xb1, 0x69,
		0x8c, 0x39, 0x5b, 0xfa, 0x99, 0x4f, 0x53, 0xdd, 0xe5, 0x35}
	inputPassword := "testpass"
	inputSalt := []byte{0x2f, 0x50, 0x25, 0x34, 0x78, 0x17, 0x1, 0x44,
		0x1d, 0xc, 0x61, 0x4f, 0x5c, 0x69, 0x65, 0x6f, 0x25, 0x66, 0x7c, 0x64}

	output, err := protocol.NativePassword(inputPassword, inputSalt)

	assert.Equal(t, expected, output)
	assert.Equal(t, nil, err)
}

func TestUpdateHeaderPayloadLength(t *testing.T) {
	// Test with a valid negative value
	expectedHeader := []byte{170, 0, 0, 0}
	inputHeader := []byte{173, 0, 0, 0}
	inputLength := int32(-3)

	output, err := protocol.UpdateHeaderPayloadLength(inputHeader, inputLength)

	assert.Equal(t, expectedHeader, output)
	assert.Equal(t, nil, err)

	// Test with a valid positive value
	expectedHeader = []byte{176, 0, 0, 0}
	inputHeader = []byte{173, 0, 0, 0}
	inputLength = int32(3)

	output, err = protocol.UpdateHeaderPayloadLength(inputHeader, inputLength)

	assert.Equal(t, expectedHeader, output)
	assert.Equal(t, nil, err)

	// Test with an invalid value for the length difference
	inputHeader = []byte{173, 0, 0, 0}
	inputLength = int32(-180)

	output, err = protocol.UpdateHeaderPayloadLength(inputHeader, inputLength)

	assert.EqualError(t, err, "Malformed packet")
}

func TestReadUint24(t *testing.T) {
	expected := uint32(173)
	input := []byte{173, 0, 0}

	output, err := protocol.ReadUint24(input)

	assert.Equal(t, expected, output)
	assert.Equal(t, nil, err)
}

func TestWriteUint24(t *testing.T) {
	expected := []byte{173, 0, 0}
	input := uint32(173)

	output := protocol.WriteUint24(input)

	assert.Equal(t, expected, output)
}
