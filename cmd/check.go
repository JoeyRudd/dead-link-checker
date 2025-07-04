package cmd

import (
	"fmt"
	"github.com/spf13/cobra"
	"github.com/your-username/dead-link-checker/internal"
)

// checkCmd represents the check command
var checkCmd = &cobra.Command{
	Use:   "check",
	Short: "Use check to check the website for dead links",
	Long:  `Use check *url* -d *depth* to check a website for dead links`,
	Run: func(cmd *cobra.Command, args []string) {
		// Get flag depth
		depth, err := cmd.Flags().GetInt("depth")
		if err != nil {
			fmt.Println("Error retreiving depth")
		}
		// Get url from command
		url := args[0]
		// Pass url into crawler
		siteURLs := internal.CrawlSite(url, depth)
		// deadURLs
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
	checkCmd.Flags().IntP("depth", "d", 2, "Depth of website to scrape")
	)
}
