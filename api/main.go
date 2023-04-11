package api

import (
	"encoding/json"
	"io/ioutil"
	"log"
	"net/http"
	"strconv"
	"time"

	db "github.com/Oluwatunmise-olat/Horus-Clone/db"
	types "github.com/Oluwatunmise-olat/Horus-Clone/types"
	"github.com/gorilla/mux"
)

type httpFunctionSig http.HandlerFunc

type wrappedHttpResponseStatus types.WrappedResponseStatus

type API struct {
	Storage *db.Store
}

var requestQueue = make(chan types.RequestLogger)

func (w *wrappedHttpResponseStatus) WriteHeader(status int) {
	w.Status = status // Keep track of response status

	w.ResponseWriter.WriteHeader(status)
}

func (w *wrappedHttpResponseStatus) Write(b []byte) (int, error) {
	if (w.Status == 0){
		w.Status = http.StatusOK
	}
	
	w.Body.Write(b) // Keep track of response body
	return w.ResponseWriter.Write(b)
}

func (next httpFunctionSig) ServeHTTP(w http.ResponseWriter, r *http.Request){
	next(w, r)
}

// For web access
func (s *API) Serve(port string) (error) {
	log.Println("🚀 Initialing Horus Clone Server...")
	router := mux.NewRouter().StrictSlash(true)

	horus := router.PathPrefix("/horus-clone").Subrouter()
	horus.HandleFunc("/logs", nil)

	// Serve as goroutine to prevent blocking of main block
	var err error = nil
	go func (){
		err = http.ListenAndServe(port, router)
	}()

	return err
}

func (s *API) Watch(next httpFunctionSig) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request){

		reqBody, err := ioutil.ReadAll(r.Body)
		if err != nil {
			log.Fatal(err)
		}
		defer r.Body.Close()

		reqHeader, err := json.Marshal(r.Header)
		if err != nil {
			log.Fatal(err)
		}
		httpResponseWriter := &wrappedHttpResponseStatus{ ResponseWriter: w }
		startTime := time.Now()
		next.ServeHTTP(httpResponseWriter, r)
		
		respBody, err := json.Marshal(httpResponseWriter.Body.String())
		if err != nil {
			log.Fatal(err)
		}

		record := &types.RequestLogger{ 
			RequestMethod: r.Method, 
			RequestHost: r.Host,
			RequestBody: string(reqBody), 
			RequestHeaders: string(reqHeader), 
			RequestPath: r.URL.Path,
			ResponseBody: string(respBody),
			ResponseStatus: strconv.Itoa(httpResponseWriter.Status),
			ReqResTime: float64(time.Since(startTime).Seconds()), // Stored in seconds
		}

		if err := s.Storage.PerformInsert(record); err != nil {
			log.Fatal(err)
		}

	}
}

