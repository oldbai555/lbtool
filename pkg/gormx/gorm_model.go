package gormx

import (
	"context"
	"fmt"
	"github.com/mitchellh/mapstructure"
	"github.com/oldbai555/gorm"
	"github.com/oldbai555/lbtool/log"
	"github.com/oldbai555/lbtool/pkg/lberr"
	"github.com/oldbai555/lbtool/utils"
	"reflect"
	"strconv"
	"strings"
)

type Scope struct {
	db  *gorm.DB
	ctx context.Context
	m   *Model

	size      uint32
	page      uint32
	skipTotal bool
}

type Result struct {
	RowsAffected int64
	Created      bool
}

func NewResult(rowsAffected int64, created bool) *Result {
	return &Result{RowsAffected: rowsAffected, Created: created}
}

type Model struct {
	Table       gorm.Tabler
	NotFoundErr error
	db          *gorm.DB
}

func NewModel(db *gorm.DB, m gorm.Tabler, err error) *Model {
	return &Model{
		Table:       m,
		NotFoundErr: err,
		db:          db,
	}
}

func (f *Model) NewScope(ctx context.Context) *Scope {
	db := f.db.WithContext(ctx).Table(f.Table.TableName())
	return &Scope{
		db:  db,
		ctx: ctx,
		m:   f,
	}
}

func (f *Model) IsNotFoundErr(err error) bool {
	return f.NotFoundErr == err
}

type Opt struct {
	orderByDesc string
	orderByAsc  string
	groupBy     string
	limit       uint32
	offset      uint32
	unScoped    bool
}

func WithOrderByDesc(s string) *Opt {
	return &Opt{
		orderByDesc: s,
	}
}

func WithOrderByAsc(s string) *Opt {
	return &Opt{
		orderByAsc: s,
	}
}

func WithGroupBy(groupBy string) *Opt {
	return &Opt{
		groupBy: groupBy,
	}
}

func WithLimit(v uint32) *Opt {
	return &Opt{
		limit: v,
	}
}

func WithOffset(v uint32) *Opt {
	return &Opt{
		offset: v,
	}
}

func WithUnScoped() *Opt {
	return &Opt{
		unScoped: true,
	}
}

func ProcessOpts(db *Scope, opts ...*Opt) {
	for _, opt := range opts {
		if len(opt.orderByDesc) > 0 {
			db.OrderByDesc(opt.orderByDesc)
			continue
		}
		if len(opt.orderByAsc) > 0 {
			db.OrderByAsc(opt.orderByAsc)
			continue
		}
		if len(opt.groupBy) > 0 {
			db.Group(opt.groupBy)
			continue
		}
		if opt.limit > 0 {
			db.db.Limit(int(opt.limit))
		}
		if opt.offset > 0 {
			db.db.Offset(int(opt.offset))
		}
		if opt.unScoped {
			db.UnScoped()
		}
	}
}

func GenSimpleSqlCond(f, op, val string) string {
	return fmt.Sprintf("`%s` %s '%s'", f, op, val)
}

//func NewList(db *gorm.DB, listOption *lb.ListOption) *Scope {
//	return &Scope{
//		db: db,
//
//		size:      listOption.GetSize(),
//		page:      listOption.GetPage(),
//		skipTotal: listOption.GetSkipTotal(),
//	}
//}

func (p *Scope) DB() *gorm.DB {
	return p.db
}

func (p *Scope) Select(fields []string) *Scope {
	p.db.Select(fields)
	return p
}

func (p *Scope) Eq(field string, v interface{}) *Scope {
	p.db.Where(fmt.Sprintf("`%s` = ?", field), v)
	return p
}

func (p *Scope) NotEq(f string, v interface{}) *Scope {
	p.db.Where(fmt.Sprintf("`%s` != ?", f), v)
	return p
}

// AndMap 示例
// key: "name" val: "hangman" sql = `name` = "hangman"
// key: "name like" val: %hangman% sql = `name` like "%hangman%"
func (p *Scope) AndMap(kv map[string]interface{}) *Scope {
	if len(kv) > 0 {
		var condList []string
		var argList []interface{}
		for k, v := range kv {
			if k == "" {
				panic(any("invalid empty key"))
			}
			split := strings.Split(k, " ")
			if len(split) == 2 {
				condList = append(condList, fmt.Sprintf("(`%s` %s ?)", split[0], split[1]))
			} else {
				condList = append(condList, fmt.Sprintf("`%s` = ?", k))
			}
			argList = append(argList, v)
		}
		cond := strings.Join(condList, " AND ")
		p.db.Where(cond, argList...)
	}
	return p
}

// OrMap 示例
// key: "name" val: "hangman" sql = `name` = "hangman"
// key: "name like" val: %hangman% sql = `name` like "%hangman%"
func (p *Scope) OrMap(kv map[string]interface{}) *Scope {
	if len(kv) > 0 {
		var condList []string
		var argList []interface{}
		for k, v := range kv {
			if k == "" {
				panic(any("invalid empty key"))
			}
			split := strings.Split(k, " ")
			if len(split) == 2 {
				condList = append(condList, fmt.Sprintf("(`%s` %s ?)", split[0], split[1]))
			} else {
				condList = append(condList, fmt.Sprintf("(`%s` = ?)", k))
			}
			argList = append(argList, v)
		}
		cond := strings.Join(condList, " OR ")
		p.db.Where(cond, argList...)
	}
	return p
}

func (p *Scope) Like(f, v string) *Scope {
	if v != "" {
		v := utils.EscapeMysqlLikeWildcardIgnore2End(v)
		p.db.Where(fmt.Sprintf("`%s` LIKE ?", f), v)
	}
	return p
}

func (p *Scope) NotLike(f, v string) *Scope {
	if v != "" {
		v = utils.EscapeMysqlLikeWildcardIgnore2End(v)
		v = utils.QuoteName(fmt.Sprintf("%%%s%%", v))
		p.db.Where(
			fmt.Sprintf("`%s` NOT LIKE %s", f, v))
	}
	return p
}

func (p *Scope) In(f string, i interface{}) *Scope {
	v := reflect.ValueOf(i)
	if v.Type().Kind() != reflect.Slice {
		panic(any("invalid input type, slice"))
	}
	if v.Len() == 0 {
		p.db.Where("1=0")
		return p
	}
	p.db.Where(fmt.Sprintf("`%s` in (?)", f), utils.UniqueSliceV2(i))
	return p
}

func (p *Scope) NotIn(f string, i interface{}) *Scope {
	v := reflect.ValueOf(i)
	// 如果不是slice，也是可以的，比如 id in (1)
	if v.Type().Kind() != reflect.Slice {
		panic(any("invalid input type, slice"))
	}
	if v.Len() == 0 {
		p.db.Where("1=0")
		return p
	}
	p.db.Where(fmt.Sprintf("`%s` not in (?)", f), utils.UniqueSliceV2(i))
	return p
}

func (p *Scope) Lt(f string, v interface{}) *Scope {
	p.db.Where(fmt.Sprintf("`%s` < ?", f), v)
	return p
}

func (p *Scope) Lte(f string, v interface{}) *Scope {
	p.db.Where(fmt.Sprintf("`%s` <= ?", f), v)
	return p
}

func (p *Scope) Gt(f string, v interface{}) *Scope {
	p.db.Where(fmt.Sprintf("`%s` > ?", f), v)
	return p
}

func (p *Scope) Gte(f string, v interface{}) *Scope {
	p.db.Where(fmt.Sprintf("`%s` >= ?", f), v)
	return p
}

func (p *Scope) Order(order string) *Scope {
	p.db.Order(order)
	return p
}

func (p *Scope) OrderByDesc(order ...string) *Scope {
	p.db.Order(fmt.Sprintf("%s DESC", strings.Join(order, ",")))
	return p
}

func (p *Scope) OrderByAsc(order ...string) *Scope {
	p.db.Order(fmt.Sprintf("%s ASC", strings.Join(order, ",")))
	return p
}

func (p *Scope) Group(group ...string) *Scope {
	for _, s := range group {
		p.db.Group(s)
	}
	return p
}

// Between 相当于 field >= min || field <= max
func (p *Scope) Between(fieldName string, min, max interface{}) *Scope {
	p.db.Where(fmt.Sprintf("%s BETWEEN ? AND ?", utils.QuoteFieldName(fieldName)), min, max)
	return p
}

// NotBetween 相当于 field < min || field > max
func (p *Scope) NotBetween(fieldName string, min, max interface{}) *Scope {
	p.db.Where(fmt.Sprintf("%s NOT BETWEEN ? AND ?", utils.QuoteFieldName(fieldName)), min, max)
	return p
}

func getFirstInvalidFieldNameCharIndex(s string) int {
	for i := 0; i < len(s); i++ {
		c := s[i]
		if !((c >= 'a' && c <= 'z') ||
			(c >= 'A' && c <= 'Z') ||
			(c >= '0' && c <= '9') ||
			c == '_') {
			return i
		}
	}
	return -1
}

func getOp(fieldName string) (newFieldName, op string) {
	op = "="
	newFieldName = fieldName
	idx := getFirstInvalidFieldNameCharIndex(fieldName)
	if idx > 0 {
		o := strings.TrimSpace(fieldName[idx:])
		newFieldName = fieldName[:idx]
		if o != "" {
			op = o
		}
	}
	return
}

func simpleTypeToStr(value interface{}, quoteSlice bool) string {
	if value == nil {
		panic("value nil")
	}
	vo := reflect.ValueOf(value)
	for vo.Kind() == reflect.Ptr || vo.Kind() == reflect.Interface {
		vo = vo.Elem()
	}
	value = vo.Interface()
	switch v := value.(type) {
	case string:
		v = utils.EscapeMysqlString(v)
		return v
	case []byte:
		s := utils.EscapeMysqlString(string(v))
		return s
	case bool:
		if v {
			return "1"
		} else {
			return "0"
		}
	}
	// 容器单独处理
	switch vo.Kind() {
	case reflect.Slice, reflect.Array:
		var elList []string
		count := vo.Len()
		for x := 0; x < count; x++ {
			el := vo.Index(x)
			elList = append(elList, simpleTypeToStr(el.Interface(), quoteSlice))
		}
		res := strings.Join(elList, ",")
		if quoteSlice {
			res = fmt.Sprintf("(%s)", res)
		}
		return res
	case reflect.Uint32, reflect.Uint64, reflect.Uint16, reflect.Uint8, reflect.Uint:
		return strconv.FormatUint(vo.Uint(), 10)
	case reflect.Int32, reflect.Int64, reflect.Int16, reflect.Int8, reflect.Int:
		return strconv.FormatInt(vo.Int(), 10)
	}
	return fmt.Sprintf("%v", value)
}

func (p *Scope) Where(args ...interface{}) *Scope {
	if len(args) == 0 {
		return p
	}
	arg0 := reflect.ValueOf(args[0])
	for arg0.Kind() == reflect.Interface || arg0.Kind() == reflect.Ptr {
		arg0 = arg0.Elem()
	}
	switch arg0.Kind() {
	case reflect.Bool:
		v := arg0.Bool()
		if v {
			p.db.Where("(?=?)", 1, 1)
		} else {
			p.db.Where("(?=?)", 1, 0)
		}
	case reflect.String:
		fieldName := arg0.String()
		if strings.HasPrefix(fieldName, "$") {
			if len(args) != 2 {
				panic(fmt.Sprintf("invalid number of args %d for $... cond, expected 2", len(args)))
			}
			p.db.Where(fieldName[1:], args[1])
			break
		}
		if strings.IndexByte(fieldName, '?') >= 0 {
			p.db.Where(fieldName, args[1:]...)
			break
		}
		var op string
		var val interface{}
		if len(args) == 2 {
			fieldName, op = getOp(fieldName)
			val = args[1]
			p.db.Where(fmt.Sprintf("%s %s ?", utils.QuoteName(fieldName), op), val)
		} else if len(args) == 3 {
			vo := reflect.ValueOf(args[1])
			if vo.Kind() == reflect.String {
				op = vo.String()
			} else if vo.Kind() == reflect.Int32 {
				// 可以支持 '>' 单括号写法
				op = strings.TrimSpace(fmt.Sprintf("%c", int(vo.Int())))
				if op == "" {
					panic(fmt.Sprintf("invalid op type with int %d", vo.Int()))
				}
			} else {
				panic(fmt.Sprintf("invalid op type %v", vo.Type()))
			}
			val = args[2]
			p.db.Where(fmt.Sprintf("%s %s ?", utils.QuoteName(fieldName), op), val)
		} else if len(args) == 1 {
			p.db.Where(fieldName)
		} else {
			panic(fmt.Sprintf("invalid number of where args %d by `string` prefix", len(args)))
		}
	case reflect.Map:
		typ := arg0.Type()
		if typ.Key().Kind() != reflect.String {
			panic(fmt.Sprintf("map key type required string, but got %v", typ.Key()))
		}
		for _, k := range arg0.MapKeys() {
			fieldName := k.String()
			val := arg0.MapIndex(k)
			if !val.IsValid() || !val.CanInterface() {
				panic(fmt.Sprintf("invalid map val for field %s", fieldName))
			}
			var op string
			fieldName, op = getOp(fieldName)
			log.Infof("val is %s", simpleTypeToStr(val, true))
			p.db.Where(fmt.Sprintf("%s %s ?", utils.QuoteName(fieldName), op), simpleTypeToStr(val, true))
		}
	case reflect.Slice, reflect.Array:
		n := arg0.Len()
		if n == 0 {
			break
		}
		p.db.Where(arg0)
	}
	return p
}

// UnScoped 去除逻辑删除条件
func (p *Scope) UnScoped() {
	p.db.Unscoped()
}

// First 查找
func (p *Scope) First(dest interface{}) error {
	err := p.db.First(dest).Error
	if gorm.ErrRecordNotFound == err {
		return p.m.NotFoundErr
	}
	return err
}

// Update 更新
func (p *Scope) Update(values map[string]interface{}) (*Result, error) {
	res := p.db.Updates(values)
	if res.Error != nil {
		log.Errorf("err is %v", res.Error)
		return nil, res.Error
	}
	return NewResult(res.RowsAffected, false), nil
}

// Delete 删除
func (p *Scope) Delete(obj interface{}) (*Result, error) {
	res := p.db.Delete(obj)
	if res.Error != nil {
		log.Errorf("err is %v", res.Error)
		return nil, res.Error
	}
	return NewResult(res.RowsAffected, false), nil
}

// Create 创建
func (p *Scope) Create(dest interface{}) (*Result, error) {
	res := p.db.Create(dest)
	if res.Error != nil {
		log.Errorf("err is %v", res.Error)
		return nil, res.Error
	}
	return NewResult(res.RowsAffected, false), nil
}

// Find 查找所有
func (p *Scope) Find(dest interface{}) error {
	err := p.db.Find(dest).Error
	if err != nil {
		log.Errorf("err is %v", err)
		return err
	}
	return nil
}

// FindPaginate 分页查找
//func (p *Scope) FindPaginate(list interface{}) (*lb.Paginate, error) {
//	var total int64
//	if !p.skipTotal {
//		err := p.db.Count(&total).Error
//		if err != nil {
//			log.Errorf("err is %v", err)
//			return nil, err
//		}
//	}
//
//	var page = uint32(0)
//	if p.page-1 > 0 {
//		page = p.page - 1
//	}
//
//	err := p.db.Limit(int(p.size)).Offset(int(page * p.size)).Find(list).Error
//	if err != nil {
//		log.Errorf("err is %v", err)
//		return nil, err
//	}
//
//	return &lb.Paginate{
//		Total: uint64(total),
//		Size:  p.size,
//		Page:  p.page,
//	}, nil
//}

func (p *Scope) Increment(field string, num uint32, candMap map[string]interface{}) error {
	if len(candMap) == 0 {
		return lberr.NewInvalidArg("candMap must be not nil")
	}
	_, err := p.AndMap(candMap).Update(map[string]interface{}{
		field: gorm.Expr("? + ?", field, num),
	})
	if err != nil {
		log.Errorf("err is %v", err)
		return err
	}
	return nil
}

func (p *Scope) Decrement(ctx context.Context, field string, num uint32, candMap map[string]interface{}) error {
	if len(candMap) == 0 {
		return lberr.NewInvalidArg("candMap must be not nil")
	}
	_, err := p.AndMap(candMap).Update(map[string]interface{}{
		field: gorm.Expr("? - ?", field, num),
	})
	if err != nil {
		log.Errorf("err is %v", err)
		return err
	}
	return nil
}

func (p *Scope) UpdateOrCreate(ctx context.Context, candMap, attrMap map[string]interface{}, out interface{}) (*Result, error) {
	err := p.AndMap(candMap).First(out)
	if err != nil && err != gorm.ErrRecordNotFound {
		log.Errorf("err:%v", err)
		return nil, err
	}

	optDb := p.m.NewScope(ctx)
	if err == gorm.ErrRecordNotFound {
		for k, v := range candMap {
			attrMap[k] = v
		}
		err = mapstructure.Decode(attrMap, out)
		if err != nil {
			log.Errorf("err is : %v", err)
			return nil, err
		}

		rows, err := optDb.Create(out)
		if err != nil {
			log.Errorf("err is %v", err)
			return nil, err
		}
		return rows, nil
	}

	rows, err := optDb.AndMap(candMap).Update(attrMap)
	if err != nil {
		log.Errorf("err:%v", err)
		return nil, err
	}

	err = mapstructure.Decode(attrMap, out)
	if err != nil {
		log.Errorf("err is : %v", err)
		return nil, err
	}

	return rows, nil
}

func (p *Scope) SetLimit(limit int) *Scope {
	p.db.Limit(limit)
	return p
}

func (p *Scope) SetOffset(offset int) *Scope {
	p.db.Offset(offset)
	return p
}
