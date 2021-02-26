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

	//	"github.com/stretchr/testify/assert"

	"github.com/ConsenSys/fc-retrieval-common/pkg/fcrcrypto"
	log "github.com/ConsenSys/fc-retrieval-common/pkg/logging"
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

	confBuilder := fcrgatewayadmin.CreateSettings()
	confBuilder.SetEstablishmentTTL(101)
	confBuilder.SetBlockchainPrivateKey(blockchainPrivateKey)
	conf := confBuilder.Build()

	gatewayAdmin := fcrgatewayadmin.NewFilecoinRetrievalGatewayAdminClient(*conf)

	gatewayRetrievalPrivateKey, err := fcrgatewayadmin.CreateKey()
	if err != nil {
		panic(err)
	}

	// TODO add a get key to check it doesn't exist

	err = gatewayAdmin.InitializeGateway("gateway", gatewayRetrievalPrivateKey);
	if err != nil {
		panic(err)
	}

	// TODO add a get key to see that the gateway has the key

	log.Info("Working!")
	gatewayAdmin.Shutdown()
}

