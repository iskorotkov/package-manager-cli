package commands

import (
	"context"
	"fmt"
	"os"

	"github.com/google/go-github/v39/github"
	"github.com/jedib0t/go-pretty/v6/table"
	"github.com/spf13/cobra"
)

//nolint:gochecknoinits
func init() {
	searchCmd := &cobra.Command{ //nolint:exhaustivestruct
		Use:   "search",
		Short: "search for a package",
		Args:  cobra.MinimumNArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			client := github.NewClient(nil)

			result, _, err := client.Search.Repositories(context.Background(), args[0], nil)
			if err != nil {
				return fmt.Errorf("error searching repositories: %w", err)
			}

			fmt.Printf("found %d repositories\n", result.GetTotal())

			printReposList(result)

			return nil
		},
	}

	rootCmd.AddCommand(searchCmd)
}

func printReposList(result *github.RepositoriesSearchResult) {
	t := table.NewWriter()

	t.Style().Options.DrawBorder = false
	t.Style().Options.SeparateColumns = false

	t.SetOutputMirror(os.Stdout)
	t.AppendHeader(table.Row{"repo", "stars", "description"})

	for _, repo := range result.Repositories[:10] {
		t.AppendRow(table.Row{
			repo.GetFullName(),
			repo.GetStargazersCount(),
			repo.GetDescription(),
		})
	}

	t.Render()
}
