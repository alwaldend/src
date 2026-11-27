package flagcompat

import (
	"strings"

	"github.com/spf13/cobra"
)

// SetArgs supplies a flat command's arguments after its flags have been
// registered. It accepts Go's single-dash flag names and ends flag parsing
// at the first operand. The caller retains ownership of args.
func SetArgs(command *cobra.Command, args []string) {
	command.InitDefaultHelpFlag()
	flags := command.Flags()
	flags.SetInterspersed(false)

	// Cobra interprets nil as a request to read the process's arguments.
	normalized := make([]string, len(args))
	copy(normalized, args)
	for i := 0; i < len(args); i++ {
		arg := args[i]
		if !strings.HasPrefix(arg, "-") || arg == "-" || arg == "--" {
			break
		}
		name, _, attached := strings.Cut(arg[1:], "=")
		name, doubleDash := strings.CutPrefix(name, "-")
		flag := flags.Lookup(name)
		if flag != nil && !doubleDash {
			normalized[i] = "-" + arg
		} else if flag == nil && !doubleDash && len(name) == 1 {
			flag = flags.ShorthandLookup(name)
		}
		if flag != nil && flag.NoOptDefVal == "" && !attached {
			i++
		}
	}
	command.SetArgs(normalized)
}
