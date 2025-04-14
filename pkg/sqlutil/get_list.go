package sqlutil

import (
	"context"
	// "fmt"
)

const sliceSize = 50

func AsArray[K any](slice []K) [][sliceSize]K {
	array := make([][sliceSize]K, 0, len(slice)/sliceSize)
	left := len(slice)
	for left > 0 {
		sliceCut := [sliceSize]K{}
		start := len(slice) - left
		end := start + sliceSize
		if end > len(slice) {
			end = len(slice)
		}
		copy(sliceCut[:], slice[start:end])
		array = append(array, sliceCut)
		left = left - sliceSize
	}
	return array
}

type GetListQuery[K any, V any] func(context.Context, [sliceSize]K) ([]V, error)
type DeleteQuery[K any] func(context.Context, [sliceSize]K) error

func GetListBy50[K any, V any](ctx context.Context, ids []K, query GetListQuery[K, V]) ([]V, error) {
	arrayOfIds := AsArray(ids)
	result := []V{}
	for _, array := range arrayOfIds {
		queryResult, err := query(ctx, array)
		if err != nil {
			return nil, err
		}
		result = append(result, queryResult[:]...)
	}
	return result, nil
}

func DeleteBy50[K any](ctx context.Context, ids []K, query DeleteQuery[K]) error {
	arrayOfIds := AsArray(ids)
	for _, array := range arrayOfIds {
		err := query(ctx, array)
		if err != nil {
			return err
		}
	}
	return nil
}

func Map[K any, V any](in []K, dto func(K) V) []V {
	out := make([]V, len(in))
	for i, v := range in {
		out[i] = dto(v)
	}
	return out
}

func Map2[K any, V any, W any](in []K, k2v func(K) V, v2w func(V) W) []W {
	out := make([]W, len(in))
	for i, k := range in {
		out[i] = v2w(k2v(k))
	}
	return out
}
