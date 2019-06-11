<<<<<<< HEAD
// Copyright 2017 gf Author(https://gitee.com/johng/gf). All Rights Reserved.
//
// This Source Code Form is subject to the terms of the MIT License.
// If a copy of the MIT was not distributed with this file,
// You can obtain one at https://gitee.com/johng/gf.

// JSON解析/封装.
// 单元测试请参考gpaser包.
package gjson

import (
    "sync"
    "strings"
    "strconv"
    "io/ioutil"
    "encoding/json"
    "gitee.com/johng/gf/g/os/gfile"
    "gitee.com/johng/gf/g/util/gconv"
    "gitee.com/johng/gf/g/util/gutil"
    "gitee.com/johng/gf/g/encoding/gxml"
    "gitee.com/johng/gf/g/encoding/gyaml"
    "gitee.com/johng/gf/g/encoding/gtoml"
    "gitee.com/johng/gf/g/util/gstr"
)

const (
    gDEFAULT_SPLIT_CHAR = '.' // 默认层级分隔符号
)

// json解析结果存放数组
type Json struct {
    mu sync.RWMutex
    p  *interface{} // 注意这是一个指针
    c  byte         // 层级分隔符，默认为"."
    vc bool         // 是否执行分隔符冲突检测(默认为true，检测会比较影响检索效率)
}

// 将变量转换为Json对象进行处理，该变量至少应当是一个map或者array，否者转换没有意义
func New(value interface{}) *Json {
    switch value.(type) {
        case map[string]interface{}:
            return &Json{
                p  : &value,
                c  : byte(gDEFAULT_SPLIT_CHAR),
                vc : true ,
            }
        case []interface{}:
            return &Json{
                p  : &value,
                c  : byte(gDEFAULT_SPLIT_CHAR),
                vc : true ,
            }
        default:
            // 这里效率会比较低
            b, _ := Encode(value)
            v, _ := Decode(b)
            return &Json{
                p  : &v,
                c  : byte(gDEFAULT_SPLIT_CHAR),
                vc : true,
            }
    }
}

// 编码go变量为json字符串，并返回json字符串指针
func Encode (v interface{}) ([]byte, error) {
    return json.Marshal(v)
}

// 解码字符串为interface{}变量
func Decode (b []byte) (interface{}, error) {
    var v interface{}
    if err := DecodeTo(b, &v); err != nil {
        return nil, err
    } else {
        return v, nil
    }
}

// 解析json字符串为go变量，注意第二个参数为指针(任意结构的变量)
func DecodeTo (b []byte, v interface{}) error {
    return json.Unmarshal(b, v)
}

// 解析json字符串为gjson.Json对象，并返回操作对象指针
func DecodeToJson (b []byte) (*Json, error) {
    if v, err := Decode(b); err != nil {
        return nil, err
    } else {
        return New(v), nil
    }
}

// 支持多种配置文件类型转换为json格式内容并解析为gjson.Json对象
func Load (path string) (*Json, error) {
    data, err := ioutil.ReadFile(path)
    if err != nil {
        return nil, err
    }
    return LoadContent(data, gfile.Ext(path))
}

// 支持的配置文件格式：xml, json, yaml/yml, toml
func LoadContent (data []byte, t string) (*Json, error) {
    var err    error
    var result interface{}
    switch t {
        case  "xml":  fallthrough
        case ".xml":
            data, err = gxml.ToJson(data)
            if err != nil {
                return nil, err
            }
        case   "yml": fallthrough
        case  "yaml": fallthrough
        case  ".yml": fallthrough
        case ".yaml":
            data, err = gyaml.ToJson(data)
            if err != nil {
                return nil, err
            }

        case  "toml": fallthrough
        case ".toml":
            data, err = gtoml.ToJson(data)
            if err != nil {
                return nil, err
            }
    }
    if err := json.Unmarshal(data, &result); err != nil {
        return nil, err
    }
    return New(result), nil
}

// 设置自定义的层级分隔符号
func (j *Json) SetSplitChar(char byte) {
    j.mu.Lock()
    j.c = char
    j.mu.Unlock()
}

// 设置自定义的层级分隔符号
func (j *Json) SetViolenceCheck(check bool) {
    j.mu.Lock()
    j.vc = check
    j.mu.Unlock()
}

// 将指定的json内容转换为指定结构返回，查找失败或者转换失败，目标对象转换为nil
// 注意第二个参数需要给的是**变量地址**
func (j *Json) GetToVar(pattern string, v interface{}) error {
    r := j.Get(pattern)
    if r != nil {
        if t, err := Encode(r); err == nil {
            return DecodeTo(t, v)
        } else {
            return err
        }
    } else {
        v = nil
    }
    return nil
}

// 获得一个键值对关联数组/哈希表，方便操作，不需要自己做类型转换
// 注意，如果获取的值不存在，或者类型与json类型不匹配，那么将会返回nil
func (j *Json) GetMap(pattern string) map[string]interface{} {
    result := j.Get(pattern)
    if result != nil {
        if r, ok := result.(map[string]interface{}); ok {
            return r
        }
    }
    return nil
}

// 将检索值转换为Json对象指针返回
func (j *Json) GetJson(pattern string) *Json {
    result := j.Get(pattern)
    if result != nil {
        return New(result)
    }
    return nil
}

// 获得一个数组[]interface{}，方便操作，不需要自己做类型转换
// 注意，如果获取的值不存在，或者类型与json类型不匹配，那么将会返回nil
func (j *Json) GetArray(pattern string) []interface{} {
    result := j.Get(pattern)
    if result != nil {
        if r, ok := result.([]interface{}); ok {
            return r
        }
    }
    return nil
}

// 返回指定json中的string
func (j *Json) GetString(pattern string) string {
    return gconv.String(j.Get(pattern))
}

// 返回指定json中的bool(false:"", 0, false, off)
func (j *Json) GetBool(pattern string) bool {
    return gconv.Bool(j.Get(pattern))
}

func (j *Json) GetInt(pattern string) int {
    return gconv.Int(j.Get(pattern))
}

func (j *Json) GetUint(pattern string) uint {
    return gconv.Uint(j.Get(pattern))
}

func (j *Json) GetFloat32(pattern string) float32 {
    return gconv.Float32(j.Get(pattern))
}

func (j *Json) GetFloat64(pattern string) float64 {
    return gconv.Float64(j.Get(pattern))
}

// 动态设置层级变量
func (j *Json) Set(pattern string, value interface{}) error {
    return j.setValue(pattern, value, false)
}

// 动态删除层级变量
func (j *Json) Remove(pattern string) error {
    return j.setValue(pattern, nil, true)
}

// 根据pattern查找并设置数据
// 注意：
// 1、写入的value为nil且removed为true时，表示删除;
// 2、里面的层级处理比较复杂，逻辑较复杂的地方在于层级检索及节点创建，叶子赋值;
=======
// Copyright 2017 gf Author(https://github.com/gogf/gf). All Rights Reserved.
//
// This Source Code Form is subject to the terms of the MIT License.
// If a copy of the MIT was not distributed with this file,
// You can obtain one at https://github.com/gogf/gf.

// Package gjson provides convenient API for JSON/XML/YAML/TOML data handling.
package gjson

import (
	"github.com/gogf/gf/g/internal/rwmutex"
	"github.com/gogf/gf/g/text/gstr"
	"github.com/gogf/gf/g/util/gconv"
	"reflect"
	"strconv"
	"strings"
)

const (
    // Separator char for hierarchical data access.
    gDEFAULT_SPLIT_CHAR = '.'
)

// The customized JSON struct.
type Json struct {
    mu *rwmutex.RWMutex
    p  *interface{} // Pointer for hierarchical data access, it's the root of data in default.
    c  byte         // Char separator('.' in default).
    vc bool         // Violence Check(false in default), which is used to access data
                    // when the hierarchical data key contains separator char.
}

// Set <value> by <pattern>.
// Notice:
// 1. If value is nil and removed is true, means deleting this value;
// 2. It's quite complicated in hierarchical data search, node creating and data assignment;
>>>>>>> upstream/master
func (j *Json) setValue(pattern string, value interface{}, removed bool) error {
    array   := strings.Split(pattern, string(j.c))
    length  := len(array)
    value    = j.convertValue(value)
    // 初始化判断
    if *j.p == nil {
        if gstr.IsNumeric(array[0]) {
            *j.p = make([]interface{}, 0)
        } else {
            *j.p = make(map[string]interface{})
        }
    }
    var pparent *interface{} = nil // 父级元素项(设置时需要根据子级的内容确定数据类型，所以必须记录父级)
    var pointer *interface{} = j.p // 当前操作层级项
    j.mu.Lock()
    defer j.mu.Unlock()
    for i:= 0; i < length; i++ {
        switch (*pointer).(type) {
            case map[string]interface{}:
                if i == length - 1 {
                    if removed && value == nil {
                        // 删除map元素
                        delete((*pointer).(map[string]interface{}), array[i])
                    } else {
                        (*pointer).(map[string]interface{})[array[i]] = value
                    }
                } else {
                    // 当键名不存在的情况这里会进行处理
                    if v, ok := (*pointer).(map[string]interface{})[array[i]]; !ok {
                        if removed && value == nil {
                            goto done
                        }
                        // 创建新节点
                        if gstr.IsNumeric(array[i + 1]) {
                            // 创建array节点
                            n, _ := strconv.Atoi(array[i + 1])
                            var v interface{} = make([]interface{}, n + 1)
                            pparent = j.setPointerWithValue(pointer, array[i], v)
                            pointer = &v
                        } else {
                            // 创建map节点
                            var v interface{} = make(map[string]interface{})
                            pparent = j.setPointerWithValue(pointer, array[i], v)
                            pointer = &v
                        }
                    } else {
                        pparent = pointer
                        pointer = &v
                    }
                }

            case []interface{}:
                // 键名与当前指针类型不符合，需要执行**覆盖操作**
                if !gstr.IsNumeric(array[i]) {
                    if i == length - 1 {
                        *pointer = map[string]interface{}{ array[i] : value }
                    } else {
                        var v interface{} = make(map[string]interface{})
                        *pointer = v
                        pparent  = pointer
                        pointer  = &v
                    }
                    continue
                }

                valn, err := strconv.Atoi(array[i])
                if err != nil {
                    return err
                }
                // 叶子节点
                if i == length - 1 {
                    if len((*pointer).([]interface{})) > valn {
                        if removed && value == nil {
                            // 删除数据元素
                            j.setPointerWithValue(pparent, array[i - 1], append((*pointer).([]interface{})[ : valn], (*pointer).([]interface{})[valn + 1 : ]...))
                        } else {
                            (*pointer).([]interface{})[valn] = value
                        }
                    } else {
                        if removed && value == nil {
                            goto done
                        }
<<<<<<< HEAD
                        j.setPointerWithValue(pointer, array[i], value)
=======
                        if pparent == nil {
                            // 表示根节点
                            j.setPointerWithValue(pointer, array[i], value)
                        } else {
                            // 非根节点
                            s := make([]interface{}, valn + 1)
                            copy(s, (*pointer).([]interface{}))
                            s[valn] = value
                            j.setPointerWithValue(pparent, array[i - 1], s)
                        }
>>>>>>> upstream/master
                    }
                } else {
                    if gstr.IsNumeric(array[i + 1]) {
                        n, _ := strconv.Atoi(array[i + 1])
                        if len((*pointer).([]interface{})) > valn {
                            (*pointer).([]interface{})[valn] = make([]interface{}, n + 1)
                            pparent                          = pointer
                            pointer                          = &(*pointer).([]interface{})[valn]
                        } else {
                            if removed && value == nil {
                                goto done
                            }
                            var v interface{} = make([]interface{}, n + 1)
                            pparent = j.setPointerWithValue(pointer, array[i], v)
                            pointer = &v
                        }
                    } else {
                        var v interface{} = make(map[string]interface{})
                        pparent = j.setPointerWithValue(pointer, array[i], v)
                        pointer = &v
                    }
                }

            // 如果当前指针指向的变量不是引用类型的，
            // 那么修改变量必须通过父级进行修改，即 pparent
            default:
                if removed && value == nil {
                    goto done
                }
                if gstr.IsNumeric(array[i]) {
                    n, _    := strconv.Atoi(array[i])
                    s       := make([]interface{}, n + 1)
                    if i == length - 1 {
                        s[n] = value
                    }
                    if pparent != nil {
                        pparent = j.setPointerWithValue(pparent, array[i - 1], s)
                    } else {
                        *pointer = s
                        pparent  = pointer
                    }
                } else {
                    var v interface{} = make(map[string]interface{})
                    if i == length - 1 {
                        v = map[string]interface{}{
                            array[i] : value,
                        }
                    }
                    if pparent != nil {
                        pparent = j.setPointerWithValue(pparent, array[i - 1], v)
                    } else {
                        *pointer = v
                        pparent  = pointer
                    }
                    pointer = &v
                }
<<<<<<< HEAD

=======
>>>>>>> upstream/master
        }
    }
done:
    return nil
}

<<<<<<< HEAD
// 数据结构转换，map参数必须转换为map[string]interface{}，数组参数必须转换为[]interface{}
=======
// Convert <value> to map[string]interface{} or []interface{},
// which can be supported for hierarchical data access.
>>>>>>> upstream/master
func (j *Json) convertValue(value interface{}) interface{} {
    switch value.(type) {
        case map[string]interface{}:
            return value
        case []interface{}:
            return value
        default:
<<<<<<< HEAD
            // 这里效率会比较低，当然比直接用反射也不会差到哪儿去
            // 为了操作的灵活性，牺牲了一定的效率
            b, _ := Encode(value)
            v, _ := Decode(b)
            return v
    }
    return value
}

// 用于Set方法中，对指针指向的内存地址进行赋值
// 返回修改后的父级指针
=======
            rv   := reflect.ValueOf(value)
            kind := rv.Kind()
            if kind == reflect.Ptr {
                rv   = rv.Elem()
                kind = rv.Kind()
            }
            switch kind {
                case reflect.Array:  return gconv.Interfaces(value)
                case reflect.Slice:  return gconv.Interfaces(value)
                case reflect.Map:    return gconv.Map(value)
                case reflect.Struct: return gconv.Map(value)
                default:
                    // Use json decode/encode at last.
                    b, _ := Encode(value)
                    v, _ := Decode(b)
                    return v
            }
    }
}

// Set <key>:<value> to <pointer>, the <key> may be a map key or slice index.
// It returns the pointer to the new value set.
>>>>>>> upstream/master
func (j *Json) setPointerWithValue(pointer *interface{}, key string, value interface{}) *interface{} {
    switch (*pointer).(type) {
        case map[string]interface{}:
            (*pointer).(map[string]interface{})[key] = value
            return &value
        case []interface{}:
            n, _ := strconv.Atoi(key)
            if len((*pointer).([]interface{})) > n {
                (*pointer).([]interface{})[n] = value
                return &(*pointer).([]interface{})[n]
            } else {
                s := make([]interface{}, n + 1)
                copy(s, (*pointer).([]interface{}))
                s[n] = value
                *pointer = s
                return &s[n]
            }
        default:
            *pointer = value
    }
    return pointer
}

<<<<<<< HEAD
// 根据约定字符串方式访问json解析数据，参数形如： "items.name.first", "list.0"
// 返回的结果类型的interface{}，因此需要自己做类型转换
// 如果找不到对应节点的数据，返回nil
func (j *Json) Get(pattern string) interface{} {
    j.mu.RLock()
    defer j.mu.RUnlock()

    var result *interface{}
    if j.vc {
        result = j.getPointerByPattern(pattern)
    } else {
        result = j.getPointerByPatternWithoutSplitCharViolenceCheck(pattern)
    }
    if result != nil {
        return *result
    }
    return nil
}

// 根据pattern层级查找**变量指针**
// 检索方式：例如检索 a.a.a ，值为1
// 1. 检索 a.a.a.a 是否存在对应map的键名；
// 2. 检索 a.a.a   是否存在对应map的键名；
// 3. 检索 a.a     是否存在对应map的键名；
// 4. 检索 a       是否存在对应map的键名，如果检索出这是一个map，假如为变量m1；
// 5. 在m1中检索 a.a.a 否存在对应map的键名；
// 6. 在m1中检索 a.a   否存在对应map的键名；
// 7. 在m1中检索 a     否存在对应map的键名，如果检索出这是一个map，假如为变量m2；
// 8. 在m2中检索 a.a   否存在对应map的键名；
// 9. 在m2中检索 a     否存在对应map的键名，检索到有值，值为1；
// 这样检索的复杂度很高，主要是为了避免键名中存在分隔符号(默认为".")的情况，避免歧义。
func (j *Json) getPointerByPattern(pattern string) *interface{} {
=======
// Get a pointer to the value by specified <pattern>.
func (j *Json) getPointerByPattern(pattern string) *interface{} {
    if j.vc {
        return j.getPointerByPatternWithViolenceCheck(pattern)
    } else {
        return j.getPointerByPatternWithoutViolenceCheck(pattern)
    }
}

// Get a pointer to the value of specified <pattern> with violence check.
func (j *Json) getPointerByPatternWithViolenceCheck(pattern string) *interface{} {
    if !j.vc {
        return j.getPointerByPatternWithoutViolenceCheck(pattern)
    }
>>>>>>> upstream/master
    index   := len(pattern)
    start   := 0
    length  := 0
    pointer := j.p
    if index == 0 {
        return pointer
    }
    for {
        if r := j.checkPatternByPointer(pattern[start:index], pointer); r != nil {
            length += index - start
            if start > 0 {
                length += 1
            }
            start = index + 1
            index = len(pattern)
            if length == len(pattern) {
                return r
            } else {
                pointer = r
            }
        } else {
<<<<<<< HEAD
            // 查找下一个分割符号的索引位置
=======
            // Get the position for next separator char.
>>>>>>> upstream/master
            index = strings.LastIndexByte(pattern[start:index], j.c)
            if index != -1 && length > 0 {
                index += length + 1
            }
        }
        if start >= index {
            break
        }
    }
    return nil
}

<<<<<<< HEAD
// 层级检索，内部不执行分隔符冲突检查，检索效率会有所提高，但是冲突需要开发者自己根据自定义的分隔符来进行解决
func (j *Json) getPointerByPatternWithoutSplitCharViolenceCheck(pattern string) *interface{} {
=======
// Get a pointer to the value of specified <pattern>, with no violence check.
func (j *Json) getPointerByPatternWithoutViolenceCheck(pattern string) *interface{} {
    if j.vc {
        return j.getPointerByPatternWithViolenceCheck(pattern)
    }
>>>>>>> upstream/master
    pointer := j.p
    if len(pattern) == 0 {
        return pointer
    }
    array := strings.Split(pattern, string(j.c))
    for k, v := range array {
        if r := j.checkPatternByPointer(v, pointer); r != nil {
            if k == len(array) - 1 {
                return r
            } else {
                pointer = r
            }
        } else {
            break
        }
    }
    return nil
}

<<<<<<< HEAD
// 判断给定的key在当前的pointer下是否有值，并返回对应的pointer
// 注意这里返回的指针都是临时变量的内存地址
=======
// Check whether there's value by <key> in specified <pointer>.
// It returns a pointer to the value.
>>>>>>> upstream/master
func (j *Json) checkPatternByPointer(key string, pointer *interface{}) *interface{} {
    switch (*pointer).(type) {
        case map[string]interface{}:
            if v, ok := (*pointer).(map[string]interface{})[key]; ok {
                return &v
            }
        case []interface{}:
            if gstr.IsNumeric(key) {
                n, err := strconv.Atoi(key)
                if err == nil && len((*pointer).([]interface{})) > n {
                    return &(*pointer).([]interface{})[n]
                }
            }
    }
    return nil
}
<<<<<<< HEAD

// 转换为map[string]interface{}类型,如果转换失败，返回nil
func (j *Json) ToMap() map[string]interface{} {
    j.mu.RLock()
    defer j.mu.RUnlock()
    switch (*(j.p)).(type) {
        case map[string]interface{}:
            return (*(j.p)).(map[string]interface{})
        default:
            return nil
    }
}

// 转换为[]interface{}类型,如果转换失败，返回nil
func (j *Json) ToArray() []interface{} {
    j.mu.RLock()
    defer j.mu.RUnlock()
    switch (*(j.p)).(type) {
        case []interface{}:
            return (*(j.p)).([]interface{})
        default:
            return nil
    }
}

func (j *Json) ToXml(rootTag...string) ([]byte, error) {
    return gxml.Encode(j.ToMap(), rootTag...)
}

func (j *Json) ToXmlIndent(rootTag...string) ([]byte, error) {
    return gxml.EncodeWithIndent(j.ToMap(), rootTag...)
}

func (j *Json) ToJson() ([]byte, error) {
    j.mu.RLock()
    defer j.mu.RUnlock()
    return Encode(*(j.p))
}

func (j *Json) ToJsonIndent() ([]byte, error) {
    j.mu.RLock()
    defer j.mu.RUnlock()
    return json.MarshalIndent(*(j.p), "", "\t")
}

func (j *Json) ToYaml() ([]byte, error) {
    j.mu.RLock()
    defer j.mu.RUnlock()
    return gyaml.Encode(*(j.p))
}

func (j *Json) ToToml() ([]byte, error) {
    j.mu.RLock()
    defer j.mu.RUnlock()
    return gtoml.Encode(*(j.p))
}

// 转换为指定的struct对象
func (j *Json) ToStruct(o interface{}) error {
    j.mu.RLock()
    defer j.mu.RUnlock()
    return gutil.MapToStruct(j.ToMap(), o)
}
=======
>>>>>>> upstream/master
