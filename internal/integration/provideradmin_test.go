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
	"github.com/ConsenSys/fc-retrieval-common/pkg/cid"
	"github.com/ConsenSys/fc-retrieval-common/pkg/fcrcrypto"
	"github.com/ConsenSys/fc-retrieval-common/pkg/logging"
	"github.com/ConsenSys/fc-retrieval-common/pkg/nodeid"
	"github.com/ConsenSys/fc-retrieval-common/pkg/register"
	"github.com/ConsenSys/fc-retrieval-itest/config"
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
	logging.Info("/*******************************************************/")
	logging.Info("/*      Start TestInitProviderAdminNoRetrievalKey	     */")
	logging.Info("/*******************************************************/")
	logging.Error("Wait two seconds for the provider to deploy and be ready for requests")
	time.Sleep(2 * time.Second)

	blockchainPrivateKey, err := fcrcrypto.GenerateBlockchainKeyPair()
	if err != nil {
		logging.ErrorAndPanic(err.Error())
	}

	confBuilder := fcrprovideradmin.CreateSettings()
	confBuilder.SetBlockchainPrivateKey(blockchainPrivateKey)
	confBuilder.SetRegisterURL(providerConfig.GetString("REGISTER_API_URL"))
	conf := confBuilder.Build()

	// Init client
	client := fcrprovideradmin.InitFilecoinRetrievalProviderAdminClient(*conf)

	// Generate private key for provider
	providerPrivKey, err := fcrcrypto.GenerateRetrievalV1KeyPair()
	if err != nil {
		logging.ErrorAndPanic(err.Error())
	}
	providerSigningKey, err := providerPrivKey.EncodePublicKey()
	if err != nil {
		logging.ErrorAndPanic(err.Error())
	}
	providerID, err := nodeid.NewRandomNodeID()
	if err != nil {
		logging.ErrorAndPanic(err.Error())
	}

	providerRegister := &register.ProviderRegister{
		NodeID:             providerID.ToString(),
		Address:            providerConfig.GetString("PROVIDER_ADDRESS"),
		RootSigningKey:     providerConfig.GetString("PROVIDER_ROOT_SIGNING_KEY"),
		SigningKey:         providerSigningKey,
		RegionCode:         providerConfig.GetString("PROVIDER_REGION_CODE"),
		NetworkInfoGateway: providerConfig.GetString("NETWORK_INFO_GATEWAY"),
		NetworkInfoClient:  providerConfig.GetString("NETWORK_INFO_CLIENT"),
		NetworkInfoAdmin:   providerConfig.GetString("NETWORK_INFO_ADMIN"),
	}

	// Initialise provider
	err = client.InitialiseProvider(providerRegister, providerPrivKey, fcrcrypto.DecodeKeyVersion(1))
	if err != nil {
		logging.ErrorAndPanic(err.Error())
	}

	logging.Info("Wait five seconds for the provider to initialise")
	time.Sleep(5 * time.Second)

	// Generate random cid offer
	contentID, _ := cid.NewRandomContentID()
	pieceCIDs := []cid.ContentID{*contentID}
	expiryDate := time.Now().Local().Add(time.Hour * time.Duration(24)).Unix()

	// Publish Group CID
	err = client.PublishGroupCID(providerID, pieceCIDs, 42, expiryDate, 42)
	if err != nil {
		logging.ErrorAndPanic(err.Error())
	}
	logging.Info("Wait 3 seconds")
	time.Sleep(3 * time.Second)

	// Get all offers
	var gatewayIDs []nodeid.NodeID
	gatewayIDs = make([]nodeid.NodeID, 0)
	logging.Info("Get all offers")
	_, cidgroupInfo, err := client.GetGroupCIDOffer(providerID, gatewayIDs)
	if err != nil {
		logging.ErrorAndPanic(err.Error())
	}
	assert.GreaterOrEqual(t, len(cidgroupInfo), 1, "Offers should be found")

	// Get offers by gatewayIDs real
	gateways, err := register.GetRegisteredGateways(providerConfig.GetString("REGISTER_API_URL"))
	if err != nil {
		logging.ErrorAndPanic(err.Error())
	}
	realNodeID, err := nodeid.NewNodeIDFromString(gateways[0].NodeID)
	if err != nil {
		logging.ErrorAndPanic(err.Error())
	}
	gatewayIDs = append(gatewayIDs, *realNodeID) // Add a gateway
	logging.Info("Get offers by gatewayID=%s", realNodeID.ToString())
	_, cidgroupInfo, err = client.GetGroupCIDOffer(providerID, gatewayIDs)
	if err != nil {
		logging.ErrorAndPanic(err.Error())
	}
	assert.GreaterOrEqual(t, len(cidgroupInfo), 1, "Offers should be found")

	// Get offers by gatewayIDs fake
	fakeNodeID, _ := nodeid.NewNodeIDFromString("101112131415161718191A1B1C1D1E1F202122232425262728292A2B2C2DFA43")
	gatewayIDs[0] = *fakeNodeID
	logging.Info("Get offers by gatewayID=%s", fakeNodeID.ToString())
	_, cidgroupInfo, err = client.GetGroupCIDOffer(providerID, gatewayIDs)
	if err != nil {
		logging.ErrorAndPanic(err.Error())
	}
	assert.Equal(t, 0, len(cidgroupInfo), "Offers should be empty")

	// The version must be 1 or more.
	assert.LessOrEqual(t, 1, 1)

	logging.Info("/*******************************************************/")
	logging.Info("/*      End TestInitProviderAdminNoRetrievalKey	       */")
	logging.Info("/*******************************************************/")
}
