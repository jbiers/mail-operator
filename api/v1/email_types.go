/*
Copyright 2024.

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

// EmailSpec defines the desired state of Email
type EmailSpec struct {
	SenderConfigRef string `json:"senderConfigRef,omitempty"`
	RecipientEmail  string `json:"recipientEmail,omitempty"`
	Subject         string `json:"subject,omitempty"`
	Body            string `json:"body,omitempty"`
}

// EmailStatus defines the observed state of Email
type EmailStatus struct {
	DeliveryStatus string `json:"deliveryStatus,omitempty"`
	MessageId      string `json:"messageId,omitempty"`
}

// +kubebuilder:object:root=true
// +kubebuilder:subresource:status
// +kubebuilder:printcolumn:name="status",type=string,JSONPath=`.status.deliveryStatus`
// +kubebuilder:printcolumn:name="recipient",type=string,JSONPath=`.spec.recipientEmail`
// +kubebuilder:printcolumn:name="sender",type=string,JSONPath=`.spec.senderConfigRef.senderEmail`

// Email is the Schema for the emails API
type Email struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec   EmailSpec   `json:"spec,omitempty"`
	Status EmailStatus `json:"status,omitempty"`
}

//+kubebuilder:object:root=true

// EmailList contains a list of Email
type EmailList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	Items           []Email `json:"items"`
}

func init() {
	SchemeBuilder.Register(&Email{}, &EmailList{})
}
