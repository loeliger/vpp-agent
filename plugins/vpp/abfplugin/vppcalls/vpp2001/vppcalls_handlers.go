//  Copyright (c) 2019 Cisco and/or its affiliates.
//
//  Licensed under the Apache License, Version 2.0 (the "License");
//  you may not use this file except in compliance with the License.
//  You may obtain a copy of the License at:
//
//      http://www.apache.org/licenses/LICENSE-2.0
//
//  Unless required by applicable law or agreed to in writing, software
//  distributed under the License is distributed on an "AS IS" BASIS,
//  WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
//  See the License for the specific language governing permissions and
//  limitations under the License.

package vpp2001

import (
	govppapi "git.fd.io/govpp.git/api"
	"github.com/ligato/cn-infra/logging"

	"github.com/ligato/vpp-agent/plugins/vpp/abfplugin/vppcalls"
	"github.com/ligato/vpp-agent/plugins/vpp/aclplugin/aclidx"
	vpp_abf "github.com/ligato/vpp-agent/plugins/vpp/binapi/vpp2001/abf"
	"github.com/ligato/vpp-agent/plugins/vpp/ifplugin/ifaceidx"
)

func init() {
	var msgs []govppapi.Message
	msgs = append(msgs, vpp_abf.AllMessages()...)

	vppcalls.Versions["vpp2001"] = vppcalls.HandlerVersion{
		Msgs: msgs,
		New: func(ch govppapi.Channel, aclIndexes aclidx.ACLMetadataIndex, ifIndexes ifaceidx.IfaceMetadataIndex, log logging.Logger) vppcalls.ABFVppAPI {
			return NewABFVppHandler(ch, aclIndexes, ifIndexes, log)
		},
	}
}

// ABFVppHandler is accessor for abf-related vppcalls methods
type ABFVppHandler struct {
	callsChannel govppapi.Channel
	aclIndexes   aclidx.ACLMetadataIndex
	ifIndexes    ifaceidx.IfaceMetadataIndex
	log          logging.Logger
}

// NewABFVppHandler returns new ABFVppHandler.
func NewABFVppHandler(calls govppapi.Channel, aclIdx aclidx.ACLMetadataIndex, ifIdx ifaceidx.IfaceMetadataIndex, log logging.Logger) *ABFVppHandler {
	return &ABFVppHandler{
		callsChannel: calls,
		aclIndexes:   aclIdx,
		ifIndexes:    ifIdx,
		log:          log,
	}
}
