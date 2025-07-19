package flags

import "github.com/spf13/cobra"



func AddDeleteFlag(cmd *cobra.Command, description string) {
	cmd.Flags().BoolP("delete", "d", false, description)
}

func ResolveDeleteFlag(cmd *cobra.Command) (bool, error){
	return cmd.Flags().GetBool("delete")
}

func AddLabelFlag(cmd *cobra.Command) {
	cmd.Flags().StringP("label", "l", "", "Label")
}

func ResolveLabelFlag(cmd *cobra.Command) (string, error){
	return cmd.Flags().GetString("label")
}
