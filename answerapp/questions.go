package answerapp

import (
	"context"
	"errors"
	"google.golang.org/appengine/datastore"
	"google.golang.org/appengine/log"
	"time"
)

type Question struct {
	Key          *datastore.Key `json:"id" datastore:"-"`
	CTime        time.Time      `json:"created"`
	Question     string         `json:"question"`
	User         UserCard       `json:"user"`
	AnswersCount int            `json:"answers_count"`
}

// 구글 클라우드 데이터스토어 데이터 가져오기
func (q Question) OK() error {
	if len(q.Question) < 10 {
		return errors.New("question is too short(질문이 너무 짧음)")
	}
	return nil
}

func (q *Question) Create(ctx context.Context) error {
	log.Debugf(ctx, "Saving question: %s", q.Question)
	if q.Key == nil {
		q.Key = datastore.NewIncompleteKey(ctx, "Question", nil)
	}
	user, err := UserFromAEUser(ctx)
	if err != nil {
		return err
	}
	q.User = user.Card()
	q.CTime = time.Now()
	q.Key, err = datastore.Put(ctx, q.Key, q)
	if err != nil {
		return err
	}
	return nil
}

func (q *Question) Update(ctx context.Context) error {
	if q.Key == nil {
		q.Key = datastore.NewIncompleteKey(ctx, "Question", nil)
	}
	var err error
	q.Key, err = datastore.Put(ctx, q.Key, q)
	if err != nil {
		return err
	}
	return nil
}

// 구글 클라우드 데이터스토어 데이터 읽기
func GetQuestion(ctx context.Context, key *datastore.Key) (*Question, error) {
	var q Question
	err := datastore.Get(ctx, key, &q)
	if err != nil {
		return nil, err
	}
	q.Key = key

	return &q, nil
}

func ToQuestion(ctx context.Context) ([]*Question, error) {
	var questions []*Question
	questionKeys, err := datastore.NewQuery("Question").
		Order("-AnswersCount").Order("-CTime").Limit(25).GetAll(ctx, &questions)
	if err != nil {
		return nil, err
	}
	for i := range questions {
		questions[i].Key = questionKeys[i]
	}
	return questions, nil
}
