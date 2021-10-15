package flowdesc_test

import (
	"free5gc/lib/flowdesc"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestIPFilterRuleEncode(t *testing.T) {
	testStr1 := "permit out ip from any to assigned 655"

	rule := flowdesc.NewIPFilterRule()
	if rule == nil {
		t.Fatal("IP Filter Rule Create Error")
	}

	if err := rule.SetAction(flowdesc.Permit); err != nil {
		assert.Nil(t, err)
	}

	if err := rule.SetDirection(flowdesc.Out); err != nil {
		assert.Nil(t, err)
	}

	if err := rule.SetProtocal(0xfc); err != nil {
		assert.Nil(t, err)
	}

	if err := rule.SetSourceIP("any"); err != nil {
		assert.Nil(t, err)
	}

	if err := rule.SetDestinationIP("assigned"); err != nil {
		assert.Nil(t, err)
	}

	if err := rule.SetDestinationPorts("655"); err != nil {
		assert.Nil(t, err)
	}

	result, err := flowdesc.Encode(rule)
	if err != nil {
		assert.Nil(t, err)
	}
	if result != testStr1 {
		t.Fatalf("Encode error, \n\t expect: %s,\n\t    get: %s", testStr1, result)
	}
}

func TestIPFilterRuleDecode(t *testing.T) {
	var testCases = map[string]struct {
		filterRule string
		action     flowdesc.Action
		dir        flowdesc.Direction
		proto      uint8
		src        string
		srcPorts   string
		dst        string
		dstPorts   string
	}{
		"fully": {
			filterRule: "permit out ip from 60.60.0.100 8080 to 60.60.0.1 80",
			action:     flowdesc.Action(flowdesc.Permit),
			dir:        flowdesc.Direction(flowdesc.Out),
			proto:      flowdesc.ProtocalNumberAny,
			src:        "60.60.0.100",
			srcPorts:   "8080",
			dst:        "60.60.0.1",
			dstPorts:   "80",
		},
		"withoutPorts": {
			filterRule: "permit out ip from 60.60.0.100 to 60.60.0.1",
			action:     flowdesc.Action(flowdesc.Permit),
			dir:        flowdesc.Direction(flowdesc.Out),
			proto:      flowdesc.ProtocalNumberAny,
			src:        "60.60.0.100",
			srcPorts:   "",
			dst:        "60.60.0.1",
			dstPorts:   "",
		},
		"withoutOnePorts": {
			filterRule: "permit out ip from 60.60.0.100 8080 to 60.60.0.1",
			action:     flowdesc.Action(flowdesc.Permit),
			dir:        flowdesc.Direction(flowdesc.Out),
			proto:      flowdesc.ProtocalNumberAny,
			src:        "60.60.0.100",
			srcPorts:   "8080",
			dst:        "60.60.0.1",
			dstPorts:   "",
		},
		"withSrcAny": {
			filterRule: "permit out ip from any to 60.60.0.1 8080",
			action:     flowdesc.Action(flowdesc.Permit),
			dir:        flowdesc.Direction(flowdesc.Out),
			proto:      flowdesc.ProtocalNumberAny,
			src:        "any",
			srcPorts:   "",
			dst:        "60.60.0.1",
			dstPorts:   "8080",
		},
		"withDstAny": {
			filterRule: "permit out ip from 60.60.0.1 8080 to any",
			action:     flowdesc.Action(flowdesc.Permit),
			dir:        flowdesc.Direction(flowdesc.Out),
			proto:      flowdesc.ProtocalNumberAny,
			src:        "60.60.0.1",
			srcPorts:   "8080",
			dst:        "any",
			dstPorts:   "",
		},
		"withAssigned": {
			filterRule: "permit out ip from assigned to 60.60.0.1 8080",
			action:     flowdesc.Action(flowdesc.Permit),
			dir:        flowdesc.Direction(flowdesc.Out),
			proto:      flowdesc.ProtocalNumberAny,
			src:        "assigned",
			srcPorts:   "",
			dst:        "60.60.0.1",
			dstPorts:   "8080",
		},
	}

	for testName, expected := range testCases {
		t.Run(testName, func(t *testing.T) {
			r := flowdesc.NewIPFilterRule()
			err := flowdesc.Decode(expected.filterRule, r)

			require.Equal(t, expected.action, r.GetAction())
			require.Equal(t, expected.dir, r.GetDirection())
			require.Equal(t, expected.proto, r.GetProtocal())
			require.Equal(t, expected.src, r.GetSourceIP())
			require.Equal(t, expected.srcPorts, r.GetSourcePorts())
			require.Equal(t, expected.dst, r.GetDestinationIP())
			require.Equal(t, expected.dstPorts, r.GetDestinationPorts())

			require.NoError(t, err)
		})
	}
}

func TestIPFilterRuleSwapSourceAndDestination(t *testing.T) {
	var testCases = map[string]struct {
		filterRule string
		action     flowdesc.Action
		dir        flowdesc.Direction
		proto      uint8
		src        string
		srcPorts   string
		dst        string
		dstPorts   string
	}{
		"fully": {
			filterRule: "permit out ip from 60.60.0.100 8080 to 60.60.0.1 80",
			action:     flowdesc.Action(flowdesc.Permit),
			dir:        flowdesc.Direction(flowdesc.Out),
			proto:      flowdesc.ProtocalNumberAny,
			src:        "60.60.0.1",
			srcPorts:   "80",
			dst:        "60.60.0.100",
			dstPorts:   "8080",
		},
		"withoutPorts": {
			filterRule: "permit out ip from 60.60.0.100 to 60.60.0.1",
			action:     flowdesc.Action(flowdesc.Permit),
			dir:        flowdesc.Direction(flowdesc.Out),
			proto:      flowdesc.ProtocalNumberAny,
			src:        "60.60.0.1",
			srcPorts:   "",
			dst:        "60.60.0.100",
			dstPorts:   "",
		},
		"withoutOnePorts": {
			filterRule: "permit out ip from 60.60.0.100 8080 to 60.60.0.1",
			action:     flowdesc.Action(flowdesc.Permit),
			dir:        flowdesc.Direction(flowdesc.Out),
			proto:      flowdesc.ProtocalNumberAny,
			src:        "60.60.0.1",
			srcPorts:   "",
			dst:        "60.60.0.100",
			dstPorts:   "8080",
		},
		"withSrcAny": {
			filterRule: "permit out ip from any to 60.60.0.1 8080",
			action:     flowdesc.Action(flowdesc.Permit),
			dir:        flowdesc.Direction(flowdesc.Out),
			proto:      flowdesc.ProtocalNumberAny,
			src:        "60.60.0.1",
			srcPorts:   "8080",
			dst:        "any",
			dstPorts:   "",
		},
		"withDstAny": {
			filterRule: "permit out ip from 60.60.0.1 8080 to any",
			action:     flowdesc.Action(flowdesc.Permit),
			dir:        flowdesc.Direction(flowdesc.Out),
			proto:      flowdesc.ProtocalNumberAny,
			src:        "any",
			srcPorts:   "",
			dst:        "60.60.0.1",
			dstPorts:   "8080",
		},
		"withAssigned": {
			filterRule: "permit out ip from assigned to 60.60.0.1 8080",
			action:     flowdesc.Action(flowdesc.Permit),
			dir:        flowdesc.Direction(flowdesc.Out),
			proto:      flowdesc.ProtocalNumberAny,
			src:        "60.60.0.1",
			srcPorts:   "8080",
			dst:        "assigned",
			dstPorts:   "",
		},
	}

	for testName, expected := range testCases {
		t.Run(testName, func(t *testing.T) {
			r := flowdesc.NewIPFilterRule()
			err := flowdesc.Decode(expected.filterRule, r)
			r.SwapSourceAndDestination()
			require.Equal(t, expected.action, r.GetAction())
			require.Equal(t, expected.dir, r.GetDirection())
			require.Equal(t, expected.proto, r.GetProtocal())
			require.Equal(t, expected.src, r.GetSourceIP())
			require.Equal(t, expected.srcPorts, r.GetSourcePorts())
			require.Equal(t, expected.dst, r.GetDestinationIP())
			require.Equal(t, expected.dstPorts, r.GetDestinationPorts())

			require.NoError(t, err)
		})
	}
}
