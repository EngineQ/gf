<<<<<<< HEAD
// Copyright 2017 gf Author(https://gitee.com/johng/gf). All Rights Reserved.
//
// This Source Code Form is subject to the terms of the MIT License.
// If a copy of the MIT was not distributed with this file,
// You can obtain one at https://gitee.com/johng/gf.
=======
// Copyright 2017 gf Author(https://github.com/gogf/gf). All Rights Reserved.
//
// This Source Code Form is subject to the terms of the MIT License.
// If a copy of the MIT was not distributed with this file,
// You can obtain one at https://github.com/gogf/gf.
>>>>>>> upstream/master

package gdb

import (
<<<<<<< HEAD
    "fmt"
    "errors"
    "strings"
    "database/sql"
    _ "github.com/lib/pq"
    _ "github.com/go-sql-driver/mysql"
    "gitee.com/johng/gf/g/util/gconv"
=======
    "database/sql"
    "fmt"
    "github.com/gogf/gf/g/text/gregex"
>>>>>>> upstream/master
    "reflect"
)

// 数据库事务对象
<<<<<<< HEAD
type Tx struct {
    db *Db
    tx *sql.Tx
}

// 事务操作，提交
func (tx *Tx) Commit() error {
=======
type TX struct {
    db     DB
    tx     *sql.Tx
    master *sql.DB
}

// 事务操作，提交
func (tx *TX) Commit() error {
>>>>>>> upstream/master
    return tx.tx.Commit()
}

// 事务操作，回滚
<<<<<<< HEAD
func (tx *Tx) Rollback() error {
=======
func (tx *TX) Rollback() error {
>>>>>>> upstream/master
    return tx.tx.Rollback()
}

// (事务)数据库sql查询操作，主要执行查询
<<<<<<< HEAD
func (tx *Tx) Query(query string, args ...interface{}) (*sql.Rows, error) {
    p         := tx.db.link.handleSqlBeforeExec(&query)
    rows, err := tx.tx.Query(*p, args ...)
    err        = tx.db.formatError(err, p, args...)
    if err == nil {
        return rows, nil
    }
    return nil, err
}

// (事务)执行一条sql，并返回执行情况，主要用于非查询操作
func (tx *Tx) Exec(query string, args ...interface{}) (sql.Result, error) {
    p      := tx.db.link.handleSqlBeforeExec(&query)
    r, err := tx.tx.Exec(*p, args ...)
    err     = tx.db.formatError(err, p, args...)
    return r, err
}

// 数据库查询，获取查询结果集，以列表结构返回
func (tx *Tx) GetAll(query string, args ...interface{}) (Result, error) {
    // 执行sql
=======
func (tx *TX) Query(query string, args ...interface{}) (rows *sql.Rows, err error) {
    return tx.db.doQuery(tx.tx, query, args...)
}

// (事务)执行一条sql，并返回执行情况，主要用于非查询操作
func (tx *TX) Exec(query string, args ...interface{}) (sql.Result, error) {
    return tx.db.doExec(tx.tx, query, args...)
}

// sql预处理，执行完成后调用返回值sql.Stmt.Exec完成sql操作
func (tx *TX) Prepare(query string) (*sql.Stmt, error) {
    return tx.db.doPrepare(tx.tx, query)
}

// 数据库查询，获取查询结果集，以列表结构返回
func (tx *TX) GetAll(query string, args ...interface{}) (Result, error) {
>>>>>>> upstream/master
    rows, err := tx.Query(query, args ...)
    if err != nil || rows == nil {
        return nil, err
    }
<<<<<<< HEAD
    // 列名称列表
    columns, err := rows.Columns()
    if err != nil {
        return nil, err
    }
    // 返回结构组装
    values   := make([]sql.RawBytes, len(columns))
    scanArgs := make([]interface{}, len(values))
    records  := make(Result, 0)
    for i := range values {
        scanArgs[i] = &values[i]
    }
    for rows.Next() {
        err = rows.Scan(scanArgs...)
        if err != nil {
            return records, err
        }
        row := make(Record)
        // 注意col字段是一个[]byte类型(slice类型本身是一个指针)，多个记录循环时该变量指向的是同一个内存地址
        for i, col := range values {
            k := columns[i]
            v := make([]byte, len(col))
            copy(v, col)
            row[k] = v
        }
        //fmt.Printf("%p\n", row["typeid"])
        records = append(records, row)
    }
    return records, nil
}

// 数据库查询，获取查询结果记录，以关联数组结构返回
func (tx *Tx) GetOne(query string, args ...interface{}) (Record, error) {
=======
    defer rows.Close()
    return tx.db.rowsToResult(rows)
}

// 数据库查询，获取查询结果记录，以关联数组结构返回
func (tx *TX) GetOne(query string, args ...interface{}) (Record, error) {
>>>>>>> upstream/master
    list, err := tx.GetAll(query, args ...)
    if err != nil {
        return nil, err
    }
    if len(list) > 0 {
        return list[0], nil
    }
    return nil, nil
}

// 数据库查询，获取查询结果记录，自动映射数据到给定的struct对象中
<<<<<<< HEAD
func (tx *Tx) GetStruct(obj interface{}, query string, args ...interface{}) error {
=======
func (tx *TX) GetStruct(obj interface{}, query string, args ...interface{}) error {
>>>>>>> upstream/master
    one, err := tx.GetOne(query, args...)
    if err != nil {
        return err
    }
    return one.ToStruct(obj)
}

<<<<<<< HEAD

// 数据库查询，获取查询字段值
func (tx *Tx) GetValue(query string, args ...interface{}) (Value, error) {
=======
// 数据库查询，查询多条记录，并自动转换为指定的slice对象, 如: []struct/[]*struct。
func (tx *TX) GetStructs(objPointerSlice interface{}, query string, args ...interface{}) error {
    all, err := tx.GetAll(query, args...)
    if err != nil {
        return err
    }
    return all.ToStructs(objPointerSlice)
}

// 将结果转换为指定的struct/*struct/[]struct/[]*struct,
// 参数应该为指针类型，否则返回失败。
// 该方法自动识别参数类型，调用Struct/Structs方法。
func (tx *TX) GetScan(objPointer interface{}, query string, args ...interface{}) error {
    t := reflect.TypeOf(objPointer)
    k := t.Kind()
    if k != reflect.Ptr {
        return fmt.Errorf("params should be type of pointer, but got: %v", k)
    }
    k = t.Elem().Kind()
    switch k {
        case reflect.Array:
        case reflect.Slice:
            return tx.db.GetStructs(objPointer, query, args ...)
        case reflect.Struct:
            return tx.db.GetStruct(objPointer, query, args ...)
        default:
            return fmt.Errorf("element type should be type of struct/slice, unsupported: %v", k)
    }
    return nil
}

// 数据库查询，获取查询字段值
func (tx *TX) GetValue(query string, args ...interface{}) (Value, error) {
>>>>>>> upstream/master
    one, err := tx.GetOne(query, args ...)
    if err != nil {
        return nil, err
    }
    for _, v := range one {
        return v, nil
    }
    return nil, nil
}

// 数据库查询，获取查询数量
<<<<<<< HEAD
func (tx *Tx) GetCount(query string, args ...interface{}) (int, error) {
    val, err := tx.GetValue(query, args ...)
    if err != nil {
        return 0, err
    }
    return gconv.Int(val), nil
}

// 数据表查询，其中tables可以是多个联表查询语句，这种查询方式较复杂，建议使用链式操作
func (tx *Tx) Select(tables, fields string, condition interface{}, groupBy, orderBy string, first, limit int, args ... interface{}) (Result, error) {
    s := fmt.Sprintf("SELECT %s FROM %s ", fields, tables)
    if condition != nil {
        s += fmt.Sprintf("WHERE %s ", tx.db.formatCondition(condition))
    }
    if len(groupBy) > 0 {
        s += fmt.Sprintf("GROUP BY %s ", groupBy)
    }
    if len(orderBy) > 0 {
        s += fmt.Sprintf("ORDER BY %s ", orderBy)
    }
    if limit > 0 {
        s += fmt.Sprintf("LIMIT %d,%d ", first, limit)
    }
    return tx.GetAll(s, args ... )
}

// sql预处理，执行完成后调用返回值sql.Stmt.Exec完成sql操作
// 记得调用sql.Stmt.Close关闭操作对象
func (tx *Tx) Prepare(query string) (*sql.Stmt, error) {
    return tx.tx.Prepare(query)
}

// insert、replace, save， ignore操作
// 0: insert:  仅仅执行写入操作，如果存在冲突的主键或者唯一索引，那么报错返回
// 1: replace: 如果数据存在(主键或者唯一索引)，那么删除后重新写入一条
// 2: save:    如果数据存在(主键或者唯一索引)，那么更新，否则写入一条新数据
// 3: ignore:  如果数据存在(主键或者唯一索引)，那么什么也不做
func (tx *Tx) insert(table string, data Map, option uint8) (sql.Result, error) {
    var keys   []string
    var values []string
    var params []interface{}
    for k, v := range data {
        keys   = append(keys,   tx.db.charl + k + tx.db.charr)
        values = append(values, "?")
        params = append(params, v)
    }
    operation := tx.db.getInsertOperationByOption(option)
    updatestr := ""
    if option == OPTION_SAVE {
        var updates []string
        for k, _ := range data {
            updates = append(updates, fmt.Sprintf("%s%s%s=VALUES(%s)", tx.db.charl, k, tx.db.charr, k))
        }
        updatestr = fmt.Sprintf(" ON DUPLICATE KEY UPDATE %s", strings.Join(updates, ","))
    }
    return tx.Exec(
        fmt.Sprintf("%s INTO %s%s%s(%s) VALUES(%s) %s",
            operation, tx.db.charl, table, tx.db.charr, strings.Join(keys, ","), strings.Join(values, ","), updatestr), params...
    )
}

// CURD操作:单条数据写入, 仅仅执行写入操作，如果存在冲突的主键或者唯一索引，那么报错返回
func (tx *Tx) Insert(table string, data Map) (sql.Result, error) {
    return tx.insert(table, data, OPTION_INSERT)
}

// CURD操作:单条数据写入, 如果数据存在(主键或者唯一索引)，那么删除后重新写入一条
func (tx *Tx) Replace(table string, data Map) (sql.Result, error) {
    return tx.insert(table, data, OPTION_REPLACE)
}

// CURD操作:单条数据写入, 如果数据存在(主键或者唯一索引)，那么更新，否则写入一条新数据
func (tx *Tx) Save(table string, data Map) (sql.Result, error) {
    return tx.insert(table, data, OPTION_SAVE)
}

// 批量写入数据
func (tx *Tx) batchInsert(table string, list List, batch int, option uint8) (sql.Result, error) {
    var keys    []string
    var values  []string
    var bvalues []string
    var params  []interface{}
    var result  sql.Result
    var size = len(list)
    // 判断长度
    if size < 1 {
        return result, errors.New("empty data list")
    }
    // 首先获取字段名称及记录长度
    for k, _ := range list[0] {
        keys   = append(keys,   k)
        values = append(values, "?")
    }
    var kstr = tx.db.charl + strings.Join(keys, tx.db.charl + "," + tx.db.charr) + tx.db.charr
    // 操作判断
    operation := tx.db.getInsertOperationByOption(option)
    updatestr := ""
    if option == OPTION_SAVE {
        var updates []string
        for _, k := range keys {
            updates = append(updates, fmt.Sprintf("%s=VALUES(%s)", tx.db.charl, k, tx.db.charr, k))
        }
        updatestr = fmt.Sprintf(" ON DUPLICATE KEY UPDATE %s", strings.Join(updates, ","))
    }
    // 构造批量写入数据格式(注意map的遍历是无序的)
    for i := 0; i < size; i++ {
        for _, k := range keys {
            params = append(params, list[i][k])
        }
        bvalues = append(bvalues, "(" + strings.Join(values, ",") + ")")
        if len(bvalues) == batch {
            r, err := tx.Exec(fmt.Sprintf("%s INTO %s%s%s(%s) VALUES%s %s", operation, tx.db.charl, table, tx.db.charr, kstr, strings.Join(bvalues, ","), updatestr), params...)
            if err != nil {
                return result, err
            }
            result  = r
            bvalues = bvalues[:0]
        }
    }
    // 处理最后不构成指定批量的数据
    if len(bvalues) > 0 {
        r, err := tx.Exec(fmt.Sprintf("%s INTO %s%s%s(%s) VALUES%s %s", operation, tx.db.charl, table, tx.db.charr, kstr, strings.Join(bvalues, ","), updatestr), params...)
        if err != nil {
            return result, err
        }
        result = r
    }
    return result, nil
}

// CURD操作:批量数据指定批次量写入
func (tx *Tx) BatchInsert(table string, list List, batch int) (sql.Result, error) {
    return tx.batchInsert(table, list, batch, OPTION_INSERT)
}

// CURD操作:批量数据指定批次量写入, 如果数据存在(主键或者唯一索引)，那么删除后重新写入一条
func (tx *Tx) BatchReplace(table string, list List, batch int) (sql.Result, error) {
    return tx.batchInsert(table, list, batch, OPTION_REPLACE)
}

// CURD操作:批量数据指定批次量写入, 如果数据存在(主键或者唯一索引)，那么更新，否则写入一条新数据
func (tx *Tx) BatchSave(table string, list List, batch int) (sql.Result, error) {
    return tx.batchInsert(table, list, batch, OPTION_SAVE)
}

// CURD操作:数据更新，统一采用sql预处理
// data参数支持字符串或者关联数组类型，内部会自行做判断处理
func (tx *Tx) Update(table string, data interface{}, condition interface{}, args ...interface{}) (sql.Result, error) {
    var params  []interface{}
    var updates string
    refValue := reflect.ValueOf(data)
    if refValue.Kind() == reflect.Map {
        var fields []string
        keys := refValue.MapKeys()
        for _, k := range keys {
            fields = append(fields, fmt.Sprintf("%s%s%s=?", tx.db.charl, k, tx.db.charr))
            params = append(params, gconv.String(refValue.MapIndex(k).Interface()))
            updates = strings.Join(fields,   ",")
        }
    } else {
        updates = gconv.String(data)
    }
    for _, v := range args {
        params = append(params, gconv.String(v))
    }
    return tx.Exec(fmt.Sprintf("UPDATE %s%s%s SET %s WHERE %s", tx.db.charl, table, tx.db.charr, updates, tx.db.formatCondition(condition)), params...)
}

// CURD操作:删除数据
func (tx *Tx) Delete(table string, condition interface{}, args ...interface{}) (sql.Result, error) {
    return tx.Exec(fmt.Sprintf("DELETE FROM %s%s%s WHERE %s", tx.db.charl, table, tx.db.charr, tx.db.formatCondition(condition)), args...)
}

=======
func (tx *TX) GetCount(query string, args ...interface{}) (int, error) {
    if !gregex.IsMatchString(`(?i)SELECT\s+COUNT\(.+\)\s+FROM`, query) {
        query, _ = gregex.ReplaceString(`(?i)(SELECT)\s+(.+)\s+(FROM)`, `$1 COUNT($2) $3`, query)
    }
    value, err := tx.GetValue(query, args ...)
    if err != nil {
        return 0, err
    }
    return value.Int(), nil
}

// CURD操作:单条数据写入, 仅仅执行写入操作，如果存在冲突的主键或者唯一索引，那么报错返回
func (tx *TX) Insert(table string, data interface{}, batch...int) (sql.Result, error) {
    return tx.db.doInsert(tx.tx, table, data, OPTION_INSERT, batch...)
}

// CURD操作:单条数据写入, 如果数据存在(主键或者唯一索引)，那么删除后重新写入一条
func (tx *TX) Replace(table string, data interface{}, batch...int) (sql.Result, error) {
    return tx.db.doInsert(tx.tx, table, data, OPTION_REPLACE, batch...)
}

// CURD操作:单条数据写入, 如果数据存在(主键或者唯一索引)，那么更新，否则写入一条新数据
func (tx *TX) Save(table string, data interface{}, batch...int) (sql.Result, error) {
    return tx.db.doInsert(tx.tx, table, data, OPTION_SAVE, batch...)
}

// CURD操作:批量数据指定批次量写入
func (tx *TX) BatchInsert(table string, list interface{}, batch...int) (sql.Result, error) {
    return tx.db.doBatchInsert(tx.tx, table, list, OPTION_INSERT, batch...)
}

// CURD操作:批量数据指定批次量写入, 如果数据存在(主键或者唯一索引)，那么删除后重新写入一条
func (tx *TX) BatchReplace(table string, list interface{}, batch...int) (sql.Result, error) {
    return tx.db.doBatchInsert(tx.tx, table, list, OPTION_REPLACE, batch...)
}

// CURD操作:批量数据指定批次量写入, 如果数据存在(主键或者唯一索引)，那么更新，否则写入一条新数据
func (tx *TX) BatchSave(table string, list interface{}, batch...int) (sql.Result, error) {
    return tx.db.doBatchInsert(tx.tx, table, list, OPTION_SAVE, batch...)
}

// CURD操作:数据更新，统一采用sql预处理,
// data参数支持字符串或者关联数组类型，内部会自行做判断处理.
func (tx *TX) Update(table string, data interface{}, condition interface{}, args ...interface{}) (sql.Result, error) {
    newWhere, newArgs := formatCondition(condition, args)
    return tx.doUpdate(table, data, newWhere, newArgs ...)
}

// 与Update方法的区别是不处理条件参数
func (tx *TX) doUpdate(table string, data interface{}, condition string, args ...interface{}) (sql.Result, error) {
    return tx.db.doUpdate(tx.tx, table, data, condition, args ...)
}

// CURD操作:删除数据
func (tx *TX) Delete(table string, condition interface{}, args ...interface{}) (sql.Result, error) {
    newWhere, newArgs := formatCondition(condition, args)
    return tx.doDelete(table, newWhere, newArgs ...)
}

// 与Delete方法的区别是不处理条件参数
func (tx *TX) doDelete(table string, condition string, args ...interface{}) (sql.Result, error) {
    return tx.db.doDelete(tx.tx, table, condition, args ...)
}



>>>>>>> upstream/master
