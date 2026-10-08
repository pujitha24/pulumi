// Copyright 2026, Pulumi Corporation.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package cmdutil

import (
	"context"
	"errors"
	"io"
	"testing"

	"github.com/pulumi/pulumi/pkg/v3/resource/plugin"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

type failingInstallRuntime struct {
	plugin.LanguageRuntime

	err error
}

func (r failingInstallRuntime) InstallDependencies(
	context.Context, plugin.InstallDependenciesRequest,
) (io.Reader, io.Reader, <-chan error, error) {
	outr, outw := io.Pipe()
	errr, errw := io.Pipe()
	done := make(chan error, 1)
	// Mirror the language host client, which closes both pipes with the same error it reports on done.
	_ = outw.CloseWithError(r.err)
	_ = errw.CloseWithError(r.err)
	done <- r.err
	close(done)
	return outr, errr, done, nil
}

func TestInstallDependenciesReportsErrorOnce(t *testing.T) {
	t.Parallel()

	want := errors.New("unable to find program: dotnet")
	err := InstallDependencies(
		t.Context(), failingInstallRuntime{err: want}, plugin.InstallDependenciesRequest{}, io.Discard, io.Discard)
	require.Error(t, err)
	assert.Equal(t, want.Error(), err.Error())
}
