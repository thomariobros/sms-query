package command

import (
	"sms-query/pkg/config"
	"sms-query/pkg/i18n"
	"sms-query/pkg/sms"
	"sms-query/pkg/util"
)

type Dispatcher struct {
	commands []*Command
}

func NewDispatcher() *Dispatcher {
	result := &Dispatcher{}
	result.commands = commands
	// append help command here to avoid initialization loop
	result.commands = append(result.commands, helpCmd)
	return result
}

func (dispatcher *Dispatcher) addCommand(cmd *Command) {
	dispatcher.commands = append(dispatcher.commands, cmd)
}

func (dispatcher *Dispatcher) findCommand(locale string, query string) (*Command, bool, map[string]string) {
	standardizedQuery := util.StandardizeStringLower(query)
	for _, cmd := range dispatcher.commands {
		if cmd.IsMatch(locale, standardizedQuery) {
			return cmd, true, cmd.GetArgs(locale, standardizedQuery)
		} else if cmd.IsMatchKeyOnly(locale, standardizedQuery) {
			// only match cmd key
			return cmd, false, nil
		}
	}
	return nil, false, nil
}

func (dispatcher *Dispatcher) Execute(from string, query string) string {
	// check aliases
	preferences := config.GetInstance().GetPreferences(from)
	alias := preferences.FindAlias(query)
	if alias != "" {
		// replace query by alias value
		query = alias
	}
	message := ""
	// find command
	if cmd, isValid, args := dispatcher.findCommand(preferences.Defaults.Locale, query); cmd != nil {
		if isValid {
			// execute command
			content, err := cmd.Run(from, args)
			// build message
			message = sms.Build(from, query, content, err)
		} else {
			// invalid message when only cmd's key matches
			message = i18n.GetTranslationForPhoneNumber(from, cmd.InvalidMessageKey)
		}
	}
	if message == "" {
		message = i18n.GetTranslationForPhoneNumber(from, "command.unknown")
	}
	return message
}

func (dispatcher *Dispatcher) ExecuteSend(from string, query string) error {
	return sms.Send(from, dispatcher.Execute(from, query))
}
