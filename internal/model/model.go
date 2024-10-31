package model

import (
	"database/sql/driver"
	"fmt"
	"time"
)

const TimeFormat = "2006-01-02 15:04:05"

// MyTime 自定义时间类型
type MyTime struct {
	time.Time
}

// Scan 实现 database/sql Scanner 接口
func (mt *MyTime) Scan(value interface{}) error {
	if value == nil {
		mt.Time = time.Time{} // 设置为零值，相当于 NULL
		return nil
	}
	v, ok := value.(time.Time)
	if !ok {
		return fmt.Errorf("failed to scan MyTime: %v", value)
	}
	mt.Time = v
	return nil
}

// Value 实现 database/sql Valuer 接口
func (mt MyTime) Value() (driver.Value, error) {
	if mt.IsZero() {
		return nil, nil // 返回 NULL
	}
	return mt.Time.Format("2006-01-02 15:04:05"), nil // 格式化为指定字符串
}

// MarshalJSON 实现 json.Marshaler 接口
func (mt MyTime) MarshalJSON() ([]byte, error) {
	if mt.IsZero() {
		return []byte("null"), nil // 返回 JSON NULL
	}
	return []byte(fmt.Sprintf("\"%s\"", mt.Format("2006-01-02 15:04:05"))), nil // 格式化为指定字符串
}

// UnmarshalJSON 实现 json.Unmarshaler 接口
func (mt *MyTime) UnmarshalJSON(b []byte) error {
	if string(b) == "null" {
		mt.Time = time.Time{} // 设置为零值
		return nil
	}
	// 去掉引号
	s := string(b)
	t, err := time.Parse("2006-01-02 15:04:05", s[1:len(s)-1])
	if err != nil {
		return err
	}
	mt.Time = t
	return nil
}
