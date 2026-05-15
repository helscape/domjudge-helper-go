package main

import (
	"flag"
	"fmt"
	"os"
	"os/signal"
	"sync"
	"syscall"
	"time"

	"github.com/helscape/domjudge-helper-go/config"
	"github.com/helscape/domjudge-helper-go/logger"
	"github.com/helscape/domjudge-helper-go/scanner"
	"github.com/helscape/domjudge-helper-go/uploader"
)

type probResult struct {
	id     string
	name   string
	total  int
	failed int
}

type probMeta struct {
	id   string
	name string
	tcs  []scanner.Testcase
}

func main() {
	configPath := flag.String("config", "config.json", "path to config file")
	flag.Parse()

	cfg, err := config.Load(*configPath)
	if err != nil {
		logger.New().Fatal(fmt.Sprintf("config error: %v", err))
	}

	log := logger.New()

	log.Header(cfg.ServerURL, "Browser (Headless Chromium)", cfg.TestcasesDir)

	client := uploader.NewClientWithTimeout(
		cfg.ServerURL,
		cfg.AdminUsername,
		cfg.AdminPassword,
		cfg.TimeoutSeconds,
	)
	defer client.CloseBrowser()

	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)
	go func() {
		<-sigChan
		log.Interrupt("Interrupt received (Ctrl+C). Force-closing Chromium...")
		client.CloseBrowser()
		os.Exit(130)
	}()

	log.Info("Scanning for test cases...")

	scanned, err := scanner.Scan(cfg.TestcasesDir)
	if err != nil {
		log.Fatal(fmt.Sprintf("scan error: %v", err))
	}
	if len(scanned) == 0 {
		log.Fatal(fmt.Sprintf("no problems found in %s", cfg.TestcasesDir))
	}

	log.Info(fmt.Sprintf("  Found %d problem(s).", len(scanned)))
	log.Info("\nResolving problem names...")

	var uploadQueue []probMeta
	for _, prob := range scanned {
		name, err := client.GetProblemName(prob.ID)
		if err != nil {
			log.Warn(fmt.Sprintf("[%s] -> could not resolve name: %v", prob.ID, err))
			name = prob.ID
		} else {
			log.Success(fmt.Sprintf("[%s] -> %s", prob.ID, name))
		}
		uploadQueue = append(uploadQueue, probMeta{
			id:   prob.ID,
			name: name,
			tcs:  prob.Testcases,
		})
	}

	// Siapkan state untuk live renderer
	items := make([]logger.ProgressItem, len(uploadQueue))
	for i, p := range uploadQueue {
		items[i] = logger.ProgressItem{
			ID:    p.id,
			Name:  p.name,
			Total: len(p.tcs),
		}
	}
	log.InitProgress(items)

	log.StartRenderer()
	log.DrawInitial()

	var allResults []probResult
	var wg sync.WaitGroup

	for _, prob := range uploadQueue {
		prob := prob
		wg.Add(1)

		go func() {
			defer wg.Done()

			tcTotal := len(prob.tcs)
			res := probResult{id: prob.id, name: prob.name, total: tcTotal}

			// Beri tahu renderer bahwa problem ini mulai diupload
			log.UpdateProgress(prob.id, 0, "uploading", 0)

			for i, tc := range prob.tcs {
				var uploadErr error
				for attempt := 1; attempt <= cfg.RetryCount; attempt++ {
					uploadErr = client.UploadTestcaseBrowser(prob.id, tc.InputFile, tc.OutputFile)
					if uploadErr == nil {
						break
					}
					if attempt < cfg.RetryCount {
						time.Sleep(time.Duration(cfg.RetryDelaySeconds) * time.Second)
					}
				}

				if uploadErr == nil {
					log.UpdateProgress(prob.id, i+1, "uploading", res.failed)
				} else {
					res.failed++
					log.UpdateProgress(prob.id, i+1, "uploading", res.failed)
				}
			}

			if res.failed == 0 {
				log.UpdateProgress(prob.id, tcTotal, "finished", 0)
			} else if res.failed == tcTotal {
				log.UpdateProgress(prob.id, tcTotal, "failed", res.failed)
			} else {
				log.UpdateProgress(prob.id, tcTotal, "partial", res.failed)
			}

			allResults = append(allResults, res)
		}()
	}

	wg.Wait()

	log.StopRenderer()

	log.SummaryStart()

	totalProbs := len(allResults)
	probsFailed := 0
	totalTC := 0
	totalFailed := 0

	for _, r := range allResults {
		totalTC += r.total
		totalFailed += r.failed
		if r.failed > 0 {
			probsFailed++
		}
		log.SummaryItem(r.name, r.id, r.total, r.failed)
	}

	log.SummaryStats(totalProbs, probsFailed, totalTC, totalFailed)
	log.SummaryEnd(totalFailed)

	if totalFailed > 0 {
		os.Exit(1)
	}
}
