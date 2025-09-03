package flags

import (
	"github.com/spf13/cobra"
)

func AddDeleteFlag(cmd *cobra.Command, description string) {
	cmd.Flags().BoolP("delete", "d", false, description)
}

func ResolveDeleteFlag(cmd *cobra.Command) (bool, error) {
	return cmd.Flags().GetBool("delete")
}

func AddLabelFlag(cmd *cobra.Command) {
	cmd.Flags().StringP("label", "l", "", "Label")
}

func ResolveLabelFlag(cmd *cobra.Command) (string, error) {
	return cmd.Flags().GetString("label")
}

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
