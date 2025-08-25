package game

import (
	"fmt"
	"toby_launcher/apperrors"
)

type GameParam interface {
	toCmdArgs() []string
	value() any
	set(value any) error
}

type VidBackend int

const (
	OpenGLBackend   VidBackend = 0
	VulkanBackend   VidBackend = 1
	OpenGLESBackend VidBackend = 2
)

var VideoBackends []VidBackend = []VidBackend{
	OpenGLBackend, VulkanBackend, OpenGLESBackend,
}

func (b VidBackend) String() string {
	switch b {
	case OpenGLBackend:
		return "OpenGL"
	case VulkanBackend:
		return "Vulkan"
	case OpenGLESBackend:
		return "OpenGL ES"
	default:
		return "unknown"
	}
}

func (b VidBackend) isValid() bool {
	switch b {
	case OpenGLBackend, VulkanBackend, OpenGLESBackend:
		return true
	default:
		return false
	}
}

type videoBackendParam struct {
	backend VidBackend
}

func newVideoBackendParam() GameParam {
	return &videoBackendParam{backend: OpenGLBackend}
}

func (p *videoBackendParam) toCmdArgs() []string {
	return []string{"+vid_preferbackend", fmt.Sprintf("%d", p.backend)}
}

func (p *videoBackendParam) value() any {
	return p.backend
}

func (p *videoBackendParam) set(value any) error {
	var backend VidBackend
	switch v := value.(type) {
	case VidBackend:
		backend = v
	case float64:
		backend = VidBackend(int(v))
	case int:
		backend = VidBackend(v)
	default:
		return apperrors.New(apperrors.Err, "invalid type for vid_preferbackend: expected VidBackend or number, got $type", map[string]any{"type": fmt.Sprintf("%T", value)})
	}
	if !backend.isValid() {
		return apperrors.New(apperrors.Err, "invalid VidBackend value: $value", map[string]any{"value": backend})
	}
	p.backend = backend
	return nil
}

type toggleSwitcher struct {
	active bool
}

func (s *toggleSwitcher) value() any {
	return s.active
}

func (s *toggleSwitcher) set(value any) error {
	status, ok := value.(bool)
	if !ok {
		return apperrors.New(apperrors.Err, "invalid type of value $value, expected bool", map[string]any{"value": value})
	}
	s.active = status
	return nil
}

type MusicParam struct{ toggleSwitcher }

func newMusicParam() GameParam {
	param := &MusicParam{}
	param.active = true
	return param
}

func (p *MusicParam) toCmdArgs() []string {
	if p.active {
		return []string{}
	}
	return []string{"-nomusic"}
}

type SoundfxParam struct{ toggleSwitcher }

func newSoundfxParam() GameParam {
	param := &SoundfxParam{}
	param.active = true
	return param
}

func (p *SoundfxParam) toCmdArgs() []string {
	if p.active {
		return []string{}
	}
	return []string{"-nosfx"}
}
