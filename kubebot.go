package main

import (
	"errors"
	"fmt"
	"regexp"
	"strings"
	"time"

	"github.com/go-chat-bot/bot"
)

type Kubebot struct {
	token    string
	admins   map[string]bool
	channels map[string]bool
	commands map[string]bool
}

const (
	forbiddenChannelMessage  string = "%s - ⚠ Channel %s forbidden (%s) for user @%s\n"
	forbiddenChannelResponse string = "Sorry @%s, but I'm not allowed to run this command here :zipper_mouth_face:"
	forbiddenCommandResponse string = "Sorry @%s, but I cannot run this command."
	forbiddenFlagResponse    string = "Sorry @%s, but I'm not allowed to run one of your flags."
	okResponse               string = "Roger that!\n@%s, this is the response to your request:\n ```\n%s\n``` "
)

var (
	ignored = map[string]map[string]bool{
		"get": map[string]bool{
			"-f":           true,
			"--filename":   true,
			"-w":           true,
			"--watch":      true,
			"--watch-only": true,
		},
		"describe": map[string]bool{
			"-f":         true,
			"--filename": true,
		},
		"create": map[string]bool{
			"-f":         true,
			"--filename": true,
		},
		"replace": map[string]bool{
			"-f":         true,
			"--filename": true,
		},
		"patch": map[string]bool{
			"-f":         true,
			"--filename": true,
		},
		"delete": map[string]bool{
			"-f":         true,
			"--filename": true,
		},
		"edit": map[string]bool{
			"-f":         true,
			"--filename": true,
		},
		"apply": map[string]bool{
			"-f":         true,
			"--filename": true,
		},
		"logs": map[string]bool{
			"-f":       true,
			"--follow": true,
		},
		"rolling-update": map[string]bool{
			"-f":         true,
			"--filename": true,
		},
		"scale": map[string]bool{
			"-f":         true,
			"--filename": true,
		},
		"attach": map[string]bool{
			"-i":      true,
			"--stdin": true,
			"-t":      true,
			"--tty":   true,
		},
		"exec": map[string]bool{
			"-i":      true,
			"--stdin": true,
			"-t":      true,
			"--tty":   true,
		},
		"run": map[string]bool{
			"--leave-stdin-open": true,
			"-i":                 true,
			"--stdin":            true,
			"--tty":              true,
		},
		"expose": map[string]bool{
			"-f":         true,
			"--filename": true,
		},
		"autoscale": map[string]bool{
			"-f":         true,
			"--filename": true,
		},
		"label": map[string]bool{
			"-f":         true,
			"--filename": true,
		},
		"annotate": map[string]bool{
			"-f":         true,
			"--filename": true,
		},
		"convert": map[string]bool{
			"-f":         true,
			"--filename": true,
		},
	}

	// Slack will automatically reformat anything that looks like a URL, including some image URIs.
	// For example: "myregistry.com/org/repository:tag" -> "<http://myregistry.com/org/repository:tag|myregistry.com/org/repository:tag>\u00a0"
	// This is what make the text look blue in the Slack UI. We need to undo this before trying to pass the message to kubectl.
	undo_urls = regexp.MustCompile(`<.*\|(.*)>`)
)

func validateFlags(arguments ...string) error {
	if len(arguments) <= 1 {
		return nil
	}

	for i := 1; i < len(arguments); i++ {
		if ignored[arguments[0]][arguments[i]] {
			return errors.New(fmt.Sprintf("Error: %s is an invalid flag", arguments[i]))
		}
	}

	return nil
}

func kubectl(command *bot.Cmd) (msg string, err error) {
	t := time.Now()
	time := t.Format(time.RFC3339)
	nickname := command.User.Nick

	if !kb.channels[command.Channel] {
		fmt.Printf(forbiddenChannelMessage, time, command.Channel, kb.channels, nickname)
		return fmt.Sprintf(forbiddenChannelResponse, nickname), nil
	}

	// command.Args gets mangled/truncated when slack thinks there are "URL"s, so we need to use command.RawArgs and split it ourselves

	// Undo any "linkification" by slack, and fix weird characters
	unlinkedArgs := undo_urls.ReplaceAllString(command.RawArgs, "${1}")
	// Also fix any instances of the special "\u00a0" character that seem to sneak in when "linkification" occurs
	fixedArgs := strings.ReplaceAll(unlinkedArgs, "\u00a0", " ")

	fixedArgsSlice := strings.Split(fixedArgs, " ")

	fmt.Printf("Running kubectl command for %s: kubectl %+v", nickname, fixedArgsSlice)

	output := execute("kubectl", fixedArgsSlice...)
	return fmt.Sprintf(okResponse, nickname, output), nil
}

func init() {
	bot.RegisterCommand(
		"kubectl",
		"Kubectl Slack integration",
		"",
		kubectl)
}
