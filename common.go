package common

import (
	"fmt"
	"log"
	"net/http"
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

// função para reindexar no ElasticSearch
func reindexToElasticSearch(id string, model string) {
	apiUrl := os.Getenv("API_URL")
	if apiUrl == "" {
		apiUrl = "http://localhost:3000/" // url da api local
	}

	http.Get(apiUrl + "reindex_to_elasticsearch/" + model + "/" + id)
}

// função para inserir um leilão no banco de dados
func InsertAuctionToDatabase(auction *Auction, property *Property, rounds *[]Round, db *gorm.DB, images []Attachment, attachments []Attachment) bool {
	if auction != nil && property != nil {
		if db == nil {
			return false
		}

		// insere a property
		result := db.Clauses(clause.OnConflict{
			Columns:   []clause.Column{{Name: "street_name"}, {Name: "street_number"}, {Name: "neighborhood"}, {Name: "city"}, {Name: "usage_type"}, {Name: "size"}, {Name: "postal_code"}, {Name: "bedrooms"}, {Name: "bathroom"}, {Name: "garage"}},
			DoUpdates: clause.AssignmentColumns([]string{
				"updated_at",
			}),
		}).Create(&property)
		if result.Error != nil {
			return false
		}

		// insere a auction
		auction.PropertyID = property.Id
		result = db.Clauses(clause.OnConflict{
			Columns: []clause.Column{{Name: "external_url"}},
			DoUpdates: clause.AssignmentColumns([]string{
				"updated_at",
				"qualified_users",
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
					Columns:   []clause.Column{{Name: "auction_id"}, {Name: "round_number"}},
					DoNothing: true,
				}).Create(&round)
				if result.Error != nil {
					return false
				}
			}
		}

		// Pós processamento
		reindexToElasticSearch(strconv.Itoa(int(auction.Id)), "Auction")
		UploadImages(property.Id, images)
		UploadAttachments(auction.Id, attachments)

		return true
	} else {
		fmt.Println("Auction or Property is nil", auction.ExternalUrl)
		return false
	}

}

// fazer upload de imagens na property
func UploadImages(id uint, images []Attachment) {
	apiUrl := os.Getenv("API_URL")
	if apiUrl == "" {
		apiUrl = "http://localhost:3000/" // url da api local
	}

	for _, image := range images {

		http.Post(apiUrl+"upload_property_image/"+strconv.Itoa(int(id)), "application/json", strings.NewReader(`{"image_url": "`+image.Url+`", "image_name": "`+image.Name+`"}`))
	}
}

// fazer upload de attachments na auction
func UploadAttachments(id uint, attachments []Attachment) {
	apiUrl := os.Getenv("API_URL")
	if apiUrl == "" {
		apiUrl = "http://localhost:3000/" // url da api local
	}

	for _, attachment := range attachments {

		http.Post(apiUrl+"upload_auction_attachment/"+strconv.Itoa(int(id)), "application/json", strings.NewReader(`{"attachment_url": "`+attachment.Url+`", "attachment_name": "`+attachment.Name+`"}`))
	}
}

// conecta ao banco de dados
func ConnectToDatabase(dsn string) *gorm.DB {
	fmt.Println("Conectando ao banco de dados...")
	if dsn == "" {
		dsn = os.Getenv("DATABASE_URL")
	}
	if dsn == "" {
		dsn = "postgresql://cardeal:cardeal@localhost:5433/cardeal" // url do docker local
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
