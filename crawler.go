package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"strconv"
	"time"

	"github.com/chromedp/cdproto/browser"
	"github.com/chromedp/chromedp"
)

type crawler struct {
	downloadTimeout   time.Duration
	collectionTimeout time.Duration
	timeBetweenSteps  time.Duration
	output            string
}

func (c crawler) crawl() error {
	// Pegar variáveis de ambiente

	// Chromedp setup.
	log.SetOutput(os.Stderr) // Enviando logs para o stderr para não afetar a execução do coletor.

	alloc, allocCancel := chromedp.NewExecAllocator(
		context.Background(),
		append(chromedp.DefaultExecAllocatorOptions[:],
			chromedp.Flag("headless", true), // mude para false para executar com navegador visível.
			chromedp.NoSandbox,
			chromedp.DisableGPU,
			chromedp.WindowSize(1920, 1080), // init with a desktop view

		)...,
	)
	defer allocCancel()

	ctx, cancel := chromedp.NewContext(
		alloc,
		chromedp.WithLogf(log.Printf), // remover comentário para depurar
	)
	defer cancel()

	url := "https://www.ribeiraopreto.sp.gov.br/portal/administracao/pesquisa-de-notas-fiscais"
	log.Printf("Acessando o site (%s)...", url)
	if err := c.acessaSite(ctx, url); err != nil {
		log.Fatalf("%w", err)
	}
	log.Printf("Site carregado com sucesso!\n")

	if err := c.realizaDownload(ctx); err != nil {
		log.Fatalf("%w", err)
	}
	log.Printf("Download realizado com sucesso!\n")

	return nil
}

func (c crawler) downloadPdf(n int) string {
	return filepath.Join(c.output, fmt.Sprintf("nf-ribeirao-preto-%d.pdf", n))
}
func (c crawler) downloadXml(n int) string {
	return filepath.Join(c.output, fmt.Sprintf("nf-ribeirao-preto-%d.xml", n))
}

// nomeiaDownload dá um nome ao último arquivo modificado dentro do diretório
// passado como parâmetro nomeiaDownload dá pega um arquivo
func nomeiaDownload(output, fName string) error {
	// Identifica qual foi o ultimo arquivo
	files, err := os.ReadDir(output)
	if err != nil {
		log.Fatalf("erro lendo diretório %s: %v", output, err)
	}
	var newestFPath string
	var newestTime int64 = 0
	for _, f := range files {
		fPath := filepath.Join(output, f.Name())
		fi, err := os.Stat(fPath)
		if err != nil {
			log.Fatalf("erro obtendo informações sobre arquivo %s: %v", fPath, err)
		}
		currTime := fi.ModTime().Unix()
		if currTime > newestTime {
			newestTime = currTime
			newestFPath = fPath
		}
	}
	// Renomeia o ultimo arquivo modificado.
	if err := os.Rename(newestFPath, fName); err != nil {
		log.Fatalf("sem planilhas baixadas: %w", err)
	}
	return nil
}

func (c crawler) acessaSite(ctx context.Context, url string) error {
	if err := chromedp.Run(ctx,
		// Acessa o site
		chromedp.Navigate(url),
		chromedp.Sleep(c.timeBetweenSteps),
		// Realiza a busca por fornecedor
		chromedp.SetValue(`//*[@id="vGCMVAN_NOM_FORNEC"]`, "VERSAO BR"),
		chromedp.Sleep(c.timeBetweenSteps),
		chromedp.DoubleClick(`//*[@id="TABLE6"]/tbody/tr/td/span[1]/span`),
		chromedp.Sleep(c.timeBetweenSteps),

		// Altera o diretório de download
		browser.SetDownloadBehavior(browser.SetDownloadBehaviorBehaviorAllowAndName).
			WithDownloadPath(c.output).
			WithEventsEnabled(true),
	); err != nil {
		log.Fatal(err)
	}

	return nil
}

func (c crawler) realizaDownload(ctx context.Context) error {
	// var nodes []*cdp.Node
	// var buf []byte
	count := 2474
	for i := 619; i <= 559; i++ {
		// Acessando a página
		if err := chromedp.Run(ctx,
			chromedp.SetValue(`//*[@id="vCURRENTPAGE"]`, strconv.Itoa(i)),
			chromedp.Sleep(c.timeBetweenSteps),
		); err != nil {
			log.Fatal(err)
		}
		// Acessando itens
		for j := 2; j <= 2; j++ {
			xpath := fmt.Sprintf(`//*[@id="vDISPLAY_000%d"]`, j)
			// Download do PDF
			if err := chromedp.Run(ctx,
				chromedp.DoubleClick(xpath),
				chromedp.Sleep(c.timeBetweenSteps),
				chromedp.DoubleClick(`/html/body/form/table/tbody/tr[4]/td/table/tbody/tr[2]/td/table/tbody/tr[1]/td/table/tbody/tr/td/span[1]/span`, chromedp.NodeVisible),
				chromedp.Sleep(c.downloadTimeout),
				// chromedp.Nodes(`/html/body/form/table/tbody/tr[4]/td/table/tbody/tr[2]/td/table/tbody/tr[1]/td/table/tbody/tr/td/span[2]/span`, &nodes, chromedp.AtLeast(0)),
			); err != nil {
				log.Fatal(err)
			}

			pdf := c.downloadPdf(count)
			if err := nomeiaDownload(c.output, pdf); err != nil {
				log.Fatal(err)
			}

			if _, err := os.Stat(pdf); os.IsNotExist(err) {
				log.Fatalf("download do arquivo de %s não realizado", pdf)
			}
			log.Println("OK: ", pdf)
			// fmt.Printf("", len(nodes))
			// if len(nodes) != 0 {
			// Download do XML
			if err := chromedp.Run(ctx,
				chromedp.DoubleClick(`/html/body/form/table/tbody/tr[4]/td/table/tbody/tr[2]/td/table/tbody/tr[1]/td/table/tbody/tr/td/span[2]/span`),
				chromedp.Sleep(c.downloadTimeout),
				chromedp.DoubleClick(`/html/body/form/table/tbody/tr[4]/td/table/tbody/tr[2]/td/table/tbody/tr[1]/td/table/tbody/tr/td/span[3]/span`),
				chromedp.Sleep(c.timeBetweenSteps),
			); err != nil {
				log.Fatal(err)
			}

			xml := c.downloadXml(count)
			if err := nomeiaDownload(c.output, xml); err != nil {
				log.Fatal(err)
			}

			if _, err := os.Stat(xml); os.IsNotExist(err) {
				log.Fatalf("download do arquivo de %s não realizado", xml)
			}
			log.Println("OK: ", xml)
			// // } else {
			// // if err := chromedp.Run(ctx,
			// // 	chromedp.DoubleClick(`/html/body/form/table/tbody/tr[4]/td/table/tbody/tr[2]/td/table/tbody/tr[1]/td/table/tbody/tr/td/span[3]/span`),
			// // 	chromedp.Sleep(c.timeBetweenSteps),
			// // 	chromedp.FullScreenshot(&buf, 90),
			// // ); err != nil {
			// // 	log.Fatal(err)
			// // }
			// // 	if err := os.WriteFile("fullScreenshot.png", buf, 0o644); err != nil {
			// // 		log.Fatal(err)
			// // 	}
			// // }
			// // break
			count++
		}
		break
	}
	return nil
}
