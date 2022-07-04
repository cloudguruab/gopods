/*
Copyright 2022.

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package v1

import (
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)


type SharePolicy string

const (
	// Stale state triggers replacing exisitng share type
	// by simple removal and creation
	Stale SharePolicy = "Replace"

	// Ok state allows existing share types to remain operational
	Ok SharePolicy = "Allow"
)

// FormatSpec defines the desired state of Format
type FormatSpec struct {

	// Share Type Instance
	Share string `json:"share"`

	// Timestamp in seconds for share types initialized
	CreationTimeSeconds *int64 `json:"creationTimeSeconds"`

	// Policy for handling stale share types
	SharePolicy SharePolicy `json:"sharePolicy,omitempty"`

	// This flag determines if controller needs to handle stale shares
	// Usecase: if we extend this to handle migration of shares on a 
	// volume then we'd want to migrate items older than certain timestmap.
	Suspend *bool `json:"suspend,omitempty"`

	// List of existing shares on cluster
	ExistingShares []string `json:"existingShares"`

	// State of existing shares
	// Valid state:
	// - Available
	// - Creating 
	ShareState interface{} `json:"shareState,omitempty"`
}

// FormatStatus defines the observed state of Format
type FormatStatus struct {

	// List of live share types - FIXME: manilav1 undefined
	Ok []manilav1.ObjectReference `json:"ok,omitempty"`

	// Most recently fetched state
	RecentShareState *metav1.Time `json:"recentShareState,omitempty"`
}

//+kubebuilder:object:root=true
//+kubebuilder:subresource:status

// Format is the Schema for the formats API
type Format struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec   FormatSpec   `json:"spec,omitempty"`
	Status FormatStatus `json:"status,omitempty"`
}

//+kubebuilder:object:root=true

// FormatList contains a list of Format
type FormatList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	Items           []Format `json:"items"`
}

func init() {
	SchemeBuilder.Register(&Format{}, &FormatList{})
}
