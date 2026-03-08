package cmd

import (
	"aozorandom-cobra/internal/read"
	"bytes"
	"fmt"
	"math/rand"
	"os"
	"strings"

	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
	Use:   "aozora [作家名または著作名]",
	Short: "青空文庫からランダムに作品を検索するCLI",
	Long:  "青空文庫からコマンドで渡された作家名または著作名からランダムに作品を検索するCLI",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		key := args[0]

		data, err := read.FetchData()
		if err != nil {
			return err
		}

		books, err := read.ParseCSV(bytes.NewReader(data))
		if err != nil {
			return err
		}

		var matched []read.Book
		for _, b := range books {
			if strings.Contains(b.Author, key) || strings.Contains(b.Title, key) {
				matched = append(matched, b)
			}
		}

		if len(matched) == 0 {
			fmt.Printf("「%s」の作品が見つかりませんでした\n", key)
			os.Exit(1)
		}

		pick := matched[rand.Intn(len(matched))]

		fmt.Printf("作家: %s\n", pick.Author)
		fmt.Printf("作品: %s\n", pick.Title)
		fmt.Printf("出版社: %s\n", pick.Publisher)

		return nil
	},
}

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		os.Exit(1)
	}
}
