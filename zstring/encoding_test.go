package zstring

import (
	"testing"

	"github.com/sohaha/zlsgo"
)

func TestBase64(t *testing.T) {
	tt := zlsgo.NewTest(t)
	str := "hi,是我"
	strbyte := []byte(str)
	s := Base64Encode(strbyte)
	deByte, err := Base64Decode(s)
	tt.EqualNil(err)
	tt.Equal(strbyte, deByte)

	s2 := Base64EncodeString(str)
	tt.Equal(Bytes2String(s), s2)

	de, err := Base64DecodeString(s2)
	tt.EqualNil(err)
	tt.Equal(str, de)

	de, _ = Base64DecodeString(string(s))
	tt.Equal(str, de)

}

type testSt struct {
	Name string
}

func TestSerialize(t *testing.T) {
	tt := zlsgo.NewTest(t)
	test := &testSt{"hi"}

	s, err := Serialize(test)
	tt.EqualNil(err)

	v, err := UnSerialize(s, &testSt{})
	tt.EqualNil(err)

	test2, ok := v.(*testSt)
	tt.EqualTrue(ok)
	tt.Equal(test.Name, test2.Name)
}
