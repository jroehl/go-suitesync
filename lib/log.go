package lib

import (
	"fmt"
	"log"
	"os"
	"sort"
	"strconv"
	"strings"

	"github.com/olekukonko/tablewriter"
	"github.com/ttacon/chalk"
)

var (
	warn   = chalk.Red.NewStyle().WithTextStyle(chalk.Bold).WithBackground(chalk.White).Style
	err    = chalk.Red.NewStyle().WithTextStyle(chalk.Bold).Style
	notice = chalk.Cyan.NewStyle().WithTextStyle(chalk.Dim).Style
	head   = chalk.Green.NewStyle().WithTextStyle(chalk.Underline).Style
	result = chalk.Green.NewStyle().Style
)

// PrHeaderf print underlined header
func PrHeaderf(format string, a ...interface{}) {
	if os.Getenv("CI") == "true" {
		fmt.Printf(format, a...)
		return
	}
	fmt.Printf(head(format), a...)
}

// PrNoticef print dimmed notice
func PrNoticef(format string, a ...interface{}) {
	if os.Getenv("CI") == "true" {
		fmt.Printf(format, a...)
		return
	}
	fmt.Printf(notice(format), a...)
}

// PrWarnf print bold warning
func PrWarnf(format string, a ...interface{}) {
	if os.Getenv("CI") == "true" {
		fmt.Printf(format, a...)
		return
	}
	fmt.Printf(warn(format), a...)
}

// PrFatalf log fatal
func PrFatalf(format string, a ...interface{}) {
	if os.Getenv("CI") == "true" {
		log.Fatalf(format, a...)
		return
	}
	log.Fatalf(err(format), a...)
}

// PrResultf print green result
func PrResultf(format string, a ...interface{}) {
	if os.Getenv("CI") == "true" {
		fmt.Printf(format, a...)
		return
	}
	fmt.Printf(result(format), a...)
}

// PrettyList output list prettified for terminal
func PrettyList(str string, list []string) {
	if len(list) > 0 {
		PrHeaderf("\n%s\n", str)
		for _, s := range list {
			fmt.Printf("  - %s\n", s)
		}
		fmt.Println()
	}
}

// PrettyHash output list prettified for terminal
func PrettyHash(str string, list []Hash) {
	fmt.Println()
	data := [][]string{}
	for i, s := range list {
		data = append(data, []string{strconv.Itoa(i + 1), s.Name, s.Hash, s.Path})
	}
	PrettyTable(str, []string{"#", "NAME", "HASH", "PATH"}, data)
	fmt.Println()
}

// PrintResponse print delete response prettified
func PrintResponse(s string, a []DeleteResult) {
	if len(a) > 0 {
		fmt.Println()
		// sort result by code
		sort.Slice(a, func(i, j int) bool {
			return a[i].Code < a[j].Code
		})
		data := [][]string{}
		for i, it := range a {
			data = append(data, []string{strconv.Itoa(i + 1), it.ID, it.Type, it.Code, it.Message})
		}
		PrettyTable(s, []string{"#", "ID", "TYPE", "CODE", "MESSAGE"}, data)
		fmt.Println()
	}
}

// PrettyTable markdown style table output
func PrettyTable(caption string, header []string, data [][]string) {
	table := tablewriter.NewWriter(os.Stdout)
	table.SetHeader(header)
	table.SetAutoWrapText(false)
	table.SetBorders(tablewriter.Border{Left: true, Top: false, Right: true, Bottom: false})
	table.SetCenterSeparator("|")
	if os.Getenv("CI") == "true" {
		table.SetCaption(true, strings.ToUpper(caption))
	} else {
		table.SetCaption(true, result(strings.ToUpper(caption)))
	}
	table.AppendBulk(data) // Add Bulk Data
	table.Render()
}
