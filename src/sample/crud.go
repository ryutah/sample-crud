package sample

import (
	"encoding/json"
	"fmt"
	"net/http"

	"golang.org/x/net/context"
	"google.golang.org/appengine"
	"google.golang.org/appengine/datastore"
)

const KindFoo = "foo"

type Foo struct {
	Foo string `json:"foo"`
	Bar string `json:"bar"`
}

func GetCRUDSampleHandler() *CRUDSample {
	return new(CRUDSample)
}

type CRUDSample struct {
	ctx context.Context
	w   http.ResponseWriter
	r   *http.Request
}

func (c *CRUDSample) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	c.ctx, c.w, c.r = appengine.NewContext(r), w, r
	switch r.Method {
	case http.MethodGet:
		c.Get()
	case http.MethodPost:
		c.Post()
	default:
		w.WriteHeader(http.StatusNotFound)
		fmt.Fprintln(w, "Not found.")
	}
}

type RespGetFoo struct {
	Id  int64  `json:"id"`
	Foo string `json:"foo"`
	Bar string `json:"bar"`
}

func (c *CRUDSample) Get() {
	var foos []Foo
	keys, err := datastore.NewQuery(KindFoo).GetAll(c.ctx, &foos)
	if err != nil {
		c.responseServerError(err)
		return
	}
	result := make([]RespGetFoo, len(keys))
	for i, f := range foos {
		result[i] = RespGetFoo{keys[i].IntID(), f.Foo, f.Bar}
	}
	c.responseJson(http.StatusOK, result)
}

func (c *CRUDSample) Post() {
	var foo Foo
	err := json.NewDecoder(c.r.Body).Decode(&foo)
	if err != nil {
		c.responseServerError(err)
		return
	}
	key := datastore.NewIncompleteKey(c.ctx, KindFoo, nil)
	newKey, err := datastore.Put(c.ctx, key, &foo)
	if err != nil {
		c.responseServerError(err)
		return
	}
	response := map[string]int64{"id": newKey.IntID()}
	c.responseJson(http.StatusCreated, response)
}

func (c *CRUDSample) responseJson(status int, v interface{}) {
	c.w.WriteHeader(status)
	json.NewEncoder(c.w).Encode(v)
}

func (c *CRUDSample) responseServerError(err error) {
	c.w.WriteHeader(http.StatusInternalServerError)
	fmt.Fprintln(c.w, err)
}
