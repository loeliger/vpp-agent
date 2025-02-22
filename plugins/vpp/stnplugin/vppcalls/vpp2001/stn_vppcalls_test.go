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

	"github.com/ligato/cn-infra/logging/logrus"
	. "github.com/onsi/gomega"

	stn "github.com/ligato/vpp-agent/api/models/vpp/stn"
	vpp_stn "github.com/ligato/vpp-agent/plugins/vpp/binapi/vpp2001/stn"
	"github.com/ligato/vpp-agent/plugins/vpp/ifplugin/ifaceidx"
	"github.com/ligato/vpp-agent/plugins/vpp/stnplugin/vppcalls"
	"github.com/ligato/vpp-agent/plugins/vpp/stnplugin/vppcalls/vpp2001"
	"github.com/ligato/vpp-agent/plugins/vpp/vppcallmock"
)

func TestAddStnRule(t *testing.T) {
	ctx, stnHandler, ifIndexes := stnTestSetup(t)
	defer ctx.TeardownTestCtx()

	ifIndexes.Put("if1", &ifaceidx.IfaceMetadata{SwIfIndex: 1})

	ctx.MockVpp.MockReply(&vpp_stn.StnAddDelRuleReply{})

	err := stnHandler.AddSTNRule(&stn.Rule{
		Interface: "if1",
		IpAddress: "10.0.0.1",
	})

	Expect(err).To(BeNil())
	vppMsg, ok := ctx.MockChannel.Msg.(*vpp_stn.StnAddDelRule)
	Expect(ok).To(BeTrue())
	Expect(vppMsg.SwIfIndex).To(BeEquivalentTo(1))
	Expect(vppMsg.IPAddress).To(BeEquivalentTo(net.ParseIP("10.0.0.1").To4()))
	Expect(vppMsg.IsIP4).To(BeEquivalentTo(1))
	Expect(vppMsg.IsAdd).To(BeEquivalentTo(1))
}

func TestAddStnRuleIPv6(t *testing.T) {
	ctx, stnHandler, ifIndexes := stnTestSetup(t)
	defer ctx.TeardownTestCtx()

	ifIndexes.Put("if1", &ifaceidx.IfaceMetadata{SwIfIndex: 1})

	ctx.MockVpp.MockReply(&vpp_stn.StnAddDelRuleReply{})

	err := stnHandler.AddSTNRule(&stn.Rule{
		Interface: "if1",
		IpAddress: "2001:db8:0:1:1:1:1:1",
	})

	Expect(err).To(BeNil())
	vppMsg, ok := ctx.MockChannel.Msg.(*vpp_stn.StnAddDelRule)
	Expect(ok).To(BeTrue())
	Expect(vppMsg.SwIfIndex).To(BeEquivalentTo(1))
	Expect(vppMsg.IPAddress).To(BeEquivalentTo(net.ParseIP("2001:db8:0:1:1:1:1:1").To16()))
	Expect(vppMsg.IsIP4).To(BeEquivalentTo(0))
	Expect(vppMsg.IsAdd).To(BeEquivalentTo(1))
}

func TestAddStnRuleInvalidIP(t *testing.T) {
	ctx, stnHandler, ifIndexes := stnTestSetup(t)
	defer ctx.TeardownTestCtx()

	ifIndexes.Put("if1", &ifaceidx.IfaceMetadata{SwIfIndex: 1})

	ctx.MockVpp.MockReply(&vpp_stn.StnAddDelRuleReply{})

	err := stnHandler.AddSTNRule(&stn.Rule{
		Interface: "if1",
		IpAddress: "invalid-ip",
	})

	Expect(err).ToNot(BeNil())
}

func TestAddStnRuleError(t *testing.T) {
	ctx, stnHandler, ifIndexes := stnTestSetup(t)
	defer ctx.TeardownTestCtx()

	ifIndexes.Put("if1", &ifaceidx.IfaceMetadata{SwIfIndex: 1})

	ctx.MockVpp.MockReply(&vpp_stn.StnAddDelRule{})

	err := stnHandler.AddSTNRule(&stn.Rule{
		Interface: "if1",
		IpAddress: "10.0.0.1",
	})

	Expect(err).ToNot(BeNil())
}

func TestAddStnRuleRetval(t *testing.T) {
	ctx, stnHandler, ifIndexes := stnTestSetup(t)
	defer ctx.TeardownTestCtx()

	ifIndexes.Put("if1", &ifaceidx.IfaceMetadata{SwIfIndex: 1})

	ctx.MockVpp.MockReply(&vpp_stn.StnAddDelRuleReply{
		Retval: 1,
	})

	err := stnHandler.AddSTNRule(&stn.Rule{
		Interface: "if1",
		IpAddress: "10.0.0.1",
	})

	Expect(err).ToNot(BeNil())
}

func TestDelStnRule(t *testing.T) {
	ctx, stnHandler, ifIndexes := stnTestSetup(t)
	defer ctx.TeardownTestCtx()

	ifIndexes.Put("if1", &ifaceidx.IfaceMetadata{SwIfIndex: 1})

	ctx.MockVpp.MockReply(&vpp_stn.StnAddDelRuleReply{})

	err := stnHandler.DeleteSTNRule(&stn.Rule{
		Interface: "if1",
		IpAddress: "10.0.0.1",
	})

	Expect(err).To(BeNil())
	vppMsg, ok := ctx.MockChannel.Msg.(*vpp_stn.StnAddDelRule)
	Expect(ok).To(BeTrue())
	Expect(vppMsg.IsAdd).To(BeEquivalentTo(0))
}

func stnTestSetup(t *testing.T) (*vppcallmock.TestCtx, vppcalls.StnVppAPI, ifaceidx.IfaceMetadataIndexRW) {
	ctx := vppcallmock.SetupTestCtx(t)
	logger := logrus.NewLogger("test-log")
	ifIndexes := ifaceidx.NewIfaceIndex(logger, "stn-if-idx")
	stnHandler := vpp2001.NewStnVppHandler(ctx.MockChannel, ifIndexes, logrus.DefaultLogger())
	return ctx, stnHandler, ifIndexes
}
