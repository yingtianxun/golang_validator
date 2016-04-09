package validator

import (
	"errors"
	"reflect"
	"strconv"
	"strings"
)

// 使用的例子：结构体定义如下：
//	type testValidator struct {
//		ItemNum int `valid:"PosNO" name:"items" ` // 验证是否是正数
//		IsCount int `valid:"PosNO"`
//	}
// valid表示要进行验证，后面跟要进行验证的函数名：如要进行是否为整数和字符串的验证
//  则可以这样写`valid:"PosNO，Str"
// name则表示当前字段的名字
// 如要进行扩展，验证函数的格式如下：
// func (this *Validator) CheckPosNO(data int, tag reflect.StructTag) error {
//	 if data <= 0 {
//		 return errors.New(tag.Get("name") + ":不能为负数!")
//	 }
//	 return nil
// }
// Check表示函数的前缀
// 而在调用的时候在标签上写：PosNO 即可

type Validator struct {
}

var Validate = newValidator()

func newValidator() *Validator {
	return &Validator{}
}
func (this *Validator) ValidateData(data interface{}) error {

	t := reflect.TypeOf(data)
	v := reflect.ValueOf(data)
	// 普通类型不解析，直接返回
	if v.Kind() != reflect.Slice && v.Kind() != reflect.Ptr && v.Kind() != reflect.Struct {
		return nil
	}
	err := this.parseParam(t, v, reflect.StructTag(""))
	return err
}
func (this *Validator) parseParam(t reflect.Type, v reflect.Value,
	tag reflect.StructTag) error {
	// 判断是否是数组
	if v.Kind() == reflect.Slice {
		for i, n := 0, v.Len(); i < n; i++ {
			t1 := reflect.TypeOf(v.Index(i).Interface())
			err := this.parseParam(t1, v.Index(i), tag)
			if err != nil {
				return err
			}
		}
		// 是否是指针
	} else if v.Kind() == reflect.Ptr {
		t = t.Elem()
		v = v.Elem()
		err := this.parseParam(t, v, tag)
		if err != nil {
			return err
		}
		// 结构体的
	} else if v.Kind() == reflect.Struct {
		for i := 0; i < t.NumField(); i++ {
			f := t.Field(i)
			t.FieldByName(f.Name)

			err := this.parseParam(f.Type, v.FieldByName(f.Name), f.Tag)
			if err != nil {
				return err
			}
		}
		// 最后普通类型解析
	} else {
		invoker := reflect.ValueOf(this)
		methodsName := strings.Split(tag.Get("valid"), ",") // 获取要进行验证的函数名
		for i, n := 0, len(methodsName); i < n; i++ {
			methodName := methodsName[i]
			if len(methodName) != 0 {
				method := invoker.MethodByName("Check" + methodName) // 调用需要验证的函数

				inVal := []reflect.Value{v, reflect.ValueOf(tag)}

				outVal := method.Call(inVal)

				outValLen := len(outVal)

				if outValLen != 1 {
					v.Set(outVal[0])
				}
				if !outVal[outValLen-1].IsNil() {
					return outVal[outValLen-1].Elem().Interface().(error) // 返回错误信息
				}
			}
		}
	}
	return nil
}

// 检查字符串
func (this *Validator) CheckStr(data string, tag reflect.StructTag) (string, error) {

	data = strings.Replace(data, "'", "", -1)
	data = strings.Replace(data, " ", "", -1)
	data = strings.Replace(data, " ", "", -1)
	data = strings.Replace(data, "\\", "", -1)
	data = strings.Replace(data, "\"", "", -1)

	var minLen, maxLen int = 0, 0

	var err error = nil

	lenStr := strings.Split(tag.Get("len"), ",")

	if lenStr[0] != "" {
		minLen, err = strconv.Atoi(lenStr[0])
		if err != nil {
			return "", errors.New(tag.Get("name") + ":验证字符串的最小长度参数输入有误!")
		}
	}
	if len(lenStr) == 2 && lenStr[1] != "" {
		maxLen, err = strconv.Atoi(lenStr[1])
		if err != nil {
			return "", errors.New(tag.Get("name") + ":验证字符串的最大长度参数输入有误!")
		}
	}

	if minLen > maxLen {
		return "", errors.New(tag.Get("name") + ":最小长度和最大长度冲突!")
	}
	dataLen := len(data)
	if dataLen < minLen {
		return data, errors.New(tag.Get("name") + ":字符串长度过短!")
	}

	if maxLen != 0 && dataLen > maxLen {
		return data, errors.New(tag.Get("name") + ":字符串过长!")
	}

	return data, nil
}

// 检正负数
func (this *Validator) CheckPosNO(data int, tag reflect.StructTag) error {
	if data <= 0 {
		return errors.New(tag.Get("name") + ":不能为负数!")
	}
	return nil
}
func (this *Validator) CheckCardType(data interface{}, tag reflect.StructTag) error {

	if reflect.TypeOf(data).Name() != "int" { // 判断是不是字符串

		return errors.New(tag.Get("name") + ":类型只能是int")
	}
	cardType := data.(int)

	if cardType < 0 || cardType > 8 {
		return errors.New(tag.Get("name") + ":卡类类型出错")
	}
	return nil
}
func (this *Validator) CheckCardState(data interface{}, tag reflect.StructTag) error {

	if reflect.TypeOf(data).Name() != "int" { // 判断是不是字符串

		return errors.New(tag.Get("name") + ":类型只能是int")
	}
	cardState := data.(int)

	if cardState < 0 || cardState > 3 {
		return errors.New(tag.Get("name") + ":卡状态出错")
	}
	return nil
}
func (this *Validator) CheckAddOrSub(data interface{}, tag reflect.StructTag) error {

	if reflect.TypeOf(data).Name() != "int" {
		return errors.New(tag.Get("name") + ":添加/减去只能是int类型")
	}
	addOrSub := data.(int)
	if addOrSub != -1 && addOrSub != 1 {
		return errors.New(tag.Get("name") + ":的值只能是1/-1")
	}
	return nil
}
func (this *Validator) CheckIntVal(data interface{}, tag reflect.StructTag) error {
	valRangeStr := tag.Get("range")
	if len(valRangeStr) == 0 {
		return nil
	}
	if reflect.TypeOf(data).Name() != "int" {
		return errors.New(tag.Get("name") + ":只能是int类型")
	}
	valRange := strings.Split(valRangeStr, ",")

	var minVal, maxVal int = 0, 0

	var err error = nil

	if valRange[0] != "" {
		minVal, err = strconv.Atoi(valRange[0])
		if err != nil {
			return errors.New(tag.Get("name") + ":解析最小值出错!")
		}
	}
	if len(valRange) == 2 && valRange[1] != "" {
		maxVal, err = strconv.Atoi(valRange[1])
		if err != nil {
			return errors.New(tag.Get("name") + ":解析最大值出错!")
		}
	}
	if minVal > maxVal {
		return errors.New(tag.Get("name") + ":最小值和最大值冲突!")
	}
	IntVal := data.(int)

	if IntVal < minVal || IntVal > maxVal {
		return errors.New(tag.Get("name") + ":值超出范围!")
	}
	return nil
}
func (this *Validator) CheckUserType(data interface{}, tag reflect.StructTag) error {
	if reflect.TypeOf(data).Name() != "int" {
		return errors.New(tag.Get("name") + ":类型只能是int")
	}
	userType := data.(int)

	if userType < 0 || userType > 3 {
		return errors.New(tag.Get("name") + ":用户类型出错")
	}
	return nil
}
func (this *Validator) CheckUserState(data interface{}, tag reflect.StructTag) error {
	if reflect.TypeOf(data).Name() != "int" {
		return errors.New(tag.Get("name") + ":类型只能是int")
	}
	userType := data.(int)

	if userType < 0 || userType > 2 {
		return errors.New(tag.Get("name") + ":用户类型出错")
	}
	return nil
}
