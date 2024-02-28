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
	"fmt"
	"os"

	mailv1 "github.com/jbiers/mail-operator/api/v1"
	"github.com/jbiers/mail-operator/internal/provider"
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
//+kubebuilder:rbac:groups=mail.my.domain,resources=emailsenderconfigs,verbs=get;list;watch;create;update;patch;delete
//+kubebuilder:rbac:groups=mail.my.domain,resources=emailsenderconfigs/status,verbs=get;update;patch
//+kubebuilder:rbac:groups=mail.my.domain,resources=emailsenderconfigs/finalizers,verbs=update

func (r *EmailReconciler) Reconcile(ctx context.Context, req ctrl.Request) (ctrl.Result, error) {
	logger := log.FromContext(ctx)
	var email mailv1.Email

	if err := r.Get(ctx, req.NamespacedName, &email); err != nil {
		logger.Error(err, "error getting Email resource.", "name", req.NamespacedName)
		return ctrl.Result{}, client.IgnoreNotFound(err)
	}

	var emailSenderConfig mailv1.EmailSenderConfig
	emailSenderConfigData := req.NamespacedName
	emailSenderConfigData.Name = email.Spec.SenderConfigRef

	if err := r.Get(ctx, emailSenderConfigData, &emailSenderConfig); err != nil {
		logger.Error(err, "error getting EmailSenderConfig resource.", "name", emailSenderConfigData)
		return ctrl.Result{}, err
	}

	apiToken := os.Getenv(emailSenderConfig.Spec.ApiToken)
	if apiToken == "" {
		err := fmt.Errorf("empty value for environment variable %s", emailSenderConfig.Spec.ApiToken)

		logger.Error(err, "error getting the API token from emailSenderConfig", "name", emailSenderConfigData)
		return ctrl.Result{}, err
	}

	// if Email resource was just created
	if email.Status.DeliveryStatus == "" {
		var p provider.EmailProvider

		switch emailSenderConfig.Spec.ApiToken {
		case "mailersend":
			p = &provider.MailerSend{}
		case "mailgun":
			p = &provider.MailGun{}
		default:
			err := fmt.Errorf("invalid value for API token")

			logger.Error(err, "supported: mailersend, mailgun", "name", emailSenderConfigData)
			return ctrl.Result{}, err
		}

		emailSender := initEmailSender(p)

		logger.Info("sending email defined in Email resource.", "name", req.NamespacedName)

		err, messageId, deliveryStatus := emailSender.sendEmail(&provider.EmailData{
			ApiToken:  apiToken,
			Text:      email.Spec.Body,
			Subject:   email.Spec.Subject,
			Sender:    emailSenderConfig.Spec.SenderEmail,
			Recipient: email.Spec.RecipientEmail,
		})

		email.Status = mailv1.EmailStatus{
			MessageId:      messageId,
			DeliveryStatus: deliveryStatus,
		}

		if updateErr := r.Status().Update(ctx, &email); err != nil {
			logger.Error(updateErr, "error updating Email resource status.", "name", req.NamespacedName)
			return ctrl.Result{}, updateErr
		}

		if err != nil || deliveryStatus != "202 Accepted" {
			m := fmt.Sprintf(`there was an issue in delivering the email %s. code: %s`, messageId, deliveryStatus)
			logger.Error(err, m, "name", req.NamespacedName)

			return ctrl.Result{}, nil
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
