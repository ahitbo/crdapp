package handlers

import (
	"context"
	"log"
	"net/http"
	"strings"

	corev1 "k8s.io/api/core/v1"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/runtime/inject"
	"sigs.k8s.io/controller-runtime/pkg/webhook/admission"
	"sigs.k8s.io/controller-runtime/pkg/webhook/admission/types"
)

type PodValidater struct {
	client  client.Client
	decoder types.Decoder
}

// podAnnotator implements admission.Handler.
var _ admission.Handler = &PodValidater{}

func (a *PodValidater) Handle(ctx context.Context, req types.Request) types.Response {
	log.Printf("Validating Webhook Handle Request %s/\n", req.AdmissionRequest)
	pod := &corev1.Pod{}
	err := a.decoder.Decode(req, pod)

	if err != nil {
		return admission.ErrorResponse(http.StatusBadRequest, err)
	}
	name := pod.GetName()
	if strings.Compare(name, "nginx-pod") == 0 {
		if pod.GetAnnotations() == nil || pod.GetAnnotations()["appid"] == "" {
			return admission.ValidationResponse(false, "appid is not allowed null")
		}
	}
	return admission.ValidationResponse(true, "ok")
}

// podAnnotator implements inject.Client.
var _ inject.Client = &PodValidater{}

// InjectClient injects the client into the podAnnotator
func (a *PodValidater) InjectClient(c client.Client) error {
	a.client = c
	return nil
}

// podAnnotator implements inject.Decoder.
var _ inject.Decoder = &PodValidater{}

// InjectDecoder injects the decoder into the podAnnotator
func (a *PodValidater) InjectDecoder(d types.Decoder) error {
	a.decoder = d
	return nil
}
