package shenmeci

import (
	"crypto/tls"
	"encoding/json"
	"fmt"
	"html/template"
	"io"
	"io/fs"
	"log"
	"mime/multipart"
	"net/http"
	"net/url"
	"strconv"
	"strings"
	"time"

	"github.com/rhcarvalho/shenmeci/internal/segmentation"
)

type QueryRecord struct {
	Query       string        `json:"query"`
	Result      []*Result     `json:"result"`
	When        time.Time     `json:"when"`
	Duration    time.Duration `json:"duration"`
	RequestInfo *RequestInfo  `json:"requestinfo"`
}

// RequestInfo is a JSON-serializable version of http.Request.
type RequestInfo struct {
	Method           string               `json:"method"`
	URL              string               `json:"url"`
	Proto            string               `json:"proto"`
	ProtoMajor       int                  `json:"protomajor"`
	ProtoMinor       int                  `json:"protominor"`
	Header           http.Header          `json:"header"`
	Body             []byte               `json:"body"`
	ContentLength    int64                `json:"contentlength"`
	TransferEncoding []string             `json:"transferencoding"`
	Close            bool                 `json:"close"`
	Host             string               `json:"host"`
	Form             url.Values           `json:"form"`
	PostForm         url.Values           `json:"postform"`
	MultipartForm    *multipart.Form      `json:"multipartform"`
	Trailer          http.Header          `json:"trailer"`
	RemoteAddr       string               `json:"remoteaddr"`
	RequestURI       string               `json:"requesturi"`
	TLS              *tls.ConnectionState `json:"tls"`
}

func requestToRequestInfo(r *http.Request) *RequestInfo {
	body, err := io.ReadAll(r.Body)
	if err != nil {
		log.Println("reading request body:", err)
	}
	return &RequestInfo{
		Method:           r.Method,
		URL:              r.URL.String(),
		Proto:            r.Proto,
		ProtoMajor:       r.ProtoMajor,
		ProtoMinor:       r.ProtoMinor,
		Header:           r.Header,
		Body:             body,
		ContentLength:    r.ContentLength,
		TransferEncoding: r.TransferEncoding,
		Close:            r.Close,
		Host:             r.Host,
		Form:             r.Form,
		PostForm:         r.PostForm,
		MultipartForm:    r.MultipartForm,
		Trailer:          r.Trailer,
		RemoteAddr:       r.RemoteAddr,
		RequestURI:       r.RequestURI,
		TLS:              r.TLS,
	}
}

type Results struct {
	R []*Result
}

type Result struct {
	Z string // Hanzi
	M string // Meaning
	P string // Pinyin
}

func (r *Results) MarshalJSON() ([]byte, error) {
	if r.R != nil {
		return json.Marshal(&map[string]interface{}{"r": r.R})
	} else {
		return json.Marshal(&map[string]interface{}{"r": []interface{}{}})
	}
}

func (r *Result) MarshalJSON() ([]byte, error) {
	return json.Marshal(
		struct {
			Z string `json:"z"`
			M string `json:"m"`
			P string `json:"p"`
		}{
			Z: template.HTMLEscapeString(r.Z),
			M: template.HTMLEscapeString(r.M),
			P: template.HTMLEscapeString(r.P),
		},
	)
}

func segmentHandler(segmenter segmentation.Segmenter) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		startTime := time.Now()
		query := r.FormValue("q")
		results := keysToResults(func() (keys []string) {
			for _, key := range segmenter.Segment([]rune(query)) {
				keys = append(keys, string(key))
			}
			return keys
		}())
		if len(results.R) == 1 && results.R[0].M == "?" {
			log.Printf("q='%v' triggers Full-Text Search", query)
			results = keysToResults(searchDB(query))
		}
		if len(results.R) == 0 {
			log.Printf("q='%v' returns no results", query)
		}
		b, _ := json.Marshal(results)
		w.Header().Set("Content-Length", strconv.FormatInt(int64(len(b)), 10))
		w.Header().Set("Content-Type", "application/json")
		w.Write(b)
		duration := time.Since(startTime)

		// Store query information in a new goroutine, without blocking the
		// request-response cycle.
		go func() {
			err := insertQueryRecord(QueryRecord{
				query,
				results.R,
				startTime,
				duration,
				requestToRequestInfo(r),
			})
			// Log insertion errors
			if err != nil {
				log.Println("insertQueryRecord:", err)
			}
		}()
	}
}

func keysToResults(keys []string) *Results {
	results := &Results{}
	var m, p []string
	for _, key := range keys {
		entry, ok := cedict.Dict[key]
		if ok {
			m = entry.definitions
			p = entry.pinyin
		} else {
			m = []string{"?"}
			p = []string{""}
		}
		results.R = append(results.R, &Result{
			Z: key,
			M: strings.Join(m, "/"),
			P: strings.Join(p, "/"),
		})
	}
	return results
}

// HTTPConfig configures an HTTP server.
type HTTPConfig struct {
	Host string
	Port int

	Static fs.FS

	Segmenter segmentation.Segmenter
}

// Serve starts an HTTP server.
func Serve(config *HTTPConfig) {
	http.Handle("/", http.FileServer(http.FS(config.Static)))
	http.HandleFunc("/segment", segmentHandler(config.Segmenter))

	addr := fmt.Sprintf("%s:%d", config.Host, config.Port)
	log.Printf("serving at http://%s", addr)
	err := http.ListenAndServe(addr, nil)
	if err != nil {
		log.Fatal(err)
	}
}
