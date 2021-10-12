package commands

import (
	"context"
	"fmt"
	"log"

	"github.com/google/go-github/v39/github"
	"github.com/iskorotkov/package-manager-cli/pkg/xlog"
	"github.com/jedib0t/go-pretty/v6/table"
	"github.com/spf13/cobra"
)

//nolint:gochecknoinits
func init() {
	searchCmd := wrapCommand(&cobra.Command{ //nolint:exhaustivestruct
		Use:   "search",
		Short: "search for a package",
		Args:  cobra.MinimumNArgs(1),
		RunE:  search,
	})

	rootCmd.AddCommand(searchCmd)
}

func search(_ *cobra.Command, args []string) error {
	packageName := args[0]

	xlog.Push(packageName)
	defer xlog.Pop()

	log.Printf("package name: %s", packageName)

	client := github.NewClient(nil)

	result, _, err := client.Search.Repositories(context.Background(), packageName, nil)
	if err != nil {
		return fmt.Errorf("error searching repositories: %w", err)
	}

	log.Printf("found repositories: %d", result.GetTotal())

	fmt.Printf("found %d repositories\n", result.GetTotal())

	printReposList(result)

	return nil
}

func printReposList(result *github.RepositoriesSearchResult) {
	t := createTable()
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
