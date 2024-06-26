/*
Copyright © 2024 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"fmt"
	"strings"

	"github.com/rayhanadev/git-esque/internal/repo"
	"github.com/spf13/cobra"
)

var branch string

// cloneCmd represents the clone command
var cloneCmd = &cobra.Command{
	Use:   "clone [repository] [directory]",
	Short: "Clone a repository",
	Args:  cobra.MinimumNArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		repoURL := args[0]
		var directory string
		if len(args) > 1 {
			directory = args[1]
		} else {
			parts := strings.Split(repoURL, "/")
			directory = strings.TrimSuffix(parts[len(parts)-1], ".git")
		}

		err := repo.Clone(repoURL, directory, branch)
		if err != nil {
			fmt.Println("Error cloning repository:", err)
		} else {
			fmt.Println("Repository cloned.")
		}
	},
}

func init() {
	rootCmd.AddCommand(cloneCmd)
	rootCmd.AddCommand(cloneCmd)
	cloneCmd.Flags().StringVarP(&branch, "branch", "b", "master", "Branch to clone")
	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// cloneCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// cloneCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
