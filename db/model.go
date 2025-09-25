package db

import (
	"github.com/RyCarlos/go-common/utils/snowflake"
	"gorm.io/gorm"
	"strconv"
	"strings"
	"time"
)

type Id int64

// UnmarshalJSON JSON返回序列化的时候转为INT4
func (id *Id) UnmarshalJSON(data []byte) error {
	// 去除引号
	str := strings.Trim(string(data), `"`)

	if str == "" || str == "null" {
		*id = 0
		return nil
	}

	// 转换为数字
	num, err := strconv.ParseInt(str, 10, 64)
	if err != nil {
		return err
	}

	*id = Id(num)
	return nil
}

// MarshalJSON int64转为字符串 避免精度丢失
func (id *Id) MarshalJSON() ([]byte, error) {
	if *id == 0 {
		return []byte("null"), nil
	}
	// 直接返回字符串格式的字节切片，不要再次使用 json.Marshal
	return []byte(`"` + strconv.FormatInt(int64(*id), 10) + `"`), nil
}

type AutoId struct {
	Id Id `gorm:"primaryKey" json:"id" gen:"-"` // 主键Id
}

type CreateTime struct {
	CreatedAt time.Time `json:"createdAt" gorm:"comment:创建时间" gen:"-"` // 创建时间
}

type UpdateTime struct {
	UpdatedAt time.Time `json:"updatedAt" gorm:"comment:更新时间" gen:"-"` // 更新时间
}

type SoftDelete struct {
	DeletedAt gorm.DeletedAt `json:"-" gorm:"index;comment:删除时间" gen:"-"` //删除时间
}

func (a *AutoId) BeforeCreate(*gorm.DB) error {
	if a.Id == 0 {
		a.Id = Id(snowflake.BuildId())
	}
	return nil
}
