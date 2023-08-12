package filters

import (
	"github.com/tempcke/rpm/entity"
)

type TenantFilter struct {
	IDs []entity.ID
}

func NewTenantFilter() TenantFilter {
	return TenantFilter{}
}
