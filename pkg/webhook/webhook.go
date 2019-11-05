package webhook

import (
	"context"

	v1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/api/errors"
	"k8s.io/client-go/rest"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/manager"
	logf "sigs.k8s.io/controller-runtime/pkg/runtime/log"
	"sigs.k8s.io/controller-runtime/pkg/webhook"
)

var log = logf.Log.WithName("webhook")

// AddToManagerFuncs is a list of functions to add all Webhooks to the Manager
var AddToManagerFuncs []func(manager.Manager) (webhook.Webhook, error)

// AddToManager adds all Webhooks to the Manager
func AddToManager(m manager.Manager) error {

	if !runningOnMinikube(m.GetConfig()) {
		log.Info("Minikube not detected; no webhooks will be configured")
		return nil
	}

	webhooks := []webhook.Webhook{}
	for _, f := range AddToManagerFuncs {
		wh, err := f(m)
		if err != nil {
			log.Error(err, "Unable to setup webhook")
			return err
		}
		webhooks = append(webhooks, wh)
	}
	if len(webhooks) == 0 {
		return nil
	}

	log.Info("Setting up webhook server")
	// This will be started when the Manager is started
	as, err := webhook.NewServer("admission-webhook-server", m, webhook.ServerOptions{
		Port:    9876,
		CertDir: "/tmp/cert",
		BootstrapOptions: &webhook.BootstrapOptions{
			Service: &webhook.Service{
				Namespace: "default",
				Name:      "admission-server-service",
				// Selectors should select the pods that runs this webhook server.
				Selectors: map[string]string{
					"app": "minikube-admission-server",
				},
			},
		},
	})
	if err != nil {
		log.Error(err, "Unable to create a new webhook server")
		return err
	}

	log.Info("Registering webhooks to the webhook server")
	err = as.Register(webhooks...)
	if err != nil {
		log.Error(err, "Unable to register webhooks in the admission server")
		return err
	}
	return nil
}

func runningOnMinikube(cfg *rest.Config) bool {
	c, err := client.New(cfg, client.Options{})
	if err != nil {
		log.Error(err, "Can't create client")
		return false
	}
	node := &v1.Node{}
	if err := c.Get(context.TODO(), client.ObjectKey{Name: "minikube"}, node); err != nil {
		if !errors.IsNotFound(err) {
			log.Error(err, "Unable to query for minikube node")
		}
		return false
	}
	return true
}
