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
	"net"
	"testing"

	ifs "github.com/ligato/vpp-agent/api/models/vpp/interfaces"
	vpp_afpacket "github.com/ligato/vpp-agent/plugins/vpp/binapi/vpp2001/af_packet"
	vpp_ifs "github.com/ligato/vpp-agent/plugins/vpp/binapi/vpp2001/interfaces"
	. "github.com/onsi/gomega"
)

func TestAddAfPacketInterface(t *testing.T) {
	ctx, ifHandler := ifTestSetup(t)
	defer ctx.TeardownTestCtx()

	ctx.MockVpp.MockReply(&vpp_afpacket.AfPacketCreateReply{})
	ctx.MockVpp.MockReply(&vpp_ifs.SwInterfaceTagAddDelReply{})

	ifIndex, err := ifHandler.AddAfPacketInterface("if1", "", &ifs.AfpacketLink{
		HostIfName: "host1",
	})

	Expect(err).To(BeNil())
	Expect(ifIndex).ToNot(BeNil())
	Expect(len(ctx.MockChannel.Msgs)).To(BeEquivalentTo(2))
	for i, msg := range ctx.MockChannel.Msgs {
		if i == 0 {
			vppMsg, ok := msg.(*vpp_afpacket.AfPacketCreate)
			Expect(ok).To(BeTrue())
			Expect(vppMsg).To(Equal(&vpp_afpacket.AfPacketCreate{
				HostIfName:      []byte("host1"),
				HwAddr:          nil,
				UseRandomHwAddr: 1,
			}))
		}
	}
}

func TestAddAfPacketInterfaceError(t *testing.T) {
	ctx, ifHandler := ifTestSetup(t)
	defer ctx.TeardownTestCtx()

	ctx.MockVpp.MockReply(&vpp_afpacket.AfPacketDeleteReply{})

	_, err := ifHandler.AddAfPacketInterface("if1", "", &ifs.AfpacketLink{
		HostIfName: "host1",
	})

	Expect(err).ToNot(BeNil())
}

func TestAddAfPacketInterfaceRetval(t *testing.T) {
	ctx, ifHandler := ifTestSetup(t)
	defer ctx.TeardownTestCtx()

	ctx.MockVpp.MockReply(&vpp_afpacket.AfPacketCreateReply{
		Retval: 1,
	})
	ctx.MockVpp.MockReply(&vpp_ifs.SwInterfaceTagAddDelReply{})

	_, err := ifHandler.AddAfPacketInterface("if1", "", &ifs.AfpacketLink{
		HostIfName: "host1",
	})

	Expect(err).ToNot(BeNil())
}

func TestDeleteAfPacketInterface(t *testing.T) {
	ctx, ifHandler := ifTestSetup(t)
	defer ctx.TeardownTestCtx()

	ctx.MockVpp.MockReply(&vpp_afpacket.AfPacketDeleteReply{})
	ctx.MockVpp.MockReply(&vpp_ifs.SwInterfaceTagAddDelReply{})

	err := ifHandler.DeleteAfPacketInterface("if1", 0, &ifs.AfpacketLink{
		HostIfName: "host1",
	})

	Expect(err).To(BeNil())
	Expect(len(ctx.MockChannel.Msgs)).To(BeEquivalentTo(2))
	for i, msg := range ctx.MockChannel.Msgs {
		if i == 0 {
			vppMsg, ok := msg.(*vpp_afpacket.AfPacketDelete)
			Expect(ok).To(BeTrue())
			Expect(vppMsg).To(Equal(&vpp_afpacket.AfPacketDelete{
				HostIfName: []byte("host1"),
			}))
		}
	}
}

func TestDeleteAfPacketInterfaceError(t *testing.T) {
	ctx, ifHandler := ifTestSetup(t)
	defer ctx.TeardownTestCtx()

	ctx.MockVpp.MockReply(&vpp_afpacket.AfPacketCreateReply{})

	err := ifHandler.DeleteAfPacketInterface("if1", 0, &ifs.AfpacketLink{
		HostIfName: "host1",
	})

	Expect(err).ToNot(BeNil())
}

func TestDeleteAfPacketInterfaceRetval(t *testing.T) {
	ctx, ifHandler := ifTestSetup(t)
	defer ctx.TeardownTestCtx()

	ctx.MockVpp.MockReply(&vpp_afpacket.AfPacketDeleteReply{
		Retval: 1,
	})
	ctx.MockVpp.MockReply(&vpp_ifs.SwInterfaceTagAddDelReply{})

	err := ifHandler.DeleteAfPacketInterface("if1", 0, &ifs.AfpacketLink{
		HostIfName: "host1",
	})

	Expect(err).ToNot(BeNil())
}

func TestAddAfPacketInterfaceMac(t *testing.T) {
	ctx, ifHandler := ifTestSetup(t)
	defer ctx.TeardownTestCtx()

	ctx.MockVpp.MockReply(&vpp_afpacket.AfPacketCreateReply{})
	ctx.MockVpp.MockReply(&vpp_ifs.SwInterfaceTagAddDelReply{})

	ifIndex, err := ifHandler.AddAfPacketInterface("if1", "a2:01:01:01:01:01", &ifs.AfpacketLink{
		HostIfName: "host1",
	})

	Expect(err).To(BeNil())
	Expect(ifIndex).ToNot(BeNil())
	Expect(len(ctx.MockChannel.Msgs)).To(BeEquivalentTo(2))

	mac, err := net.ParseMAC("a2:01:01:01:01:01")
	Expect(err).To(BeNil())

	for i, msg := range ctx.MockChannel.Msgs {
		if i == 0 {
			vppMsg, ok := msg.(*vpp_afpacket.AfPacketCreate)
			Expect(ok).To(BeTrue())
			Expect(vppMsg).To(Equal(&vpp_afpacket.AfPacketCreate{
				HostIfName:      []byte("host1"),
				HwAddr:          mac,
				UseRandomHwAddr: 0,
			}))
		}
	}
}
