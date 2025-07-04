package cmd

import (
	"os"

	"github.com/spf13/cobra"
)

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:   "dead-link-checker [command]",
	Short: "Fast, reliable dead link detection for websites",
	Long: `Dead Link Checker is a command-line tool that crawls websites to detect broken links.

It recursively follows internal links to map your entire site structure and identifies:
• HTTP errors (404, 500, etc.)
• Timeout and network issues  
• Unreachable resources

Perfect for web developers, SEO audits, and site maintenance.`,
	Example: `  # Check a website with default depth (2 levels)
  dead-link-checker check https://example.com
  
  # Check with custom depth
  dead-link-checker check https://example.com -d 5`,
}

func Execute() {
	err := rootCmd.Execute()
	if err != nil {
		os.Exit(1)
	}
}

func init() {
	
}
