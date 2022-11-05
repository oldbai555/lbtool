package webtool

// SortField 排序字段
// id
// created_at
// updated_at
type SortField string

const (
	SortFieldID        SortField = "id"
	SortFieldCreatedAt SortField = "created_at"
	SortFieldUpdatedAt SortField = "updated_at"
)

// SortType 排序方式，升序，降序，默认降序
// asc
// desc
type SortType string

const (
	SortTypeAsc  SortType = "asc"
	SortTypeDesc SortType = "desc"
)

type Sorter struct {
	// SortField 排序字段,id created_at updated_at sort_weight total today_join_cnt today_drop_out_cnt relation_delete_at relation_create_at msg_time
	SortField SortField `form:"sort_field" json:"sort_field" gorm:"-" validate:"omitempty,oneof=id created_at updated_at sort_weight add_customer_count total today_join_cnt today_drop_out_cnt createtime customer_delete_staff_at relation_delete_at relation_create_at in_connection_time_range order today_join_member_num today_quit_member_num create_time msg_time"`
	// SortType 排序类型,asc desc
	SortType SortType `form:"sort_type" json:"sort_type" gorm:"-" validate:"omitempty,oneof=asc desc"`
}

func (o *Sorter) SetDefault() *Sorter {
	if o.SortField == "" {
		o.SortField = SortFieldCreatedAt
	}
	if o.SortType == "" {
		o.SortType = SortTypeDesc
	}
	return o
}
