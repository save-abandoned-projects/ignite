// +build !ignore_autogenerated

// Code generated by deepcopy-gen. DO NOT EDIT.

package v1alpha1

import (
	net "net"
)

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *DMID) DeepCopyInto(out *DMID) {
	*out = *in
	return
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new DMID.
func (in *DMID) DeepCopy() *DMID {
	if in == nil {
		return nil
	}
	out := new(DMID)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in IPAddresses) DeepCopyInto(out *IPAddresses) {
	{
		in := &in
		*out = make(IPAddresses, len(*in))
		for i := range *in {
			if (*in)[i] != nil {
				in, out := &(*in)[i], &(*out)[i]
				*out = make(net.IP, len(*in))
				copy(*out, *in)
			}
		}
		return
	}
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new IPAddresses.
func (in IPAddresses) DeepCopy() IPAddresses {
	if in == nil {
		return nil
	}
	out := new(IPAddresses)
	in.DeepCopyInto(out)
	return *out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *OCIContentID) DeepCopyInto(out *OCIContentID) {
	*out = *in
	return
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new OCIContentID.
func (in *OCIContentID) DeepCopy() *OCIContentID {
	if in == nil {
		return nil
	}
	out := new(OCIContentID)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *OCIImageRef) DeepCopyInto(out *OCIImageRef) {
	*out = *in
	return
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new OCIImageRef.
func (in *OCIImageRef) DeepCopy() *OCIImageRef {
	if in == nil {
		return nil
	}
	out := new(OCIImageRef)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *PortMapping) DeepCopyInto(out *PortMapping) {
	*out = *in
	if in.BindAddress != nil {
		in, out := &in.BindAddress, &out.BindAddress
		*out = make(net.IP, len(*in))
		copy(*out, *in)
	}
	return
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new PortMapping.
func (in *PortMapping) DeepCopy() *PortMapping {
	if in == nil {
		return nil
	}
	out := new(PortMapping)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in PortMappings) DeepCopyInto(out *PortMappings) {
	{
		in := &in
		*out = make(PortMappings, len(*in))
		for i := range *in {
			(*in)[i].DeepCopyInto(&(*out)[i])
		}
		return
	}
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new PortMappings.
func (in PortMappings) DeepCopy() PortMappings {
	if in == nil {
		return nil
	}
	out := new(PortMappings)
	in.DeepCopyInto(out)
	return *out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *Size) DeepCopyInto(out *Size) {
	*out = *in
	return
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new Size.
func (in *Size) DeepCopy() *Size {
	if in == nil {
		return nil
	}
	out := new(Size)
	in.DeepCopyInto(out)
	return out
}
