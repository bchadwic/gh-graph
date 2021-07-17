package main

import (
	"bytes"
	"fmt"
	"net/http"
	"os"
	"os/exec"
	"strconv"

	"github.com/PuerkitoBio/goquery"
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
	DaysOfWeek   = 7
	WeeksInAYear = 52
)

type GraphOptions struct {
	Username string
	Matrix   bool
	Solid    bool
}

type Stats struct {
	TotalContributions int
	LongestStreak      int
	BestDay            int
	AveragePerDay      float32
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

func printGraph(opts *GraphOptions, graph [][]int, stats *Stats) {
	colors := GetColors(stats)
	// TODO
	if colors == nil {
		return
	}
	return
}

func getGraph(doc *goquery.Document) ([][]int, *Stats) {
	graph := make([][]int, DaysOfWeek)
	for i := 0; i < DaysOfWeek; i++ {
		graph[i] = make([]int, WeeksInAYear+1)
	}

	stats := &Stats{
		TotalContributions: 0,
		LongestStreak:      0,
		AveragePerDay:      0,
	}
	k := -1
	count := 0
	curr := 0
	doc.Find(".js-calendar-graph rect[data-count]").Each(func(i int, s *goquery.Selection) {
		cell, exists := s.Attr("data-count")

		j := i % DaysOfWeek
		if j == 0 {
			k++
		}

		if exists {
			contribution, _ := strconv.Atoi(cell)
			count++
			if contribution > 0 {
				stats.TotalContributions += contribution
				curr++
			} else {
				if curr > stats.LongestStreak {
					stats.LongestStreak = curr
				}
				if contribution > stats.BestDay {
					stats.BestDay = contribution
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
