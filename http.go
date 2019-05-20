package errors

//noinspection GoUnusedExportedFunction
func HttpStatus(errType string) int {
    switch errType {
    case BadRequestType:
        return 400
    case UnauthorizedType:
        return 401
    case ForbiddenType:
        return 403
    case NotFoundType:
        return 404
    case TimeoutType:
        return 441
    case InternalErrorType:
        return 500
    case NotImplementType:
        return 501
    }

    return 520
}
