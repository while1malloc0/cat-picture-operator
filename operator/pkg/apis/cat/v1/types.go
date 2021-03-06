package v1

import (
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object

type CatPictureList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata"`
	Items           []CatPicture `json:"items"`
}

// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object

type CatPicture struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata"`
	Spec              CatPictureSpec   `json:"spec"`
	Status            CatPictureStatus `json:"status,omitempty"`
}

type CatPictureSpec struct {
	// The number of containers to have running in the ReplicaSet
	Num int32 `json:"num,omitempty"`
	// The size of the served cat pictures. One of "small" or "full". Defaults to "small".
	Size string `json:"size,omitempty"`
	// The format that the cat pictures are served in.
	Format string `json:"format,omitempty"`
}

type CatPictureStatus struct {
	// Fill me
}
