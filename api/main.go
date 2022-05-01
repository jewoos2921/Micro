package main

import (
	"context"
	"flag"
	"gopkg.in/mgo.v2"
	"log"
	"net/http"
)

// RESTful API 설계
//====================================================================================================
// HTTP 메소드는 수행할 작업의 종류를 설명한다. 예로, GET 메소드는 데이터를 읽는 반면 POST 요청은 무언가를 생성한다.
// 데이터는 리소스의 모음으로 표현된다.
// 액션은 데이터 변화로 표현된다.
// URL 은 특정 데이터를 참조하는 데 사용된다.
// HTTP 헤더는 서버와 주고받는 표현의 종류를 설명하는 데 사용된다.
//====================================================================================================
func main() {
	var (
		addr  = flag.String("addr", ":8080", "endpoint address")
		mongo = flag.String("mongo", "localhost", "mongodb address")
	)
	log.Println("Dialing mongo", *mongo)
	db, err := mgo.Dial(*mongo)
	if err != nil {
		log.Fatalln("failed to connect ot mongo:", err)
	}
	defer db.Close()
	s := &Server{db: db}
	mux := http.NewServeMux()
	mux.HandleFunc("/polls/", withCORS(withAPIKey(s.handlePolls)))
	/*
		1. 적절한 헤더를 설정하는 withCORS함수가 호출된다.
		2. withAPIKey함수가 다음을 호출해 APIKey에 대한 요청을 확인하고 유효하지 않는 경우 중단하거나
			다음 핸들러 함수를 호출한다.
		3. 그러면 respond.go의 헬퍼 함수를 사용해 클라이언트에 응답을 쓸 수 있는 handlePolls 함수가 호출된다.
		4. 실행은 APIKey로 돌아가고 종료된다.
		5. 실행은 마침내 종료된 withCORS로 다시 돌아간다.
	*/
	log.Println("Starting web server on", *addr)
	http.ListenAndServe(":8080", mux)
	log.Println("Stopping...")
}

// API 서버다.
type Server struct {
	db *mgo.Session
}

type contextKey struct {
	name string
}

var contextKeyAPIKey = &contextKey{"api-key"}

func APIKey(ctx context.Context) (string, bool) {
	key, ok := ctx.Value(contextKeyAPIKey).(string)
	return key, ok
}

func withAPIKey(fn http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		key := r.URL.Query().Get("key")
		if !isValidAPIKey(key) {
			respondErr(w, r, http.StatusUnauthorized, "invalid API key")
			return
		}
		ctx := context.WithValue(r.Context(), contextKeyAPIKey, key)
		fn(w, r.WithContext(ctx))
	}
}

func isValidAPIKey(key string) bool {
	return key == "abc123"
}
func withCORS(fn http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Access-Control-Expose-Headers", "Location")
		fn(w, r)
	}
}
