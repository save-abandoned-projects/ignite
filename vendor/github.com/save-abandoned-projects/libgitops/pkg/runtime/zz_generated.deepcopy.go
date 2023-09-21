// +build !ignore_autogenerated

// Code generated by deepcopy-gen. DO NOT EDIT.

package runtime

import (
	pkgruntime "k8s.io/apimachinery/pkg/runtime"
)

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *Metav1NameIdentifierFactory) DeepCopyInto(out *Metav1NameIdentifierFactory) {
	*out = *in
	return
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new Metav1NameIdentifierFactory.
func (in *Metav1NameIdentifierFactory) DeepCopy() *Metav1NameIdentifierFactory {
	if in == nil {
		return nil
	}
	out := new(Metav1NameIdentifierFactory)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *ObjectUIDIdentifierFactory) DeepCopyInto(out *ObjectUIDIdentifierFactory) {
	*out = *in
	return
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new ObjectUIDIdentifierFactory.
func (in *ObjectUIDIdentifierFactory) DeepCopy() *ObjectUIDIdentifierFactory {
	if in == nil {
		return nil
	}
	out := new(ObjectUIDIdentifierFactory)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *PartialObjectImpl) DeepCopyInto(out *PartialObjectImpl) {
	*out = *in
	out.TypeMeta = in.TypeMeta
	in.ObjectMeta.DeepCopyInto(&out.ObjectMeta)
	return
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new PartialObjectImpl.
func (in *PartialObjectImpl) DeepCopy() *PartialObjectImpl {
	if in == nil {
		return nil
	}
	out := new(PartialObjectImpl)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyObject is an autogenerated deepcopy function, copying the receiver, creating a new pkgruntime.Object.
func (in *PartialObjectImpl) DeepCopyObject() pkgruntime.Object {
	if c := in.DeepCopy(); c != nil {
		return c
	}
	return nil
}
