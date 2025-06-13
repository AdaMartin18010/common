package component

var (
	_ Cmder = (*cmder)(nil)
)

type cmder struct {
	cmds map[any]func(map[any]any) any
}

// NewCmder returns a new Commander.
func NewCmder() Cmder {
	return &cmder{
		cmds: make(map[any]func(map[any]any) any),
	}
}

// Command returns the command interface when passed a valid command name
func (c *cmder) Cmd(name any) (command func(map[any]any) any) {
	command = c.cmds[name]
	return
}

// Commands returns the entire map of valid commands
func (c *cmder) Cmds() map[any]func(map[any]any) any {
	return c.cmds
}

// AddCommand adds a new command, when passed a command name and the command interface.
func (c *cmder) AddCmd(name any, command func(map[any]any) any) {
	c.cmds[name] = command
}
