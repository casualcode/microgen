// This file was automatically generated by "microgen 0.7.0b" utility.
// Please, do not edit.
package middleware

import (
	context "context"
	generated "github.com/devimteam/microgen/example/generated"
	entity "github.com/devimteam/microgen/example/svc/entity"
	log "github.com/go-kit/kit/log"
)

// ServiceErrorLogging writes to logger any error, if it is not nil.
func ServiceErrorLogging(logger log.Logger) Middleware {
	return func(next generated.StringService) generated.StringService {
		return &serviceErrorLogging{
			logger: logger,
			next:   next,
		}
	}
}

type serviceErrorLogging struct {
	logger log.Logger
	next   generated.StringService
}

func (L *serviceErrorLogging) Uppercase(ctx context.Context, stringsMap map[string]string) (ans string, err error) {
	defer func() {
		if err != nil {
			L.logger.Log("method", "Uppercase", "message", err)
		}
	}()
	return L.next.Uppercase(ctx, stringsMap)
}

func (L *serviceErrorLogging) Count(ctx context.Context, text string, symbol string) (count int, positions []int, err error) {
	defer func() {
		if err != nil {
			L.logger.Log("method", "Count", "message", err)
		}
	}()
	return L.next.Count(ctx, text, symbol)
}

func (L *serviceErrorLogging) TestCase(ctx context.Context, comments []*entity.Comment) (tree map[string]int, err error) {
	defer func() {
		if err != nil {
			L.logger.Log("method", "TestCase", "message", err)
		}
	}()
	return L.next.TestCase(ctx, comments)
}
