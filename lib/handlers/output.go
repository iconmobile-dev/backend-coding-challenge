// Package handlers contains shared router handlers and middleware
package handlers

import (
	"bytes"
	"compress/gzip"
	"fmt"
	"io"
	"math"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/iconmobile-dev/go-core/errors"
	jsoniter "github.com/json-iterator/go"
)

var json = jsoniter.ConfigFastest

// JSONMsgStr string message response.
type JSONMsgStr struct {
	Msg string
}

// JSONMsg returns an HTTP response as JSON message with given status code
// if v is a string then it is sent as "Msg" property value.
// Otherwise it encodes v as JSON
// is preferred over `json.NewEncoder` since it sets the correct Content-Type
func JSONMsg(w http.ResponseWriter, r *http.Request, status int, v interface{}) {
	var body []byte

	// get the type of the interface as string value
	interfaceType := fmt.Sprintf("%T", v)

	t2 := time.Now()
	if interfaceType == "string" {
		// return the value as "Msg" property
		//json.NewEncoder(w).Encode(struct{ Msg string }{Msg: v.(string)}
		data := JSONMsgStr{}
		data.Msg = v.(string)

		var err error
		body, err = json.Marshal(&data)
		if err != nil {
			log.Error(err)
		}
	} else {
		// encode the interface as json
		//json.NewEncoder(w).Encode(v)
		var err error
		body, err = json.Marshal(&v)
		if err != nil {
			log.Error(err)
		}
	}
	if ms := ElapsedMilliseconds(t2); ms > 0 {
		log.Verbose("JSONMsg: JSON Marshal took", ms, "ms")
	}
	t3 := time.Now()

	// add the time headers after the marshaling
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.Header().Set("X-Content-Type-Options", "nosniff")

	val := r.Context().Value(ctxKeyTime)
	if val != nil {
		t1 := val.(time.Time)
		AddTimeHeaders(w, r, t1)
	}
	w.WriteHeader(status)
	_, err := io.Copy(w, bytes.NewReader(body))
	if err != nil {
		log.Error(err)
	}
	if ms := ElapsedMilliseconds(t3); ms > 0 {
		log.Verbose("JSONMsg: writing body took", ms, "ms")
	}
}

// RawJSON is a utility function to forward a response during proxy
// It returns the body as GZIP if possible
func RawJSON(w http.ResponseWriter, req *http.Request, status int, reader io.Reader) {
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.Header().Set("X-Content-Type-Options", "nosniff")

	// during proxy calls the ctxKeyTime somehow stays nil?!
	// it is manually set by Worf in proxy handling route
	val := req.Context().Value(ctxKeyTime)
	if val != nil {
		t1 := val.(time.Time)
		AddTimeHeaders(w, req, t1)
	}

	if strings.Contains(req.Header.Get("Accept-Encoding"), "gzip") {
		// GZIP the response
		w.Header().Set("Content-Encoding", "gzip")
		w.WriteHeader(status)

		// write gzipped reader data to writer
		t1 := time.Now()
		var b bytes.Buffer
		_, err := b.ReadFrom(reader)
		if err != nil {
			log.Error(err)
		}

		gzipWriter, _ := gzip.NewWriterLevel(w, gzip.BestSpeed) // gzip.BestCompression
		_, err = gzipWriter.Write(b.Bytes())
		if err != nil {
			log.Error(err)
		}

		err = gzipWriter.Close()
		if err != nil {
			log.Error(err)
		}

		if ms := ElapsedMilliseconds(t1); ms > 0 {
			log.Verbose("RawJSON: GZIP compression took", ms, "ms")
		}
	} else {
		// normal response
		w.WriteHeader(status)
		_, err := io.Copy(w, reader)
		if err != nil {
			log.Error(err)
		}
	}
}

// AddTimeHeaders add Action-Time and Request-Time headers
// if not already set using the time context value
func AddTimeHeaders(w http.ResponseWriter, r *http.Request, startTime time.Time) {
	duration := strconv.FormatInt(ElapsedMilliseconds(startTime), 10)
	if duration == "0" {
		duration = "1"
	}
	//log.Debug("t1:", TimeToMilliseconds(startTime), "ms")
	//log.Debug("now:", TimeToMilliseconds(time.Now()), "ms")
	//log.Debug("duration:", duration, "ms")

	if w.Header().Get("Request-Time") == "" {
		w.Header().Set("Request-Time", duration)
	}
	if w.Header().Get("Action-Time") == "" {
		w.Header().Set("Action-Time", duration)
	}
}

// TimeToMilliseconds returns a Golang time as integer milliseconds
func TimeToMilliseconds(t time.Time) int64 {
	return int64(math.Round(float64(t.UnixNano() / 1000000)))
}

// ElapsedMilliseconds returns time duration to Now in Milliseconds
func ElapsedMilliseconds(startTime time.Time) int64 {
	nanoseconds := time.Now().UnixNano() - startTime.UnixNano()
	return int64(math.Round(float64(nanoseconds / 1000000)))
}

// JSONMsgErr returns an HTTP response as JSON message with given status code
// if v is a string then it is sent as "Msg" property value.
// Otherwise it encodes v as JSON
// is preferred over `json.NewEncoder` since it sets the correct Content-Type
func JSONMsgErr(w http.ResponseWriter, r *http.Request, err error, msg string) {
	var body []byte
	var status int
	data := JSONMsgStr{}

	if msg != "" {
		data.Msg = msg
	}

	// set custom app err Message
	appErr, ok := err.(*errors.Error)
	if !ok {
		status = 500
		data.Msg = fmt.Sprintf("%s: Internal Server Error", data.Msg)
	} else {
		switch appErr.Kind {
		case errors.Unauthorized:
			status = 401
		case errors.Forbidden:
			status = 403
		case errors.NotFound:
			status = 404
		case errors.Conflict:
			status = 409
		case errors.Unprocessable:
			status = 422
		case errors.Internal:
			status = 500
		default:
			status = 500
		}

		errMsg := errors.ToHTTPResponse(appErr)
		if errMsg != "" {
			data.Msg = fmt.Sprintf("%s: %s", data.Msg, errMsg)
		}
	}

	body, _ = json.Marshal(&data)

	// add the time headers after the marshaling
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.Header().Set("X-Content-Type-Options", "nosniff")

	w.WriteHeader(status)
	_, err = io.Copy(w, bytes.NewReader(body))
	if err != nil {
		log.Error(err)
	}
}
