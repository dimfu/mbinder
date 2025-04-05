package cmd

import "github.com/spf13/cobra"

type TagsCommand struct {
	cmd *cobra.Command
}

func NewTagsCommand() *TagsCommand {
	tc := &TagsCommand{}
	tc.cmd = &cobra.Command{
		Use:   "tags",
		Short: "Get media files by tags",
		Args:  cobra.MatchAll(cobra.MinimumNArgs(1), cobra.OnlyValidArgs),
		Run: func(cmd *cobra.Command, args []string) {
		},
	}
	return tc
}

func init() {
	tagsCmd := NewTagsCommand()
	rootCmd.cmd.AddCommand(tagsCmd.cmd)
}
