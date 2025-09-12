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

package scalableobject

import (
	"fmt"

	corev1 "k8s.io/api/core/v1"
	cbclient "k8s.io/autoscaler/cluster-autoscaler/capacitybuffer/client"
)

// Kinds of the supported objects
const (
	DeploymentKind            = "Deployment"
	ReplicaSetKind            = "ReplicaSet"
	StatefulSetKind           = "StatefulSet"
	ReplicationControllerKind = "ReplicationController"
	JobKind                   = "Job"
)

// ScalableObjectResolver resolves a scalable object to pod spec and replicas
type ScalableObjectResolver interface {
	GetTemplateAndReplicas(namespace, name string) (*corev1.PodTemplateSpec, *int32, error)
	GetKind() string
}

// GetSupportedScalableObjectResolvers returns the default ScalableObjectResolvers
func GetSupportedScalableObjectResolvers(client *cbclient.CapacityBufferClient) []ScalableObjectResolver {
	return []ScalableObjectResolver{
		&deployment{scalable{client: client, kind: DeploymentKind}},
		&replicaSet{scalable{client: client, kind: ReplicaSetKind}},
		&statefulSet{scalable{client: client, kind: StatefulSetKind}},
		&replicationController{scalable{client: client, kind: ReplicationControllerKind}},
		&job{scalable{client: client, kind: JobKind}},
	}
}

type scalable struct {
	client *cbclient.CapacityBufferClient
	kind   string
}

// GetKind returns the kind of the sclable object
func (s *scalable) GetKind() string {
	return s.kind
}

type deployment struct{ scalable }

// GetTemplateAndReplicas returns the pod spec template of the passed object name and namespace
func (s *deployment) GetTemplateAndReplicas(namespace, name string) (*corev1.PodTemplateSpec, *int32, error) {
	obj, err := s.client.GetDeployment(namespace, name)
	if err != nil {
		return nil, nil, fmt.Errorf("failed to get Deployment: %w", err)
	}
	return obj.Spec.Template.DeepCopy(), obj.Spec.Replicas, nil
}

type replicaSet struct{ scalable }

// GetTemplateAndReplicas returns the pod spec template of the passed object name and namespace
func (s *replicaSet) GetTemplateAndReplicas(namespace, name string) (*corev1.PodTemplateSpec, *int32, error) {
	obj, err := s.client.GetReplicaSet(namespace, name)
	if err != nil {
		return nil, nil, fmt.Errorf("failed to get Deployment: %w", err)
	}
	return obj.Spec.Template.DeepCopy(), obj.Spec.Replicas, nil
}

type statefulSet struct{ scalable }

// GetTemplateAndReplicas returns the pod spec template of the passed object name and namespace
func (s *statefulSet) GetTemplateAndReplicas(namespace, name string) (*corev1.PodTemplateSpec, *int32, error) {
	obj, err := s.client.GetStatefulSet(namespace, name)
	if err != nil {
		return nil, nil, fmt.Errorf("failed to get Deployment: %w", err)
	}
	return obj.Spec.Template.DeepCopy(), obj.Spec.Replicas, nil
}

type replicationController struct{ scalable }

// GetTemplateAndReplicas returns the pod spec template of the passed object name and namespace
func (s *replicationController) GetTemplateAndReplicas(namespace, name string) (*corev1.PodTemplateSpec, *int32, error) {
	obj, err := s.client.GetReplicationController(namespace, name)
	if err != nil {
		return nil, nil, fmt.Errorf("failed to get Deployment: %w", err)
	}
	return obj.Spec.Template.DeepCopy(), obj.Spec.Replicas, nil
}

type job struct{ scalable }

// GetTemplateAndReplicas returns the pod spec template of the passed object name and namespace
func (s *job) GetTemplateAndReplicas(namespace, name string) (*corev1.PodTemplateSpec, *int32, error) {
	obj, err := s.client.GetJob(namespace, name)
	if err != nil {
		return nil, nil, fmt.Errorf("failed to get Deployment: %w", err)
	}
	return obj.Spec.Template.DeepCopy(), obj.Spec.Parallelism, nil
}
