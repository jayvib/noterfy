package middleware

import (
	"github.com/didip/tollbooth"
	"github.com/didip/tollbooth/limiter"
	"github.com/gorilla/mux"
	"net/http"
	"noterfy/api"
	"time"
)

// NewRateLimitMiddleware takes conf rate-limit config an return an instance
// of named rate-limit middleware
func NewRateLimitMiddleware(conf RateLimitConfig) api.NamedMiddleware {
	return api.NewNamedMiddleware("RateLimit", RateLimit(conf))
}

// RateLimitConfig contains all the necessary
// configuration for the rate-limit middleware.
type RateLimitConfig struct {
	// DefaultExpirationTTL is the duration ttl of the token
	// in the token-bucket algorithm.
	DefaultExpirationTTL time.Duration
	// How frequently expire job triggers. How long to
	// retain the token from the memory cache after expiration.
	ExpireJobInterval time.Duration
	// MaxBurst is the maximum threshold of the rate-limit. Meaning
	// when request peak at MaxBurst in a time period succeeding
	// request will return a http.StatusTooManyRequests.
	MaxBurst float64
	// IPLookUps list of Headers to lookup the IP.
	// Default is: "RemoteAddr", "X-Forwarded-For", "X-Real-IP"
	IPLookUps []string
	// LimitMethods are the methods that will only apply the rate-limit.
	LimitMethods []string
	// LimitBasicAuthUsers are the list of users from the Basic Authentication
	// to apply the rate-limit.
	LimitBasicAuthUsers []string
}

// RateLimit do a rate-limiting request middleware based on the
// conf rate-limit configuration.
func RateLimit(conf RateLimitConfig) mux.MiddlewareFunc {
	lmt := tollbooth.NewLimiter(
		conf.MaxBurst,
		&limiter.ExpirableOptions{
			DefaultExpirationTTL: conf.DefaultExpirationTTL,
			ExpireJobInterval:    conf.ExpireJobInterval,
		},
	)

	if len(conf.IPLookUps) > 1 {
		lmt.SetIPLookups(conf.IPLookUps)
	}

	if len(conf.LimitMethods) > 1 {
		lmt.SetMethods(conf.LimitMethods)
	}

	lmt.SetMessage("You have reached maximum request limit.")
	lmt.SetMessageContentType("text/plain;charset=utf-8")

	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

			httpError := tollbooth.LimitByRequest(lmt, w, r)
			if httpError != nil {
				lmt.ExecOnLimitReached(w, r)
				w.Header().Add("Content-Type", lmt.GetMessageContentType())
				w.WriteHeader(httpError.StatusCode)
				_, _ = w.Write([]byte(httpError.Message))
				return
			}

			next.ServeHTTP(w, r)
		})
	}
}
