package ss_ge

import (
	"github.com/PuerkitoBio/goquery"
	"github.com/golang-module/carbon/v2"
	"github.com/kontsevoye/rentaflat/cmd/parser"
	"go.uber.org/zap"
	"net/http"
	"strconv"
	"strings"
)

func NewParser(logger *zap.SugaredLogger) *Parser {
	p := &Parser{
		logger: logger,
	}

	return p
}

type Parser struct {
	logger *zap.SugaredLogger
}

func trimmer(s string) string {
	trimmed := strings.ReplaceAll(s, "\n", "")
	trimmed = strings.ReplaceAll(trimmed, " ", "")

	return trimmed
}

func (p *Parser) Parse(url string) []parser.Flat {
	res, err := http.Get(url)
	if err != nil {
		p.logger.Fatal(err)
	}
	defer res.Body.Close()
	if res.StatusCode != 200 {
		p.logger.Fatalf("status code error: %d %s", res.StatusCode, res.Status)
	}

	doc, err := goquery.NewDocumentFromReader(res.Body)
	if err != nil {
		p.logger.Fatal(err)
	}

	var flats []parser.Flat
	urls := doc.Find(".latest_article_each").Map(func(i int, s *goquery.Selection) string {
		rawId, _ := s.Attr("data-id")

		return "https://ss.ge/en/real-estate/" + rawId
	})
	queue := make(chan parser.Flat, len(urls))
	defer close(queue)
	for _, url := range urls {
		go func(url string) {
			res, err := http.Get(url)
			if err != nil {
				p.logger.Fatal(err)
			}
			defer res.Body.Close()
			if res.StatusCode != 200 {
				p.logger.Fatalf("status code error: %d %s", res.StatusCode, res.Status)
			}

			doc, err := goquery.NewDocumentFromReader(res.Body)
			if err != nil {
				p.logger.Fatal(err)
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

			queue <- parser.Flat{
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
		}(url)
	}

	for range urls {
		flat := <-queue
		flats = append(flats, flat)
	}

	return flats
}
