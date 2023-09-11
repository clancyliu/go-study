// Code generated by goctl. DO NOT EDIT.

package model

import (
	"context"
	"database/sql"
	"fmt"
	"strings"

	"github.com/zeromicro/go-zero/core/stores/builder"
	"github.com/zeromicro/go-zero/core/stores/sqlc"
	"github.com/zeromicro/go-zero/core/stores/sqlx"
	"github.com/zeromicro/go-zero/core/stringx"
)

var (
	testTabFieldNames          = builder.RawFieldNames(&TestTab{})
	testTabRows                = strings.Join(testTabFieldNames, ",")
	testTabRowsExpectAutoSet   = strings.Join(stringx.Remove(testTabFieldNames, "`id`", "`create_at`", "`create_time`", "`created_at`", "`update_at`", "`update_time`", "`updated_at`"), ",")
	testTabRowsWithPlaceHolder = strings.Join(stringx.Remove(testTabFieldNames, "`id`", "`create_at`", "`create_time`", "`created_at`", "`update_at`", "`update_time`", "`updated_at`"), "=?,") + "=?"
)

type (
	testTabModel interface {
		Insert(ctx context.Context, data *TestTab) (sql.Result, error)
		FindOne(ctx context.Context, id int64) (*TestTab, error)
		Update(ctx context.Context, data *TestTab) error
		Delete(ctx context.Context, id int64) error
	}

	defaultTestTabModel struct {
		conn  sqlx.SqlConn
		table string
	}

	TestTab struct {
		Id         int64        `db:"id"`
		Deleted    BitBool      `db:"deleted"`
		CreateTime sql.NullTime `db:"create_time"`
		UpdateTime sql.NullTime `db:"update_time"`
	}
)

func newTestTabModel(conn sqlx.SqlConn) *defaultTestTabModel {
	return &defaultTestTabModel{
		conn:  conn,
		table: "`test_tab`",
	}
}

func (m *defaultTestTabModel) Delete(ctx context.Context, id int64) error {
	query := fmt.Sprintf("delete from %s where `id` = ?", m.table)
	_, err := m.conn.ExecCtx(ctx, query, id)
	return err
}

func (m *defaultTestTabModel) FindOne(ctx context.Context, id int64) (*TestTab, error) {
	query := fmt.Sprintf("select %s from %s where `id` = ? limit 1", testTabRows, m.table)
	var resp TestTab
	err := m.conn.QueryRowCtx(ctx, &resp, query, id)
	switch err {
	case nil:
		return &resp, nil
	case sqlc.ErrNotFound:
		return nil, ErrNotFound
	default:
		return nil, err
	}
}

func (m *defaultTestTabModel) Insert(ctx context.Context, data *TestTab) (sql.Result, error) {
	query := fmt.Sprintf("insert into %s (%s) values (?)", m.table, testTabRowsExpectAutoSet)
	ret, err := m.conn.ExecCtx(ctx, query, data.Deleted)
	return ret, err
}

func (m *defaultTestTabModel) Update(ctx context.Context, data *TestTab) error {
	query := fmt.Sprintf("update %s set %s where `id` = ?", m.table, testTabRowsWithPlaceHolder)
	_, err := m.conn.ExecCtx(ctx, query, data.Deleted, data.Id)
	return err
}

func (m *defaultTestTabModel) tableName() string {
	return m.table
}
