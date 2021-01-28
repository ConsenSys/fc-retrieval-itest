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

	"github.com/stretchr/testify/assert"
	

	"github.com/ConsenSys/fc-retrieval-gateway/pkg/logging"

	"github.com/ConsenSys/fc-retrieval-client/pkg/fcrclient"
	"github.com/ConsenSys/fc-retrieval-gateway/pkg/fcrcrypto"
)


func TestGetEmptyOffer(t *testing.T) {
	var pieceCIDToFind [32]byte

	client := InitClient()
	offers := client.FindBestOffers(pieceCIDToFind, 1000, 1000)
    assert.Equal(t, 0, len(offers))
	CloseClient(client)
}

