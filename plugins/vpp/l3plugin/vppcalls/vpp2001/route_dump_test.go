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
	l3 "github.com/ligato/vpp-agent/api/models/vpp/l3"
	"testing"

	"github.com/ligato/vpp-agent/plugins/vpp/l3plugin/vrfidx"

	"github.com/ligato/cn-infra/logging/logrus"
	netallock_mock "github.com/ligato/vpp-agent/plugins/netalloc/mock"
	vpp_ip "github.com/ligato/vpp-agent/plugins/vpp/binapi/vpp2001/ip"
	vpp_vpe "github.com/ligato/vpp-agent/plugins/vpp/binapi/vpp2001/vpe"
	"github.com/ligato/vpp-agent/plugins/vpp/ifplugin/ifaceidx"
	"github.com/ligato/vpp-agent/plugins/vpp/vppcallmock"
	. "github.com/onsi/gomega"
)

// Test dumping routes
func TestDumpStaticRoutes(t *testing.T) {
	ctx := vppcallmock.SetupTestCtx(t)
	defer ctx.TeardownTestCtx()
	ifIndexes := ifaceidx.NewIfaceIndex(logrus.NewLogger("test-if"), "test-if")
	vrfIndexes := vrfidx.NewVRFIndex(logrus.NewLogger("test-vrf"), "test-vrf")
	l3handler := NewRouteVppHandler(ctx.MockChannel, ifIndexes, vrfIndexes, netallock_mock.NewMockNetAlloc(),
		logrus.DefaultLogger())

	vrfIndexes.Put("vrf1-ipv4", &vrfidx.VRFMetadata{Index: 0, Protocol: l3.VrfTable_IPV4})
	vrfIndexes.Put("vrf1-ipv6", &vrfidx.VRFMetadata{Index: 0, Protocol: l3.VrfTable_IPV6})
	ifIndexes.Put("if1", &ifaceidx.IfaceMetadata{SwIfIndex: 1})
	ifIndexes.Put("if2", &ifaceidx.IfaceMetadata{SwIfIndex: 2})

	ctx.MockVpp.MockReply(&vpp_ip.IPRouteDetails{
		Route: vpp_ip.IPRoute{
			Prefix: vpp_ip.Prefix{
				Address: vpp_ip.Address{
					Af: vpp_ip.ADDRESS_IP4,
					Un: vpp_ip.AddressUnionIP4([4]uint8{10, 0, 0, 1}),
				},
			},
			Paths: []vpp_ip.FibPath{
				{
					SwIfIndex: 2,
				},
			},
		},
	})
	ctx.MockVpp.MockReply(&vpp_vpe.ControlPingReply{})
	ctx.MockVpp.MockReply(&vpp_ip.IPRouteDetails{
		Route: vpp_ip.IPRoute{
			Prefix: vpp_ip.Prefix{
				Address: vpp_ip.Address{
					Af: vpp_ip.ADDRESS_IP6,
					Un: vpp_ip.AddressUnionIP6([16]uint8{255, 255, 10, 1}),
				},
			},
			Paths: []vpp_ip.FibPath{
				{
					SwIfIndex: 1,
				},
			},
		},
	})
	ctx.MockVpp.MockReply(&vpp_vpe.ControlPingReply{})

	rtDetails, err := l3handler.DumpRoutes()
	Expect(err).To(Succeed())
	Expect(rtDetails).To(HaveLen(2))
	Expect(rtDetails[0].Route.OutgoingInterface).To(Equal("if2"))
	Expect(rtDetails[1].Route.OutgoingInterface).To(Equal("if1"))
}
