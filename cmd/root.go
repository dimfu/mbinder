package cmd

import (
	"os"
	"path"

	"github.com/dimfu/mbinder/models"
	"github.com/glebarez/sqlite"
	"github.com/spf13/cobra"
	"gorm.io/gorm"
)

type RootCmd struct {
	cmd *cobra.Command
	DB  *gorm.DB
}

func NewRootCmd() *RootCmd {
	return &RootCmd{
		DB: nil,
		cmd: &cobra.Command{
			Use:   "mbinder",
			Short: "Utility to organize all kind of media, preferrably image and video.",
		},
	}
}

var (
	rootCmd = NewRootCmd()
)

func Execute() (*RootCmd, error) {
	return rootCmd, rootCmd.cmd.Execute()
}

func init() {
	cobra.OnInitialize(initConfig)
}

func initConfig() {
	dir, err := os.UserHomeDir()
	cobra.CheckErr(err)
	p := path.Join(dir, ".mbinder")
	_, err = os.Stat(p)
	if os.IsNotExist(err) {
		if err := os.MkdirAll(p, os.ModePerm); err != nil {
			cobra.CheckErr(err)
		}
	} else if err != nil {
		cobra.CheckErr(err)
	}
	db, err := gorm.Open(sqlite.Open(path.Join(p, "database.db")), &gorm.Config{})
	cobra.CheckErr(err)
	rootCmd.DB = db

	err = rootCmd.DB.AutoMigrate(&models.Item{}, &models.Collection{}, &models.Tag{})
	cobra.CheckErr(err)
}
