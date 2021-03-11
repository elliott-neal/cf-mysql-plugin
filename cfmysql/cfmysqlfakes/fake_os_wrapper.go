// Code generated by counterfeiter. DO NOT EDIT.
package cfmysqlfakes

import (
	"os"
	"sync"

	"github.com/elliott-neal/cf-mysql-plugin/cfmysql"
)

type FakeOsWrapper struct {
	LookupEnvStub        func(key string) (string, bool)
	lookupEnvMutex       sync.RWMutex
	lookupEnvArgsForCall []struct {
		key string
	}
	lookupEnvReturns struct {
		result1 string
		result2 bool
	}
	lookupEnvReturnsOnCall map[int]struct {
		result1 string
		result2 bool
	}
	NameStub        func(file *os.File) string
	nameMutex       sync.RWMutex
	nameArgsForCall []struct {
		file *os.File
	}
	nameReturns struct {
		result1 string
	}
	nameReturnsOnCall map[int]struct {
		result1 string
	}
	RemoveStub        func(name string) error
	removeMutex       sync.RWMutex
	removeArgsForCall []struct {
		name string
	}
	removeReturns struct {
		result1 error
	}
	removeReturnsOnCall map[int]struct {
		result1 error
	}
	WriteStringStub        func(file *os.File, s string) (n int, err error)
	writeStringMutex       sync.RWMutex
	writeStringArgsForCall []struct {
		file *os.File
		s    string
	}
	writeStringReturns struct {
		result1 int
		result2 error
	}
	writeStringReturnsOnCall map[int]struct {
		result1 int
		result2 error
	}
	invocations      map[string][][]interface{}
	invocationsMutex sync.RWMutex
}

func (fake *FakeOsWrapper) LookupEnv(key string) (string, bool) {
	fake.lookupEnvMutex.Lock()
	ret, specificReturn := fake.lookupEnvReturnsOnCall[len(fake.lookupEnvArgsForCall)]
	fake.lookupEnvArgsForCall = append(fake.lookupEnvArgsForCall, struct {
		key string
	}{key})
	fake.recordInvocation("LookupEnv", []interface{}{key})
	fake.lookupEnvMutex.Unlock()
	if fake.LookupEnvStub != nil {
		return fake.LookupEnvStub(key)
	}
	if specificReturn {
		return ret.result1, ret.result2
	}
	return fake.lookupEnvReturns.result1, fake.lookupEnvReturns.result2
}

func (fake *FakeOsWrapper) LookupEnvCallCount() int {
	fake.lookupEnvMutex.RLock()
	defer fake.lookupEnvMutex.RUnlock()
	return len(fake.lookupEnvArgsForCall)
}

func (fake *FakeOsWrapper) LookupEnvArgsForCall(i int) string {
	fake.lookupEnvMutex.RLock()
	defer fake.lookupEnvMutex.RUnlock()
	return fake.lookupEnvArgsForCall[i].key
}

func (fake *FakeOsWrapper) LookupEnvReturns(result1 string, result2 bool) {
	fake.LookupEnvStub = nil
	fake.lookupEnvReturns = struct {
		result1 string
		result2 bool
	}{result1, result2}
}

func (fake *FakeOsWrapper) LookupEnvReturnsOnCall(i int, result1 string, result2 bool) {
	fake.LookupEnvStub = nil
	if fake.lookupEnvReturnsOnCall == nil {
		fake.lookupEnvReturnsOnCall = make(map[int]struct {
			result1 string
			result2 bool
		})
	}
	fake.lookupEnvReturnsOnCall[i] = struct {
		result1 string
		result2 bool
	}{result1, result2}
}

func (fake *FakeOsWrapper) Name(file *os.File) string {
	fake.nameMutex.Lock()
	ret, specificReturn := fake.nameReturnsOnCall[len(fake.nameArgsForCall)]
	fake.nameArgsForCall = append(fake.nameArgsForCall, struct {
		file *os.File
	}{file})
	fake.recordInvocation("Name", []interface{}{file})
	fake.nameMutex.Unlock()
	if fake.NameStub != nil {
		return fake.NameStub(file)
	}
	if specificReturn {
		return ret.result1
	}
	return fake.nameReturns.result1
}

func (fake *FakeOsWrapper) NameCallCount() int {
	fake.nameMutex.RLock()
	defer fake.nameMutex.RUnlock()
	return len(fake.nameArgsForCall)
}

func (fake *FakeOsWrapper) NameArgsForCall(i int) *os.File {
	fake.nameMutex.RLock()
	defer fake.nameMutex.RUnlock()
	return fake.nameArgsForCall[i].file
}

func (fake *FakeOsWrapper) NameReturns(result1 string) {
	fake.NameStub = nil
	fake.nameReturns = struct {
		result1 string
	}{result1}
}

func (fake *FakeOsWrapper) NameReturnsOnCall(i int, result1 string) {
	fake.NameStub = nil
	if fake.nameReturnsOnCall == nil {
		fake.nameReturnsOnCall = make(map[int]struct {
			result1 string
		})
	}
	fake.nameReturnsOnCall[i] = struct {
		result1 string
	}{result1}
}

func (fake *FakeOsWrapper) Remove(name string) error {
	fake.removeMutex.Lock()
	ret, specificReturn := fake.removeReturnsOnCall[len(fake.removeArgsForCall)]
	fake.removeArgsForCall = append(fake.removeArgsForCall, struct {
		name string
	}{name})
	fake.recordInvocation("Remove", []interface{}{name})
	fake.removeMutex.Unlock()
	if fake.RemoveStub != nil {
		return fake.RemoveStub(name)
	}
	if specificReturn {
		return ret.result1
	}
	return fake.removeReturns.result1
}

func (fake *FakeOsWrapper) RemoveCallCount() int {
	fake.removeMutex.RLock()
	defer fake.removeMutex.RUnlock()
	return len(fake.removeArgsForCall)
}

func (fake *FakeOsWrapper) RemoveArgsForCall(i int) string {
	fake.removeMutex.RLock()
	defer fake.removeMutex.RUnlock()
	return fake.removeArgsForCall[i].name
}

func (fake *FakeOsWrapper) RemoveReturns(result1 error) {
	fake.RemoveStub = nil
	fake.removeReturns = struct {
		result1 error
	}{result1}
}

func (fake *FakeOsWrapper) RemoveReturnsOnCall(i int, result1 error) {
	fake.RemoveStub = nil
	if fake.removeReturnsOnCall == nil {
		fake.removeReturnsOnCall = make(map[int]struct {
			result1 error
		})
	}
	fake.removeReturnsOnCall[i] = struct {
		result1 error
	}{result1}
}

func (fake *FakeOsWrapper) WriteString(file *os.File, s string) (n int, err error) {
	fake.writeStringMutex.Lock()
	ret, specificReturn := fake.writeStringReturnsOnCall[len(fake.writeStringArgsForCall)]
	fake.writeStringArgsForCall = append(fake.writeStringArgsForCall, struct {
		file *os.File
		s    string
	}{file, s})
	fake.recordInvocation("WriteString", []interface{}{file, s})
	fake.writeStringMutex.Unlock()
	if fake.WriteStringStub != nil {
		return fake.WriteStringStub(file, s)
	}
	if specificReturn {
		return ret.result1, ret.result2
	}
	return fake.writeStringReturns.result1, fake.writeStringReturns.result2
}

func (fake *FakeOsWrapper) WriteStringCallCount() int {
	fake.writeStringMutex.RLock()
	defer fake.writeStringMutex.RUnlock()
	return len(fake.writeStringArgsForCall)
}

func (fake *FakeOsWrapper) WriteStringArgsForCall(i int) (*os.File, string) {
	fake.writeStringMutex.RLock()
	defer fake.writeStringMutex.RUnlock()
	return fake.writeStringArgsForCall[i].file, fake.writeStringArgsForCall[i].s
}

func (fake *FakeOsWrapper) WriteStringReturns(result1 int, result2 error) {
	fake.WriteStringStub = nil
	fake.writeStringReturns = struct {
		result1 int
		result2 error
	}{result1, result2}
}

func (fake *FakeOsWrapper) WriteStringReturnsOnCall(i int, result1 int, result2 error) {
	fake.WriteStringStub = nil
	if fake.writeStringReturnsOnCall == nil {
		fake.writeStringReturnsOnCall = make(map[int]struct {
			result1 int
			result2 error
		})
	}
	fake.writeStringReturnsOnCall[i] = struct {
		result1 int
		result2 error
	}{result1, result2}
}

func (fake *FakeOsWrapper) Invocations() map[string][][]interface{} {
	fake.invocationsMutex.RLock()
	defer fake.invocationsMutex.RUnlock()
	fake.lookupEnvMutex.RLock()
	defer fake.lookupEnvMutex.RUnlock()
	fake.nameMutex.RLock()
	defer fake.nameMutex.RUnlock()
	fake.removeMutex.RLock()
	defer fake.removeMutex.RUnlock()
	fake.writeStringMutex.RLock()
	defer fake.writeStringMutex.RUnlock()
	copiedInvocations := map[string][][]interface{}{}
	for key, value := range fake.invocations {
		copiedInvocations[key] = value
	}
	return copiedInvocations
}

func (fake *FakeOsWrapper) recordInvocation(key string, args []interface{}) {
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

var _ cfmysql.OsWrapper = new(FakeOsWrapper)
