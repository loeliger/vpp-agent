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

package vpp2001_test

import (
	"testing"

	"github.com/ligato/vpp-agent/plugins/vpp/l3plugin/vrfidx"

	"github.com/ligato/cn-infra/logging/logrus"
	l3 "github.com/ligato/vpp-agent/api/models/vpp/l3"
	netallock_mock "github.com/ligato/vpp-agent/plugins/netalloc/mock"
	vpp_ip "github.com/ligato/vpp-agent/plugins/vpp/binapi/vpp2001/ip"
	"github.com/ligato/vpp-agent/plugins/vpp/ifplugin/ifaceidx"
	ifvppcalls "github.com/ligato/vpp-agent/plugins/vpp/ifplugin/vppcalls"
	ifvpp2001 "github.com/ligato/vpp-agent/plugins/vpp/ifplugin/vppcalls/vpp2001"
	"github.com/ligato/vpp-agent/plugins/vpp/l3plugin/vppcalls"
	"github.com/ligato/vpp-agent/plugins/vpp/l3plugin/vppcalls/vpp2001"
	"github.com/ligato/vpp-agent/plugins/vpp/vppcallmock"
	. "github.com/onsi/gomega"
)

var routes = []*l3.Route{
	{
		VrfId:             1,
		DstNetwork:        "192.168.10.21/24",
		NextHopAddr:       "192.168.30.1",
		OutgoingInterface: "iface1",
	},
	{
		VrfId:       2,
		DstNetwork:  "10.0.0.1/24",
		NextHopAddr: "192.168.30.1",
	},
	{
		VrfId:             2,
		DstNetwork:        "10.11.0.1/16",
		NextHopAddr:       "192.168.30.1",
		OutgoingInterface: "iface3",
	},
}

// Test adding routes
func TestAddRoute(t *testing.T) {
	ctx, _, rtHandler := routeTestSetup(t)
	defer ctx.TeardownTestCtx()

	ctx.MockVpp.MockReply(&vpp_ip.IPRouteAddDelReply{})
	err := rtHandler.VppAddRoute(routes[0])
	Expect(err).To(Succeed())

	ctx.MockVpp.MockReply(&vpp_ip.IPRouteAddDelReply{})
	err = rtHandler.VppAddRoute(routes[2])
	Expect(err).To(Not(BeNil())) // unknown interface
}

// Test deleting routes
func TestDeleteRoute(t *testing.T) {
	ctx, _, rtHandler := routeTestSetup(t)
	defer ctx.TeardownTestCtx()

	ctx.MockVpp.MockReply(&vpp_ip.IPRouteAddDelReply{})
	err := rtHandler.VppDelRoute(routes[0])
	Expect(err).To(Succeed())

	ctx.MockVpp.MockReply(&vpp_ip.IPRouteAddDelReply{})
	err = rtHandler.VppDelRoute(routes[1])
	Expect(err).To(Succeed())

	ctx.MockVpp.MockReply(&vpp_ip.IPRouteAddDelReply{Retval: 1})
	err = rtHandler.VppDelRoute(routes[0])
	Expect(err).To(Not(BeNil()))
}

func routeTestSetup(t *testing.T) (*vppcallmock.TestCtx, ifvppcalls.InterfaceVppAPI, vppcalls.RouteVppAPI) {
	ctx := vppcallmock.SetupTestCtx(t)
	log := logrus.NewLogger("test-log")
	ifHandler := ifvpp2001.NewInterfaceVppHandler(ctx.MockChannel, log)
	ifIndexes := ifaceidx.NewIfaceIndex(logrus.NewLogger("test-if"), "test-if")
	vrfIndexes := vrfidx.NewVRFIndex(logrus.NewLogger("test-vrf"), "test-vrf")
	ifIndexes.Put("iface1", &ifaceidx.IfaceMetadata{
		SwIfIndex: 1,
	})
	rtHandler := vpp2001.NewRouteVppHandler(ctx.MockChannel, ifIndexes, vrfIndexes, netallock_mock.NewMockNetAlloc(), log)
	return ctx, ifHandler, rtHandler
}
