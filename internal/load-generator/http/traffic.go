package http

import (
	"context"
	"encoding/json"
	"io"
	"net/http"
	"strconv"
	"strings"
	"sync"
	"time"

	"tpiii.local/daniel-tpiii/internal/load-generator/results"
)

var TrafficHeader = []string{"timestamp_ms", "service", "code", "instance_id", "latency_ms"}

type Traffic struct {
	client   *http.Client
	url      string
	service  string
	interval time.Duration
	writer   *results.Writer
}

func NewTraffic(
	loadBalancerURL string,
	service string,
	interval time.Duration,
	timeout time.Duration,
	writer *results.Writer,
) *Traffic {
	return &Traffic{
		client:   &http.Client{Timeout: timeout},
		url:      strings.TrimSuffix(loadBalancerURL, "/") + "/" + service + "/",
		service:  service,
		interval: interval,
		writer:   writer,
	}
}

func (t *Traffic) Run(ctx context.Context) {
	ticker := time.NewTicker(t.interval)
	defer ticker.Stop()

	var wg sync.WaitGroup
	defer wg.Wait()

	for {
		select {
		case <-ticker.C:
			wg.Go(func() {
				t.send(ctx)
			})

		case <-ctx.Done():
			return
		}
	}
}

func (t *Traffic) send(ctx context.Context) {
	start := time.Now()
	code, instanceID := t.get(ctx)
	latency := time.Since(start)

	if ctx.Err() != nil {
		return
	}

	t.writer.Write(
		results.FormatTimestamp(start),
		t.service,
		strconv.Itoa(code),
		instanceID,
		results.FormatMilliseconds(latency),
	)
}

func (t *Traffic) get(ctx context.Context) (int, string) {
	request, err := http.NewRequestWithContext(ctx, http.MethodGet, t.url, nil)
	if err != nil {
		return 0, ""
	}

	response, err := t.client.Do(request)
	if err != nil {
		return 0, ""
	}
	defer response.Body.Close()

	var body struct {
		InstanceID string `json:"instance_id"`
	}

	if response.StatusCode == http.StatusOK {
		_ = json.NewDecoder(response.Body).Decode(&body)
	}

	_, _ = io.Copy(io.Discard, response.Body)

	return response.StatusCode, body.InstanceID
}
