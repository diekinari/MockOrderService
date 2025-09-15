package http

import (
	"MockOrderService/internal/domain/model"
	"context"
	"encoding/json"
	"html/template"
	"io"
	"net/http"
	"net/url"
	"path/filepath"
	"runtime"
	"strings"
	"time"

	"go.uber.org/zap"
)

type WebServer struct {
	server *http.Server
}

// StartWebServer starts client server in a separate goroutine.
// handles order requests and shuts down gracefully on context cancel.
// If server fails to start or shutdown, logs error.
// Notice: we could reuse api server, but let's assume it's a separate service
func (ws *WebServer) StartWebServer(sugar *zap.SugaredLogger) error {
	mux := http.NewServeMux()
	mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		handleRequest(w, r, sugar)
	})
	srv := &http.Server{
		Addr:    ":8082",
		Handler: mux,
	}
	ws.server = srv
	sugar.Infow("started client server at :8082")
	err := srv.ListenAndServe()
	if err != nil {
		sugar.Errorw("client server failed", "error", err)
		return err
	}
	return nil
}

func (ws *WebServer) Shutdown(ctx context.Context) error {
	if ws.server != nil {
		return ws.server.Shutdown(ctx)
	}
	return nil

}

type APIError struct {
	Error string `json:"Error"`
}

// templateData передаётся в HTML-шаблон
type templateData struct {
	Query string
	Order *model.Order
	Error string
	Now   time.Time
}

// getTemplateDir возвращает директорию с HTML-шаблонами
func getTemplateDir() string {
	_, filename, _, _ := runtime.Caller(0)
	root := filepath.Join(filepath.Dir(filename), "..", "..", "..")
	return filepath.Join(root, "web")
}

// getDashboardTemplate возвращает путь к HTML-шаблону
func getDashboardTemplate() string {
	return filepath.Join(getTemplateDir(), "dashboard.html")
}

var dashboardTmpl = template.Must(template.New("dashboard.html").Funcs(template.FuncMap{
	"nl2br": func(s string) template.HTML {
		// безопасный простой перевод newlines -> <br>
		return template.HTML(strings.ReplaceAll(template.HTMLEscapeString(s), "\n", "<br>"))
	},
	"formatTime": func(t *time.Time) string {
		if t == nil {
			return ""
		}
		return t.Format(time.RFC3339)
	},
}).ParseFiles(getDashboardTemplate()))

// handleRequest обрабатывает форму и делает запрос к API, парсит JSON и рендерит шаблон.
func handleRequest(w http.ResponseWriter, r *http.Request, sugar *zap.SugaredLogger) {
	ctx := r.Context()
	q := strings.TrimSpace(r.URL.Query().Get("orderUID"))

	data := templateData{
		Query: q,
		Now:   time.Now(),
	}

	if q != "" {
		// build API URL
		apiURL := "http://localhost:8081/api/order/" + url.PathEscape(q)

		// HTTP client with timeout
		client := &http.Client{Timeout: 5 * time.Second}

		req, err := http.NewRequestWithContext(ctx, http.MethodGet, apiURL, nil)
		if err != nil {
			sugar.Errorw("failed to create request", "error", err)
			data.Error = "internal: failed to create request"
			_ = dashboardTmpl.Execute(w, data)
			return
		}

		resp, err := client.Do(req)
		if err != nil {
			sugar.Errorw("request to API failed", "error", err, "apiURL", apiURL)
			data.Error = "failed to reach API: " + err.Error()
			_ = dashboardTmpl.Execute(w, data)
			return
		}
		defer func() {
			if cerr := resp.Body.Close(); cerr != nil {
				sugar.Warnw("failed to close response body", "error", cerr, "apiURL", apiURL)
			}
		}()

		// Read and decode into Order (primary) or APIError (fallback)
		var order model.Order
		dec := json.NewDecoder(resp.Body)
		if err := dec.Decode(&order); err == nil && order.OrderUID != "" {
			// success: we got an order
			data.Order = &order
		} else {
			// try to decode API error (rewind by re-fetching body not possible, so re-request or decode via bytes)

			// reset by reading raw bytes from a fresh request
			if cerr := resp.Body.Close(); cerr != nil {
				sugar.Warnw("failed to close response body", "error", cerr, "apiURL", apiURL)
			}
			resp2, err2 := client.Get(apiURL)
			if err2 != nil {
				sugar.Errorw("request retry failed", "error", err2)
				data.Error = "failed to read API response"
				_ = dashboardTmpl.Execute(w, data)
				return
			}
			defer func() {
				if cerr := resp2.Body.Close(); cerr != nil {
					sugar.Warnw("failed to close response body", "error", cerr, "apiURL", apiURL)
				}
			}()
			raw, err3 := io.ReadAll(resp2.Body)
			if err3 != nil {
				sugar.Errorw("failed to read body", "error", err3)
				data.Error = "failed to read API response"
				_ = dashboardTmpl.Execute(w, data)
				return
			}

			// try unmarshal into Order one more time
			if err := json.Unmarshal(raw, &order); err == nil && order.OrderUID != "" {
				data.Order = &order
			} else {
				// try APIError
				var apiErr APIError
				if err := json.Unmarshal(raw, &apiErr); err == nil && apiErr.Error != "" {
					data.Error = apiErr.Error
				} else {
					// fallback: show raw body as error
					data.Error = "unexpected API response: " + string(raw)
				}
			}
		}
	}

	// Render template
	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	if err := dashboardTmpl.Execute(w, data); err != nil {
		sugar.Errorw("failed to execute template", "error", err)
		http.Error(w, "internal template error", http.StatusInternalServerError)
	}
}
