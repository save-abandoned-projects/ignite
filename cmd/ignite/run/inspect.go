package run

import (
	"bytes"
	"fmt"
	api "github.com/save-abandoned-projects/ignite/pkg/apis/ignite"
	"github.com/save-abandoned-projects/ignite/pkg/apis/ignite/scheme"
	"github.com/save-abandoned-projects/ignite/pkg/providers"
	"github.com/save-abandoned-projects/libgitops/pkg/filter"
	"github.com/save-abandoned-projects/libgitops/pkg/runtime"
	"github.com/save-abandoned-projects/libgitops/pkg/serializer"
	"strings"
	"text/template"
)

// InspectFlags contains the flags supported by inspect.
type InspectFlags struct {
	OutputFormat   string
	TemplateFormat string
}

type InspectOptions struct {
	*InspectFlags
	object runtime.Object
}

// NewInspectOptions constructs and returns InspectOptions with the given kind
// and object ID.
func (i *InspectFlags) NewInspectOptions(k, objectMatch string) (*InspectOptions, error) {
	var err error
	var kind api.Kind
	io := &InspectOptions{InspectFlags: i}

	switch strings.ToLower(k) {
	case api.KindImage.Lower():
		kind = api.KindImage
	case api.KindKernel.Lower():
		kind = api.KindKernel
	case api.KindVM.Lower():
		kind = api.KindVM
	default:
		return nil, fmt.Errorf("unrecognized kind: %q", k)
	}

	if io.object, err = providers.Client.Dynamic(kind).Find(filter.NameFilter{Name: objectMatch}); err != nil {
		return nil, err
	}

	return io, nil
}

// Inspect renders the result of inspect in different formats based on the
// InspectOptions.
func Inspect(io *InspectOptions) error {
	var b bytes.Buffer
	var err error

	// If a template format is specified, render the template.
	if io.TemplateFormat != "" {
		output := &bytes.Buffer{}
		tmpl, err := template.New("").Parse(io.TemplateFormat)
		if err != nil {
			return fmt.Errorf("failed to parse template: %v", err)
		}
		if err := tmpl.Execute(output, io.object); err != nil {
			return fmt.Errorf("failed rendering template: %v", err)
		}
		fmt.Println(output.String())
		return nil
	}

	// Select the encoder and encode the object with it
	switch io.OutputFormat {
	case "json":
		err = scheme.Serializer.Encoder().Encode(serializer.NewJSONFrameWriter(&b), io.object)
	case "yaml":
		err = scheme.Serializer.Encoder().Encode(serializer.NewYAMLFrameWriter(&b), io.object)
	default:
		err = fmt.Errorf("unrecognized output format: %q", io.OutputFormat)
	}

	if err != nil {
		return err
	}

	// Print the encoded object
	fmt.Println(string(bytes.TrimSpace(b.Bytes())))
	return nil
}
