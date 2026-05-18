// Package ingest coordinates the collection of deployment events from one or
// more CI/CD source adapters and their ingestion into the shared audit
// timeline.
//
// # Overview
//
// A Processor is created with a source Registry and a Timeline. Named
// Fetcher implementations are then registered against source names that
// exist in the Registry. Calling Run triggers all registered Fetchers
// concurrently-safe and appends the returned events to the Timeline.
//
// # Usage
//
//	reg := source.NewRegistry()
//	_ = reg.Register("github-actions", source.Meta{DisplayName: "GitHub Actions"})
//
//	tl := timeline.New()
//	p := ingest.NewProcessor(reg, tl)
//	_ = p.Register("github-actions", myGitHubFetcher)
//
//	if err := p.Run(ctx); err != nil {
//		log.Fatal(err)
//	}
package ingest
