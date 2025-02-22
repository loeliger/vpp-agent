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

	vpp_ifs "github.com/ligato/vpp-agent/plugins/vpp/binapi/vpp2001/interfaces"
	. "github.com/onsi/gomega"
)

func TestSetInterfaceMtu(t *testing.T) {
	ctx, ifHandler := ifTestSetup(t)
	defer ctx.TeardownTestCtx()

	ctx.MockVpp.MockReply(&vpp_ifs.HwInterfaceSetMtuReply{})

	err := ifHandler.SetInterfaceMtu(1, 1500)

	Expect(err).To(BeNil())
	vppMsg, ok := ctx.MockChannel.Msg.(*vpp_ifs.HwInterfaceSetMtu)
	Expect(ok).To(BeTrue())
	Expect(vppMsg.SwIfIndex).To(BeEquivalentTo(1))
	Expect(vppMsg.Mtu).To(BeEquivalentTo(1500))
}

func TestSetInterfaceMtuError(t *testing.T) {
	ctx, ifHandler := ifTestSetup(t)
	defer ctx.TeardownTestCtx()

	ctx.MockVpp.MockReply(&vpp_ifs.HwInterfaceSetMtu{})

	err := ifHandler.SetInterfaceMtu(1, 1500)

	Expect(err).ToNot(BeNil())
}

func TestSetInterfaceMtuRetval(t *testing.T) {
	ctx, ifHandler := ifTestSetup(t)
	defer ctx.TeardownTestCtx()

	ctx.MockVpp.MockReply(&vpp_ifs.HwInterfaceSetMtuReply{
		Retval: 1,
	})

	err := ifHandler.SetInterfaceMtu(1, 1500)

	Expect(err).ToNot(BeNil())
}
