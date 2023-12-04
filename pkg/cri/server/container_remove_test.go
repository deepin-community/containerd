/*
   Copyright The containerd Authors.

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

package server

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"

	containerstore "github.com/containerd/containerd/pkg/cri/store/container"
)

// TestSetContainerRemoving tests setContainerRemoving sets removing
// state correctly.
func TestSetContainerRemoving(t *testing.T) {
	testID := "test-id"
	for desc, test := range map[string]struct {
		status    containerstore.Status
		expectErr bool
	}{
		"should return error when container is in running state": {
			status: containerstore.Status{
				CreatedAt: time.Now().UnixNano(),
				StartedAt: time.Now().UnixNano(),
			},
			expectErr: true,
		},
		"should return error when container is in starting state": {
			status: containerstore.Status{
				CreatedAt: time.Now().UnixNano(),
				Starting:  true,
			},
			expectErr: true,
		},
		"should return error when container is in removing state": {
			status: containerstore.Status{
				CreatedAt:  time.Now().UnixNano(),
				StartedAt:  time.Now().UnixNano(),
				FinishedAt: time.Now().UnixNano(),
				Removing:   true,
			},
			expectErr: true,
		},
		"should not return error when container is not running and removing": {
			status: containerstore.Status{
				CreatedAt:  time.Now().UnixNano(),
				StartedAt:  time.Now().UnixNano(),
				FinishedAt: time.Now().UnixNano(),
			},
			expectErr: false,
		},
	} {
		t.Logf("TestCase %q", desc)
		container, err := containerstore.NewContainer(
			containerstore.Metadata{ID: testID},
			containerstore.WithFakeStatus(test.status),
		)
		assert.NoError(t, err)
		err = setContainerRemoving(container)
		if test.expectErr {
			assert.Error(t, err)
			assert.Equal(t, test.status, container.Status.Get(), "metadata should not be updated")
		} else {
			assert.NoError(t, err)
			assert.True(t, container.Status.Get().Removing, "removing should be set")
			assert.NoError(t, resetContainerRemoving(container))
			assert.False(t, container.Status.Get().Removing, "removing should be reset")
		}
	}
}
