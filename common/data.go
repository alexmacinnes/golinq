package common

type KeyValuePair[T_Key comparable, T_Value any] struct {
	Key   T_Key
	Value T_Value
}
