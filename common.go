package common

import (
	"fmt"
	"log"
	"os"
	"regexp"
	"strconv"
	"strings"
	"time"

	"github.com/ztrue/tracerr"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"

	"github.com/gocolly/colly"
)

// função para extrair dominio de uma url
// exemplo: https://www.google.com -> www.google.com
func ExtractDomainList(urls []string) []string {
	domains := []string{}
	for _, url := range urls {
		domain := strings.Replace(url, "https://", "", -1)
		domain = strings.Replace(domain, "http://", "", -1)
		domain = strings.Split(domain, "/")[0]
		domains = append(domains, domain)
	}
	return domains
}

// função pra visitar uma lista de urls e esperar a coleta terminar
func VisitPageUrls(collector *colly.Collector, urls []string) {
	for _, url := range urls {
		collector.Visit(url)
	}
	collector.Wait()
}

func MoneyStringToFloat(moneyString string) float64 {
	moneyString = strings.Replace(moneyString, "R$", "", -1)
	moneyString = strings.Replace(moneyString, ".", "", -1)
	moneyString = strings.Replace(moneyString, ",", ".", -1)
	moneyString = strings.TrimSpace(moneyString)
	money, _ := strconv.ParseFloat(moneyString, 64)
	return money
}

func GetProcessNumber(text string) string {
	// exemplo de texto: <b>processo: </b>1006075-64.2019.8.26.0554 asdasdasd
	// exemplo de processo: 1006075-64.2019.8.26.0554

	processNumber := regexp.MustCompile(`\d{7}-\d{2}\.\d{4}\.\d{1}\.\d{2}\.\d{4}`).FindString(text)
	return processNumber
}

// função para inserir um leilão no banco de dados
func InsertAuctionToDatabase(auction *Auction, property *Property, rounds *[]Round) bool {
	if auction != nil && property != nil && (rounds != nil && len(*rounds) > 0) {
		// conecta ao banco de dados
		db := ConnectToDatabase()
		if db == nil {
			return false
		}

		// insere a property
		result := db.Clauses(clause.OnConflict{
			Columns: []clause.Column{{Name: "street_name"}, {Name: "street_number"}, {Name: "neighborhood"}, {Name: "city"}, {Name: "usage_type"}, {Name: "size"}, {Name: "postal_code"}, {Name: "bedrooms"}, {Name: "bathroom"}, {Name: "garage"}},
			DoUpdates: clause.AssignmentColumns([]string{
				"bedrooms",
				"size",
				"garage",
				"bathroom",
				"floor",
				"neighborhood",
				"city",
				"latitude",
				"longitude",
				"street_number",
				"street_name",
				"postal_code",
				"updated_at",
			}),
		}).Create(&property)
		if result.Error != nil {
			return false
		}

		// insere a auction
		auction.PropertyID = int(property.Id)
		result = db.Clauses(clause.OnConflict{
			Columns: []clause.Column{{Name: "origin"}, {Name: "external_id"}},
			DoUpdates: clause.AssignmentColumns([]string{
				"title",
				"updated_at",
				"external_url",
				"auctioneer_comission",
				"auctioneer_views",
				"price_sold",
				"qualified_users",
				"status",
				"description",
				"views_count",
			}),
		}).Create(&auction)
		if result.Error != nil {
			return false
		}

		for _, round := range *rounds {

			// insere o round
			round.AuctionId = auction.Id
			if !(round.RoundNumber == 0 && round.MinPrice == 0) {
				result = db.Clauses(clause.OnConflict{
					Columns: []clause.Column{{Name: "auction_id"}, {Name: "round_number"}},
					DoUpdates: clause.AssignmentColumns([]string{
						"discount",
						"end_date",
						"start_date",
						"increment_value",
						"min_price",
						"round_number",
						"updated_at",
					}),
				}).Create(&round)
				if result.Error != nil {
					return false
				}
			}
		}

		return true
	}

	return false
}

// conecta ao banco de dados
func ConnectToDatabase() *gorm.DB {
	fmt.Println("Conectando ao banco de dados...")
	dsn := os.Getenv("DATABASE_URL")
	if dsn == "" {
		dsn = "postgres://postgres:postgres@localhost:5432/cardeal_app_development" // url do docker local
	}

	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{
		NowFunc: func() time.Time {
			return time.Now().UTC()
		},
	})
	if err != nil {
		fmt.Println("Erro ao conectar ao banco de dados")
	}

	return db
}

// função para transformar texto em float
func StringToFloat(text string) float64 {
	text = strings.Replace(text, ".", "", -1)
	text = strings.Replace(text, ",", ".", -1)
	text = strings.TrimSpace(text)
	number, _ := strconv.ParseFloat(text, 64)
	return number
}

func StringToInt(text string) int {
	text = strings.TrimSpace(text)
	number, _ := strconv.Atoi(text)
	return number
}

// Função para tratar erros
func HandleError() {

	if err := recover(); err != nil {
		e := tracerr.Wrap(err.(error))
		frame := tracerr.StackTrace(e)
		for index, f := range frame {
			if index == 3 && e.Error() != "runtime error: invalid memory address or nil pointer dereference" {
				log.Printf(fmt.Sprintf(" AI CE ME QUEBRA \n==========\n erro:	%s\n file: %s\n linha: %d\n funcao: %s\n==========\n", e.Error(), f.Path, f.Line, f.Func))
			}
		}
	}
}
