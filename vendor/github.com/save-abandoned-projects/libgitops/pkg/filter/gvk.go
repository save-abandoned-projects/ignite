package filter

import (
	"fmt"
	"github.com/save-abandoned-projects/libgitops/pkg/runtime"
	"k8s.io/apimachinery/pkg/runtime/schema"
)

// NameFilter implements ObjectFilter and ListOption.
var _ ObjectFilter = GvkFilter{}
var _ ListOption = GvkFilter{}

// GvkFilter is an ObjectFilter that compares runtime.Object.GetObjectKind().GroupVersionKind()
// to the Gvk field by either equality
type GvkFilter struct {
	// Gvk matches the object's gvk
	// +required
	Gvk *schema.GroupVersionKind
}

// Filter implements ObjectFilter
func (g GvkFilter) Filter(obj runtime.Object) (bool, error) {
	// Require g.Gvk to always be set.
	if g.Gvk == nil {
		return false, fmt.Errorf("the GvkFilter.Gvk field must not be nil: %w", ErrInvalidFilterParams)
	}

	// Otherwise, just use an equality check
	return g.Gvk.String() == obj.GetObjectKind().GroupVersionKind().String(), nil
}

// ApplyToListOptions implements ListOption, and adds itself converted to
// a ListFilter to ListOptions.Filters.
func (g GvkFilter) ApplyToListOptions(target *ListOptions) error {
	target.Filters = append(target.Filters, ObjectToListFilter(g))
	return nil
}
