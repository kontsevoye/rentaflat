package flat_parser

import (
	"errors"
	"fmt"
	"github.com/PuerkitoBio/goquery"
	"github.com/golang-module/carbon/v2"
	"github.com/kontsevoye/rentaflat/internal/config"
	"go.uber.org/zap"
	"net/http"
	"strconv"
	"strings"
)

func NewSsGeParser(logger *zap.Logger, c *config.AppConfig) *SsGeParser {
	p := &SsGeParser{
		logger:      logger,
		workerCount: c.WorkerCount,
	}

	return p
}

type SsGeParser struct {
	logger      *zap.Logger
	workerCount int
}

func trimmer(s string) string {
	trimmed := strings.ReplaceAll(s, "\n", "")
	trimmed = strings.ReplaceAll(trimmed, " ", "")

	return trimmed
}

func parseSingleFlat(url string) (*Flat, error) {
	res, err := http.Get(url)
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()
	if res.StatusCode != 200 {
		return nil, fmt.Errorf("status code error: %d %s", res.StatusCode, res.Status)
	}

	doc, err := goquery.NewDocumentFromReader(res.Body)
	if err != nil {
		return nil, err
	}

	paramBlockNodes := doc.Find(".EAchParamsBlocks text")
	rawArea := paramBlockNodes.Eq(0).Text()
	rawArea = strings.Replace(rawArea, "m²", "", 1)
	area, _ := strconv.Atoi(trimmer(rawArea))
	rawRooms := paramBlockNodes.Eq(1).Text()
	rooms, _ := strconv.Atoi(trimmer(rawRooms))
	rawFloor := paramBlockNodes.Eq(3).Text()
	splitFloor := strings.Split(trimmer(rawFloor), "/")
	floor, _ := strconv.Atoi(splitFloor[0])

	rawPrice := trimmer(doc.Find(".desktopPriceBlockDet .article_right_price").Text())
	price, _ := strconv.Atoi(rawPrice)

	rawPhone, _ := doc.Find(".phone-row-top .UserMObileNumbersBlock a").Attr("href")
	phone := strings.Replace(rawPhone, "tel:", "", 1)

	photos := doc.Find(".slider-full-content .carousel-image img").Map(func(i int, s *goquery.Selection) string {
		src, _ := s.Attr("src")

		return src
	})

	rawTime := trimmer(doc.Find(".add_date_block").Text())
	pubTime := carbon.SetTimezone("Asia/Tbilisi").ParseByFormat(rawTime, "d.m.Y/H:i")

	isAgency := false
	rawContact := strings.Replace(trimmer(doc.Find(".author_type").Text()), "", "", 1)
	splitContact := strings.Split(rawContact, "allclassifiedads")
	splitContact = strings.Split(splitContact[0], "Agencyallads")
	if len(splitContact) > 1 {
		isAgency = true
	}
	splitContact = strings.Split(splitContact[0], "Agent'sallapplications")
	if len(splitContact) > 1 {
		isAgency = true
	}
	contact := splitContact[0]

	flat := &Flat{
		Id:          trimmer(doc.Find(".article_item_id span").Text()),
		Url:         url,
		PhotoUrls:   photos,
		Title:       doc.Find(".article_in_title h1").Text(),
		Description: doc.Find(".article_item_desc_body .details_text").Text(),
		Area:        area,
		Rooms:       rooms,
		Floor:       floor,
		Price:       price,
		ContactName: contact,
		Phone:       phone,
		IsAgency:    isAgency,
		PublishedAt: pubTime.Carbon2Time(),
	}

	return flat, nil
}

func parseFlatList(url string) ([]string, error) {
	res, err := http.Get(url)
	if err != nil {
		return []string{}, err
	}
	defer res.Body.Close()
	if res.StatusCode != 200 {
		return []string{}, fmt.Errorf("status code error: %d %s", res.StatusCode, res.Status)
	}

	doc, err := goquery.NewDocumentFromReader(res.Body)
	if err != nil {
		return []string{}, err
	}

	urls := doc.Find(".latest_article_each").Map(func(i int, s *goquery.Selection) string {
		rawId, _ := s.Attr("data-id")

		return "https://ss.ge/en/real-estate/" + rawId
	})

	return urls, nil
}

func (p *SsGeParser) Supports(url string) bool {
	return strings.Contains(url, "ss.ge/en/real-estate/l/")
}

func (p *SsGeParser) Parse(url string) (<-chan Flat, <-chan error) {
	results := make(chan Flat, 20)
	errs := make(chan error, 20)

	go func(results chan<- Flat, errs chan<- error) {
		defer func() {
			p.logger.Debug("closing result channels")
			close(results)
			close(errs)
		}()
		if !p.Supports(url) {
			errs <- errors.New("Unsupported url: " + url)
			return
		}
		p.logger.Debug("flat list fetching started")
		urls, err := parseFlatList(url)
		p.logger.Debug("flat list fetching finished")

		if err != nil {
			errs <- err
			return
		}

		urlsCount := len(urls)
		jobs := make(chan string, urlsCount)
		internalResults := make(chan Flat)
		internalErrs := make(chan error)

		for w := 1; w <= p.workerCount; w++ {
			go func(urls <-chan string, results chan<- Flat, errs chan<- error) {
				for url := range urls {
					p.logger.Debug("worker started job " + url)
					flat, err := parseSingleFlat(url)
					if err != nil {
						errs <- err
						return
					}
					results <- *flat
					p.logger.Debug("worker finished job " + url)
				}
			}(jobs, internalResults, internalErrs)
		}

		for _, url := range urls {
			jobs <- url
		}
		close(jobs)

		remainingResults := urlsCount
		for range urls {
			p.logger.Debug(strconv.Itoa(remainingResults) + " jobs left")
			select {
			case flat := <-internalResults:
				remainingResults -= 1
				results <- flat
			case err := <-internalErrs:
				errs <- err
			}
		}
	}(results, errs)

	return results, errs
}
