package cmd

import (
	"regexp"

	"github.com/dimfu/mbinder/models"
	"github.com/dimfu/mbinder/services/http_server"
	"github.com/spf13/cobra"
)

type ShowCommand struct {
	cmd        *cobra.Command
	collection string
	tags       string
}

func NewShowCommand() *ShowCommand {
	sc := &ShowCommand{}
	sc.cmd = &cobra.Command{
		Use:   "show",
		Short: "Get media files by tags",
		Run: func(cmd *cobra.Command, args []string) {
			re := regexp.MustCompile(`\s*,\s*`)
			tags := re.Split(sc.tags, -1)
			sc.collect(tags)
		},
	}
	sc.cmd.Flags().StringVarP(&sc.tags, "tags", "t", "", "tags to be added (eg; 'blue, outside, beach')")
	sc.cmd.Flags().StringVarP(&sc.collection, "collection", "c", "", "collection name, default to current folder name")
	return sc
}

func init() {
	show := NewShowCommand()
	rootCmd.cmd.AddCommand(show.cmd)
}

func (sc *ShowCommand) collect(tags []string) {
	var items []*models.Item
	err := rootCmd.DB.
		Joins("JOIN items_tags ON items.id = items_tags.item_id").
		Joins("JOIN tags ON tags.id = items_tags.tag_id").
		Where("tags.name IN ?", tags).
		Preload("Tags").
		Find(&items).Error

	if err != nil {
		cobra.CheckErr(err)
	}

	http_server.Run(items)
}
