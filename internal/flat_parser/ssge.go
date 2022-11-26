package flat_parser

import (
	"fmt"
	"github.com/PuerkitoBio/goquery"
	"github.com/golang-module/carbon/v2"
	"github.com/kontsevoye/rentaflat/internal/config"
	"go.uber.org/zap"
	"net/http"
	neturl "net/url"
	"strconv"
	"strings"
	"sync"
)

func NewSsGeParser(logger *zap.Logger, c *config.AppConfig, ff FlatFactory) *SsGeParser {
	return &SsGeParser{
		logger:      logger,
		flatFactory: ff,
		workerCount: c.WorkerCount,
	}
}

type SsGeParser struct {
	logger      *zap.Logger
	flatFactory FlatFactory
	workerCount int
}

func trimmer(s string) string {
	trimmed := strings.ReplaceAll(s, "\n", "")
	trimmed = strings.ReplaceAll(trimmed, " ", "")

	return trimmed
}

func (p *SsGeParser) parseSingleFlat(url string) (*Flat, error) {
	parsedUrl, err := neturl.Parse(url)
	if err != nil {
		return nil, err
	}

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
	var photoUrls []*neturl.URL
	for _, photo := range photos {
		parsedPhotoUrl, err := neturl.Parse(photo)
		if err != nil {
			return nil, err
		}
		photoUrls = append(photoUrls, parsedPhotoUrl)
	}

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

	flat, err := p.flatFactory.NewFlat(
		trimmer(doc.Find(".article_item_id span").Text()),
		parsedUrl,
		photoUrls,
		doc.Find(".article_in_title h1").Text(),
		doc.Find(".article_item_desc_body .details_text").Text(),
		uint(area),
		uint(rooms),
		uint(floor),
		uint(price),
		contact,
		phone,
		isAgency,
		pubTime.Carbon2Time(),
	)

	return flat, err
}

func (p *SsGeParser) parseFlatList(url string) ([]string, error) {
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

func (p *SsGeParser) generateUrl() *neturl.URL {
	url, _ := neturl.Parse("https://ss.ge/en/real-estate/l/Flat/For-Rent")
	query := url.Query()
	query.Set("Page", "1")
	query.Set("CurrencyId", "1")
	query.Set("WithImageOnly", "true")
	query.Set("WithImageOnly", "true")
	query.Set("Sort.SortExpression", "\"OrderDate\" DESC")
	url.RawQuery = query.Encode()

	return url
}

func (p *SsGeParser) work(flatListUrl *neturl.URL, request Request, jobUrls chan<- string, results chan<- Flat) error {
	p.logger.Debug(
		"flat list fetching started",
		zap.String("url", flatListUrl.String()),
		zap.String("lastId", request.LastId),
	)
	flatUrls, err := p.parseFlatList(flatListUrl.String())
	p.logger.Debug(
		"flat list fetching finished",
		zap.String("url", flatListUrl.String()),
		zap.String("lastId", request.LastId),
	)
	if err != nil {
		return err
	}

	for _, flatUrl := range flatUrls {
		if request.LastId != "0" && strings.Contains(flatUrl, request.LastId) {
			p.logger.Debug(
				"nothing new for me",
				zap.String("url", flatListUrl.String()),
				zap.String("lastId", request.LastId),
				zap.String("flatUrl", flatUrl),
			)
			return nil
		}
		jobUrls <- flatUrl
	}

	containsLastId := false
	if request.LastId == "0" {
		containsLastId = true
	} else {
		for _, flatUrl := range flatUrls {
			if strings.Contains(flatUrl, request.LastId) {
				containsLastId = true
				break
			}
		}
	}

	if !containsLastId {
		query := flatListUrl.Query()
		currentPage, err := strconv.Atoi(query.Get("Page"))
		if err != nil {
			return err
		}

		query.Set("Page", strconv.Itoa(currentPage+1))
		flatListUrl.RawQuery = query.Encode()
		err = p.work(flatListUrl, request, jobUrls, results)
		if err != nil {
			return err
		}
	}

	return nil
}

func (p *SsGeParser) spawnWorker(workerId int, wg *sync.WaitGroup, jobs <-chan string, results chan<- Flat, errs chan<- error) {
	defer wg.Done()
	p.logger.Debug("worker spawned", zap.Int("id", workerId))
	for url := range jobs {
		p.logger.Debug("worker started job", zap.Int("id", workerId), zap.String("url", url))
		flat, err := p.parseSingleFlat(url)
		if err != nil {
			errs <- err
			return
		}
		results <- *flat
		p.logger.Debug("worker finished job", zap.Int("id", workerId), zap.String("url", url))
	}
	p.logger.Debug("worker stopped", zap.Int("id", workerId))
}

func (p *SsGeParser) Parse(request Request) (<-chan Flat, <-chan error) {
	results := make(chan Flat, 20)
	errs := make(chan error, 20)
	jobs := make(chan string, 20)

	url := p.generateUrl()
	p.logger.Debug("url generated", zap.String("url", url.String()))

	wg := &sync.WaitGroup{}
	for w := 1; w <= p.workerCount; w++ {
		wg.Add(1)
		go p.spawnWorker(w, wg, jobs, results, errs)
	}

	go func() {
		err := p.work(url, request, jobs, results)
		close(jobs)
		if err != nil {
			errs <- err
		}

		wg.Wait()
		p.logger.Debug("closing results & errs channels")
		close(results)
		close(errs)
	}()

	return results, errs
}
