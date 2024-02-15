package dynamo

import (
	"context"

	"github.com/aws/aws-sdk-go-v2/feature/dynamodb/attributevalue"
	awstypes "github.com/aws/aws-sdk-go-v2/service/dynamodb/types"
	sqstypes "github.com/aws/aws-sdk-go-v2/service/sqs/types"
	"github.com/entegral/gobox/types"
)

type RowGen[T types.Linkable] struct {
	Object T
	row    Row
}

func (r RowGen[T]) Type() string {
	return r.Object.Type()
}

func (r RowGen[T]) TableName(ctx context.Context) string {
	return r.Object.TableName(ctx)
}

func (r RowGen[T]) Keys(gsi int) (string, string, error) {
	return r.Object.Keys(gsi)
}

func (r RowGen[T]) MaxShard() int {
	return r.Object.MaxShard()
}

func NewRow[T types.Linkable](object T) *RowGen[T] {
	return &RowGen[T]{Object: object}
}

func (r RowGen[T]) Put(ctx context.Context) (oldItem T, err error) {
	err = r.row.Put(ctx, r.Object)
	if err != nil {
		return oldItem, err
	}
	err = attributevalue.UnmarshalMap(r.row.OldPutValues(), &oldItem)
	return oldItem, err
}

func (r RowGen[T]) Get(ctx context.Context) (bool, error) {
	return r.row.Get(ctx, r.Object)
}

func (r RowGen[T]) Delete(ctx context.Context) (oldItem T, err error) {
	err = r.row.Delete(ctx, r.Object)
	if err != nil {
		return oldItem, err
	}
	err = attributevalue.UnmarshalMap(r.row.OldDeleteValues(), &oldItem)
	return oldItem, err
}

func (r RowGen[T]) WasPutSuccessful() bool {
	return r.row.WasPutSuccessful()
}

func (r RowGen[T]) WasGetSuccessful() bool {
	return r.row.WasGetSuccessful()
}

func (r RowGen[T]) OldPutValues() map[string]awstypes.AttributeValue {
	return r.row.OldPutValues()
}

func (r RowGen[T]) OldDeleteValues() map[string]awstypes.AttributeValue {
	return r.row.OldDeleteValues()
}

func (r RowGen[T]) LoadFromMessage(ctx context.Context, message sqstypes.Message) (bool, error) {
	return r.row.LoadFromMessage(ctx, message, r.Object)
}

func (r *RowGen[T]) WithTableName(tableName string) *RowGen[T] {
	r.row.SetTableName(tableName)
	return r
}
