package filter

import (
	"fmt"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/types"
	"testing"
	"time"

	api "github.com/save-abandoned-projects/ignite/pkg/apis/ignite"
	"github.com/stretchr/testify/assert"
)

func TestMetaFiltering(t *testing.T) {
	t.Run("SuccessCPUsEqual", func(t *testing.T) {
		oMeta := &api.VM{
			ObjectMeta: metav1.ObjectMeta{
				UID:               types.UID("myuid"),
				CreationTimestamp: metav1.Time{},
				Labels: map[string]string{
					"first":  "f_value",
					"second": "s_value",
				},
			},
			Spec: api.VMSpec{
				CPUs: uint64(2),
			},
		}

		f := metaFilter{
			identifier:    "{{.Spec.CPUs}}",
			expectedValue: "2",
			operator:      "=",
		}

		res, err := f.isExpected(oMeta)
		assert.Nil(t, err)
		assert.True(t, res)

	})
	t.Run("SuccessNameEqual", func(t *testing.T) {
		oMeta := &api.VM{
			ObjectMeta: metav1.ObjectMeta{
				Name:              "success_object",
				UID:               types.UID("myuid"),
				CreationTimestamp: metav1.Time{},
				Labels: map[string]string{
					"first":  "f_value",
					"second": "s_value",
				},
			},
		}

		f := metaFilter{
			identifier:    "{{.ObjectMeta.Name}}",
			expectedValue: "success_object",
			operator:      "=",
		}

		res, err := f.isExpected(oMeta)
		assert.Nil(t, err)
		assert.True(t, res)
	})
	t.Run("SuccessNameDiff", func(t *testing.T) {
		oMeta := &api.VM{
			ObjectMeta: metav1.ObjectMeta{
				Name:              "success_object",
				UID:               types.UID("myuid"),
				CreationTimestamp: metav1.Time{},
				Labels: map[string]string{
					"first":  "f_value",
					"second": "s_value",
				},
			},
		}

		f := metaFilter{
			identifier:    "{{.Name}}",
			expectedValue: "success_object_diff",
			operator:      "!=",
		}

		res, err := f.isExpected(oMeta)
		assert.Nil(t, err)
		assert.True(t, res)
	})
	t.Run("FailNameEqual", func(t *testing.T) {
		oMeta := &api.VM{
			ObjectMeta: metav1.ObjectMeta{
				Name: "fail_object",
			},
		}

		f := metaFilter{
			identifier:    "{{.Name}}",
			expectedValue: "success_object",
			operator:      "=",
		}

		res, err := f.isExpected(oMeta)
		assert.Nil(t, err)
		assert.False(t, res)
	})
	t.Run("FailNameDiff", func(t *testing.T) {
		oMeta := &api.VM{
			ObjectMeta: metav1.ObjectMeta{
				Name: "fail_object",
			},
		}

		f := metaFilter{
			identifier:    "{{.Name}}",
			expectedValue: "fail_object",
			operator:      "!=",
		}

		res, err := f.isExpected(oMeta)
		assert.Nil(t, err)
		assert.False(t, res)
	})
	t.Run("SuccessNameContains", func(t *testing.T) {
		oMeta := &api.VM{
			ObjectMeta: metav1.ObjectMeta{
				Name: "fail_object",
			},
		}

		f := metaFilter{
			identifier:    "{{.Name}}",
			expectedValue: "object",
			operator:      "=~",
		}

		res, err := f.isExpected(oMeta)
		assert.Nil(t, err)
		assert.True(t, res)
	})
	t.Run("SuccessNameNotContains", func(t *testing.T) {
		oMeta := &api.VM{
			ObjectMeta: metav1.ObjectMeta{
				Name: "fail_object",
			},
		}

		f := metaFilter{
			identifier:    "{{.Name}}",
			expectedValue: "object2",
			operator:      "!~",
		}

		res, err := f.isExpected(oMeta)
		assert.Nil(t, err)
		assert.True(t, res)
	})
	t.Run("FailNameContains", func(t *testing.T) {
		oMeta := &api.VM{
			ObjectMeta: metav1.ObjectMeta{
				Name: "fail_object",
			},
		}

		f := metaFilter{
			identifier:    "{{.Name}}",
			expectedValue: "object2",
			operator:      "=~",
		}

		res, err := f.isExpected(oMeta)
		assert.Nil(t, err)
		assert.False(t, res)
	})
	t.Run("FailNameNotContains", func(t *testing.T) {
		oMeta := &api.VM{
			ObjectMeta: metav1.ObjectMeta{
				Name: "fail_object",
			},
		}

		f := metaFilter{
			identifier:    "{{.Name}}",
			expectedValue: "object",
			operator:      "!~",
		}

		res, err := f.isExpected(oMeta)
		assert.Nil(t, err)
		assert.False(t, res)
	})
	t.Run("SuccessUID", func(t *testing.T) {
		oMeta := &api.VM{
			ObjectMeta: metav1.ObjectMeta{
				UID: types.UID("myuid"),
			},
		}

		f := metaFilter{
			identifier:    "{{.UID}}",
			expectedValue: "myuid",
			operator:      "=",
		}

		res, err := f.isExpected(oMeta)
		assert.Nil(t, err)
		assert.True(t, res)
	})
	t.Run("FailUID", func(t *testing.T) {
		oMeta := &api.VM{
			ObjectMeta: metav1.ObjectMeta{
				UID: "failuid",
			},
		}

		f := metaFilter{
			identifier:    "{{.UID}}",
			expectedValue: "myuid",
			operator:      "=",
		}

		res, err := f.isExpected(oMeta)
		assert.Nil(t, err)
		assert.False(t, res)
	})
	t.Run("SuccessCreated", func(t *testing.T) {
		nowtime := metav1.Time{
			Time: time.Now().UTC(),
		}
		oMeta := &api.VM{
			ObjectMeta: metav1.ObjectMeta{
				CreationTimestamp: nowtime,
			},
		}

		f := metaFilter{
			identifier:    "{{.CreationTimestamp}}",
			expectedValue: nowtime.String(),
			operator:      "=",
		}

		res, err := f.isExpected(oMeta)
		assert.Nil(t, err)
		assert.True(t, res)
	})
	t.Run("FailCreated", func(t *testing.T) {
		nowtime := metav1.Time{
			Time: time.Now().UTC(),
		}
		oMeta := &api.VM{
			ObjectMeta: metav1.ObjectMeta{
				CreationTimestamp: nowtime,
			},
		}

		othertime := nowtime.Add(time.Duration(5))
		f := metaFilter{
			identifier:    "{{.CreationTimestamp}}",
			expectedValue: othertime.String(),
			operator:      "=",
		}

		res, err := f.isExpected(oMeta)
		assert.Nil(t, err)
		assert.False(t, res)
	})
	t.Run("SuccessLabels", func(t *testing.T) {
		oMeta := &api.VM{
			ObjectMeta: metav1.ObjectMeta{
				Labels: map[string]string{
					"foo": "bar",
				},
			},
		}

		f := metaFilter{
			identifier:    "{{.Labels.foo}}",
			expectedValue: "bar",
			operator:      "=",
		}

		res, err := f.isExpected(oMeta)
		assert.Nil(t, err)
		assert.True(t, res)
	})
	t.Run("FailLabels", func(t *testing.T) {
		oMeta := &api.VM{
			ObjectMeta: metav1.ObjectMeta{
				Labels: map[string]string{
					"foo": "bar2",
				},
			},
		}

		f := metaFilter{
			identifier:    "{{.Labels.foo}}",
			expectedValue: "bar",
			operator:      "=",
		}

		res, err := f.isExpected(oMeta)
		assert.Nil(t, err)
		assert.False(t, res)
	})
}

func TestExtractKeyValueFiltering(t *testing.T) {
	tests := []struct {
		name string
		str  string
		key  string
		val  string
		op   string
		err  error
	}{
		{
			name: "Success1",
			str:  "{{.Name}}=t/a-r:g_et",
			key:  "{{.Name}}",
			val:  "t/a-r:g_et",
			op:   "=",
			err:  nil,
		},
		{
			name: "Success2",
			str:  "{{.Name}}!=ta-rg_et",
			key:  "{{.Name}}",
			val:  "ta-rg_et",
			op:   "!=",
			err:  nil,
		},
		{
			name: "Success3",
			str:  "{{.Name}}==ta-rg_et",
			key:  "{{.Name}}",
			val:  "ta-rg_et",
			op:   "==",
			err:  nil,
		},
		{
			name: "Success4",
			str:  "{{.Name}}=~ta-rg_et",
			key:  "{{.Name}}",
			val:  "ta-rg_et",
			op:   "=~",
			err:  nil,
		},
		{
			name: "Success5",
			str:  "{{.Name}}!~ta-rg_et",
			key:  "{{.Name}}",
			val:  "ta-rg_et",
			op:   "!~",
			err:  nil,
		},
		{
			name: "Success6",
			str:  "{{.Name}}=8",
			key:  "{{.Name}}",
			val:  "8",
			op:   "=",
			err:  nil,
		},
		{
			name: "FailEqualBadPlace",
			str:  "{{.Name=}}target",
			key:  "",
			val:  "",
			op:   "",
			err:  fmt.Errorf("expected error"),
		},
		{
			name: "FailEqualBadPlace2",
			str:  "={{.Name}}target",
			key:  "",
			val:  "",
			op:   "",
			err:  fmt.Errorf("expected error"),
		},
		{
			name: "FailEqualBadPlace3",
			str:  "{{.Name}}tar=get",
			key:  "",
			val:  "",
			op:   "",
			err:  fmt.Errorf("expected error"),
		},
	}
	for _, utest := range tests {
		t.Run(utest.name, func(t *testing.T) {
			key, val, op, err := extractKeyValueFiltering(utest.str)
			if utest.err == nil {
				assert.Nil(t, err)
			} else {
				assert.NotNil(t, err)
			}
			assert.Equal(t, utest.key, key)
			assert.Equal(t, utest.val, val)
			assert.Equal(t, utest.op, op)
		})
	}
}

func TestExtractMultipleKeyValueFiltering(t *testing.T) {
	tests := []struct {
		name string
		str  string
		res  []metaFilter
		err  error
	}{
		{
			name: "Success",
			str:  "{{.Name}}=target1,{{.Age}}=38",
			res: []metaFilter{

				{
					identifier:    "{{.Name}}",
					expectedValue: "target1",
					operator:      "=",
				},
				{
					identifier:    "{{.Age}}",
					expectedValue: "38",
					operator:      "=",
				},
			},
			err: nil,
		},
		{
			name: "FailWithoutSeparator",
			str:  "{{.Name}}=target1{{.Age}}=38",
			res:  nil,
			err:  fmt.Errorf("expected error"),
		},
		{
			name: "FailBadFormat",
			str:  "{{.Name}}=target1{{.Age}}38",
			res:  nil,
			err:  fmt.Errorf("expected error"),
		},
	}
	for _, utest := range tests {
		t.Run(utest.name, func(t *testing.T) {
			res, err := extractMultipleKeyValueFiltering(utest.str)
			if err != nil {
				assert.NotNil(t, err)
			} else {
				assert.Nil(t, err)
			}
			assert.Equal(t, utest.res, res)
		})
	}
}

func TestMultipleMetaFilter(t *testing.T) {
	tests := []struct {
		name     string
		str      string
		object   *api.VM
		expected bool
		err      error
	}{
		{
			name: "SuccessOneFilter",
			str:  "{{.Name}}=hello",
			object: &api.VM{
				ObjectMeta: metav1.ObjectMeta{
					Name: "hello",
					UID:  "123",
				},
			},
			expected: true,
			err:      nil,
		},
		{
			name: "SuccessTwoFilter",
			str:  "{{.Name}}=hello,{{.UID}}=123",
			object: &api.VM{
				ObjectMeta: metav1.ObjectMeta{
					Name: "hello",
					UID:  "123",
				},
			},
			expected: true,
			err:      nil,
		},
		{
			name: "SuccessOneValueDiffer",
			str:  "{{.Name}}=hello,{{.UID}}=1234",
			object: &api.VM{
				ObjectMeta: metav1.ObjectMeta{
					Name: "hello",
					UID:  "123",
				},
			},
			expected: false,
			err:      nil,
		},
		{
			name: "FailBadFormat",
			str:  "{{.Name}}=hello,{{.Unexisting}}=1234",
			object: &api.VM{
				ObjectMeta: metav1.ObjectMeta{
					Name: "hello",
					UID:  "123",
				},
			},
			expected: false,
			err:      fmt.Errorf("expected error"),
		},
	}

	for _, utest := range tests {
		t.Run(utest.name, func(t *testing.T) {
			mmf, err := GenerateMultipleMetadataFiltering(utest.str)
			assert.Nil(t, err)
			expected, err := mmf.AreExpected(utest.object)
			if utest.err != nil {
				assert.NotNil(t, err)
			} else {
				assert.Nil(t, err)
			}
			assert.Equal(t, utest.expected, expected)
		})
	}
}
