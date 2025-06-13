package component

import (
	mdl "common/model"
)

// *************************** Component ***************************
// Component is the interface which describes the behaviour for a Component
// which composite itself and its kind.
type KindName string
type IdName string

type Cpt interface {
	// for predicate purposes
	Id() IdName
	Kind() KindName
	//for logger
	CmptInfo() string
	//for runtime status
	IsRunning() bool

	//for golang ctx waitgroup etc. control structure
	//WithCtx(ctx context.Context) Component
	Ctrl() *mdl.CtrlSt

	// for basic control
	Start() error
	Stop() error
}

type CptRoot interface {
	Cpt
	Finalize() error
}

type CptsOperator interface {
	AddCpts(...Cpt)
	RemoveCpts(...Cpt)
	Cpt(IdName) Cpt
	Each(f func(Cpt))
}

type CptComposite interface {
	Cpt
	CptsOperator
}

// *************************** Commander ******************************
// Commander is the interface which describes the behaviour for a Component
// which exposes API commands.

// for structure api supports
type Cmder interface {
	// Cmd returns a command given a name. Returns nil if the command is not found.
	Cmd(any) (command func(map[any]any) any)
	// Commands returns a map of commands.
	Cmds() (commands map[any]func(map[any]any) any)
	// AddCommand adds a command given a name.
	AddCmd(name any, command func(map[any]any) any)
}
