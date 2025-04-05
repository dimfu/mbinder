package cmd

import (
	"fmt"
	"io/fs"
	"log"
	"os"
	"path/filepath"
	"regexp"
	"slices"
	"strings"

	"github.com/dimfu/mbinder/models"
	"github.com/spf13/cobra"
)

type AddCommand struct {
	cmd        *cobra.Command
	collection string
	tags       string
}

func NewAddCommand() *AddCommand {
	ac := &AddCommand{}
	ac.cmd = &cobra.Command{
		Use:   "add",
		Short: "Add one or more media files to be organized",
		Args:  cobra.MatchAll(cobra.MinimumNArgs(1), cobra.OnlyValidArgs),
		Run: func(cmd *cobra.Command, args []string) {
			p, err := ac.scan(args)
			cobra.CheckErr(err)
			if len(p) == 0 {
				fmt.Fprintf(os.Stderr, "error: could not find media files\n")
				os.Exit(1)
			}
			ac.modify(p)
		},
	}
	ac.cmd.Flags().StringVarP(&ac.tags, "tags", "t", "", "tags to be added (eg; 'blue, outside, beach')")
	ac.cmd.Flags().StringVarP(&ac.collection, "collection", "c", "", "collection name, default to current folder name")
	return ac
}

func init() {
	addCmd := NewAddCommand()
	rootCmd.cmd.AddCommand(addCmd.cmd)
}

func (a *AddCommand) isMedia(ext string) bool {
	allowedExtensions := []string{
		".jpg", ".jpeg", ".png", ".gif", ".bmp", ".tiff", ".webp", ".heic", ".avif",
		".mp4", ".mov", ".avi", ".mkv", ".webm", ".flv", ".wmv", ".3gp",
	}
	return slices.Contains(allowedExtensions, ext)
}

func (a *AddCommand) _tags() []*models.Tag {
	var tags []*models.Tag
	re := regexp.MustCompile(`\s*,\s*`)
	splitTags := re.Split(a.tags, -1)
	for _, t := range splitTags {
		var tag *models.Tag
		err := rootCmd.DB.FirstOrCreate(&tag, models.Tag{
			Name: t,
		}).Error
		if err != nil {
			fmt.Println(err)
			continue
		}
		tags = append(tags, tag)
	}
	return tags
}

func (a *AddCommand) scan(args []string) (map[string][]*models.Item, error) {
	paths := make(map[string][]*models.Item)
	tags := a._tags()
	walk := func(path string, info fs.FileInfo, err error) error {
		if info.IsDir() {
			return nil
		}

		ext := strings.ToLower(filepath.Ext(info.Name()))
		if !a.isMedia(ext) {
			return nil
		}

		abs, err := filepath.Abs(path)
		cobra.CheckErr(err)
		folder := filepath.Base(filepath.Dir(abs))
		if _, exists := paths[folder]; !exists {
			paths[folder] = []*models.Item{}
		}

		paths[folder] = append(paths[folder], &models.Item{
			Path:      abs,
			MediaType: ext,
			Tags:      tags,
		})

		return err
	}

	for _, p := range args {
		err := filepath.Walk(p, walk)
		cobra.CheckErr(err)
	}
	return paths, nil
}

func (a *AddCommand) modify(payload map[string][]*models.Item) {
	for key, items := range payload {
		err := rootCmd.DB.FirstOrCreate(&models.Collection{}, models.Collection{
			Name: key,
		}).Error
		cobra.CheckErr(err)
		result := rootCmd.DB.CreateInBatches(items, 50)
		if result.Error != nil {
			log.Println("insert failed:", result.Error)
			continue
		}
		log.Println("inserted:", result.RowsAffected, "records")
	}
}
