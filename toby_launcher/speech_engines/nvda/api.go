//go:build windows

package nvda

import (
	"fmt"
	"golang.org/x/sys/windows"
	"toby_launcher/utils/file_utils"
	"unsafe"
)

var nvda_procs = []string{
	"nvdaController_testIfRunning",
	"nvdaController_speakText",
	"nvdaController_cancelSpeech",
}

var nvda_dll string = "nvdaControllerClient.dll"

type nvdaSynthesizer struct {
	dll   *windows.DLL
	procs map[string]*windows.Proc
}

func newNvdaSynthesizer() (*nvdaSynthesizer, error) {
	dll, err := file_utils.LoadDLL(nvda_dll)
	if err != nil {
		return nil, err
	}
	procs := make(map[string]*windows.Proc, len(nvda_procs))
	for _, p := range nvda_procs {
		proc, err := dll.FindProc(p)
		if err != nil {
			dll.Release()
			return nil, err
		}
		procs[p] = proc
	}
	return &nvdaSynthesizer{
		dll:   dll,
		procs: procs,
	}, nil
}

func (s *nvdaSynthesizer) free() {
	if s.dll != nil {
		s.dll.Release()
	}
}

func (s *nvdaSynthesizer) callProc(procName string, args ...uintptr) error {
	proc, exists := s.procs[procName]
	if !exists {
		return fmt.Errorf("Nvda function %s is not loaded.\r\n", procName)
	}
	ret, _, err := proc.Call(args...)
	if err != nil && err.Error() != "The operation completed successfully." {
		return err
	}
	status := uint32(ret)
	if status != 0 {
		return fmt.Errorf("Nvda error %d.\r\n", status)
	}
	return nil
}

func (s *nvdaSynthesizer) checkRunning() error {
	if err := s.callProc("nvdaController_testIfRunning"); err != nil {
		return err
	}
	return nil
}

func (s *nvdaSynthesizer) speak(text string) error {
	if err := s.checkRunning(); err != nil {
		return err
	}
	p, err := windows.UTF16PtrFromString(text)
	if err != nil {
		return err
	}
	return s.callProc("nvdaController_speakText", uintptr(unsafe.Pointer(p)))
}

func (s *nvdaSynthesizer) stop() error {
	if err := s.checkRunning(); err != nil {
		return err
	}
	return s.callProc("nvdaController_cancelSpeech")
}
