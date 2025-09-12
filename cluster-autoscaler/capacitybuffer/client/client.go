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

package client

import (
	"context"
	"fmt"
	"time"

	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"

	corev1 "k8s.io/api/core/v1"

	"k8s.io/apimachinery/pkg/labels"
	v1 "k8s.io/autoscaler/cluster-autoscaler/apis/capacitybuffer/autoscaling.x-k8s.io/v1"
	capacitybuffer "k8s.io/autoscaler/cluster-autoscaler/apis/capacitybuffer/client/clientset/versioned"
	"k8s.io/autoscaler/cluster-autoscaler/apis/capacitybuffer/client/informers/externalversions"
	bufferslisters "k8s.io/autoscaler/cluster-autoscaler/apis/capacitybuffer/client/listers/autoscaling.x-k8s.io/v1"
	kubernetes "k8s.io/client-go/kubernetes"
	appsv1listers "k8s.io/client-go/listers/apps/v1"
	batchv1lister "k8s.io/client-go/listers/batch/v1"
	corev1listers "k8s.io/client-go/listers/core/v1"

	appsv1 "k8s.io/api/apps/v1"
	batchv1 "k8s.io/api/batch/v1"
	"k8s.io/client-go/informers"
	"k8s.io/client-go/rest"
	klog "k8s.io/klog/v2"
)

// DefaultNamespace is a const for namespace default
const DefaultNamespace = "default"

// CapacityBufferClient represents client for v1 capacitybuffer CRD.
type CapacityBufferClient struct {
	buffersClient         capacitybuffer.Interface
	kubernetesClient      kubernetes.Interface
	buffersLister         bufferslisters.CapacityBufferLister
	podTemplateLister     corev1listers.PodTemplateLister
	replicaSetsLister     appsv1listers.ReplicaSetLister
	statefulSetsLister    appsv1listers.StatefulSetLister
	jobsLister            batchv1lister.JobLister
	deploymentLister      appsv1listers.DeploymentLister
	replicationContLister corev1listers.ReplicationControllerLister
}

// NewCapacityBufferClient returns a capacityBufferClient.
func NewCapacityBufferClient(buffersClient capacitybuffer.Interface, kubernetesClient kubernetes.Interface, buffersLister bufferslisters.CapacityBufferLister,
	podTemplateLister corev1listers.PodTemplateLister, replicaSetsLister appsv1listers.ReplicaSetLister, statefulSetsLister appsv1listers.StatefulSetLister,
	jobsLister batchv1lister.JobLister, deploymentLister appsv1listers.DeploymentLister, replicationContLister corev1listers.ReplicationControllerLister) (*CapacityBufferClient, error) {
	return &CapacityBufferClient{
		buffersClient:         buffersClient,
		kubernetesClient:      kubernetesClient,
		buffersLister:         buffersLister,
		podTemplateLister:     podTemplateLister,
		replicaSetsLister:     replicaSetsLister,
		statefulSetsLister:    statefulSetsLister,
		jobsLister:            jobsLister,
		deploymentLister:      deploymentLister,
		replicationContLister: replicationContLister,
	}, nil
}

// NewCapacityBufferClientFromConfig configures and returns a CapacityBufferClient.
func NewCapacityBufferClientFromConfig(kubeConfig *rest.Config) (*CapacityBufferClient, error) {
	buffersClient, err := capacitybuffer.NewForConfig(kubeConfig)
	if err != nil {
		return nil, fmt.Errorf("Failed to create capacity buffer client: %v", err)
	}
	kubernetesClient, err := kubernetes.NewForConfig(kubeConfig)
	if err != nil {
		return nil, fmt.Errorf("Failed to create kubernetes client: %v", err)
	}
	return NewCapacityBufferClientFromClients(buffersClient, kubernetesClient)
}

// NewCapacityBufferClientFromClients returns a CapacityBufferClient based on the passed clients
func NewCapacityBufferClientFromClients(buffersClient capacitybuffer.Interface, kubernetesClient kubernetes.Interface) (*CapacityBufferClient, error) {
	stopChannel := make(chan struct{})
	buffersLister, err := newBuffersLister(buffersClient, stopChannel, 5*time.Second)
	if err != nil {
		return nil, err
	}
	factory := informers.NewSharedInformerFactory(kubernetesClient, 1*time.Minute)
	bufferClient := &CapacityBufferClient{
		buffersClient:         buffersClient,
		kubernetesClient:      kubernetesClient,
		buffersLister:         buffersLister,
		podTemplateLister:     factory.Core().V1().PodTemplates().Lister(),
		replicaSetsLister:     factory.Apps().V1().ReplicaSets().Lister(),
		statefulSetsLister:    factory.Apps().V1().StatefulSets().Lister(),
		jobsLister:            factory.Batch().V1().Jobs().Lister(),
		deploymentLister:      factory.Apps().V1().Deployments().Lister(),
		replicationContLister: factory.Core().V1().ReplicationControllers().Lister(),
	}
	factory.Start(stopChannel)
	informersSynced := factory.WaitForCacheSync(stopChannel)
	for _, synced := range informersSynced {
		if !synced {
			return nil, fmt.Errorf("Can't initiate informer factory syncer")
		}
	}
	return bufferClient, nil
}

// newBuffersLister creates a lister for the buffers in the cluster.
func newBuffersLister(client capacitybuffer.Interface, stopChannel <-chan struct{}, defaultResync time.Duration) (bufferslisters.CapacityBufferLister, error) {
	factory := externalversions.NewSharedInformerFactory(client, defaultResync)
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

// ListCapacityBuffers lists all Capacity buffer.
func (c *CapacityBufferClient) ListCapacityBuffers() ([]*v1.CapacityBuffer, error) {
	buffers, err := c.buffersLister.List(labels.Everything())
	if err != nil {
		return nil, fmt.Errorf("Error fetching capacity buffers: %v", err)
	}
	buffersCopy := []*v1.CapacityBuffer{}
	for _, buffer := range buffers {
		buffersCopy = append(buffersCopy, buffer.DeepCopy())
	}
	return buffersCopy, nil
}

// GetPodTemplate returns pod template with the passed name
func (c *CapacityBufferClient) GetPodTemplate(name string) (*corev1.PodTemplate, error) {
	if c.podTemplateLister != nil {
		template, err := c.podTemplateLister.PodTemplates(DefaultNamespace).Get(name)
		if err == nil {
			return template.DeepCopy(), nil
		}
	}
	if c.kubernetesClient != nil {
		template, err := c.kubernetesClient.CoreV1().PodTemplates(DefaultNamespace).Get(context.TODO(), name, metav1.GetOptions{})
		if err == nil {
			return template.DeepCopy(), nil
		}
	}
	return nil, fmt.Errorf("Capacity buffer client can't get pod template, clients are not configured correctly")
}

// UpdateCapacityBuffer fetches the cached object using a lister object in the client
func (c *CapacityBufferClient) UpdateCapacityBuffer(buffer *v1.CapacityBuffer) (*v1.CapacityBuffer, error) {
	buffer, err := c.buffersClient.AutoscalingV1().CapacityBuffers(DefaultNamespace).UpdateStatus(context.TODO(), buffer, metav1.UpdateOptions{})
	if err == nil {
		return buffer.DeepCopy(), nil
	}
	return nil, err
}

// CreatePodTemplate creates a pod template
func (c *CapacityBufferClient) CreatePodTemplate(podTemplate *corev1.PodTemplate) (*corev1.PodTemplate, error) {
	template, err := c.kubernetesClient.CoreV1().PodTemplates(DefaultNamespace).Create(context.TODO(), podTemplate, metav1.CreateOptions{})
	if err == nil {
		return template.DeepCopy(), nil
	}
	return nil, err
}

// UpdatePodTemplate updates the pod template
func (c *CapacityBufferClient) UpdatePodTemplate(podTemplate *corev1.PodTemplate) (*corev1.PodTemplate, error) {
	template, err := c.kubernetesClient.CoreV1().PodTemplates(DefaultNamespace).Update(context.TODO(), podTemplate, metav1.UpdateOptions{})
	if err == nil {
		return template.DeepCopy(), nil
	}
	return nil, err
}

// GetDeployment fetches the cached object using a lister object in the client
func (c *CapacityBufferClient) GetDeployment(namespace, name string) (*appsv1.Deployment, error) {
	obj, err := c.deploymentLister.Deployments(namespace).Get(name)
	if err != nil {
		return nil, fmt.Errorf("failed to get Deployment: %w", err)
	}
	return obj, nil
}

// GetReplicaSet fetches the cached object using a lister object in the client
func (c *CapacityBufferClient) GetReplicaSet(namespace, name string) (*appsv1.ReplicaSet, error) {
	obj, err := c.replicaSetsLister.ReplicaSets(namespace).Get(name)
	if err != nil {
		return nil, fmt.Errorf("failed to get ReplicaSet: %w", err)
	}
	return obj, nil
}

// GetJob fetches the cached object using a lister object in the client
func (c *CapacityBufferClient) GetJob(namespace, name string) (*batchv1.Job, error) {
	obj, err := c.jobsLister.Jobs(namespace).Get(name)
	if err != nil {
		return nil, fmt.Errorf("failed to get Job: %w", err)
	}
	return obj, nil
}

// GetReplicationController fetches the cached object using a lister object in the client
func (c *CapacityBufferClient) GetReplicationController(namespace, name string) (*corev1.ReplicationController, error) {
	obj, err := c.replicationContLister.ReplicationControllers(namespace).Get(name)
	if err != nil {
		return nil, fmt.Errorf("failed to get ReplicationController: %w", err)
	}
	return obj, nil
}

// GetStatefulSet fetches the cached object using a lister object in the client
func (c *CapacityBufferClient) GetStatefulSet(namespace, name string) (*appsv1.StatefulSet, error) {
	obj, err := c.statefulSetsLister.StatefulSets(namespace).Get(name)
	if err != nil {
		return nil, fmt.Errorf("failed to get StatefulSet: %w", err)
	}
	return obj, nil
}
