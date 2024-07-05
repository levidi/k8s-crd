package types

import (
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/runtime/schema"
	"k8s.io/client-go/kubernetes/scheme"
)

// BucketSpec defines the spec for Bucket
type BucketSpec struct {
	BucketName string `json:"bucketName"`
	Region     string `json:"region"`
}

// Bucket defines the structure of our custom resource
type Bucket struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec BucketSpec `json:"spec"`
}

// DeepCopyObject implements runtime.Object interface
func (in *Bucket) DeepCopyObject() runtime.Object {
	out := &Bucket{
		TypeMeta:   in.TypeMeta,
		ObjectMeta: *in.ObjectMeta.DeepCopy(),
	}
	out.Spec.BucketName = in.Spec.BucketName
	out.Spec.Region = in.Spec.Region
	return out
}

// BucketList is a list of Bucket resources
type BucketList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	Items           []Bucket `json:"items"`
}

// DeepCopyObject implements runtime.Object interface
func (in *BucketList) DeepCopyObject() runtime.Object {
	out := &BucketList{
		TypeMeta: in.TypeMeta,
		ListMeta: *in.ListMeta.DeepCopy(),
	}
	out.Items = make([]Bucket, len(in.Items))
	for i := range in.Items {
		out.Items[i] = *in.Items[i].DeepCopyObject().(*Bucket)
	}
	return out
}

// AddToScheme adds the Bucket and BucketList types to the given scheme.
func AddToScheme(scheme *runtime.Scheme) error {
	scheme.AddKnownTypes(schema.GroupVersion{
		Group:   "levi.com", // API group
		Version: "v1",       // API version
	},
		&Bucket{},
		&BucketList{},
	)
	metav1.AddToGroupVersion(scheme, schema.GroupVersion{
		Group:   "levi.com", // API group
		Version: "v1",       // API version
	})
	return nil
}

func init() {
	// Add the schema for the Bucket type to the global scheme
	AddToScheme(scheme.Scheme)
}
