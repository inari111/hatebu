package cmd

import (
	"encoding/xml"
	"fmt"
	"io/ioutil"
	"net/http"
	"strconv"
	"strings"
	"unicode/utf8"

	"github.com/fatih/color"
	"github.com/mattn/go-runewidth"
	"github.com/spf13/cobra"
)

type HotEntry struct {
	Items []*Item `xml:"item"`
}

type Item struct {
	Title         string `xml:"title"`
	Link          string `xml:"link"`
	Description   string `xml:"description"`
	Date          string `xml:"date"`
	BookmarkCount int    `xml:"bookmarkcount"`
}

var RootCmd = &cobra.Command{
	Use:   "hatebu",
	Short: "hatebu is a CLI",
	Long:  "hatebu is a CLI",
	Run: func(cmd *cobra.Command, args []string) {
		data := httpGet("http://b.hatena.ne.jp/hotentry/it.rss")

		result := HotEntry{}
		err := xml.Unmarshal([]byte(data), &result)
		if err != nil {
			fmt.Printf("error: %v", err)
			return
		}

		bookmarkWidth := 8
		bookmarkFmt := fmt.Sprintf("%%-%ds", bookmarkWidth)

		titleWidth := maxTitleWidth(result.Items)
		titleFmt := fmt.Sprintf("%%-%ds", titleWidth)

		urlWidth := maxURLWidth(result.Items)
		urlFmt := fmt.Sprintf("%%-%ds", urlWidth)

		fmt.Fprintf(color.Output, " %s | %s | %s \n",
			color.YellowString(fmt.Sprintf(bookmarkFmt, "Bookmark")),
			color.BlueString(titleFmt, "Title"),
			fmt.Sprintf(urlFmt, "Url"),
		)

		fmt.Println(strings.Repeat("-", bookmarkWidth+titleWidth+urlWidth))
		for _, bookmark := range result.Items {
			title := bookmark.Title
			link := bookmark.Link
			fmt.Fprintf(
				color.Output,
				" %s | %s | %s \n",
				color.YellowString(fmt.Sprintf(bookmarkFmt, strconv.Itoa(bookmark.BookmarkCount))),
				color.CyanString(runewidth.FillRight(title, titleWidth)),
				fmt.Sprintf(urlFmt, link),
			)
		}

	},
}

func init() {
	cobra.OnInitialize()
}

func httpGet(url string) string {
	response, _ := http.Get(url)
	body, _ := ioutil.ReadAll(response.Body)
	defer response.Body.Close()
	return string(body)
}

func maxTitleWidth(entries []*Item) int {
	width := 0
	for _, e := range entries {
		count := runewidth.StringWidth(e.Title)
		if count > width {
			width = count
		}
	}
	return width
}

func maxURLWidth(entries []*Item) int {
	width := 0
	for _, e := range entries {
		count := utf8.RuneCountInString(e.Link)
		if count > width {
			width = count
		}
	}
	if width > 100 {
		return 100
	}
	return width
}
