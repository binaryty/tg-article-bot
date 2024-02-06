package fetcher

import (
	"context"
	"fmt"
	"github.com/PuerkitoBio/goquery"
	"github.com/binaryty/tg-bot/internal/models"
	"log"
	"net/http"
	"net/url"
	"sync"
)

const sourceUrl = "https://habr.com/ru/hubs/go/articles/page"

type ArticleSaver interface {
	Save(context.Context, models.Article) error
}

type Fetcher struct {
	articles []models.Article
	storage  ArticleSaver
}

func New(s ArticleSaver) *Fetcher {
	arts := make([]models.Article, 0)

	return &Fetcher{
		articles: arts,
		storage:  s,
	}
}

// Start fetcher.
func (f *Fetcher) Start(ctx context.Context) error {
	if err := f.Fetch(ctx); err != nil {
		return err
	}

	return nil
}

// Fetch articles.
func (f *Fetcher) Fetch(ctx context.Context) error {

	var wg sync.WaitGroup

	for i := 1; i < 50; i++ {
		wg.Add(1)

		link := fmt.Sprintf("%s%d", sourceUrl, i)

		go func(source string) {
			defer wg.Done()

			err := f.fetch(link)
			if err != nil {
				log.Printf("[ERROR] can't fetch source %s: %v", link, err)
				return
			}
		}(link)
	}

	wg.Wait()

	return nil
}

// fetch from url.
func (f *Fetcher) fetch(link string) error {
	resp, err := http.Get(link)
	if err != nil {
		return err
	}

	defer func() { _ = resp.Body.Close() }()

	if resp.StatusCode != http.StatusOK {
		log.Printf("[ERROR] can't pase %s status: %s", link, resp.Status)
		return fmt.Errorf("[ERROR] can't parse %s status: %s", link, resp.Status)
	}

	doc, err := goquery.NewDocumentFromReader(resp.Body)
	if err != nil {
		return err
	}

	f.parseDoc(doc)

	return nil
}

// parseDoc parse a HTML document.
func (f *Fetcher) parseDoc(doc *goquery.Document) {

	doc.Find(".tm-articles-list__item").Each(func(i int, s *goquery.Selection) {
		title, _ := s.Find("h2").Find("span").Html()
		l, _ := s.Find("h2").Find("a").Attr("href")
		thumbUrl, _ := s.Find(".tm-article-body").Find("img").Attr("src")
		published, _ := s.Find("time").Attr("title")

		link, err := url.JoinPath("https://habr.com/", l)
		if err != nil {
			log.Printf("[ERROR] cant't join url: %v", err)
		}

		article := models.Article{
			Title:       title,
			Link:        link,
			ThumbUrl:    thumbUrl,
			PublishedAt: published,
		}

		f.articles = append(f.articles, article)
	})
}

// ProcessArticles save parsed articles to storage.
func (f *Fetcher) ProcessArticles(ctx context.Context) error {
	for _, art := range f.articles {
		if err := f.storage.Save(ctx, art); err != nil {
			return err
		}
	}

	return nil
}
