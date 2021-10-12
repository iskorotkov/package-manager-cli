package commands

import (
	"context"
	"fmt"
	"log"
	"os"
	"strings"

	"github.com/google/go-github/v39/github"
	"github.com/iskorotkov/package-manager-cli/pkg/xlog"
	"github.com/jedib0t/go-pretty/v6/table"
	"github.com/spf13/cobra"
)

//nolint:gochecknoinits
func init() {
	infoCmd := wrapCommand(&cobra.Command{ //nolint:exhaustivestruct
		Use:   "info",
		Short: "show info about package",
		Args:  cobra.MinimumNArgs(1),
		RunE:  info,
	})

	rootCmd.AddCommand(infoCmd)
}

func info(_ *cobra.Command, args []string) error {
	packageName := args[0]

	xlog.Push(packageName)
	defer xlog.Pop()

	log.Printf("package name: %s", packageName)

	client := github.NewClient(nil)

	result, _, err := client.Search.Repositories(context.Background(), packageName, nil)
	if err != nil {
		return fmt.Errorf("error searching repositories: %w", err)
	}

	log.Printf("found repositories: %d", len(result.Repositories))

	if len(result.Repositories) == 0 {
		return fmt.Errorf("no results")
	}

	repo := result.Repositories[0]

	log.Printf("got repo info: %s", repo.GetFullName())

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

	log.Printf("got releases: %d", len(releases))

	printReleasesList(releases)

	return nil
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
