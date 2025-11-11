package server

import (
	"encoding/json"
	"log"
	"net/http"
	"strings"

	"search-service/internal/config"

	gql "github.com/graphql-go/graphql"
)

// Server wraps the HTTP transport for the GraphQL schema.
type Server struct {
	conf   *config.Config
	schema gql.Schema
}

// New creates a Server instance.
func New(conf *config.Config, schema gql.Schema) *Server {
	return &Server{
		conf:   conf,
		schema: schema,
	}
}

// Run starts the HTTP server.
func (s *Server) Run() error {
	mux := http.NewServeMux()
	mux.Handle("/graphql", http.HandlerFunc(s.handleGraphQL))
	mux.HandleFunc("/healthz", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte(`{"status":"ok"}`))
	})

	addr := ":" + s.conf.AppPort
	log.Printf("search-service listening on %s (playground=%t)", addr, s.conf.EnablePlayground)
	return http.ListenAndServe(addr, mux)
}

func (s *Server) handleGraphQL(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodOptions:
		w.WriteHeader(http.StatusNoContent)
		return
	case http.MethodGet:
		s.handleGraphQLGet(w, r)
	case http.MethodPost:
		s.handleGraphQLPost(w, r)
	default:
		http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
	}
}

func (s *Server) handleGraphQLGet(w http.ResponseWriter, r *http.Request) {
	query := strings.TrimSpace(r.URL.Query().Get("query"))
	if query == "" {
		if s.conf.EnablePlayground {
			w.Header().Set("Content-Type", "text/html; charset=utf-8")
			_, _ = w.Write([]byte(playgroundHTML))
			return
		}
		http.Error(w, "query is required", http.StatusBadRequest)
		return
	}

	vars, err := parseVariables(r.URL.Query().Get("variables"))
	if err != nil {
		http.Error(w, "invalid variables payload", http.StatusBadRequest)
		return
	}

	s.writeResult(w, gql.Params{
		Schema:         s.schema,
		Context:        r.Context(),
		RequestString:  query,
		VariableValues: vars,
		OperationName:  r.URL.Query().Get("operationName"),
	})
}

func (s *Server) handleGraphQLPost(w http.ResponseWriter, r *http.Request) {
	defer r.Body.Close()
	var req graphQLRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "invalid JSON payload", http.StatusBadRequest)
		return
	}
	if strings.TrimSpace(req.Query) == "" {
		http.Error(w, "query is required", http.StatusBadRequest)
		return
	}

	vars, err := req.variablesMap()
	if err != nil {
		http.Error(w, "invalid variables payload", http.StatusBadRequest)
		return
	}

	s.writeResult(w, gql.Params{
		Schema:         s.schema,
		Context:        r.Context(),
		RequestString:  req.Query,
		VariableValues: vars,
		OperationName:  req.OperationName,
	})
}

func (s *Server) writeResult(w http.ResponseWriter, params gql.Params) {
	result := gql.Do(params)
	w.Header().Set("Content-Type", "application/json")
	status := http.StatusOK
	if len(result.Errors) > 0 {
		status = http.StatusBadRequest
	}
	w.WriteHeader(status)
	_ = json.NewEncoder(w).Encode(result)
}

func (s *Server) wrapCORS(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Access-Control-Allow-Origin", "http://localhost:3000")
		w.Header().Set("Access-Control-Allow-Methods", "GET,POST,OPTIONS")
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization, X-Member-ID, x-member-id")
		w.Header().Set("Access-Control-Max-Age", "3600")

		if r.Method == http.MethodOptions {
			w.WriteHeader(http.StatusNoContent)
			return
		}

		next.ServeHTTP(w, r)
	})
}

type graphQLRequest struct {
	Query         string      `json:"query"`
	OperationName string      `json:"operationName"`
	Variables     interface{} `json:"variables"`
}

func (r graphQLRequest) variablesMap() (map[string]interface{}, error) {
	switch v := r.Variables.(type) {
	case nil:
		return nil, nil
	case map[string]interface{}:
		return v, nil
	case string:
		return parseVariables(v)
	default:
		raw, err := json.Marshal(v)
		if err != nil {
			return nil, err
		}
		return parseVariablesBytes(raw)
	}
}

func parseVariables(raw string) (map[string]interface{}, error) {
	raw = strings.TrimSpace(raw)
	if raw == "" {
		return nil, nil
	}
	return parseVariablesBytes([]byte(raw))
}

func parseVariablesBytes(raw []byte) (map[string]interface{}, error) {
	var vars map[string]interface{}
	if len(raw) == 0 || string(raw) == "null" {
		return nil, nil
	}
	if err := json.Unmarshal(raw, &vars); err != nil {
		return nil, err
	}
	return vars, nil
}

const playgroundHTML = `<!DOCTYPE html>
<html>
<head>
  <meta charset="utf-8" />
  <title>TeamUp Search Service</title>
  <link rel="stylesheet" href="https://cdn.jsdelivr.net/npm/graphql-playground-react@1.7.42/build/static/css/index.css" />
  <link rel="shortcut icon" href="https://graphql.org/img/favicon.png" />
  <script src="https://cdn.jsdelivr.net/npm/graphql-playground-react@1.7.42/build/static/js/middleware.js"></script>
</head>
<body>
  <div id="root" />
  <script>
    window.addEventListener('load', function () {
      GraphQLPlayground.init(document.getElementById('root'), {
        endpoint: '/graphql'
      })
    })
  </script>
</body>
</html>`
