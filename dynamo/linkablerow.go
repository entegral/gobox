package dynamo

import (
	"context"

	"github.com/aws/aws-sdk-go-v2/feature/dynamodb/attributevalue"
	awstypes "github.com/aws/aws-sdk-go-v2/service/dynamodb/types"
	sqstypes "github.com/aws/aws-sdk-go-v2/service/sqs/types"
	"github.com/entegral/gobox/types"
)

type LinkableRow[T types.Linkable] struct {
	Object T
	row    Row
}

func (r LinkableRow[T]) Type() string {
	return r.Object.Type()
}

func (r LinkableRow[T]) TableName(ctx context.Context) string {
	return r.Object.TableName(ctx)
}

func (r LinkableRow[T]) Keys(gsi int) (string, string, error) {
	return r.Object.Keys(gsi)
}

func (r LinkableRow[T]) MaxShard() int {
	return r.Object.MaxShard()
}

func NewLinkableRow[T types.Linkable](object T) *LinkableRow[T] {
	return &LinkableRow[T]{Object: object}
}

func (r LinkableRow[T]) Put(ctx context.Context) (oldItem T, err error) {
	err = r.row.Put(ctx, r.Object)
	if err != nil {
		return oldItem, err
	}
	err = attributevalue.UnmarshalMap(r.row.OldPutValues(), &oldItem)
	return oldItem, err
}

func (r LinkableRow[T]) Get(ctx context.Context) (bool, error) {
	return r.row.Get(ctx, r.Object)
}

func (r LinkableRow[T]) Delete(ctx context.Context) (oldItem T, err error) {
	err = r.row.Delete(ctx, r.Object)
	if err != nil {
		return oldItem, err
	}
	err = attributevalue.UnmarshalMap(r.row.OldDeleteValues(), &oldItem)
	return oldItem, err
}

func (r LinkableRow[T]) WasPutSuccessful() bool {
	return r.row.WasPutSuccessful()
}

func (r LinkableRow[T]) WasGetSuccessful() bool {
	return r.row.WasGetSuccessful()
}

func (r LinkableRow[T]) OldPutValues() map[string]awstypes.AttributeValue {
	return r.row.OldPutValues()
}

func (r LinkableRow[T]) OldDeleteValues() map[string]awstypes.AttributeValue {
	return r.row.OldDeleteValues()
}

func (r LinkableRow[T]) LoadFromMessage(ctx context.Context, message sqstypes.Message) (bool, error) {
	return r.row.LoadFromMessage(ctx, message, r.Object)
}

func (r *LinkableRow[T]) WithTableName(tableName string) *LinkableRow[T] {
	r.row.SetTableName(tableName)
	return r
}
