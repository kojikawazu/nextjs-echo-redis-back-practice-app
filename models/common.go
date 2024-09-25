package models

import (
	"fmt"
	"time"
)

// カスタムの時間フォーマットを扱うための型
type CustomTime struct {
	time.Time
}

// time.Time から CustomTime を生成する
func NewCustomTime(t time.Time) CustomTime {
	return CustomTime{Time: t}
}

// CustomTime を JSON にフォーマット
func (ct CustomTime) MarshalJSON() ([]byte, error) {
	formatted := fmt.Sprintf("\"%s\"", ct.Format("2006-01-02 15:04:05.999999"))
	return []byte(formatted), nil
}

// JSON から CustomTime をパース
func (ct *CustomTime) UnmarshalJSON(b []byte) error {
	// クオートを除去
	strInput := string(b)
	strInput = strInput[1 : len(strInput)-1]

	// フォーマットを2つ用意
	const layoutWithT = "2006-01-02T15:04:05.999999"
	const layoutWithSpace = "2006-01-02 15:04:05.999999"

	// "T" を含むフォーマットでまずパース
	parsedTime, err := time.Parse(layoutWithT, strInput)
	if err != nil {
		// スペース区切りのフォーマットでもパースを試みる
		parsedTime, err = time.Parse(layoutWithSpace, strInput)
		if err != nil {
			return err
		}
	}

	ct.Time = parsedTime
	return nil
}
