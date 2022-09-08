package custom_resources

import (
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
)

type DeploymentObejctSpec struct {
	Kind      string `json:"kind"`
	NReplicas int32  `json:"nReplicas"`
}

type DeploymentObject struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec DeploymentObejctSpec `json:"spec"`
}

type DeploymentObjectList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`

	Items []DeploymentObject `json:"items"`
}

// DeepCopyInto copies all properties of this object into another object of the
// same type that is provided as a pointer.
func (in *DeploymentObject) DeepCopyInto(out *DeploymentObject) {
	out.TypeMeta = in.TypeMeta
	out.ObjectMeta = in.ObjectMeta
	out.Spec = DeploymentObejctSpec{
		Kind:      in.Spec.Kind,
		NReplicas: in.Spec.NReplicas,
	}
}

// DeepCopyObject returns a generically typed copy of an object
func (in *DeploymentObject) DeepCopyObject() runtime.Object {
	out := DeploymentObject{}
	in.DeepCopyInto(&out)

	return &out
}

// DeepCopyObject returns a generically typed copy of an object
func (in *DeploymentObjectList) DeepCopyObject() runtime.Object {
	out := DeploymentObjectList{}
	out.TypeMeta = in.TypeMeta
	out.ListMeta = in.ListMeta

	if in.Items != nil {
		out.Items = make([]DeploymentObject, len(in.Items))
		for i := range in.Items {
			in.Items[i].DeepCopyInto(&out.Items[i])
		}
	}

	return &out
}
