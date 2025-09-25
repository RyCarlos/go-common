package sqls

import (
	"github.com/RyCarlos/go-common/log"
	"gorm.io/gorm"
)

type Query struct {
	QueryFields []string
	Conditions  []Condition
	Orders      []Order
	Paging      *Paging
}

type Order struct {
	Column string
	Asc    bool
}

type Condition struct {
	Expr string
	Args []interface{}
}

func NewQuery() *Query {
	return &Query{}
}

func (q *Query) Fields(fields ...string) *Query {
	if len(fields) > 0 {
		q.QueryFields = append(q.QueryFields, fields...)
	}
	return q
}

func (q *Query) Eq(column string, value interface{}) *Query {
	q.Where(getColumn(column)+" = ?", value)
	return q
}

func (q *Query) NotEq(column string, value interface{}) *Query {
	q.Where(getColumn(column)+" <> ?", value)
	return q
}

func (q *Query) Gt(column string, value interface{}) *Query {
	q.Where(getColumn(column)+" > ?", value)
	return q
}

func (q *Query) Gte(column string, value interface{}) *Query {
	q.Where(getColumn(column)+" >= ?", value)
	return q
}

func (q *Query) Lt(column string, value interface{}) *Query {
	q.Where(getColumn(column)+" < ?", value)
	return q
}

func (q *Query) Lte(column string, value interface{}) *Query {
	q.Where(getColumn(column)+" <= ?", value)
	return q
}

func (q *Query) In(column string, value interface{}) *Query {
	q.Where(getColumn(column)+" IN (?)", value)
	return q
}

func (q *Query) NotIn(column string, value interface{}) *Query {
	q.Where(getColumn(column)+" NOT IN (?)", value)
	return q
}

func (q *Query) Or(column string, args ...interface{}) *Query {
	q.Where(getColumn(column)+" OR ?", args)
	return q
}

func (q *Query) Like(column string, value string) *Query {
	q.Where(getColumn(column)+" LIKE ?", "%"+value+"%")
	return q
}

func (q *Query) LeftLike(column string, value string) *Query {
	q.Where(getColumn(column)+" LIKE ?", "%"+value)
	return q
}

func (q *Query) RightLike(column string, value string) *Query {
	q.Where(getColumn(column)+" LIKE ?", value+"%")
	return q
}

func (q *Query) Between(column string, start interface{}, end interface{}) *Query {
	q.Where(getColumn(column)+" BETWEEN ? AND ?", start, end)
	return q
}

func (q *Query) Where(expr string, args ...interface{}) {
	q.Conditions = append(q.Conditions, Condition{Expr: expr, Args: args})
}

func (q *Query) Limit(limit int) *Query {
	q.Page(1, limit)
	return q
}

func (q *Query) Page(page, limit int) *Query {
	if q.Paging == nil {
		q.Paging = &Paging{Page: page, Limit: limit}
	} else {
		q.Paging.Page = page
		q.Paging.Limit = limit
	}
	return q
}

func (q *Query) Find(db *gorm.DB, out interface{}) {
	if err := q.Build(db).Find(out).Error; err != nil {
		log.Error("Find query err", err)
	}
}

func (q *Query) FindOne(db *gorm.DB, out interface{}) error {
	if err := q.Build(db).Limit(1).First(out).Error; err != nil {
		return err
	}
	return nil
}

func (q *Query) OrderAsc(column string) *Query {
	q.Order(column, true)
	return q
}

func (q *Query) OrderDesc(column string) *Query {
	q.Order(column, false)
	return q
}

func (q *Query) Order(column string, isAsc bool) {
	q.Orders = append(q.Orders, Order{getColumn(column), isAsc})
}

func (q *Query) Count(db *gorm.DB, model interface{}) int64 {
	_db := db.Model(model)
	// where 条件
	_db = q.buildWhere(_db)
	var count int64
	if err := _db.Count(&count).Error; err != nil {
		log.Error("count err", err)
	}
	return count
}

func (q *Query) Build(db *gorm.DB) *gorm.DB {
	_db := db
	// 查询字段
	if len(q.QueryFields) > 0 {
		_db = _db.Select(q.QueryFields)
	}
	// where 条件
	_db = q.buildWhere(_db)
	// 排序
	if len(q.Orders) > 0 {
		for _, order := range q.Orders {
			if order.Asc {
				_db = _db.Order(order.Column + " ASC")
			} else {
				_db = _db.Order(order.Column + " DESC")
			}
		}
	}
	// 分页
	if q.Paging != nil {
		_db = _db.Offset(q.Paging.Offset()).Limit(q.Paging.Limit)
	}
	return _db
}

func (q *Query) buildWhere(db *gorm.DB) *gorm.DB {
	_db := db
	if len(q.Conditions) > 0 {
		for _, condition := range q.Conditions {
			_db = _db.Where(condition.Expr, condition.Args...)
		}
	}
	return _db
}

func getColumn(column string) string {
	return "`" + column + "`"
}
