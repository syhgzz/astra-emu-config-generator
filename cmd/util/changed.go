package util

import "github.com/spf13/cobra"

func SetIfChanged[T any](cmd *cobra.Command, flag string, setter func(T)) {
	if cmd.Flags().Changed(flag) {
		var v any
		switch any(*new(T)).(type) {
		case int:
			val, _ := cmd.Flags().GetInt(flag)
			v = val
		case int32:
			val, _ := cmd.Flags().GetInt32(flag)
			v = val
		case float64:
			val, _ := cmd.Flags().GetFloat64(flag)
			v = val
		case string:
			val, _ := cmd.Flags().GetString(flag)
			v = val
		case bool:
			val, _ := cmd.Flags().GetBool(flag)
			v = val
		}
		setter(v.(T))
	}
}
