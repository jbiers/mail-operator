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

package controller

import (
	"context"

	mailv1 "github.com/jbiers/mail-operator/api/v1"
	"github.com/jbiers/mail-operator/provider"
	"k8s.io/apimachinery/pkg/runtime"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/log"
)

// EmailReconciler reconciles a Email object
type EmailReconciler struct {
	client.Client
	Scheme *runtime.Scheme
}

//+kubebuilder:rbac:groups=mail.my.domain,resources=emails,verbs=get;list;watch;create;update;patch;delete
//+kubebuilder:rbac:groups=mail.my.domain,resources=emails/status,verbs=get;update;patch
//+kubebuilder:rbac:groups=mail.my.domain,resources=emails/finalizers,verbs=update

// Reconcile is part of the main kubernetes reconciliation loop which aims to
// move the current state of the cluster closer to the desired state.
// TODO(user): Modify the Reconcile function to compare the state specified by
// the Email object against the actual cluster state, and then
// perform operations to make the cluster state reflect the state specified by
// the user.
//
// For more details, check Reconcile and its Result here:
// - https://pkg.go.dev/sigs.k8s.io/controller-runtime@v0.17.0/pkg/reconcile
func (r *EmailReconciler) Reconcile(ctx context.Context, req ctrl.Request) (ctrl.Result, error) {
	logger := log.FromContext(ctx)
	var email mailv1.Email

	err := r.Get(ctx, req.NamespacedName, &email)
	if err != nil {
		logger.Error(err, "error getting Email resource.", "name", req.NamespacedName)
		return ctrl.Result{}, client.IgnoreNotFound(nil)
	}

	if email.Status.DeliveryStatus == "" {
		err, messageId, deliveryStatus := provider.SendEmail(&provider.EmailData{
			ApiKey:    email.Spec.SenderConfigRef, // This will actually be data gotten from the senderConfig
			Text:      email.Spec.Body,
			Subject:   email.Spec.Subject,
			Sender:    "anon@juliacodes.net", // This will actually be data gotten from the senderConfig
			Recipient: email.Spec.RecipientEmail,
		})

		email.Status = mailv1.EmailStatus{
			MessageId:      messageId,
			DeliveryStatus: deliveryStatus,
		}

		if err != nil || email.Status.DeliveryStatus != "202 Accepted" {
			logger.Error(err, "error sending email.", "name", req.NamespacedName)
			return ctrl.Result{}, client.IgnoreNotFound(nil)
		}

		err = r.Status().Update(ctx, &email)
		if err != nil {
			logger.Error(err, "error updating Email resource status.", "name", req.NamespacedName)
			return ctrl.Result{}, client.IgnoreNotFound(nil)
		}
	}

	return ctrl.Result{}, nil
}

// SetupWithManager sets up the controller with the Manager.
func (r *EmailReconciler) SetupWithManager(mgr ctrl.Manager) error {
	return ctrl.NewControllerManagedBy(mgr).
		For(&mailv1.Email{}).
		Complete(r)
}
