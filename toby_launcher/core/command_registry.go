package core

import (
	"sort"
	"strings"
)

type CommandRegistry struct {
	globalCommands []Command
	localCommands  []Command
	globalAliasMap map[string]string
	localAliasMap  map[string]string
}

func NewCommandRegistry() *CommandRegistry {
	return &CommandRegistry{
		globalAliasMap: make(map[string]string),
		localAliasMap:  make(map[string]string),
	}
}

func (r *CommandRegistry) UpdateAliases() {
	r.updateGlobalAliases()
	r.updateLocalAliases()
}

func (r *CommandRegistry) updateGlobalAliases() {
	r.globalAliasMap = make(map[string]string)
	for _, cmd := range r.globalCommands {
		cmdName := cmd.Name()
		r.globalAliasMap[strings.ToLower(cmdName)] = cmdName
		for _, alias := range cmd.Aliases() {
			r.globalAliasMap[strings.ToLower(alias)] = cmdName
		}
	}
}

func (r *CommandRegistry) updateLocalAliases() {
	r.localAliasMap = make(map[string]string)
	for _, cmd := range r.localCommands {
		cmdName := cmd.Name()
		r.localAliasMap[strings.ToLower(cmdName)] = cmdName
		for _, alias := range cmd.Aliases() {
			r.localAliasMap[strings.ToLower(alias)] = cmdName
		}
	}
}

func (r *CommandRegistry) containsCommand(cmds []Command, cmd Command) bool {
	for _, c := range cmds {
		if c.Name() == cmd.Name() {
			return true
		}
	}
	return false
}

func (r *CommandRegistry) findCommandByName(cmds []Command, name string) Command {
	for _, cmd := range cmds {
		if cmd.Name() == name {
			return cmd
		}
	}
	return nil
}

func (r *CommandRegistry) RegisterGlobalCommands(cmds []Command) {
	r.globalCommands = cmds
	r.sortCommands(r.globalCommands)
	r.updateGlobalAliases()
}

func (r *CommandRegistry) RegisterLocalCommands(cmds []Command) {
	r.localCommands = cmds
	r.sortCommands(r.localCommands)
	r.updateLocalAliases()
}

func (r *CommandRegistry) GetGlobalCommands() []Command {
	return r.globalCommands
}

func (r *CommandRegistry) GetLocalCommands() []Command {
	return r.localCommands
}

func (r *CommandRegistry) sortCommands(cmds []Command) {
	sort.Slice(cmds, func(i, j int) bool {
		return cmds[i].Name() < cmds[j].Name()
	})
}

func (r *CommandRegistry) GetCommand(input string) Command {
	input = strings.ToLower(input)
	if cmdId, exists := r.localAliasMap[input]; exists {
		return r.findCommandByName(r.localCommands, cmdId)
	}
	if cmdId, exists := r.globalAliasMap[input]; exists {
		return r.findCommandByName(r.globalCommands, cmdId)
	}
	if cmd := r.findCommandOrAliasByPrefix(r.localCommands, r.localAliasMap, input); cmd != nil {
		return cmd
	}
	if cmd := r.findCommandOrAliasByPrefix(r.globalCommands, r.globalAliasMap, input); cmd != nil {
		return cmd
	}
	return nil
}

func (r *CommandRegistry) findCommandOrAliasByPrefix(cmds []Command, aliasMap map[string]string, prefix string) Command {
	// 1. Проверяем команды по префиксу
	if cmd := r.findCommandByPrefix(cmds, prefix); cmd != nil {
		return cmd
	}
	// 2. Проверяем алиасы по префиксу
	for alias, cmdId := range aliasMap {
		if strings.HasPrefix(alias, prefix) {
			return r.findCommandByName(cmds, cmdId)
		}
	}
	return nil
}

func (r *CommandRegistry) findCommandByPrefix(cmds []Command, prefix string) Command {
	index := sort.Search(len(cmds), func(i int) bool {
		name := cmds[i].Name()
		return name >= prefix
	})
	if index < len(cmds) {
		name := cmds[index].Name()
		if strings.HasPrefix(name, prefix) {
			return cmds[index]
		}
	}
	return nil
}

func (r *CommandRegistry) ParseInput(input string) (Command, []string) {
	if input == "" {
		return nil, []string{}
	}
	args := strings.Fields(input)
	cmdPart := args[0]
	args[0] = strings.ToLower(cmdPart)
	if cmd := r.GetCommand(cmdPart); cmd != nil {
		return cmd, args
	}
	return nil, []string{}
}
