package pinbase

import (
	"testing"
)

type BackendMock struct {
	PinProcessorBumpFn      func() <-chan struct{}
	PinProcessorBumpInvoked bool

	PinRequirementsFn      func() map[Hash]bool
	PinRequirementsInvoked bool

	NotifyPinFn      func(Hash, *PinBackendState)
	NotifyPinInvoked bool
}

func (m *BackendMock) PinProcessorBump() <-chan struct{} {
	m.PinProcessorBumpInvoked = true
	return m.PinProcessorBumpFn()
}

func (m *BackendMock) PinRequirements() map[Hash]bool {
	m.PinRequirementsInvoked = true
	return m.PinRequirementsFn()
}

func (m *BackendMock) NotifyPin(h Hash, s *PinBackendState) {
	m.NotifyPinInvoked = true
	m.NotifyPinFn(h, s)
}

var _ PinBackend = &BackendMock{}

func TestProcessPins(t *testing.T) {

}
