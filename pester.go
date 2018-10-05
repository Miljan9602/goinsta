package goinsta

import (
	"github.com/sethgrid/pester"
)

type PesterOptions struct {
	Concurrency int
	MaxRetries  int
	Backoff     pester.BackoffStrategy
}


func DefaultOptions() *PesterOptions {
	return &PesterOptions{
		Concurrency : pester.DefaultClient.Concurrency,
		MaxRetries : pester.DefaultClient.MaxRetries,
		Backoff :  pester.DefaultClient.Backoff,
	}
}