package main

import (
	"gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
	"net/http"
)

// 엔드포인트 헨들링

type poll struct {
	ID      bson.ObjectId  `bson:"_id" json:"id"`
	Title   string         `json:"title"`
	Options []string       `json:"options"`
	Results map[string]int `json:"results,omitempty"`
	APIKey  string         `json:"apikey"`
}

func (s *Server) handlePolls(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case "GET":
		s.handlePollsGet(w, r)
		return
	case "POST":
		s.handlePollsPost(w, r)
		return
	case "DELETE":
		s.handlePollsDelete(w, r)
		return
	case "OPTIONS":
		// CORS 지원
		w.Header().Add("Access-Control-Allow-Methods", "DELETE")
		respond(w, r, http.StatusOK, nil)
		return
	}
	// 찾지 못함
	respondHTTPErr(w, r, http.StatusNotFound)
}

// 투표 읽기
func (s *Server) handlePollsGet(w http.ResponseWriter, r *http.Request) {
	session := s.db.Copy()
	defer session.Close()

	c := session.DB("ballots").C("polls")

	var q *mgo.Query
	p := NewPath(r.URL.Path)
	if p.HasID() {
		// 특정 투표 가져오기
		q = c.FindId(bson.ObjectIdHex(p.ID))
	} else {
		// 모든 투표 가져오기
		q = c.Find(nil)
	}

	var result []*poll

	if err := q.All(&result); err != nil {
		respondErr(w, r, http.StatusInternalServerError, err)
		return
	}

	respond(w, r, http.StatusOK, nil)
}

// 투표 생성
func (s *Server) handlePollsPost(w http.ResponseWriter, r *http.Request) {
	session := s.db.Copy()
	defer session.Close()

	c := session.DB("ballots").C("polls")

	var p poll

	if err := decodeBody(r, &p); err != nil {
		respondErr(w, r, http.StatusBadRequest, "failed to read poll from request", err)
		return
	}

	apikey, ok := APIKey(r.Context())
	if ok {
		p.APIKey = apikey
	}

	p.ID = bson.NewObjectId()
	if err := c.Insert(p); err != nil {
		respondErr(w, r, http.StatusInternalServerError, "failed to read poll from request", err)
		return
	}

	w.Header().Set("Location", "polls/"+p.ID.Hex())
	respond(w, r, http.StatusCreated, nil)
}

// 투표 삭제
func (s *Server) handlePollsDelete(w http.ResponseWriter, r *http.Request) {
	session := s.db.Copy()
	defer session.Close()

	c := session.DB("ballots").C("polls")

	p := NewPath(r.URL.Path)
	if !p.HasID() {
		respondErr(w, r, http.StatusMethodNotAllowed, "Cannot delete all polls.")
		return
	}

	if err := c.RemoveId(bson.ObjectIdHex(p.ID)); err != nil {
		respondErr(w, r, http.StatusInternalServerError, "failed to delete poll", err)
		return
	}

	respond(w, r, http.StatusOK, nil) // ok
}
