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

	"github.com/ConsenSys/fc-retrieval-gateway/pkg/logging"
	"github.com/ConsenSys/fc-retrieval-gateway/pkg/nodeid"
	"github.com/ConsenSys/fc-retrieval-gateway/pkg/register"

	"github.com/ConsenSys/fc-retrieval-itest/internal/provider"
	prov "github.com/ConsenSys/fc-retrieval-provider/pkg/provider"
	"github.com/stretchr/testify/assert"
)

// TestProviderPublishMessage test a provider publishing messages
func TestProviderPublishMessage(t *testing.T) {
	p := InitProvider()


	gateways, err := register.GetRegisteredGateways(p.Conf.GetString("REGISTER_API_URL"))
	if err != nil {
		panic(err)
	}
	for _, gw := range gateways {
		message := provider.GenerateDummyMessage()
		logging.Info("Message: %v", message)
		gatewayID, err := nodeid.NewNodeIDFromString(gw.NodeID)
		if err != nil {
			logging.Error("Error with nodeID %v: %v", gw.NodeID, err)
			continue
		}
		p.GatewayCommPool.RegisteredNodeMap[gw.NodeID] = &gw
		prov.SendMessageToGateway(message, gatewayID, p.GatewayCommPool)
		logging.Info("Message sent to gateway %s", gw.NetworkProviderInfo)
	}
	CloseProvider(p)

	assert.LessOrEqual(t, 10, 1)
}
