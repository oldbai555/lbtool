package topdf

import (
	"context"
	"fmt"
	"github.com/oldbai555/lbtool/pkg/html"
	"log"
	"os"

	"github.com/chromedp/cdproto/page"
	"github.com/chromedp/chromedp"
)

func Do(url string) {
	// create context
	ctx, cancel := chromedp.NewContext(context.Background())
	defer cancel()

	// capture pdf
	var buf []byte
	if err := chromedp.Run(ctx, printToPDF(fmt.Sprintf(`%s`, url), &buf)); err != nil {
		log.Fatal(err)
	}
	byUrl, err := html.GetHtmlResultByUrl(url)
	if err != nil {
		log.Fatal(err)
	}
	if err := os.WriteFile(fmt.Sprintf("%s.pdf", byUrl.Title), buf, 0o644); err != nil {
		log.Fatal(err)
	}
	fmt.Printf("write %s.pdf ok \n", byUrl.Title)
}

// print a specific pdf page.
func printToPDF(urlstr string, res *[]byte) chromedp.Tasks {
	return chromedp.Tasks{
		chromedp.Navigate(urlstr),
		chromedp.ActionFunc(func(ctx context.Context) error {
			buf, _, err := page.PrintToPDF().WithPrintBackground(false).Do(ctx)
			if err != nil {
				return err
			}
			*res = buf
			return nil
		}),
	}
}
