package dynamo

import (
	"reflect"
	"testing"
)

type ExtractionUser struct {
	Row
	Email string
}

func (u *ExtractionUser) Keys(index int) (string, string, error) {
	return "userPk", "userSk", nil
}

func (u *ExtractionUser) Type() string {
	return "User"
}

type House struct {
	Row
	Address string
}

func (h *House) Keys(index int) (string, string, error) {
	return "housePk", "houseSk", nil
}

func (h *House) Type() string {
	return "House"
}

type Title struct {
	DiLink[*ExtractionUser, *House]
	PointerUser *ExtractionUser
}

func (t *Title) Type() string {
	return "Title"
}

func BenchmarkAll(b *testing.B) {
	b.Run("Reflect", BenchmarkReflect)
	b.Run("Direct", BenchmarkDirect)
}

func BenchmarkReflect(b *testing.B) {
	title := Title{}

	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		_ = reflect.ValueOf(title.Entity0).IsNil()
	}
}

func BenchmarkDirect(b *testing.B) {
	title := Title{}

	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		_ = title.PointerUser == nil
	}
}

func BenchmarkInterfaceNilCheck(b *testing.B) {
	var title = Title{}

	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		_ = title.Entity0 == nil
	}
}
