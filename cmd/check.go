package cmd

import (
	"fmt"
	"github.com/spf13/cobra"
	"github.com/your-username/dead-link-checker/internal"
)

// checkCmd represents the check command
var checkCmd = &cobra.Command{
	Use:   "check <url>",
	Short: "Crawl a website and detect broken links",
	Long: `Crawl a website starting from the specified URL and detect all broken links.

The crawler will:
• Follow internal links recursively up to the specified depth
• Check all discovered links for HTTP errors
• Report broken links with their status codes
• Preserve relative paths for complete site coverage

Use the --depth flag to control how deep the crawler goes into your site.`,
	Example: `  # Check homepage and one level deep
  dead-link-checker check https://example.com -d 1
  
  # Check with default depth (2 levels)
  dead-link-checker check https://example.com
  
  # Deep crawl (5 levels)
  dead-link-checker check https://mysite.com --depth 5`,
	Args: cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		// Get flag depth
		depth, err := cmd.Flags().GetInt("depth")
		if err != nil {
			fmt.Println("Error retreiving depth")
		}
		// Get url and pass into crawler then retreive deadlinks
		url := args[0]
		fmt.Println("Checking " + url)
		siteURLs := internal.CrawlSite(url, depth)
		fmt.Println("Collecting dead URLs:")
		deadURLs := internal.GetDeadLinks(&siteURLs)
		for _, deadURL := range deadURLs {
			fmt.Println(deadURL)
		}
		fmt.Println("check called")
	},
}

func init() {
	rootCmd.AddCommand(checkCmd)
	// Depth flag
	checkCmd.Flags().IntP("depth", "d", 2, "Maximum crawl depth")
}
