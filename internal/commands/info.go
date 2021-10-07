package commands

import (
	"context"
	"fmt"
	"os"
	"strings"

	"github.com/google/go-github/v39/github"
	"github.com/jedib0t/go-pretty/v6/table"
	"github.com/spf13/cobra"
)

//nolint:gochecknoinits
func init() {
	infoCmd := &cobra.Command{ //nolint:exhaustivestruct
		Use:   "info",
		Short: "show info about package",
		Args:  cobra.MinimumNArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			packageName := args[0]

			client := github.NewClient(nil)

			result, _, err := client.Search.Repositories(context.Background(), packageName, nil)
			if err != nil {
				return fmt.Errorf("error searching repositories: %w", err)
			}

			if len(result.Repositories) == 0 {
				return fmt.Errorf("no results")
			}

			repo := result.Repositories[0]

			printRepoInfo(repo)

			releases, _, err := client.Repositories.ListReleases(
				context.Background(),
				repo.GetOwner().GetLogin(),
				repo.GetName(),
				nil,
			)
			if err != nil {
				return fmt.Errorf("error getting releases: %w", err)
			}

			printReleasesList(releases)

			return nil
		},
	}

	rootCmd.AddCommand(infoCmd)
}

func printReleasesList(releases []*github.RepositoryRelease) {
	t := table.NewWriter()

	t.Style().Options.DrawBorder = false
	t.Style().Options.SeparateColumns = false

	t.SetOutputMirror(os.Stdout)
	t.AppendHeader(table.Row{"version", "name", "url"})

	for _, release := range releases {
		t.AppendRow(table.Row{
			release.GetTagName(),
			release.GetName(),
			release.GetHTMLURL(),
		})
	}

	t.Render()
}

func printRepoInfo(repo *github.Repository) {
	fmt.Printf("name: %s\n", repo.GetFullName())
	fmt.Printf("stars: %d\n", repo.GetStargazersCount())
	fmt.Printf("description: %s\n", repo.GetDescription())
	fmt.Println("-----")
	fmt.Printf("homepage: %s\n", repo.GetHomepage())
	fmt.Printf("url: %s\n", repo.GetURL())
	fmt.Printf("language: %s\n", repo.GetLanguage())
	fmt.Printf("forks: %d\n", repo.GetForksCount())
	fmt.Printf("topics: %s\n", strings.Join(repo.Topics, ", "))
	fmt.Println("-----")
}
