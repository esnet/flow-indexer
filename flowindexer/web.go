package flowindexer

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
)

type fiHandler struct {
	fi *FlowIndexer
}

func (fh *fiHandler) handleSearch(w http.ResponseWriter, req *http.Request) {
	indexerParam := req.FormValue("i")
	query := req.FormValue("q")
	if indexerParam == "" {
		http.Error(w, "Missing parameter: i", http.StatusBadRequest)
		return
	}
	if query == "" {
		http.Error(w, "Missing parameter: q", http.StatusBadRequest)
		return
	}

	indexer, err := fh.fi.GetIndexer(indexerParam)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	docs, err := indexer.QueryString(query)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	for _, doc := range docs {
		_, err := fmt.Fprintln(w, doc)
		if err != nil {
			return
		}
	}
}
func (fh *fiHandler) handleStats(w http.ResponseWriter, req *http.Request) {
	indexerParam := req.FormValue("i")
	query := req.FormValue("q")
	if indexerParam == "" {
		http.Error(w, "Missing parameter: i", http.StatusBadRequest)
		return
	}
	if query == "" {
		http.Error(w, "Missing parameter: q", http.StatusBadRequest)
		return
	}

	indexer, err := fh.fi.GetIndexer(indexerParam)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	stats, err := indexer.Stats(query)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	json.NewEncoder(w).Encode(stats)
}

func startWeb(fi *FlowIndexer) {
	fh := &fiHandler{fi: fi}
	http.HandleFunc("/search", fh.handleSearch)
	http.HandleFunc("/stats", fh.handleStats)

	bind := fi.config.HTTP.Bind
	log.Printf("Listening on %q\n", bind)
	log.Fatal(http.ListenAndServe(bind, nil))
}