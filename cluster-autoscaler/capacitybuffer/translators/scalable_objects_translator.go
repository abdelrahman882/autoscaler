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

	cbclient "k8s.io/autoscaler/cluster-autoscaler/capacitybuffer/client"
	"k8s.io/autoscaler/cluster-autoscaler/capacitybuffer/common"

	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	apiv1 "k8s.io/autoscaler/cluster-autoscaler/apis/capacitybuffer/autoscaling.x-k8s.io/v1"
	scalableobject "k8s.io/autoscaler/cluster-autoscaler/capacitybuffer/translators/scalable_objects"
)

// ScalableObjectsTranslator translates buffers processors into pod capacity.
type ScalableObjectsTranslator struct {
	client                           *cbclient.CapacityBufferClient
	supportedScalableObjectResolvers map[string]scalableobject.ScalableObjectResolver
}

// NewDefaultScalableObjectsTranslator creates an instance of ScalableObjectsTranslator.
func NewDefaultScalableObjectsTranslator(client *cbclient.CapacityBufferClient) *ScalableObjectsTranslator {
	supportedScalableObjects := map[string]scalableobject.ScalableObjectResolver{}
	for _, scalableObject := range scalableobject.GetSupportedScalableObjectResolvers(client) {
		supportedScalableObjects[scalableObject.GetKind()] = scalableObject
	}
	return &ScalableObjectsTranslator{
		client:                           client,
		supportedScalableObjectResolvers: supportedScalableObjects,
	}
}

// Translate translates buffers processors into pod capacity.
func (t *ScalableObjectsTranslator) Translate(buffers []*apiv1.CapacityBuffer) []error {

	errors := []error{}
	for _, buffer := range buffers {
		if isScalableObjectBuffer(buffer) {
			scalableObj, found := t.supportedScalableObjectResolvers[buffer.Spec.ScalableRef.Kind]
			if !found {
				err := fmt.Errorf("Couldn't get pod template spec for scalable object buffer %v: kind %v is not supported", buffer.Name, buffer.Spec.ScalableRef.Kind)
				common.SetBufferAsNotReadyForProvisioning(buffer, nil, nil, nil, err)
				errors = append(errors, err)
				continue
			}
			podTemplateSpec, replicasFromScalable, err := scalableObj.GetTemplateAndReplicas(corev1.NamespaceDefault, buffer.Spec.ScalableRef.Name)
			if err != nil {
				err := fmt.Errorf("Couldn't get pod template spec for scalable object buffer %v for referenced object %v, with error: %v", buffer.Name, buffer.Spec.ScalableRef.Name, err.Error())
				common.SetBufferAsNotReadyForProvisioning(buffer, nil, nil, nil, err)
				errors = append(errors, err)
				continue
			}
			podTemplate := t.getPodTemplateFromSpecs(podTemplateSpec, buffer)
			createdPodTemplate, err := t.client.GetPodTemplate(podTemplate.Name)
			if err != nil {
				createdPodTemplate, err = t.client.CreatePodTemplate(podTemplate)
			} else {
				createdPodTemplate, err = t.client.UpdatePodTemplate(podTemplate)
			}
			if err != nil {
				err := fmt.Errorf("Failed to create pod template object for buffer %v with error: %v", buffer.Name, err.Error())
				common.SetBufferAsNotReadyForProvisioning(buffer, nil, nil, nil, err)
				errors = append(errors, err)
				continue
			}
			numberOfPods := t.getBufferNumberOfPods(buffer, replicasFromScalable)
			if numberOfPods == nil {
				common.SetBufferAsNotReadyForProvisioning(buffer, &apiv1.LocalObjectRef{Name: createdPodTemplate.Name}, &createdPodTemplate.Generation, nil, fmt.Errorf("Couldn't get number of replicas for buffer %v, replicas and percentage are not defined", buffer.Name))
				continue
			}

			common.SetBufferAsReadyForProvisioning(buffer, &apiv1.LocalObjectRef{Name: createdPodTemplate.Name}, &createdPodTemplate.Generation, numberOfPods)

		}
	}
	return errors
}

func (t *ScalableObjectsTranslator) getBufferNumberOfPods(buffer *apiv1.CapacityBuffer, scalableReplicas *int32) *int32 {

	var numberOfPodsFromPercentage *int32
	var numberOfPodsFromReplicas *int32

	if buffer.Spec.Percentage != nil {
		if scalableReplicas != nil {
			percentValue := buffer.Spec.Percentage
			numberOfPods := max(0, int32(int32(*percentValue)*(*scalableReplicas)/100.0))
			numberOfPodsFromPercentage = &numberOfPods
		}
	}
	if buffer.Spec.Replicas != nil {
		numberOfPods := max(0, int32(*buffer.Spec.Replicas))
		numberOfPodsFromReplicas = &numberOfPods
	}
	if numberOfPodsFromPercentage != nil && numberOfPodsFromReplicas != nil {
		numberOfPods := min(*numberOfPodsFromPercentage, *numberOfPodsFromReplicas)
		return &numberOfPods
	} else if numberOfPodsFromPercentage != nil {
		return numberOfPodsFromPercentage
	}
	return numberOfPodsFromReplicas
}

func (t *ScalableObjectsTranslator) getPodTemplateFromSpecs(pts *corev1.PodTemplateSpec, buffer *apiv1.CapacityBuffer) *corev1.PodTemplate {
	return &corev1.PodTemplate{
		ObjectMeta: metav1.ObjectMeta{
			Name:      fmt.Sprintf("capacitybuffer-%v-pod-template", buffer.Name),
			Namespace: cbclient.DefaultNamespace,
			OwnerReferences: []metav1.OwnerReference{
				{
					APIVersion:         common.CapacityBufferApiVersion,
					Kind:               common.CapacityBufferKind,
					Name:               buffer.Name,
					UID:                buffer.UID,
					Controller:         pointerBool(true),
					BlockOwnerDeletion: pointerBool(true),
				},
			},
		},
		Template: corev1.PodTemplateSpec{
			ObjectMeta: metav1.ObjectMeta{
				Labels:      pts.Labels,
				Annotations: pts.Annotations,
			},
			Spec: pts.Spec,
		},
	}
}

func isScalableObjectBuffer(buffer *apiv1.CapacityBuffer) bool {
	return buffer.Spec.ScalableRef != nil
}

// CleanUp cleans up the translator's internal structures.
func (t *ScalableObjectsTranslator) CleanUp() {
}

func pointerBool(b bool) *bool {
	return &b
}
