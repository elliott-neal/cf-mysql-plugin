// This file was generated by counterfeiter
package cfmysqlfakes

import (
	"os/exec"
	"sync"

	"github.com/andreasf/cf-mysql-plugin/cfmysql"
)

type FakeExec struct {
	LookPathStub        func(file string) (string, error)
	lookPathMutex       sync.RWMutex
	lookPathArgsForCall []struct {
		file string
	}
	lookPathReturns struct {
		result1 string
		result2 error
	}
	RunStub        func(*exec.Cmd) error
	runMutex       sync.RWMutex
	runArgsForCall []struct {
		arg1 *exec.Cmd
	}
	runReturns struct {
		result1 error
	}
	invocations      map[string][][]interface{}
	invocationsMutex sync.RWMutex
}

func (fake *FakeExec) LookPath(file string) (string, error) {
	fake.lookPathMutex.Lock()
	fake.lookPathArgsForCall = append(fake.lookPathArgsForCall, struct {
		file string
	}{file})
	fake.recordInvocation("LookPath", []interface{}{file})
	fake.lookPathMutex.Unlock()
	if fake.LookPathStub != nil {
		return fake.LookPathStub(file)
	} else {
		return fake.lookPathReturns.result1, fake.lookPathReturns.result2
	}
}

func (fake *FakeExec) LookPathCallCount() int {
	fake.lookPathMutex.RLock()
	defer fake.lookPathMutex.RUnlock()
	return len(fake.lookPathArgsForCall)
}

func (fake *FakeExec) LookPathArgsForCall(i int) string {
	fake.lookPathMutex.RLock()
	defer fake.lookPathMutex.RUnlock()
	return fake.lookPathArgsForCall[i].file
}

func (fake *FakeExec) LookPathReturns(result1 string, result2 error) {
	fake.LookPathStub = nil
	fake.lookPathReturns = struct {
		result1 string
		result2 error
	}{result1, result2}
}

func (fake *FakeExec) Run(arg1 *exec.Cmd) error {
	fake.runMutex.Lock()
	fake.runArgsForCall = append(fake.runArgsForCall, struct {
		arg1 *exec.Cmd
	}{arg1})
	fake.recordInvocation("Run", []interface{}{arg1})
	fake.runMutex.Unlock()
	if fake.RunStub != nil {
		return fake.RunStub(arg1)
	} else {
		return fake.runReturns.result1
	}
}

func (fake *FakeExec) RunCallCount() int {
	fake.runMutex.RLock()
	defer fake.runMutex.RUnlock()
	return len(fake.runArgsForCall)
}

func (fake *FakeExec) RunArgsForCall(i int) *exec.Cmd {
	fake.runMutex.RLock()
	defer fake.runMutex.RUnlock()
	return fake.runArgsForCall[i].arg1
}

func (fake *FakeExec) RunReturns(result1 error) {
	fake.RunStub = nil
	fake.runReturns = struct {
		result1 error
	}{result1}
}

func (fake *FakeExec) Invocations() map[string][][]interface{} {
	fake.invocationsMutex.RLock()
	defer fake.invocationsMutex.RUnlock()
	fake.lookPathMutex.RLock()
	defer fake.lookPathMutex.RUnlock()
	fake.runMutex.RLock()
	defer fake.runMutex.RUnlock()
	return fake.invocations
}

func (fake *FakeExec) recordInvocation(key string, args []interface{}) {
	fake.invocationsMutex.Lock()
	defer fake.invocationsMutex.Unlock()
	if fake.invocations == nil {
		fake.invocations = map[string][][]interface{}{}
	}
	if fake.invocations[key] == nil {
		fake.invocations[key] = [][]interface{}{}
	}
	fake.invocations[key] = append(fake.invocations[key], args)
}

var _ cfmysql.Exec = new(FakeExec)