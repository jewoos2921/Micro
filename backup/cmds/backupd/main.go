package main

import (
	"Micro/backup"
	"encoding/json"
	"errors"
	"flag"
	"fmt"
	"github.com/matryer/filedb"
	"log"
	"os"
	"os/signal"
	"syscall"
	"time"
)

type path struct {
	Path string
	Hash string
}

func main() {
	// 데몬 백업 툴
	var fatalErr error
	defer func() {
		if fatalErr != nil {
			log.Fatalln(fatalErr)
		}
	}()
	var (
		// 폴더가 변경됐는지 여부를 확인
		interval = flag.Duration("interval", 10*time.Second, "interval between check")
		// ZIP파일이 저장될 보관 위치의 경로
		archive = flag.String("archive", "archive", "path to archive location")
		// backup 명령이 있는 동일한 filedb 데이터베이스의 상호작용하는 경로
		dbpath = flag.String("db", "./db", "path to filedb database")
	)
	flag.Parse()
	m := &backup.Monitor{Destination: *archive,
		Archiver: backup.ZIP, Paths: make(map[string]string),
	}

	db, err := filedb.Dial(*dbpath)
	if err != nil {
		fatalErr = err
		return
	}

	defer db.Close()
	col, err := db.C("paths")
	if err != nil {
		fatalErr = err
		return
	}
	// 데이터 캐싱
	var path path
	col.ForEach(func(_ int, data []byte) bool {
		if err := json.Unmarshal(data, &path); err != nil {
			fatalErr = err
			return true
		}
		m.Paths[path.Path] = path.Hash
		return false // 계속한다.
	})
	if fatalErr != nil {
		return
	}
	if len(m.Paths) < 1 {
		fatalErr = errors.New("no paths - use backup tool to add at least one")
		return
	}
	check(m, col)
	signalChan := make(chan os.Signal, 1)
	signal.Notify(signalChan, syscall.SIGINT, syscall.SIGTERM)
	for {
		select {
		case <-time.After(*interval):
			check(m, col)
		case <-signalChan:
			// 중지
			fmt.Println()
			log.Printf("Stopping...")
			return
		}
	}
}

// filedb 레코드 업데이트
func check(m *backup.Monitor, col *filedb.C) {
	log.Println("Checking...")
	counter, err := m.Now()
	if err != nil {
		log.Fatalln("failed to backup:", err)
	}
	if counter > 0 {
		log.Printf(" Archived %d directories\n", counter)
		// 해시 업데이트
		var path path
		col.SelectEach(func(_ int, data []byte) (bool, []byte, bool) {
			if err := json.Unmarshal(data, &path); err != nil {
				log.Println("failed to unmarshal data (skipping):", err)
				return true, data, false
			}
			path.Hash, _ = m.Paths[path.Path]
			newdata, err := json.Marshal(&path)
			if err != nil {
				log.Println("failed to marshal data (skipping):", err)
				return true, data, false
			}
			return true, newdata, false
		})
	} else {
		log.Println("No changes")
	}
}
