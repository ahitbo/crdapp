package handlers

import (
	"context"
	"log"
	"net/http"

	corev1 "k8s.io/api/core/v1"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/runtime/inject"
	"sigs.k8s.io/controller-runtime/pkg/webhook/admission"
	"sigs.k8s.io/controller-runtime/pkg/webhook/admission/types"
)

type PodAnnotator struct {
	client  client.Client
	decoder types.Decoder
}

// PodAnnotator implements admission.Handler.
var _ admission.Handler = &PodAnnotator{}

func (a *PodAnnotator) Handle(ctx context.Context, req types.Request) types.Response {
	log.Printf("Webhook Handle Request %s/\n", req.AdmissionRequest)
	pod := &corev1.Pod{}

	err := a.decoder.Decode(req, pod)
	if err != nil {
		return admission.ErrorResponse(http.StatusBadRequest, err)
	}
	copy := pod.DeepCopy()

	err = a.mutatePodsFn(ctx, copy)
	if err != nil {
		return admission.ErrorResponse(http.StatusInternalServerError, err)
	}
	// admission.PatchResponse generates a Response containing patches.
	return admission.PatchResponse(pod, copy)
}

func (a *PodAnnotator) mutatePodsFn(ctx context.Context, pod *corev1.Pod) error {
	if pod.Annotations == nil {
		pod.Annotations = map[string]string{}
	}
	pod.Annotations["example-mutating-admission-webhook"] = "foo"
	return nil
}

// PodAnnotator implements inject.Client.
var _ inject.Client = &PodAnnotator{}

// InjectClient injects the client into the PodAnnotator
func (a *PodAnnotator) InjectClient(c client.Client) error {
	a.client = c
	return nil
}

// PodAnnotator implements inject.Decoder.
var _ inject.Decoder = &PodAnnotator{}

// InjectDecoder injects the decoder into the PodAnnotator
func (a *PodAnnotator) InjectDecoder(d types.Decoder) error {
	a.decoder = d
	return nil
}
