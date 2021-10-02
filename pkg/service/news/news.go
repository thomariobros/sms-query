package news

import (
	"github.com/mmcdole/gofeed"

	"sms-query/pkg/util"
)

func Rss(rss string, max int) ([]gofeed.Item, error) {
	fp := gofeed.NewParser()
	feed, err := fp.ParseURL(rss)
	if err != nil {
		return nil, err
	}
	if feed.Items != nil {
		var result []gofeed.Item
		for index, item := range feed.Items {
			if index == max {
				break
			}
			// standardize title and description
			item.Title = util.StandardizeString(item.Title, false)
			item.Description = util.StandardizeString(item.Description, false)
			result = append(result, *item)
		}
		return result, nil
	}
	return nil, nil
}
