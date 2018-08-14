package grpcError

import (
	"github.com/onedaycat/errors"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

var mapGrpcCode = map[int]codes.Code{
	400: 3,
	401: 16,
	403: 7,
	404: 5,
	441: 4,
	500: 13,
	501: 12,
	503: 14,
	520: 2,
}

func ToGRPC(err error) error {
	if err == nil {
		return nil
	}

	herr, ok := err.(*errors.AppError)
	if !ok {
		return status.Error(2, err.Error())
	}

	c := mapGrpcCode[herr.Status]
	return status.Error(c, herr.Message)
}
