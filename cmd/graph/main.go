package main

import (
	"bytes"
	"fmt"
	"net/http"
	"os"
	"os/exec"
	"strconv"
	"strings"

	"github.com/PuerkitoBio/goquery"
	"github.com/bchadwic/gh-graph/pkg/color"
	"github.com/bchadwic/gh-graph/pkg/stats"
	lg "github.com/charmbracelet/lipgloss"
	"github.com/cli/cli/git"
	"github.com/spf13/cobra"
)

func main() {
	graph := NewCmdGraph()

	if err := graph.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

const (
	WeeksInAYear = 52
	DaysInAWeek  = 7
)

type GraphOptions struct {
	Username string
	Matrix   bool
	Solid    bool
}

func NewCmdGraph() *cobra.Command {
	opts := &GraphOptions{}
	cmd := &cobra.Command{
		Use:   "graph",
		Short: "Display your GitHub contribution graph",
		Long:  "Display your GitHub contribution graph in the terminal",
		RunE: func(cmd *cobra.Command, args []string) error {
			return runGraph(opts)
		},
		SilenceErrors: true,
		SilenceUsage:  true,
	}

	cmd.Flags().StringVarP(&opts.Username, "username", "u", "", "Specify a user")
	cmd.Flags().BoolVarP(&opts.Matrix, "matrix", "m", false, "Set cells to matrix digital rain")
	cmd.Flags().BoolVarP(&opts.Solid, "solid", "s", false, "Set cells to solid blocks")
	return cmd
}

func runGraph(opts *GraphOptions) error {
	if opts.Username == "" {
		_, err := git.GitCommand()
		if err != nil {
			return err
		}
		if output, err := exec.Command("git", "config", "user.name").Output(); err != nil {
			return err
		} else {
			opts.Username = firstLine(output)
		}
	}

	resp, err := http.Get("https://github.com/" + opts.Username)
	if err != nil {
		return err
	}

	defer resp.Body.Close()

	if resp.StatusCode != 200 {
		if resp.StatusCode == 404 {
			return fmt.Errorf("user %s was not found on GitHub, choose a new user with -u / --user", opts.Username)
		}
		return fmt.Errorf("unexpected status code: %d", resp.StatusCode)
	}

	doc, err := goquery.NewDocumentFromReader(resp.Body)
	if err != nil {
		return err
	}

	graph, stats := getGraph(doc)
	printGraph(opts, graph, stats)
	return nil
}

func printGraph(opts *GraphOptions, graph [][]int, stats *stats.Stats) {
	cp := &color.ColorPalette{}
	cp.Initialize(stats)
	DaysOfTheWeek := []string{"    ", "Mon ", "    ", "Wed ", "    ", "Fri ", "    "}

	b := strings.Builder{}
	for i, x := range graph {
		b.WriteString(DaysOfTheWeek[i])
		for _, y := range x {
			s := lg.Style{}
			if y != 0 {
				s = lg.NewStyle().SetString("#").Foreground(lg.Color(cp.GetColor(y)))
			} else {
				s = lg.NewStyle().SetString(" ")
			}
			b.WriteString(s.String())
		}
		b.WriteString("\n")
	}
	b.WriteString(
		fmt.Sprintf("%s\ncontributions in the last year: %d\nlongest streak: %d, average: %.3f/day, best day: %d",
			"github.com/"+opts.Username, stats.TotalContributions, stats.LongestStreak, stats.AveragePerDay, stats.BestDay))

	dialogBoxStyle := lg.NewStyle().SetString(b.String()).
		Border(lg.RoundedBorder()).
		Margin(1).
		BorderTop(true).
		BorderLeft(true).
		BorderRight(true).
		BorderBottom(true)
	fmt.Println(dialogBoxStyle)

	return
}

func getGraph(doc *goquery.Document) ([][]int, *stats.Stats) {
	graph := make([][]int, DaysInAWeek)
	for i := 0; i < DaysInAWeek; i++ {
		graph[i] = make([]int, WeeksInAYear+1)
	}

	stats := &stats.Stats{
		TotalContributions: 0,
		LongestStreak:      0,
		AveragePerDay:      0,
		BestDay:            0,
	}
	k := -1
	count := 0
	curr := 0
	doc.Find(".js-calendar-graph rect[data-count]").Each(func(i int, s *goquery.Selection) {
		cell, exists := s.Attr("data-count")

		j := i % DaysInAWeek
		if j == 0 {
			k++
		}

		if exists {
			contribution, _ := strconv.Atoi(cell)
			count++
			if contribution > 0 {
				stats.TotalContributions += contribution
				curr++
				if contribution > stats.BestDay {
					stats.BestDay = contribution
				}
			} else {
				if curr > stats.LongestStreak {
					stats.LongestStreak = curr
				}
				curr = 0
			}
			graph[j][k] = contribution
		}
	})
	stats.AveragePerDay = float32(stats.TotalContributions) / float32(count)
	return graph, stats
}

func firstLine(output []byte) string {
	if i := bytes.IndexAny(output, "\n"); i >= 0 {
		return string(output)[0:i]
	}
	return string(output)
}
