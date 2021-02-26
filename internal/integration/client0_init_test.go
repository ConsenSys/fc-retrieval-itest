package integration

/*
 * Copyright 2021 ConsenSys Software Inc.
 *
 * Licensed under the Apache License, Version 2.0 (the "License"); you may not use this file except in compliance with
 * the License. You may obtain a copy of the License at
 *
 * http://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing, software distributed under the License is distributed on
 * an "AS IS" BASIS, WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied. See the License for the
 * specific language governing permissions and limitations under the License.
 *
 * SPDX-License-Identifier: Apache-2.0
 */

import (
	"strconv"
	"testing"
	"time"

	"github.com/ConsenSys/fc-retrieval-client/pkg/fcrclient"
	"github.com/ConsenSys/fc-retrieval-common/pkg/fcrcrypto"
	"github.com/ConsenSys/fc-retrieval-common/pkg/nodeid"
	"github.com/stretchr/testify/assert"
)

// Tests in this file use the Client API, but don't need the rest of the system to be
// configured. These tests need to be run prior to the other client tests.

func TestGetClientVersion(t *testing.T) {
	versionInfo := fcrclient.GetVersion()
	// Verify that the client version is an integer number.
	ver, err := strconv.Atoi(versionInfo.Version)
	if err != nil {
		panic(err)
	}

	// The version must be 1 or more.
	assert.LessOrEqual(t, 1, ver)
}

// Test that the client can be initialized without causing an error.
func TestInitClientNoRetrievalKey(t *testing.T) {
	blockchainPrivateKey, err := fcrcrypto.GenerateBlockchainKeyPair()
	if err != nil {
		panic(err)
	}

	confBuilder := fcrclient.CreateSettings()
	confBuilder.SetEstablishmentTTL(101)
	confBuilder.SetBlockchainPrivateKey(blockchainPrivateKey)
	conf := confBuilder.Build()

	client, err := fcrclient.NewFilecoinRetrievalClient(*conf)
	if err != nil {
		panic(err)
	}

	client.Shutdown()
}

func TestNoConfiguredGateways(t *testing.T) {
	// The current configuration means that there should only be one connected gateway
	client := InitClient()
	gateways := client.ConnectedGateways()
	assert.Equal(t, 0, len(gateways), "Unexpected number of gateways returned")
	if len(gateways) > 0 {
		assert.Equal(t, "http://gateway:9011/", gateways[0])
	}
	CloseClient(client)

}

func TestUnknownGatewayAdded(t *testing.T) {
	client := InitClient()
	randomGatewayID, err := nodeid.NewRandomNodeID()
	if err != nil {
		panic(err)
	}
	newGatwaysToBeAdded := make([]*nodeid.NodeID, 0)
	newGatwaysToBeAdded = append(newGatwaysToBeAdded, randomGatewayID)
	numAdded := client.AddGateways(newGatwaysToBeAdded)
	assert.Equal(t, 1, numAdded)
	gws := client.GetGateways()
	assert.Equal(t, 1, len(gws))

	// Give the client time to fail the look-up for the random gateway id
	time.Sleep(500 * time.Millisecond)

	gateways := client.ConnectedGateways()
	assert.Equal(t, 0, len(gateways), "Unexpected number of gateways returned")
	if len(gateways) > 0 {
		assert.Equal(t, "http://gateway:9011/", gateways[0])
	}
	CloseClient(client)
}
