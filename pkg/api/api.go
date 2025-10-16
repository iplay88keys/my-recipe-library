package api

import (
	"encoding/json"
	"fmt"
	"log"
	"net"
	"net/http"
	"strconv"
	"time"

	"github.com/iplay88keys/my-recipe-library/pkg/token"
)

type validate func(r *http.Request) (*token.AccessDetails, error)
type retrieveAccessDetails func(details *token.AccessDetails) (int64, error)

type Config struct {
	Port                  string
	StaticDir             string
	Validate              validate
	RetrieveAccessDetails retrieveAccessDetails
	Endpoints             []*Endpoint
}

type API struct {
	Config *Config
	Server *http.Server
}

func New(config *Config) *API {
	server := &API{
		Config: config,
		Server: &http.Server{
			Addr:         net.JoinHostPort("", config.Port),
			ReadTimeout:  15 * time.Second,
			WriteTimeout: 15 * time.Second,
		},
	}

	fmt.Println("Registering endpoints:")

	mux := http.NewServeMux()

	// Register API endpoints
	for _, endpoint := range config.Endpoints {
		pattern := fmt.Sprintf("%s /api/v1/%s", endpoint.Method, endpoint.Path)
		fmt.Printf("%s %s\n", endpoint.Method, pattern)
		mux.Handle(pattern, server.createHandler(endpoint))
	}

	// Handle API 404s
	mux.Handle("/api/v1/", http.HandlerFunc(notFoundHandler))

	// SPA handler for all other routes
	spa, err := NewSPAHandler(config.StaticDir, "index.html")
	if err != nil {
		log.Fatalf("Unable to create SPA handler: %s", err)
	}

	mux.Handle("/", spa)

	server.Server.Handler = mux

	return server
}

func (a *API) Start() (shutdown func()) {
	go a.Server.ListenAndServe()

	return func() {
		a.Server.Close()
	}
}

func (a *API) createHandler(e *Endpoint) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		resp := &Response{}
		valid := true
		userID := int64(-1)

		req := &Request{
			Req: r,
		}

		startTime := time.Now()

		if e.Auth {
			userID, valid = a.ValidateUserToken(r)
		}

		if valid {
			req.UserID = userID

			resp = e.Handle(req)
		} else {
			resp.StatusCode = http.StatusUnauthorized
		}

		writeResponse(w, resp)
		logRequest(startTime, time.Now(), req, resp)
	})
}

func writeResponse(w http.ResponseWriter, resp *Response) {
	var body []byte
	status := resp.StatusCode

	if resp.Body != nil {
		var err error
		body, err = json.Marshal(resp.Body)
		if err != nil {
			fmt.Printf("Failed to marshal response body: %s\n", err.Error())
			status = http.StatusInternalServerError
		} else {
			w.Header().Set("Content-type", "application/json")
		}
	}

	w.WriteHeader(status)
	_, err := w.Write(body)
	if err != nil {
		fmt.Printf("Error writing body: %s", err.Error())
	}
}

func (a *API) ValidateUserToken(r *http.Request) (int64, bool) {
	details, err := a.Config.Validate(r)
	if err != nil {
		fmt.Printf("Failed to validate token: %s\n", err.Error())
		return -1, false
	}

	userID, err := a.Config.RetrieveAccessDetails(details)
	if err != nil {
		fmt.Printf("Failed to retrieve token details from redis: %s\n", err.Error())
		return -1, false
	}

	return userID, true
}

func logRequest(startTime, endTime time.Time, req *Request, resp *Response) {
	start := startTime.Format(time.RFC3339)
	end := endTime.Format(time.RFC3339)

	latency := endTime.Sub(startTime)

	var userID string
	if req.UserID == int64(-1) {
		userID = "<none>"
	} else {
		userID = strconv.FormatInt(req.UserID, 10)
	}

	fmt.Printf("[API] %s [path=%s] [userID=%s] [start=%s] [end=%s] [latency=%s] [status=%d]\n",
		req.Req.Method,
		req.Req.URL.Path,
		userID,
		start,
		end,
		latency,
		resp.StatusCode,
	)
}

func notFoundHandler(w http.ResponseWriter, r *http.Request) {
	resp := NewResponse(http.StatusNotFound, nil)

	logRequest(
		time.Now(),
		time.Now(),
		&Request{
			Req:    r,
			UserID: -1,
		},
		resp,
	)

	writeResponse(w, resp)
}
