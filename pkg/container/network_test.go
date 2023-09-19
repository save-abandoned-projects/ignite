package container

import (
	"testing"

	api "github.com/save-abandoned-projects/ignite/pkg/apis/ignite"
	"gotest.tools/assert"
)

func TestParseExtraIntfs(t *testing.T) {
	cases := []struct {
		name        string
		annotations map[string]string
		wantIntfs   map[string]string
	}{
		{
			name:      "empty object",
			wantIntfs: make(map[string]string),
		},
		{
			name: "wrong annotations",
			annotations: map[string]string{
				"foo":                                 "bar",
				"ignite.weave.works/interface/":       "dhcp-bridge",
				"ignite.weave.works/interface/eth123": "foo",
			},
			wantIntfs: make(map[string]string),
		},
		{
			name: "one interface",
			annotations: map[string]string{
				"foo":                                 "bar",
				"ignite.weave.works/interface/":       "dhcp-bridge",
				"ignite.weave.works/interface/eth123": "tc-redirect",
			},
			wantIntfs: map[string]string{
				"eth123": "tc-redirect",
			},
		},
		{
			name: "many interfaces",
			annotations: map[string]string{
				"foo":                                 "bar",
				"ignite.weave.works/interface/eth0":   "dhcp-bridge",
				"ignite.weave.works/interface/eth123": "tc-redirect",
			},
			wantIntfs: map[string]string{
				"eth0":   "dhcp-bridge",
				"eth123": "tc-redirect",
			},
		},
	}

	for _, rt := range cases {
		t.Run(rt.name, func(t *testing.T) {
			vm := &api.VM{}
			vm.SetAnnotations(rt.annotations)
			parsedIntfs := parseExtraIntfs(vm)

			assert.DeepEqual(t, parsedIntfs, rt.wantIntfs)

		})
	}
}
