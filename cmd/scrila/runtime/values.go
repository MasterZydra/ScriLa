package runtime

type ValueType string

const (
	BoolValueType ValueType = "bool"
	IntValueType  ValueType = "int"
	NullValueType ValueType = "null"
	StrValueType  ValueType = "str"
)

type IRuntimeVal interface {
	GetType() ValueType
}

type RuntimeVal struct {
	valueType ValueType
}

func (self *RuntimeVal) GetType() ValueType {
	return self.valueType
}

type INullVal interface {
	IRuntimeVal
	GetValue() string
}

type NullVal struct {
	valueType ValueType
	value     string
}

func NewNullVal() *NullVal {
	return &NullVal{
		valueType: NullValueType,
		value:     "null",
	}
}

func (self *NullVal) GetType() ValueType {
	return self.valueType
}

func (self *NullVal) GetValue() string {
	return self.value
}

type IIntVal interface {
	IRuntimeVal
	GetValue() int64
}

type IntVal struct {
	valueType ValueType
	value     int64
}

func NewIntVal(value int64) *IntVal {
	return &IntVal{
		valueType: IntValueType,
		value:     value,
	}
}

func (self *IntVal) GetType() ValueType {
	return self.valueType
}

func (self *IntVal) GetValue() int64 {
	return self.value
}

type IBoolVal interface {
	IRuntimeVal
	GetValue() bool
}

type BoolVal struct {
	valueType ValueType
	value     bool
}

func NewBoolVal(value bool) *BoolVal {
	return &BoolVal{
		valueType: BoolValueType,
		value:     value,
	}
}

func (self *BoolVal) GetType() ValueType {
	return self.valueType
}

func (self *BoolVal) GetValue() bool {
	return self.value
}

// TODO BoolVal
