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

package filter

import (
	v1 "k8s.io/autoscaler/cluster-autoscaler/apis/capacitybuffer/autoscaling.x-k8s.io/v1"
)

// bufferGenerationFilter filters out buffers with the defined conditions
type bufferGenerationFilter struct {
	buffersGenerations map[string]int64
}

// NewBufferGenerationFilter creates an instance of bufferGenerationFilter that filters the buffers that needs to be updated.
func NewBufferGenerationFilter() *bufferGenerationFilter {
	return &bufferGenerationFilter{
		buffersGenerations: map[string]int64{},
	}
}

// Filter filters the passed buffers based on buffer status conditions
func (f *bufferGenerationFilter) Filter(buffersToFilter []*v1.CapacityBuffer) ([]*v1.CapacityBuffer, []*v1.CapacityBuffer) {
	var buffers []*v1.CapacityBuffer
	var filteredOutBuffers []*v1.CapacityBuffer

	for _, buffer := range buffersToFilter {
		if f.generationChanged(buffer) {
			buffers = append(buffers, buffer)
		} else {
			filteredOutBuffers = append(filteredOutBuffers, buffer)
		}
	}
	return buffers, filteredOutBuffers
}

func (f *bufferGenerationFilter) generationChanged(buffer *v1.CapacityBuffer) bool {
	newGeneration := buffer.Generation
	oldGeneration, found := f.buffersGenerations[buffer.Name]
	f.buffersGenerations[buffer.Name] = newGeneration

	if found && oldGeneration == newGeneration {
		return false
	}

	// we will return true -process the buffer- in cases of a) not found (new buffer or CA restarted) b) found or different generation
	return true
}

// CleanUp cleans up the filter's internal structures.
func (f *bufferGenerationFilter) CleanUp() {
}
