package webhook

import (
	admissionregistrationv1beta1 "k8s.io/api/admissionregistration/v1beta1"
	corev1 "k8s.io/api/core/v1"
	"sigs.k8s.io/controller-runtime/pkg/manager"
	"sigs.k8s.io/controller-runtime/pkg/webhook"
	"sigs.k8s.io/controller-runtime/pkg/webhook/admission/builder"

	"crdapp/pkg/webhook/handlers"

	apitypes "k8s.io/apimachinery/pkg/types"
)

func init() {
	// AddToManagerFuncs is a list of functions to create controllers and add them to a manager.
	AddToManagerFuncs = append(AddToManagerFuncs, Add)
}

func Add(mgr manager.Manager) error {
	return add(mgr)
}

func add(mgr manager.Manager) error {
	vd, err := builder.NewWebhookBuilder().
		Validating().
		Operations(admissionregistrationv1beta1.Create).
		ForType(&corev1.Pod{}).
		Handlers(&handlers.PodValidater{}).
		WithManager(mgr).
		Build()
	if err != nil {
		return err
	}
	//创建Webhook//规则在这里设置
	wh, err := builder.NewWebhookBuilder().
		Mutating().
		Operations(admissionregistrationv1beta1.Create).
		ForType(&corev1.Pod{}).
		Handlers(&handlers.PodAnnotator{}).
		WithManager(mgr).
		Build()

	if err != nil {
		return err
	}
	//创建Webhook Server
	svr, err := webhook.NewServer("application-admission-server", mgr, webhook.ServerOptions{
		CertDir: "/tmp/cert",
		BootstrapOptions: &webhook.BootstrapOptions{
			Secret: &apitypes.NamespacedName{
				Namespace: "crdapp-system",
				Name:      "crdapp-webhook-server-secret",
			},

			Service: &webhook.Service{
				Namespace: "crdapp-system",
				Name:      "crdapp-admission-server-service",
				// Selectors should select the pods that runs this webhook server.
				Selectors: map[string]string{
					"control-plane":           "controller-manager",
					"controller-tools.k8s.io": "1.0",
				},
			},
		},
	})

	if err != nil {
		return err
	}

	//注册
	svr.Register(vd, wh)
	return nil
}
