package core

type Command interface {
	Name() string
	Description() string
	Aliases() []string
	Execute(ctx *AppContext, ui *UiContext, args []string) (State, error)
}

type BaseCommand struct{}

func (c *BaseCommand) Execute(ctx *AppContext, ui *UiContext, args []string) (State, error) {
	ui.DisplayText("Unknown command action.")
	return ctx.GetCurrentState()
}

func (c *BaseCommand) Name() string {
	return "unknown"
}

func (c *BaseCommand) Description() string {
	return "The description for this command is not defined."
}

func (c *BaseCommand) Aliases() []string {
	return []string{}
}

func DefaultGlobalCommands() []Command {
	return []Command{
		&HelpCommand{},
		&QuitCommand{},
		&VersionCommand{},
	}
}
