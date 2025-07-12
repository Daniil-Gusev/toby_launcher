package core

type State interface {
	Name() string
	Description() string
	Commands() []Command
	RequiresInput() bool
	Init(ctx *AppContext, ui *UiContext) (State, error)
	Handle(ctx *AppContext, ui *UiContext, input string) (State, error)
	Display(ctx *AppContext, ui *UiContext)
}

type BaseState struct{}

func (b *BaseState) Name() string {
	return "unknown"
}

func (b *BaseState) Description() string {
	return "The description for this state is not defined."
}

func (b *BaseState) Commands() []Command {
	return []Command{}
}

func (b *BaseState) RequiresInput() bool {
	return true
}

func (b *BaseState) Init(ctx *AppContext, ui *UiContext) (State, error) {
	return b, nil
}

func (b *BaseState) Display(ctx *AppContext, ui *UiContext) {}

func (b *BaseState) Handle(ctx *AppContext, ui *UiContext, input string) (State, error) {
	return b, nil
}
