package main

import (
	"fmt"
	"net/http"
	"sync"
	"time"

	"github.com/tomasen/realip"
	"golang.org/x/time/rate"
)

type Client struct {
	limiter      *rate.Limiter
	lastSeenTime time.Time
}

func (app *application) recoverPanic(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		defer func() {
			if err := recover(); err != nil {
				w.Header().Set("Connection", "Close")
				app.serverErrorResponse(w, r, fmt.Errorf("%s", err))
			}
		}()
		next.ServeHTTP(w, r)
	})
}

func (app *application) rateLimit(next http.Handler) http.Handler {

	var (
		mutex   sync.Mutex
		clients map[string]*Client = make(map[string]*Client)
	)

	// background cleanup check runs every
	// minute deleting old inactive clients
	go func() {
		for {
			time.Sleep(time.Minute)

			mutex.Lock()

			for key, client := range clients {
				// client have not been seen in more than 3
				// minutes ago eliminate it from the map
				if time.Since(client.lastSeenTime) > 3*time.Minute {
					delete(clients, key)
				}
			}

			mutex.Unlock()
		}
	}()

	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

		if app.config.Limiter.Enabled {
			ip := realip.FromRequest(r)

			mutex.Lock() // guarding this entire checking process, 1 goroutine at a time

			if _, ok := clients[ip]; !ok {
				clients[ip] = &Client{limiter: rate.NewLimiter(rate.Limit(app.config.Limiter.RPS), app.config.Limiter.Burst)}
			}

			clients[ip].lastSeenTime = time.Now()

			if !clients[ip].limiter.Allow() {
				mutex.Unlock() // unlock, because checking is done
				app.rateLimitExceededResponse(w, r)
				return
			}

			mutex.Unlock() // unlock, because checking is done
		}
		next.ServeHTTP(w, r)
	})
}

func (app *application) enableCORS(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Add("Vary", "Origin")
		w.Header().Add("Vary", "Access-Control-Request-Method")
		origin := r.Header.Get("Origin")

		if origin != "" && len(app.config.CORS.TrustedOrigins) > 0 {
			for _, trustedOrigin := range app.config.CORS.TrustedOrigins {
				if origin == trustedOrigin { // make sure it matches trustedOrigin exactly, no partial matches
					w.Header().Set("Access-Control-Allow-Origin", trustedOrigin)
					w.Header().Set("Access-Control-Allow-Credentials", "true")

					// if it is a pre-flight CORS request
					if r.Method == http.MethodOptions && r.Header.Get("Access-Control-Request-Method") != "" {

						w.Header().Set("Access-Control-Allow-Credentials", "true")
						w.Header().Set("Access-Control-Allow-Methods", "OPTIONS, PUT, PATCH, DELETE")
						w.Header().Set("Access-Control-Allow-Headers", "Authorization, Content-Type, X-CSRF-Token")

						w.WriteHeader(http.StatusOK)
						return
					}
				}
			}
		}

		next.ServeHTTP(w, r)
	})
}

// func (app *application) metrics(next http.Handler) http.Handler {
// 	totalRequestsReceived := expvar.NewInt("total_requests_received")
// 	totalResponsesSent := expvar.NewInt("total_responses_sent")
// 	totalProcessingTimeMicroseconds := expvar.NewInt("total_processing_time_μs")
// 	totalActiveRequests := expvar.NewInt("total_active_requests")
// 	totalResponsesSentByStatus := expvar.NewMap("total_responses_sent_by_status")

// 	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

// 		totalRequestsReceived.Add(1)

// 		metrics := httpsnoop.CaptureMetrics(next, w, r)

// 		totalResponsesSent.Add(1)

// 		totalProcessingTimeMicroseconds.Add(metrics.Duration.Microseconds())
// 		totalResponsesSentByStatus.Add(strconv.Itoa(metrics.Code), 1)
// 		totalActiveRequests.Set(totalRequestsReceived.Value() - totalResponsesSent.Value())
// 	})
// }
