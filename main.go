package main

import (
	"app/cmd/admin"
	"app/cmd/public"
	"app/config"
	"fmt"
	"os"
	"strings"
)

func main() {
	if err := run(os.Args[1:]); err != nil {
		fmt.Fprintln(os.Stderr, "error : ", err)
		os.Exit(1)
	}
}

func run(args []string) error {

	actions := map[string]func(config.Parameters) error{
		admin.Action:  admin.Run,
		public.Action: public.Run,
	}

	c, err := config.New(args)
	if err != nil {
		return err
	}

	run, ok := actions[c.Action]
	if ok {
		return run(c)
	}

	listOfAction := make([]string, 0, len(actions))
	for a := range actions {
		listOfAction = append(listOfAction, a)
	}
	return fmt.Errorf("invalid action '%s' valid are [%s]", c.Action, strings.Join(listOfAction, " | "))
}
