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

	"github.com/stretchr/testify/assert"

	"github.com/ConsenSys/fc-retrieval-gateway/pkg/fcrcrypto"
	log "github.com/ConsenSys/fc-retrieval-gateway/pkg/logging"
	"github.com/ConsenSys/fc-retrieval-gateway-admin/pkg/fcrgatewayadmin"
)

// Test the Client API.

// func TestGetClientVersion(t *testing.T) {
// 	versionInfo := fcrclient.GetVersion()
// 	// Verify that the client version is an integer number.
// 	ver, err := strconv.Atoi(versionInfo.Version)
// 	if err != nil {
// 		panic(err)
// 	}

// 	// The version must be 1 or more.
// 	assert.LessOrEqual(t, 1, ver)
// }

func TestInitGateway(t *testing.T) {
	blockchainPrivateKey, err := fcrcrypto.GenerateBlockchainKeyPair()
	if err != nil {
		log.Panic(err.Error())
	}
	gatewayRetrievalPrivateKey, err := fcrcrypto.GenerateRetrievalV1KeyPair()
	if err != nil {
		log.Panic(err.Error())
	}

	confBuilder := fcrgatewayadmin.CreateSettings()
	confBuilder.SetEstablishmentTTL(101)
	confBuilder.SetBlockchainPrivateKey(blockchainPrivateKey)
	conf := confBuilder.Build()

	gatewayAdmin := fcrgatewayadmin.InitFilecoinRetrievalGatewayAdminClient(*conf)
	err = gatewayAdmin.SendKeyToGateway("gateway", gatewayRetrievalPrivateKey);
	gatewayAdmin.Shutdown()}

func TestOneConnectedGateway(t *testing.T) {
	// The current configuration means that there should only be one connected gateway
	client := InitClient()
	gateways := client.ConnectedGateways()
	assert.Equal(t, 1, len(gateways), "Unexpected number of gateways returned")
	if len(gateways) > 0 {
		assert.Equal(t, "http://gateway:9011/", gateways[0])
	}
	CloseClient(client)

}
