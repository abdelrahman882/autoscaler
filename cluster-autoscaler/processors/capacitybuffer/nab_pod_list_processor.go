package capacitybufferpodlister

import (
	"fmt"
	"time"

	"github.com/google/uuid"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/labels"
	"k8s.io/apimachinery/pkg/types"
	"k8s.io/klog/v2"

	apiv1 "k8s.io/api/core/v1"
	api_v1 "k8s.io/autoscaler/cluster-autoscaler/apis/capacitybuffer/autoscaling.x-k8s.io/v1"
	client "k8s.io/autoscaler/cluster-autoscaler/capacitybuffer/client"
	"k8s.io/autoscaler/cluster-autoscaler/capacitybuffer/common"
	buffersfilter "k8s.io/autoscaler/cluster-autoscaler/capacitybuffer/filters"
	"k8s.io/autoscaler/cluster-autoscaler/context"
	"k8s.io/autoscaler/cluster-autoscaler/simulator/fake"
)

// CapacityBufferPodListProcessor processes the pod lists before scale up
// and adds buffres api virtual pods.
type CapacityBufferPodListProcessor struct {
	client            *client.CapacityBufferClient
	toProvisionFilter buffersfilter.Filter
}

// NewCapacityRequestPodListProcessor creates a new CapacityRequestPodListProcessor.
func NewCapacityBufferPodListProcessor(client *client.CapacityBufferClient) *CapacityBufferPodListProcessor {

	return &CapacityBufferPodListProcessor{
		client: client,
		toProvisionFilter: buffersfilter.NewStatusFilter(map[string]string{
			common.ReadyForProvisioningCondition: common.ConditionTrue,
			common.ProvisioningCondition:         common.ConditionTrue,
		}),
	}
}

// Process updates unschedulablePods by injecting fake pods to match target replica count
func (p *CapacityBufferPodListProcessor) Process(ctx *context.AutoscalingContext, unschedulablePods []*apiv1.Pod) ([]*apiv1.Pod, error) {
	buffers, err := p.client.CapacityBuffers()
	if err != nil {
		klog.Errorf("CapacityBufferPodListProcessor failed to list buffers with error: %v", err.Error())
		return unschedulablePods, nil
	}
	_, filteredOutBuffers := p.toProvisionFilter.Filter(buffers)

	totalFakePods := []*apiv1.Pod{}
	for _, buffer := range filteredOutBuffers {
		fakePods := p.provision(buffer)
		totalFakePods = append(totalFakePods, fakePods...)
	}

	unschedulablePods = append(unschedulablePods, totalFakePods...)
	return unschedulablePods, nil
}

func (p *CapacityBufferPodListProcessor) provision(buffer *api_v1.CapacityBuffer) []*apiv1.Pod {
	if buffer.Status.PodTemplateRef == nil || buffer.Status.Replicas == nil {
		return []*apiv1.Pod{}
	}

	podTemplateName := buffer.Status.PodTemplateRef.Name
	replicas := buffer.Status.Replicas
	podTemplate, err := common.GetPodTemplate(p.client.PodTemplateClient, podTemplateName)

	if err != nil {
		p.updateBufferStatusToFailedProvisioing(buffer, fmt.Sprintf("failed to get pod template with error: %v", err.Error()))
		return []*apiv1.Pod{}
	}

	fakePods, err := makeFakePods(buffer.UID, &podTemplate.Template, int(*replicas))
	if err != nil {
		p.updateBufferStatusToFailedProvisioing(buffer, fmt.Sprintf("failed to get pod template with error: %v", err.Error()))
		return []*apiv1.Pod{}
	}

	p.updateBufferStatusToSuccessfullyProvisioing(buffer)
	return fakePods
}

// CleanUp is called at CA termination
func (p *CapacityBufferPodListProcessor) CleanUp() {
}

func (p *CapacityBufferPodListProcessor) updateBufferStatusToFailedProvisioing(buffer *api_v1.CapacityBuffer, errorMessage string) {
	buffer.Status.Conditions = []metav1.Condition{metav1.Condition{
		Type:               common.ProvisioningCondition,
		Status:             common.ConditionFalse,
		Message:            errorMessage,
		Reason:             "error",
		LastTransitionTime: metav1.Time{Time: time.Now()},
	}}
	err := common.UpdateBufferStatus(p.client.BuffersClient, buffer)
	if err != nil {
		klog.Errorf("Failed to update buffer status for buffer %v", buffer.Name)
	}
}

func (p *CapacityBufferPodListProcessor) updateBufferStatusToSuccessfullyProvisioing(buffer *api_v1.CapacityBuffer) {
	buffer.Status.Conditions = []metav1.Condition{metav1.Condition{
		Type:               common.ProvisioningCondition,
		Status:             common.ConditionTrue,
		Reason:             "provisioning",
		LastTransitionTime: metav1.Time{Time: time.Now()},
	}}
	err := common.UpdateBufferStatus(p.client.BuffersClient, buffer)
	if err != nil {
		klog.Errorf("Failed to update buffer status for buffer %v", buffer.Name)
	}
}

// makeFakePods creates podCount number of copies of the sample pod
func makeFakePods(ownerUid types.UID, samplePodTemplate *apiv1.PodTemplateSpec, podCount int) ([]*apiv1.Pod, error) {
	var fakePods []*apiv1.Pod
	samplePod := getPodFromTemplate(samplePodTemplate)

	for i := 1; i <= podCount; i++ {
		newPod := fake.WithFakePodAnnotation(samplePod)
		newPod.Name = fmt.Sprintf("NAB-INJECTED-POD-%s-%d", ownerUid, i)
		newPod.UID = types.UID(fmt.Sprintf("%s-%d", string(ownerUid), i))
		newPod.Spec.NodeName = ""
		fakePods = append(fakePods, newPod)
	}
	return fakePods, nil
}

func getPodFromTemplate(template *apiv1.PodTemplateSpec) *apiv1.Pod {
	desiredLabels := getPodsLabelSet(template)
	desiredFinalizers := getPodsFinalizers(template)
	desiredAnnotations := getPodsAnnotationSet(template)

	pod := &apiv1.Pod{
		ObjectMeta: metav1.ObjectMeta{
			Labels:       desiredLabels,
			Namespace:    template.Namespace,
			Annotations:  desiredAnnotations,
			GenerateName: uuid.NewString(),
			Finalizers:   desiredFinalizers,
		},
	}

	pod.Spec = template.Spec
	return pod
}

func getPodsLabelSet(template *apiv1.PodTemplateSpec) labels.Set {
	desiredLabels := make(labels.Set)
	for k, v := range template.Labels {
		desiredLabels[k] = v
	}
	return desiredLabels
}

func getPodsFinalizers(template *apiv1.PodTemplateSpec) []string {
	desiredFinalizers := make([]string, len(template.Finalizers))
	copy(desiredFinalizers, template.Finalizers)
	return desiredFinalizers
}

func getPodsAnnotationSet(template *apiv1.PodTemplateSpec) labels.Set {
	desiredAnnotations := make(labels.Set)
	for k, v := range template.Annotations {
		desiredAnnotations[k] = v
	}
	return desiredAnnotations
}
