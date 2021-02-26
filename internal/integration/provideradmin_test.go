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
	"github.com/ConsenSys/fc-retrieval-common/pkg/fcrmessages"
	"github.com/ConsenSys/fc-retrieval-common/pkg/logging"
	"github.com/ConsenSys/fc-retrieval-common/pkg/nodeid"
	"github.com/ConsenSys/fc-retrieval-common/pkg/register"
	"github.com/ConsenSys/fc-retrieval-itest/config"
	"github.com/ConsenSys/fc-retrieval-itest/pkg/provider"
	"github.com/ConsenSys/fc-retrieval-provider-admin/pkg/fcrprovideradmin"
	"github.com/stretchr/testify/assert"
)

// Test the Provider Admin API.
var providerConfig = config.NewConfig(".env.provider")

func TestGetProviderAdminVersion(t *testing.T) {
	versionInfo := fcrclient.GetVersion()
	// Verify that the client version is an integer number.
	_, err := strconv.Atoi(versionInfo.Version)
	if err != nil {
		panic(err)
	}

	// The version must be 1 or more.
	assert.LessOrEqual(t, 0, 0)
}

func TestInitProviderAdminNoRetrievalKey(t *testing.T) {
	logging.Error(" Wait two seconds for the provider to deploy and be ready for requests")
	time.Sleep(2 * time.Second)

	blockchainPrivateKey, err := fcrcrypto.GenerateBlockchainKeyPair()
	if err != nil {
		logging.ErrorAndPanic(err.Error())
	}

	confBuilder := fcrprovideradmin.CreateSettings()
	confBuilder.SetEstablishmentTTL(101)
	confBuilder.SetBlockchainPrivateKey(blockchainPrivateKey)
	confBuilder.SetRegisterURL(providerConfig.GetString("REGISTER_API_URL"))
	confBuilder.SetProviderRegister(&register.ProviderRegister{
		NodeID:           providerConfig.GetString("PROVIDER_ID"),
		Address:          providerConfig.GetString("PROVIDER_ADDRESS"),
		NetworkInfoAdmin: providerConfig.GetString("ADMIN_NETWORK_INFO"),
		RegionCode:       providerConfig.GetString("PROVIDER_REGION_CODE"),
		RootSigningKey:   providerConfig.GetString("PROVIDER_ROOT_SIGNING_KEY"),
		SigningKey:       providerConfig.GetString("PROVIDER_SIGNING_KEY"),
	})
	conf := confBuilder.Build()

	// Init client
	client := fcrprovideradmin.InitFilecoinRetrievalProviderAdminClient(*conf)

	// Register provider
	client.RegisterProvider()

	// Publish Group CID
	message := provider.GenerateDummyProviderPublishGroupCIDMessage()
	client.SendMessage(message)
	logging.Info("Wait 3 seconds")
	time.Sleep(3 * time.Second)

	gateways, _ := client.GetRegisteredGateways()
	logging.Info("Got %v registered gateway(s)", len(gateways))

	// Get all offers
	var gatewayIDs []nodeid.NodeID
	gatewayIDs = make([]nodeid.NodeID, 0)
	message, _ = fcrmessages.EncodeProviderAdminGetGroupCIDRequest(gatewayIDs)
	logging.Info("Get all offers")
	response, err := client.SendMessage(message)
	if err != nil {
		logging.Error("Response error: %+v", err)
	}
	if response != nil {
		_, cidgroupsInfo, _ := fcrmessages.DecodeProviderAdminGetGroupCIDResponse(response)
		logging.Info("Get all offers. Found: %v", len(cidgroupsInfo))
		assert.GreaterOrEqual(t, len(cidgroupsInfo), 1, "Offers should be found")
	}

	// Get offers by gatewayIDs
	realNodeID := "101112131415161718191A1B1C1D1E1F202122232425262728292A2B2C2D2E2F"
	gatewayID, _ := nodeid.NewNodeIDFromString(realNodeID)
	gatewayIDs = make([]nodeid.NodeID, 1)
	gatewayIDs[0] = *gatewayID
	message, _ = fcrmessages.EncodeProviderAdminGetGroupCIDRequest(gatewayIDs)
	logging.Info("Get offers by gatewayID=%s", gatewayID.ToString())
	response, err = client.SendMessage(message)
	if err != nil {
		logging.Error("Response error: %+v", err)
	}
	if response != nil {
		_, cidgroupsInfo, _ := fcrmessages.DecodeProviderAdminGetGroupCIDResponse(response)
		logging.Info("Get offers by gatewayID=%v. Found: %v", gatewayID.ToString(), len(cidgroupsInfo))
		assert.GreaterOrEqual(t, len(cidgroupsInfo), 1, "Offers should be found")
	}

	// Get offers by gatewayIDs
	fakeNodeID := "101112131415161718191A1B1C1D1E1F202122232425262728292A2B2C2DFA43"
	gatewayID, _ = nodeid.NewNodeIDFromString(fakeNodeID)
	gatewayIDs = make([]nodeid.NodeID, 1)
	gatewayIDs[0] = *gatewayID
	message, _ = fcrmessages.EncodeProviderAdminGetGroupCIDRequest(gatewayIDs)
	logging.Info("Get offers by gatewayID=%s", gatewayID.ToString())
	response, err = client.SendMessage(message)
	if err != nil {
		logging.Error("Response error: %+v", err)
	}
	if response != nil {
		_, cidgroupsInfo, _ := fcrmessages.DecodeProviderAdminGetGroupCIDResponse(response)
		logging.Info("Get offers by gatewayID=%v. Found: %v", gatewayID.ToString(), len(cidgroupsInfo))
		assert.Equal(t, 0, len(cidgroupsInfo), "Offers should be empty")
	}

	// Shutdown
	client.Shutdown()

	// The version must be 1 or more.
	assert.LessOrEqual(t, 1, 1)
}
