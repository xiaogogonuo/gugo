package gugo

import (
	"sync/atomic"
)

type module struct {
	calledCount    uint64 // 代表请求调用的计数
	failedCount    uint64 // 代表请求失败的计数
	acceptedCount  uint64 // 代表请求被接受的计数
	interceptCount uint64 // 代表请求被拦截的计数
	completedCount uint64 // 代表请求成功完成的计数
	handlingNumber uint64 // 代表请求实时处理的计数
}

func (m *module) IncrCalledCount() {
	atomic.AddUint64(&m.calledCount, 1)
}

func (m *module) DecrCalledCount() {
	atomic.AddUint64(&m.calledCount, ^uint64(0))
}

func (m *module) IncrFailedCount() {
	atomic.AddUint64(&m.failedCount, 1)
}

func (m *module) DecrFailedCount() {
	atomic.AddUint64(&m.failedCount, ^uint64(0))
}

func (m *module) IncrAcceptedCount() {
	atomic.AddUint64(&m.acceptedCount, 1)
}

func (m *module) DecrAcceptedCount() {
	atomic.AddUint64(&m.acceptedCount, ^uint64(0))
}

func (m *module) IncrInterceptCount() {
	atomic.AddUint64(&m.interceptCount, 1)
}

func (m *module) DecrInterceptCount() {
	atomic.AddUint64(&m.interceptCount, ^uint64(0))
}

func (m *module) IncrCompletedCount() {
	atomic.AddUint64(&m.completedCount, 1)
}

func (m *module) DecrCompletedCount() {
	atomic.AddUint64(&m.completedCount, ^uint64(0))
}

func (m *module) IncrHandlingNumber() {
	atomic.AddUint64(&m.handlingNumber, 1)
}

func (m *module) DecrHandlingNumber() {
	atomic.AddUint64(&m.handlingNumber, ^uint64(0))
}

func (m *module) CalledCount() uint64 {
	return atomic.LoadUint64(&m.calledCount)
}

func (m *module) FailedCount() uint64 {
	return atomic.LoadUint64(&m.failedCount)
}

func (m *module) AcceptedCount() uint64 {
	return atomic.LoadUint64(&m.acceptedCount)
}

func (m *module) InterceptCount() uint64 {
	return atomic.LoadUint64(&m.interceptCount)
}

func (m *module) CompletedCount() uint64 {
	return atomic.LoadUint64(&m.completedCount)
}

func (m *module) HandlingNumber() uint64 {
	return atomic.LoadUint64(&m.handlingNumber)
}

func (m *module) Clear() {
	atomic.StoreUint64(&m.calledCount, 0)
	atomic.StoreUint64(&m.failedCount, 0)
	atomic.StoreUint64(&m.acceptedCount, 0)
	atomic.StoreUint64(&m.interceptCount, 0)
	atomic.StoreUint64(&m.completedCount, 0)
	atomic.StoreUint64(&m.handlingNumber, 0)
}
