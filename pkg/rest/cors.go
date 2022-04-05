package rest

import (
	"net/http"
	"regexp"
	"strconv"
	"strings"
	"time"
)

const (
	headerAllowOrigin      = "Access-Control-Allow-Origin"
	headerAllowCredentials = "Access-Control-Allow-Credentials"
	headerAllowHeaders     = "Access-Control-Allow-Headers"
	headerAllowMethods     = "Access-Control-Allow-Methods"
	headerExposeHeaders    = "Access-Control-Expose-Headers"
	headerMaxAge           = "Access-Control-Max-Age"
	headerOrigin           = "Origin"
	headerRequestMethod    = "Access-Control-Request-Method"
	headerRequestHeaders   = "Access-Control-Request-Headers"
)

var (
	defaultAllowHeaders = []string{"Origin", "Accept", "Content-Type", "Authorization"}
	// Regex patterns are generated from AllowOrigins. These are used and generated internally.
	allowOriginPatterns = []string{}
)

type Options struct {
	AllowAllOrigins   bool                                        // If set, all origins are allowed.
	AllowOrigins      []string                                    // A list of allowed origins. Wild cards and FQDNs are supported.
	ShouldAllowOrigin func(origin string, req *http.Request) bool // A func for determining if `origin` is allowed at request time
	AllowCredentials  bool                                        // If set, allows to share auth credentials such as cookies.
	AllowMethods      []string                                    // A list of allowed HTTP methods.
	AllowHeaders      []string                                    // A list of allowed HTTP headers.
	ExposeHeaders     []string                                    // A list of exposed HTTP headers.
	MaxAge            time.Duration                               // Max age of the CORS headers.
}

func (opt *Options) Header(origin string, req *http.Request) (headers map[string]string) {
	headers = make(map[string]string)
	if !opt.AllowAllOrigins && !opt.IsOriginAllowed(origin, req) {
		return
	}

	headers[headerAllowOrigin] = origin
	headers[headerAllowCredentials] = strconv.FormatBool(opt.AllowCredentials)

	if len(opt.AllowMethods) > 0 {
		headers[headerAllowMethods] = strings.Join(opt.AllowMethods, ",")
	}

	if len(opt.AllowHeaders) > 0 {
		headers[headerAllowHeaders] = strings.Join(opt.AllowHeaders, ",")
	}

	if len(opt.ExposeHeaders) > 0 {
		headers[headerExposeHeaders] = strings.Join(opt.ExposeHeaders, ",")
	}

	if opt.MaxAge > time.Duration(0) {
		headers[headerMaxAge] = strconv.FormatInt(int64(opt.MaxAge/time.Second), 10)
	}

	return
}

func (opt *Options) IsOriginAllowed(origin string, req *http.Request) (allowed bool) {
	if opt.ShouldAllowOrigin != nil {
		return opt.ShouldAllowOrigin(origin, req)
	}
	for _, pattern := range allowOriginPatterns {
		allowed, _ = regexp.MatchString(pattern, origin)
		if allowed {
			return
		}
	}
	return
}

func (opt *Options) Handler(next http.Handler) http.HandlerFunc {
	if len(opt.AllowHeaders) == 0 {
		opt.AllowHeaders = defaultAllowHeaders
	}

	for _, origin := range opt.AllowOrigins {
		pattern := regexp.QuoteMeta(origin)
		pattern = strings.Replace(pattern, "\\*", ".*", -1)
		pattern = strings.Replace(pattern, "\\?", ".", -1)
		allowOriginPatterns = append(allowOriginPatterns, "^"+pattern+"$")
	}

	return func(w http.ResponseWriter, req *http.Request) {
		if origin := req.Header.Get(headerOrigin); origin != "" {
			for key, value := range opt.Header(origin, req) {
				w.Header().Set(key, value)
			}
			if req.Method == "OPTIONS" {
				w.WriteHeader(200)
				return
			}
		}
		next.ServeHTTP(w, req)
	}
}
