// Code generated by go-mockgen 1.3.3; DO NOT EDIT.
//
// This file was generated by running `sg generate` (or `go-mockgen`) at the root of
// this repository. To add additional mocks to this or another package, add a new entry
// to the mockgen.yaml file in the root of this repository.

package app

import (
	"context"
	"sync"

	github "github.com/google/go-github/v41/github"
)

// MockGithubClient is a mock implementation of the githubClient interface
// (from the package
// github.com/sourcegraph/sourcegraph/enterprise/cmd/frontend/internal/app)
// used for unit testing.
type MockGithubClient struct {
	// GetAppInstallationFunc is an instance of a mock function object
	// controlling the behavior of the method GetAppInstallation.
	GetAppInstallationFunc *GithubClientGetAppInstallationFunc
}

// NewMockGithubClient creates a new mock of the githubClient interface. All
// methods return zero values for all results, unless overwritten.
func NewMockGithubClient() *MockGithubClient {
	return &MockGithubClient{
		GetAppInstallationFunc: &GithubClientGetAppInstallationFunc{
			defaultHook: func(context.Context, int64) (r0 *github.Installation, r1 error) {
				return
			},
		},
	}
}

// NewStrictMockGithubClient creates a new mock of the githubClient
// interface. All methods panic on invocation, unless overwritten.
func NewStrictMockGithubClient() *MockGithubClient {
	return &MockGithubClient{
		GetAppInstallationFunc: &GithubClientGetAppInstallationFunc{
			defaultHook: func(context.Context, int64) (*github.Installation, error) {
				panic("unexpected invocation of MockGithubClient.GetAppInstallation")
			},
		},
	}
}

// surrogateMockGithubClient is a copy of the githubClient interface (from
// the package
// github.com/sourcegraph/sourcegraph/enterprise/cmd/frontend/internal/app).
// It is redefined here as it is unexported in the source package.
type surrogateMockGithubClient interface {
	GetAppInstallation(context.Context, int64) (*github.Installation, error)
}

// NewMockGithubClientFrom creates a new mock of the MockGithubClient
// interface. All methods delegate to the given implementation, unless
// overwritten.
func NewMockGithubClientFrom(i surrogateMockGithubClient) *MockGithubClient {
	return &MockGithubClient{
		GetAppInstallationFunc: &GithubClientGetAppInstallationFunc{
			defaultHook: i.GetAppInstallation,
		},
	}
}

// GithubClientGetAppInstallationFunc describes the behavior when the
// GetAppInstallation method of the parent MockGithubClient instance is
// invoked.
type GithubClientGetAppInstallationFunc struct {
	defaultHook func(context.Context, int64) (*github.Installation, error)
	hooks       []func(context.Context, int64) (*github.Installation, error)
	history     []GithubClientGetAppInstallationFuncCall
	mutex       sync.Mutex
}

// GetAppInstallation delegates to the next hook function in the queue and
// stores the parameter and result values of this invocation.
func (m *MockGithubClient) GetAppInstallation(v0 context.Context, v1 int64) (*github.Installation, error) {
	r0, r1 := m.GetAppInstallationFunc.nextHook()(v0, v1)
	m.GetAppInstallationFunc.appendCall(GithubClientGetAppInstallationFuncCall{v0, v1, r0, r1})
	return r0, r1
}

// SetDefaultHook sets function that is called when the GetAppInstallation
// method of the parent MockGithubClient instance is invoked and the hook
// queue is empty.
func (f *GithubClientGetAppInstallationFunc) SetDefaultHook(hook func(context.Context, int64) (*github.Installation, error)) {
	f.defaultHook = hook
}

// PushHook adds a function to the end of hook queue. Each invocation of the
// GetAppInstallation method of the parent MockGithubClient instance invokes
// the hook at the front of the queue and discards it. After the queue is
// empty, the default hook function is invoked for any future action.
func (f *GithubClientGetAppInstallationFunc) PushHook(hook func(context.Context, int64) (*github.Installation, error)) {
	f.mutex.Lock()
	f.hooks = append(f.hooks, hook)
	f.mutex.Unlock()
}

// SetDefaultReturn calls SetDefaultHook with a function that returns the
// given values.
func (f *GithubClientGetAppInstallationFunc) SetDefaultReturn(r0 *github.Installation, r1 error) {
	f.SetDefaultHook(func(context.Context, int64) (*github.Installation, error) {
		return r0, r1
	})
}

// PushReturn calls PushHook with a function that returns the given values.
func (f *GithubClientGetAppInstallationFunc) PushReturn(r0 *github.Installation, r1 error) {
	f.PushHook(func(context.Context, int64) (*github.Installation, error) {
		return r0, r1
	})
}

func (f *GithubClientGetAppInstallationFunc) nextHook() func(context.Context, int64) (*github.Installation, error) {
	f.mutex.Lock()
	defer f.mutex.Unlock()

	if len(f.hooks) == 0 {
		return f.defaultHook
	}

	hook := f.hooks[0]
	f.hooks = f.hooks[1:]
	return hook
}

func (f *GithubClientGetAppInstallationFunc) appendCall(r0 GithubClientGetAppInstallationFuncCall) {
	f.mutex.Lock()
	f.history = append(f.history, r0)
	f.mutex.Unlock()
}

// History returns a sequence of GithubClientGetAppInstallationFuncCall
// objects describing the invocations of this function.
func (f *GithubClientGetAppInstallationFunc) History() []GithubClientGetAppInstallationFuncCall {
	f.mutex.Lock()
	history := make([]GithubClientGetAppInstallationFuncCall, len(f.history))
	copy(history, f.history)
	f.mutex.Unlock()

	return history
}

// GithubClientGetAppInstallationFuncCall is an object that describes an
// invocation of method GetAppInstallation on an instance of
// MockGithubClient.
type GithubClientGetAppInstallationFuncCall struct {
	// Arg0 is the value of the 1st argument passed to this method
	// invocation.
	Arg0 context.Context
	// Arg1 is the value of the 2nd argument passed to this method
	// invocation.
	Arg1 int64
	// Result0 is the value of the 1st result returned from this method
	// invocation.
	Result0 *github.Installation
	// Result1 is the value of the 2nd result returned from this method
	// invocation.
	Result1 error
}

// Args returns an interface slice containing the arguments of this
// invocation.
func (c GithubClientGetAppInstallationFuncCall) Args() []interface{} {
	return []interface{}{c.Arg0, c.Arg1}
}

// Results returns an interface slice containing the results of this
// invocation.
func (c GithubClientGetAppInstallationFuncCall) Results() []interface{} {
	return []interface{}{c.Result0, c.Result1}
}
