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
	"testing"
	"time"

	"github.com/ConsenSys/fc-retrieval-common/pkg/nodeid"
	"github.com/stretchr/testify/assert"

	"github.com/ConsenSys/fc-retrieval-gateway-admin/pkg/fcrgatewayadmin"
)

// Tests in this file check the ability to do node discovery.

func TestOneGateway(t *testing.T) {
	gwAdmin := InitGatewayAdmin()
	gatewayRetrievalPrivateKey, err := fcrgatewayadmin.CreateKey()
	if err != nil {
		panic(err)
	}

	// TODO fix this hard coded domain name
	err = gwAdmin.InitializeGateway("gateway", gatewayRetrievalPrivateKey)
	if err != nil {
		panic(err)
	}

	gwID, err := nodeid.NewNodeIDFromPublicKey(gatewayRetrievalPrivateKey)
	if err != nil {
		panic(err)
	}

	client := InitClient()
	newGatwaysToBeAdded := make([]*nodeid.NodeID, 0)
	newGatwaysToBeAdded = append(newGatwaysToBeAdded, gwID)
	numAdded := client.AddGateways(newGatwaysToBeAdded)
	assert.Equal(t, 1, numAdded)
	gws := client.GetGateways()
	assert.Equal(t, 1, len(gws))

	// Give the client time to talk to the register and get the gateway.
	time.Sleep(500 * time.Millisecond)

	client.ConnectedGateways()
	// TODO UNCOMMENT WHEN GATEWAY REGISTRATION IS WORKING
	//	gateways := client.ConnectedGateways()
	//	assert.Equal(t, 1, len(gateways), "Unexpected number of gateways returned")
	CloseClient(client)
	CloseGatewayAdmin(gwAdmin)
}
