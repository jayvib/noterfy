package rest

import (
	"context"
	"github.com/go-kit/kit/endpoint"
	"net/http"
	"noterfy/note"
	"strconv"
)

type fetchRequest struct {
	Pagination *note.Pagination
}

type fetchResponse struct {
	Notes      []*note.Note `json:"notes"`
	TotalCount uint64       `json:"total_count"`
	TotalPage  uint64       `json:"total_page"`
}

func decodeFetchRequest(_ context.Context, r *http.Request) (response interface{}, err error) {

	page := r.URL.Query().Get("page")
	size := r.URL.Query().Get("size")
	sortBy := r.URL.Query().Get("sort_by")

	response = fetchRequest{
		Pagination: &note.Pagination{
			Size:   convertAtoU(size),
			Page:   convertAtoU(page),
			SortBy: note.GetSortBy(sortBy),
			Ascend: false,
		},
	}

	return
}

func makeFetchEndpoint(svc fetchService) endpoint.Endpoint {
	return func(ctx context.Context, req interface{}) (resp interface{}, err error) {
		request := req.(fetchRequest)
		iter, err := svc.Fetch(ctx, request.Pagination)
		if err != nil {
			return newErrorWrapper(err), nil
		}

		defer func() {
			cerr := iter.Close()
			if cerr != nil && err == nil {
				err = cerr
			}
		}()

		var notes []*note.Note
		for iter.Next() {
			n := iter.Note()
			notes = append(notes, n)
		}

		if iter.Error() != nil {
			return newErrorWrapper(err), nil
		}

		resp = fetchResponse{
			Notes:      notes,
			TotalCount: iter.TotalCount(),
			TotalPage:  iter.TotalPage(),
		}

		return
	}
}

func convertAtoU(s string) uint64 {
	v, _ := strconv.Atoi(s)
	return uint64(v)
}
