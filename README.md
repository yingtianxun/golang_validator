# golang_validator
可以自由定制的golang参数验证器

关于该组件：
本组件是可高度定制的验证，例子如下：

	type testValidator struct {
		ItemNum int `valid:"PosNO" name:"items" ` // 验证是否是正数
		IsCount int `valid:"PosNO"`
	}
 valid表示要进行验证，后面跟要进行验证的函数名：如要进行是否为整数和字符串的验证
  则可以这样写`valid:"PosNO，Str"
 name则表示当前字段的名字
 如要进行扩展，验证函数的格式如下：
 
	 func (this *Validator) CheckPosNO(data int, tag reflect.StructTag) error {
		 if data <= 0 {
			 return errors.New(tag.Get("name") + ":不能为负数!")
		 }
		 return nil
	 }
 
 Check表示函数的前缀
 而在调用的时候在标签上写：PosNO 即可
 
 关于该组件，目前golang的很多web框架都带有，但是为什么作者还要自己写呢，理由如下：
 虽然很多框架有提供现成的，但是支持定制，使用本人根据项目的需求就写了一个。
 在使用方面上更加易用
