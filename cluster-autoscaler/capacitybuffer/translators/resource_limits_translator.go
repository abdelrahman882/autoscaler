/*
Copyright 2025 The Kubernetes Authors.

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

package translator

import (
	"fmt"
	"math"

	corev1 "k8s.io/api/core/v1"
	api_v1 "k8s.io/autoscaler/cluster-autoscaler/apis/capacitybuffer/autoscaling.x-k8s.io/v1"
	"k8s.io/autoscaler/cluster-autoscaler/capacitybuffer/common"
	"k8s.io/client-go/kubernetes"
)

// resourceLimitsTranslator translates buffers processors into pod capacity.
type resourceLimitsTranslator struct {
	client kubernetes.Interface
}

// NewResourceLimitsTranslator creates an instance of resourceLimitsTranslator.
func NewResourceLimitsTranslator(client kubernetes.Interface) *resourceLimitsTranslator {
	return &resourceLimitsTranslator{
		client: client,
	}
}

// Translate translates buffers processors into pod capacity.
func (t *resourceLimitsTranslator) Translate(buffers []*api_v1.CapacityBuffer) []error {
	errors := []error{}
	for _, buffer := range buffers {
		if isResourcesLimitsDefinedInBuffer(buffer) {

			if buffer.Status.PodTemplateRef == nil {
				errors = append(errors, fmt.Errorf("Can't get pod template, PodTemplateRef is nil"))
				continue
			}
			podTemplate, err := common.GetPodTemplate(t.client, buffer.Status.PodTemplateRef.Name)
			if err != nil {
				errors = append(errors, fmt.Errorf("Couldn't get pod template, error: %v", err))
				continue
			}

			numberOfPods, err := limitNumberOfPodsForResource(podTemplate, *buffer.Spec.Limits)
			if err != nil {
				errors = append(errors, fmt.Errorf("Skipping applying resource limits on buffer %v, with error: %v", buffer.Name, err.Error()))
				continue
			}
			if buffer.Status.Replicas != nil {
				numberOfPods = min(*buffer.Status.Replicas, numberOfPods)
			}
			setBufferAsReadyForProvisioning(buffer, buffer.Status.PodTemplateRef.Name, podTemplate.Generation, numberOfPods)
		}
	}
	return errors
}

func limitNumberOfPodsForResource(podTemplate *corev1.PodTemplate, limits api_v1.ResourceList) (int32, error) {
	maximumNumberOfPods := math.MaxInt
	podResourcesValues := map[string]int64{}

	for _, container := range podTemplate.Template.Spec.Containers {
		for resourceName, quantity := range container.Resources.Requests {
			if _, found := limits[api_v1.ResourceName(resourceName.String())]; found {
				podResourcesValues[resourceName.String()] += quantity.MilliValue()
			}
		}
	}
	for resourceName, quantity := range podResourcesValues {
		if quantity <= 0 {
			continue
		}
		if limitQuantity, found := limits[api_v1.ResourceName(resourceName)]; found {
			maxPods := limitQuantity.MilliValue() / quantity
			if maxPods < 0 {
				continue
			}
			maximumNumberOfPods = min(maximumNumberOfPods, int(maxPods))
		}
	}

	return int32(max(0, maximumNumberOfPods)), nil
}

func isResourcesLimitsDefinedInBuffer(buffer *api_v1.CapacityBuffer) bool {
	return buffer.Spec.Limits != nil
}

// CleanUp cleans up the translator's internal structures.
func (t *resourceLimitsTranslator) CleanUp() {
}
