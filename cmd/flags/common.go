package flags

import (
	"github.com/spf13/cobra"
)

func AddVerboseFlag(cmd *cobra.Command) {
	cmd.Flags().BoolP("verbose", "v", false, "Verbose output")
}

func ResolveVerboseFlag(cmd *cobra.Command) (bool, error) {
	return cmd.Flags().GetBool("verbose")
}

func AddJsonFlag(cmd *cobra.Command) {
	cmd.Flags().BoolP("json", "j", false, "Output in JSON format")
}

func ResolveJsonFlag(cmd *cobra.Command) (bool, error) {
	return cmd.Flags().GetBool("json")
}

func AddShellFlag(cmd *cobra.Command) {
	cmd.Flags().BoolP("shell", "s", false, "Shell friendly output")
}

func ResolveShellFlag(cmd *cobra.Command) (bool, error) {
	return cmd.Flags().GetBool("shell")
}
