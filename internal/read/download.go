package read

import (
	"aozorandom-cobra/internal/cache"
	"archive/zip"
	"bytes"
	"encoding/csv"
	"fmt"
	"io"
	"net/http"
	"strings"

	"golang.org/x/text/encoding/japanese"
	"golang.org/x/text/transform"
)

type Book struct {
	Author    string
	Title     string
	Publisher string
}

const URL = "https://www.aozora.gr.jp/index_pages/list_person_all.zip"

func downloadData() ([]byte, error) {
	resp, err := http.Get(URL)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	zipReader, err := zip.NewReader(bytes.NewReader(body), int64(len(body)))
	if err != nil {
		return nil, err
	}

	for _, f := range zipReader.File {
		if !strings.HasSuffix(f.Name, ".csv") {
			continue
		}

		rc, err := f.Open()
		if err != nil {
			return nil, err
		}
		defer rc.Close()

		return io.ReadAll(rc)
	}

	return nil, fmt.Errorf("ZIPの中にファイルが存在しませんでした。")
}

// FetchData は、データを取得します
func FetchData() ([]byte, error) {
	path, err := cache.CachePath()
	if err != nil {
		return nil, err
	}

	if cache.IsCacheValid(path) {
		fmt.Println("キャッシュから読み込み中")
		return cache.LoadCache(path)
	}

	fmt.Println("作品リストをダウンロード中")
	data, err := downloadData()
	if err != nil {
		return nil, err
	}

	if err := cache.SaveData(path, data); err != nil {
		return nil, err
	}

	return data, nil
}

// ParseCSV は、CSVファイルからデータを抽出し必要なデータを返却します
func ParseCSV(r io.Reader) ([]Book, error) {
	reader := transform.NewReader(r, japanese.ShiftJIS.NewDecoder())
	csvReader := csv.NewReader(reader)
	csvReader.FieldsPerRecord = -1

	if _, err := csvReader.Read(); err != nil {
		return nil, err
	}

	var books []Book
	for {
		record, err := csvReader.Read()
		if err == io.EOF {
			break
		}
		if err != nil {
			return nil, err
		}
		books = append(books, Book{
			Author:    record[1],
			Title:     record[3],
			Publisher: record[11],
		})
	}

	return books, nil
}