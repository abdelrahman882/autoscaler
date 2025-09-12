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
	"testing"

	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/api/resource"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	fakeClient "k8s.io/client-go/kubernetes/fake"

	"github.com/stretchr/testify/assert"
	v1 "k8s.io/autoscaler/cluster-autoscaler/apis/capacitybuffer/autoscaling.x-k8s.io/v1"
	"k8s.io/autoscaler/cluster-autoscaler/capacitybuffer/testutil"
)

func TestResourceLimitsTranslator(t *testing.T) {
	podTemp4mem100cpu := getPodTemplateWithResources(testutil.SomePodTemplateRefName, corev1.ResourceList{
		"memory": resource.MustParse("4Gi"),
		"cpu":    resource.MustParse("100m"),
	})
	podTemp8mem200cpu := getPodTemplateWithResources(testutil.AnotherPodTemplateRefName, corev1.ResourceList{
		"memory": resource.MustParse("8Gi"),
		"cpu":    resource.MustParse("200m"),
	})
	fakeClient := fakeClient.NewSimpleClientset(podTemp4mem100cpu, podTemp8mem200cpu)
	tests := []struct {
		name                   string
		buffers                []*v1.CapacityBuffer
		expectedStatus         []*v1.CapacityBufferStatus
		expectedNumberOfErrors int
	}{
		{
			name: "Limits set to nil, buffer filtered out",
			buffers: []*v1.CapacityBuffer{
				getTestBufferWithLimits(nil, nil, nil),
			},
			expectedStatus: []*v1.CapacityBufferStatus{
				testutil.GetBufferStatus(nil, nil, nil, nil),
			},
			expectedNumberOfErrors: 0,
		},
		{
			name: "Limits exist, podTemplateRef is nil",
			buffers: []*v1.CapacityBuffer{
				getTestBufferWithLimits(nil, nil, &v1.ResourceList{
					"cpu": resource.MustParse("500m"),
				}),
			},
			expectedStatus: []*v1.CapacityBufferStatus{
				testutil.GetBufferStatus(nil, nil, nil, nil),
			},
			expectedNumberOfErrors: 1,
		},
		{
			name: "Limits exist and no replicas, buffer filtered",
			buffers: []*v1.CapacityBuffer{
				getTestBufferWithLimits(&v1.LocalObjectRef{Name: podTemp4mem100cpu.Name}, nil, &v1.ResourceList{
					"memory": resource.MustParse("5Gi"),
					"cpu":    resource.MustParse("200m"),
				}),
			},
			expectedStatus: []*v1.CapacityBufferStatus{
				getTestBufferStatusWithReplicas(&v1.LocalObjectRef{Name: podTemp4mem100cpu.Name}, 1),
			},
			expectedNumberOfErrors: 0,
		},
		{
			name: "Limits exist and no replicas, buffer filtered",
			buffers: []*v1.CapacityBuffer{
				getTestBufferWithLimits(&v1.LocalObjectRef{Name: podTemp4mem100cpu.Name}, nil, &v1.ResourceList{
					"memory": resource.MustParse("9Gi"),
					"cpu":    resource.MustParse("200m"),
				}),
			},
			expectedStatus: []*v1.CapacityBufferStatus{
				getTestBufferStatusWithReplicas(&v1.LocalObjectRef{Name: podTemp4mem100cpu.Name}, 2),
			},
			expectedNumberOfErrors: 0,
		},
		{
			name: "Limits exist and with bigger replicas",
			buffers: []*v1.CapacityBuffer{
				getTestBufferWithLimits(&v1.LocalObjectRef{Name: podTemp4mem100cpu.Name}, int32Ptr(3), &v1.ResourceList{
					"memory": resource.MustParse("9Gi"),
					"cpu":    resource.MustParse("200m"),
				}),
			},
			expectedStatus: []*v1.CapacityBufferStatus{
				getTestBufferStatusWithReplicas(&v1.LocalObjectRef{Name: podTemp4mem100cpu.Name}, 2),
			},
			expectedNumberOfErrors: 0,
		},
		{
			name: "Limits exist and with smaller replicas",
			buffers: []*v1.CapacityBuffer{
				getTestBufferWithLimits(&v1.LocalObjectRef{Name: podTemp4mem100cpu.Name}, int32Ptr(1), &v1.ResourceList{
					"memory": resource.MustParse("9Gi"),
					"cpu":    resource.MustParse("200m"),
				}),
			},
			expectedStatus: []*v1.CapacityBufferStatus{
				getTestBufferStatusWithReplicas(&v1.LocalObjectRef{Name: podTemp4mem100cpu.Name}, 1),
			},
			expectedNumberOfErrors: 0,
		},
		{
			name: "Limits exist and no replicas, buffer filtered",
			buffers: []*v1.CapacityBuffer{
				getTestBufferWithLimits(&v1.LocalObjectRef{Name: podTemp4mem100cpu.Name}, int32Ptr(5), &v1.ResourceList{
					"memory": resource.MustParse("10Gi"),
					"cpu":    resource.MustParse("1000m"),
				}),
				getTestBufferWithLimits(&v1.LocalObjectRef{Name: podTemp8mem200cpu.Name}, int32Ptr(3), &v1.ResourceList{
					"memory": resource.MustParse("100Gi"),
					"cpu":    resource.MustParse("10000m"),
				}),
			},
			expectedStatus: []*v1.CapacityBufferStatus{
				getTestBufferStatusWithReplicas(&v1.LocalObjectRef{Name: podTemp4mem100cpu.Name}, 2),
				getTestBufferStatusWithReplicas(&v1.LocalObjectRef{Name: podTemp8mem200cpu.Name}, 3),
			},
			expectedNumberOfErrors: 0,
		},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			translator := NewResourceLimitsTranslator(fakeClient)
			errors := translator.Translate(test.buffers)
			assert.Equal(t, len(errors), test.expectedNumberOfErrors)
			assert.ElementsMatch(t, test.expectedStatus, testutil.SanitizeBuffersStatus(test.buffers))
		})
	}
}

func getTestBufferWithLimits(podTemplateRef *v1.LocalObjectRef, replicas *int32, limits *v1.ResourceList) *v1.CapacityBuffer {
	return testutil.GetBuffer(nil, nil, nil, podTemplateRef, replicas, nil, nil, limits)
}

func getTestBufferStatusWithReplicas(podTemplateRef *v1.LocalObjectRef, replicas int32) *v1.CapacityBufferStatus {
	return testutil.GetBufferStatus(podTemplateRef, &replicas, int64Ptr(1), testutil.GetConditionReady())
}

func getPodTemplateWithResources(name string, resources corev1.ResourceList) *corev1.PodTemplate {
	return &corev1.PodTemplate{
		ObjectMeta: metav1.ObjectMeta{
			Name:       name,
			Namespace:  "default",
			Generation: 1,
		},
		Template: corev1.PodTemplateSpec{
			Spec: corev1.PodSpec{
				Containers: []corev1.Container{
					{
						Resources: corev1.ResourceRequirements{
							Requests: resources,
						},
					},
				},
			},
		},
	}
}

func int64Ptr(i int64) *int64 { return &i }

func int32Ptr(i int32) *int32 { return &i }
