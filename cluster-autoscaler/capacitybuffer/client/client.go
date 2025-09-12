/*
Copyright 2023 The Kubernetes Authors.

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

package client

import (
	"context"
	"fmt"
	"time"

	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"

	"k8s.io/apimachinery/pkg/labels"
	v1 "k8s.io/autoscaler/cluster-autoscaler/apis/capacitybuffer/autoscaling.x-k8s.io/v1"
	"k8s.io/autoscaler/cluster-autoscaler/apis/capacitybuffer/client/clientset/versioned"
	"k8s.io/autoscaler/cluster-autoscaler/apis/capacitybuffer/client/informers/externalversions"
	listers "k8s.io/autoscaler/cluster-autoscaler/apis/capacitybuffer/client/listers/autoscaling.x-k8s.io/v1"
	"k8s.io/client-go/informers"
	"k8s.io/client-go/kubernetes"
	client "k8s.io/client-go/kubernetes"
	corev1 "k8s.io/client-go/listers/core/v1"
	"k8s.io/client-go/rest"

	klog "k8s.io/klog/v2"
)

const (
	clientCallTimeout = 4 * time.Second
)

// CapacityBufferClient represents client for v1 capacitybuffer CRD.
type CapacityBufferClient struct {
	BuffersClient     versioned.Interface
	BuffersLister     listers.CapacityBufferLister
	PodTemplateClient client.Interface
}

// NewCapacityBufferClient configures and returns a capacityBufferClient.
func NewCapacityBufferClient(kubeConfig *rest.Config) (*CapacityBufferClient, error) {
	client, err := newClient(kubeConfig)
	if err != nil {
		return nil, fmt.Errorf("failed to create Provisioning Request client: %v", err)
	}

	buffersLister, err := newBuffersLister(client, make(chan struct{}))
	if err != nil {
		return nil, err
	}

	podTemplateClient, err := kubernetes.NewForConfig(kubeConfig)
	if err != nil {
		return nil, fmt.Errorf("failed to create Pod Template client: %v", err)
	}

	return &CapacityBufferClient{
		BuffersClient:     client,
		BuffersLister:     buffersLister,
		PodTemplateClient: podTemplateClient,
	}, nil
}

// CapacityBuffers gets all CapacityBuffer CRs.
func (c *CapacityBufferClient) CapacityBuffers() ([]*v1.CapacityBuffer, error) {
	buffers, err := c.BuffersLister.List(labels.Everything())
	if err != nil {
		return nil, fmt.Errorf("error fetching capacityBuffers: %w", err)
	}
	return buffers, nil
}

// UpdateCapacityBuffer updates the given CapacityBuffer CR by propagating the changes using the CapacityBufferInterface and returns the updated instance or the original one in case of an error.
func (c *CapacityBufferClient) UpdateCapacityBuffer(buffer *v1.CapacityBuffer) (*v1.CapacityBuffer, error) {
	ctx, cancel := context.WithTimeout(context.Background(), clientCallTimeout)
	defer cancel()

	updatedBuffer, err := c.BuffersClient.AutoscalingV1().CapacityBuffers(buffer.Namespace).UpdateStatus(ctx, buffer, metav1.UpdateOptions{})
	if err != nil {
		return updatedBuffer, err
	}
	klog.V(4).Infof("Updated CapacityBuffer %s/%s,  status: %v,", updatedBuffer.Namespace, updatedBuffer.Name, updatedBuffer.Status)
	return updatedBuffer, nil
}

// newClient creates a new Provisioning Request client from the given config.
func newClient(kubeConfig *rest.Config) (*versioned.Clientset, error) {
	return versioned.NewForConfig(kubeConfig)
}

// newBuffersLister creates a lister for the buffers in the cluster.
func newBuffersLister(client versioned.Interface, stopChannel <-chan struct{}) (listers.CapacityBufferLister, error) {
	factory := externalversions.NewSharedInformerFactory(client, 1*time.Minute)
	buffersLister := factory.Autoscaling().V1().CapacityBuffers().Lister()
	factory.Start(stopChannel)
	informersSynced := factory.WaitForCacheSync(stopChannel)
	for _, synced := range informersSynced {
		if !synced {
			return nil, fmt.Errorf("can't create buffers lister")
		}
	}
	klog.V(2).Info("Successful initial buffers sync")
	return buffersLister, nil
}

// newPodTemplatesLister creates a lister for the Pod Templates in the cluster.
func newPodTemplatesLister(client *kubernetes.Clientset, stopChannel <-chan struct{}) (corev1.PodTemplateLister, error) {
	factory := informers.NewSharedInformerFactory(client, 1*time.Hour)
	podTemplLister := factory.Core().V1().PodTemplates().Lister()
	factory.Start(stopChannel)
	informersSynced := factory.WaitForCacheSync(stopChannel)
	for _, synced := range informersSynced {
		if !synced {
			return nil, fmt.Errorf("can't create Pod Template lister")
		}
	}
	klog.V(2).Info("Successful initial Pod Template sync")
	return podTemplLister, nil
}
