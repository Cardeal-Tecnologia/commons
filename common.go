package common

import (
	"fmt"
	"log"
	"regexp"
	"strconv"
	"strings"

	"github.com/ztrue/tracerr"

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
func InsertAuctionToDatabase(auction *Auction, property *Property, rounds *[]Round) {
	if auction != nil && property != nil && (rounds != nil && len(*rounds) > 0) {
		// insere no banco de dados
	}
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
