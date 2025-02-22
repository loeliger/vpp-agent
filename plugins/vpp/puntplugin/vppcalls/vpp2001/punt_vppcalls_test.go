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

	"github.com/ligato/cn-infra/logging/logrus"
	. "github.com/onsi/gomega"

	punt "github.com/ligato/vpp-agent/api/models/vpp/punt"
	vpp_ip "github.com/ligato/vpp-agent/plugins/vpp/binapi/vpp2001/ip"
	vpp_punt "github.com/ligato/vpp-agent/plugins/vpp/binapi/vpp2001/punt"
	"github.com/ligato/vpp-agent/plugins/vpp/ifplugin/ifaceidx"
	"github.com/ligato/vpp-agent/plugins/vpp/puntplugin/vppcalls"
	"github.com/ligato/vpp-agent/plugins/vpp/puntplugin/vppcalls/vpp2001"
	"github.com/ligato/vpp-agent/plugins/vpp/vppcallmock"
)

// TODO test below temporary disabled (re-enable with set_punt)
/*
func TestAddPunt(t *testing.T) {
	ctx, puntHandler, _ := puntTestSetup(t)
	defer ctx.TeardownTestCtx()

	ctx.MockVpp.MockReply(&ba_punt.SetPuntReply{})

	err := puntHandler.AddPunt(&punt.ToHost{
		L3Protocol: punt.L3Protocol_IPv4,
		L4Protocol: punt.L4Protocol_UDP,
		Port:       9000,
	})

	Expect(err).To(BeNil())
	vppMsg, ok := ctx.MockChannel.Msg.(*ba_punt.SetPunt)
	Expect(ok).To(BeTrue())
	Expect(vppMsg.IsAdd).To(Equal(uint8(1)))
	Expect(vppMsg.Punt.IPv).To(Equal(uint8(4)))
	Expect(vppMsg.Punt.L4Protocol).To(Equal(uint8(17)))
	Expect(vppMsg.Punt.L4Port).To(Equal(uint16(9000)))
}

func TestDeletePunt(t *testing.T) {
	ctx, puntHandler, _ := puntTestSetup(t)
	defer ctx.TeardownTestCtx()

	ctx.MockVpp.MockReply(&ba_punt.SetPuntReply{})

	err := puntHandler.DeletePunt(&punt.ToHost{
		L3Protocol: punt.L3Protocol_IPv4,
		L4Protocol: punt.L4Protocol_UDP,
		Port:       9000,
	})

	Expect(err).To(BeNil())
	vppMsg, ok := ctx.MockChannel.Msg.(*ba_punt.SetPunt)
	Expect(ok).To(BeTrue())
	Expect(vppMsg.IsAdd).To(Equal(uint8(0)))
	Expect(vppMsg.Punt.IPv).To(Equal(uint8(4)))
	Expect(vppMsg.Punt.L4Protocol).To(Equal(uint8(17)))
	Expect(vppMsg.Punt.L4Port).To(Equal(uint16(9000)))
}
*/

func TestRegisterPuntSocket(t *testing.T) {
	ctx, puntHandler, _ := puntTestSetup(t)
	defer ctx.TeardownTestCtx()

	ctx.MockVpp.MockReply(&vpp_punt.PuntSocketRegisterReply{
		Pathname: []byte("/othersock"),
	})

	path, err := puntHandler.RegisterPuntSocket(&punt.ToHost{
		L3Protocol: punt.L3Protocol_IPv4,
		L4Protocol: punt.L4Protocol_UDP,
		Port:       9000,
		SocketPath: "/test/path/socket",
	})

	Expect(err).To(BeNil())
	vppMsg, ok := ctx.MockChannel.Msg.(*vpp_punt.PuntSocketRegister)
	Expect(ok).To(BeTrue())
	Expect(vppMsg.HeaderVersion).To(Equal(uint32(1)))
	Expect(vppMsg.Punt.Punt.GetL4().Af).To(Equal(vpp_punt.ADDRESS_IP4))
	Expect(vppMsg.Punt.Punt.GetL4().Protocol).To(Equal(vpp_punt.IP_API_PROTO_UDP))
	Expect(vppMsg.Punt.Punt.GetL4().Port).To(Equal(uint16(9000)))
	Expect(path).To(Equal("/othersock"))
}

func TestRegisterPuntSocketAllIPv4(t *testing.T) {
	ctx, puntHandler, _ := puntTestSetup(t)
	defer ctx.TeardownTestCtx()

	ctx.MockVpp.MockReply(&vpp_punt.PuntSocketRegisterReply{
		Pathname: []byte("/othersock"),
	})
	ctx.MockVpp.MockReply(&vpp_punt.PuntSocketRegisterReply{
		Pathname: []byte("/othersock"),
	})

	path, err := puntHandler.RegisterPuntSocket(&punt.ToHost{
		L3Protocol: punt.L3Protocol_ALL,
		L4Protocol: punt.L4Protocol_UDP,
		Port:       9000,
		SocketPath: "/test/path/socket",
	})

	Expect(err).To(BeNil())
	for _, msg := range ctx.MockChannel.Msgs {
		vppMsg, ok := msg.(*vpp_punt.PuntSocketRegister)
		Expect(ok).To(BeTrue())

		if vppMsg.Punt.Punt.GetL4().Af == vpp_punt.ADDRESS_IP4 {
			Expect(vppMsg.HeaderVersion).To(Equal(uint32(1)))
			Expect(vppMsg.Punt.Punt.GetL4().Protocol).To(Equal(vpp_punt.IP_API_PROTO_UDP))
			Expect(vppMsg.Punt.Punt.GetL4().Port).To(Equal(uint16(9000)))
			Expect(path).To(Equal("/othersock"))
		}
		if vppMsg.Punt.Punt.GetL4().Af == vpp_punt.ADDRESS_IP6 {
			Expect(vppMsg.HeaderVersion).To(Equal(uint32(1)))
			Expect(vppMsg.Punt.Punt.GetL4().Protocol).To(Equal(vpp_punt.IP_API_PROTO_UDP))
			Expect(vppMsg.Punt.Punt.GetL4().Port).To(Equal(uint16(9000)))
			Expect(path).To(Equal("/othersock"))
		}
	}
}

func TestAddIPRedirect(t *testing.T) {
	ctx, puntHandler, ifIndexes := puntTestSetup(t)
	defer ctx.TeardownTestCtx()

	ctx.MockVpp.MockReply(&vpp_ip.IPPuntRedirectReply{})

	ifIndexes.Put("if1", &ifaceidx.IfaceMetadata{SwIfIndex: 1})
	ifIndexes.Put("if2", &ifaceidx.IfaceMetadata{SwIfIndex: 2})

	err := puntHandler.AddPuntRedirect(&punt.IPRedirect{
		L3Protocol:  punt.L3Protocol_IPv4,
		RxInterface: "if1",
		TxInterface: "if2",
		NextHop:     "10.0.0.1",
	})

	Expect(err).To(BeNil())
	vppMsg, ok := ctx.MockChannel.Msg.(*vpp_ip.IPPuntRedirect)
	Expect(ok).To(BeTrue())
	Expect(vppMsg.IsAdd).To(Equal(uint8(1)))
	Expect(vppMsg.Punt.Nh.Af).To(Equal(vpp_ip.ADDRESS_IP4))
	Expect(vppMsg.Punt.RxSwIfIndex).To(Equal(uint32(1)))
	Expect(vppMsg.Punt.TxSwIfIndex).To(Equal(uint32(2)))
	//Expect(vppMsg.Nh).To(Equal([]uint8(net.ParseIP("10.0.0.1").To4())))
}

func TestAddIPRedirectAll(t *testing.T) {
	ctx, puntHandler, ifIndexes := puntTestSetup(t)
	defer ctx.TeardownTestCtx()

	ctx.MockVpp.MockReply(&vpp_ip.IPPuntRedirectReply{})

	ifIndexes.Put("if1", &ifaceidx.IfaceMetadata{SwIfIndex: 1})

	err := puntHandler.AddPuntRedirect(&punt.IPRedirect{
		L3Protocol:  punt.L3Protocol_IPv4,
		TxInterface: "if1",
		NextHop:     "30.0.0.1",
	})

	Expect(err).To(BeNil())
	vppMsg, ok := ctx.MockChannel.Msg.(*vpp_ip.IPPuntRedirect)
	Expect(ok).To(BeTrue())
	Expect(vppMsg.IsAdd).To(Equal(uint8(1)))
	//Expect(vppMsg.IsIP6).To(Equal(uint8(0)))
	Expect(vppMsg.Punt.Nh.Af).To(Equal(vpp_ip.ADDRESS_IP4))
	Expect(vppMsg.Punt.RxSwIfIndex).To(Equal(^uint32(0)))
	Expect(vppMsg.Punt.TxSwIfIndex).To(Equal(uint32(1)))
	//Expect(vppMsg.Nh).To(Equal([]uint8(net.ParseIP("30.0.0.1").To4())))
}

func TestDeleteIPRedirect(t *testing.T) {
	ctx, puntHandler, ifIndexes := puntTestSetup(t)
	defer ctx.TeardownTestCtx()

	ctx.MockVpp.MockReply(&vpp_ip.IPPuntRedirectReply{})

	ifIndexes.Put("if1", &ifaceidx.IfaceMetadata{SwIfIndex: 1})
	ifIndexes.Put("if2", &ifaceidx.IfaceMetadata{SwIfIndex: 2})

	err := puntHandler.DeletePuntRedirect(&punt.IPRedirect{
		L3Protocol:  punt.L3Protocol_IPv4,
		RxInterface: "if1",
		TxInterface: "if2",
		NextHop:     "10.0.0.1",
	})

	Expect(err).To(BeNil())
	vppMsg, ok := ctx.MockChannel.Msg.(*vpp_ip.IPPuntRedirect)
	Expect(ok).To(BeTrue())
	Expect(vppMsg.IsAdd).To(Equal(uint8(0)))
	//Expect(vppMsg.IsIP6).To(Equal(uint8(0)))
	Expect(vppMsg.Punt.Nh.Af).To(Equal(vpp_ip.ADDRESS_IP4))
	Expect(vppMsg.Punt.RxSwIfIndex).To(Equal(uint32(1)))
	Expect(vppMsg.Punt.TxSwIfIndex).To(Equal(uint32(2)))
	//Expect(vppMsg.Nh).To(Equal([]uint8(net.ParseIP("10.0.0.1").To4())))
}

func puntTestSetup(t *testing.T) (*vppcallmock.TestCtx, vppcalls.PuntVppAPI, ifaceidx.IfaceMetadataIndexRW) {
	ctx := vppcallmock.SetupTestCtx(t)
	logger := logrus.NewLogger("test-log")
	ifIndexes := ifaceidx.NewIfaceIndex(logger, "punt-if-idx")
	puntHandler := vpp2001.NewPuntVppHandler(ctx.MockChannel, ifIndexes, logrus.DefaultLogger())
	return ctx, puntHandler, ifIndexes
}
