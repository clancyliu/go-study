package model

import (
	"github.com/zeromicro/go-zero/core/stores/sqlx"
)

var _ TestTabModel = (*customTestTabModel)(nil)

type (
	// TestTabModel is an interface to be customized, add more methods here,
	// and implement the added methods in customTestTabModel.
	TestTabModel interface {
		testTabModel
	}

	customTestTabModel struct {
		*defaultTestTabModel
	}
)

// NewTestTabModel returns a model for the database table.
func NewTestTabModel(conn sqlx.SqlConn) TestTabModel {
	return &customTestTabModel{
		defaultTestTabModel: newTestTabModel(conn),
	}
}
