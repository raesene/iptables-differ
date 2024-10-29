package main

import (
	"bufio"
	"flag"
	"fmt"
	"os"
	"strings"

	"github.com/fatih/color"
)

const usage = `IPTables Rules Comparison Utility

This tool compares two sets of iptables rules and shows the differences between them.

To generate the input files, use the following iptables commands:
  Before changes: iptables-save > rules-before.txt
  After changes:  iptables-save > rules-after.txt

Usage:
  iptables-diff -before <before-file> -after <after-file>

Example:
  iptables-diff -before rules-before.txt -after rules-after.txt

Output will be color-coded:
  - Red for removed rules
  - Green for added rules
  - Yellow for table changes
`

// Color printers
var (
	red    = color.New(color.FgRed).PrintfFunc()
	green  = color.New(color.FgGreen).PrintfFunc()
	yellow = color.New(color.FgYellow).PrintfFunc()
	cyan   = color.New(color.FgCyan).PrintfFunc()
)

func main() {
	beforeFile := flag.String("before", "", "File containing the initial iptables rules")
	afterFile := flag.String("after", "", "File containing the modified iptables rules")
	noColor := flag.Bool("no-color", false, "Disable color output")
	flag.Usage = func() {
		fmt.Fprintf(os.Stderr, usage)
	}
	flag.Parse()

	if *beforeFile == "" || *afterFile == "" {
		flag.Usage()
		os.Exit(1)
	}

	// Handle color disable flag
	if *noColor {
		color.NoColor = true
	}

	before, err := loadRules(*beforeFile)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error loading before rules: %v\n", err)
		os.Exit(1)
	}

	after, err := loadRules(*afterFile)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error loading after rules: %v\n", err)
		os.Exit(1)
	}

	compareRules(before, after)
}

func loadRules(filename string) (map[string][]string, error) {
	file, err := os.Open(filename)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	rules := make(map[string][]string)
	var currentChain string
	scanner := bufio.NewScanner(file)

	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())
		if line == "" || strings.HasPrefix(line, "#") {
			continue
		}

		if strings.HasPrefix(line, "*") {
			// Table declaration (e.g., *filter)
			currentChain = line
			rules[currentChain] = []string{}
		} else if strings.HasPrefix(line, "COMMIT") {
			currentChain = ""
		} else {
			if currentChain != "" {
				rules[currentChain] = append(rules[currentChain], line)
			}
		}
	}

	return rules, scanner.Err()
}

func compareRules(before, after map[string][]string) {
	// Compare tables
	allTables := make(map[string]bool)
	for table := range before {
		allTables[table] = true
	}
	for table := range after {
		allTables[table] = true
	}

	for table := range allTables {
		cyan("\nComparing table %s:\n", table)

		beforeRules := before[table]
		afterRules := after[table]

		if _, exists := before[table]; !exists {
			yellow("! Table %s has been added\n", table)
			continue
		}

		if _, exists := after[table]; !exists {
			yellow("! Table %s has been removed\n", table)
			continue
		}

		// Find removed rules
		for _, rule := range beforeRules {
			if !containsRule(afterRules, rule) {
				red("- %s\n", rule)
			}
		}

		// Find added rules
		for _, rule := range afterRules {
			if !containsRule(beforeRules, rule) {
				green("+ %s\n", rule)
			}
		}
	}
}

func containsRule(rules []string, rule string) bool {
	for _, r := range rules {
		if r == rule {
			return true
		}
	}
	return false
}
