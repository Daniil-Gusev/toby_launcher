package game

import (
	"toby_launcher/apperrors"
	"toby_launcher/config"
)

var (
	videoBackendParamKey string = "vid_preferbackend"
	musicParamKey        string = "music"
	soundfxParamKey      string = "sound_fx"
)

type GameParams struct {
	params map[string]GameParam
	config config.GzdoomParams
}

func newGameParams(cfg config.GzdoomParams) (*GameParams, error) {
	params := map[string]GameParam{
		videoBackendParamKey: newVideoBackendParam(),
		musicParamKey:        newMusicParam(),
		soundfxParamKey:      newSoundfxParam(),
	}
	gp := &GameParams{
		params: params,
		config: cfg,
	}
	var warn error
	if err := gp.ApplyConfig(cfg); err != nil {
		warn = apperrors.New(apperrors.Err, "Warning: incorrect GZDoom params in configuration:\n$errors", map[string]any{"errors": err})
	}
	gp.updateConfig()
	if warn != nil {
		return gp, warn
	}
	return gp, nil
}

func (p *GameParams) ApplyConfig(cfg config.GzdoomParams) error {
	errs := apperrors.NewErrors(nil)
	for key, val := range cfg {
		if _, exists := p.params[key]; !exists {
			err := apperrors.New(apperrors.Err, "param with key $key does not exist", map[string]any{"key": key})
			errs.Add(err)
			continue
		}
		if err := p.params[key].set(val); err != nil {
			errs.Add(err)
			continue
		}
	}
	if errs.Count() > 0 {
		return errs
	}
	return nil
}

func (p *GameParams) updateConfig() {
	for key, param := range p.params {
		if _, exists := p.config[key]; !exists {
			p.config[key] = param.value()
		}
	}
}

func (p *GameParams) save() {
	for key, param := range p.params {
		p.config[key] = param.value()
	}
}

func (p *GameParams) toCmdArgs() []string {
	args := make([]string, 0, 2*len(p.params))
	for _, param := range p.params {
		args = append(args, param.toCmdArgs()...)
	}
	return args
}

func (p *GameParams) Params() map[string]GameParam {
	return p.params
}

func (p *GameParams) set(key string, value any) error {
	if _, exists := p.params[key]; !exists {
		return apperrors.New(apperrors.Err, "param with key $key does not exist", map[string]any{"key": key})
	}
	if err := p.params[key].set(value); err != nil {
		return err
	}
	p.config[key] = value
	return nil
}

func (p *GameParams) VideoBackend() VidBackend {
	return p.params[videoBackendParamKey].value().(VidBackend)
}

func (p *GameParams) SetVideoBackend(backend VidBackend) error {
	return p.set(videoBackendParamKey, backend)
}

func (p *GameParams) MusicPtr() *bool {
	return &p.params[musicParamKey].(*MusicParam).active
}

func (p *GameParams) SetMusic(status bool) error {
	return p.set(musicParamKey, status)
}

func (p *GameParams) SoundfxPtr() *bool {
	return &p.params[soundfxParamKey].(*SoundfxParam).active
}

func (p *GameParams) SetSoundfx(status bool) error {
	return p.set(soundfxParamKey, status)
}
