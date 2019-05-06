package middleware

import (
	"context"
	"time"

	"github.com/go-kit/kit/log"
	"github.com/mitchell/selfpass/credentials/types"
)

func NewServiceLogger(l log.Logger, next types.Service) ServiceLogger {
	return ServiceLogger{
		l:    l,
		next: next,
	}
}

type ServiceLogger struct {
	l    log.Logger
	next types.Service
}

func (svc ServiceLogger) GetAllMetadata(ctx context.Context, sourceHost string) (output <-chan types.Metadata, errch chan error) {
	defer func(begin time.Time) {
		_ = svc.l.Log(
			"service", "Credentials",
			"method", "GetAllMetadata",
			"input", sourceHost,
			"output", "channel",
			"err", "channel",
			"took", time.Since(begin),
		)
	}(time.Now())

	return svc.next.GetAllMetadata(ctx, sourceHost)
}

func (svc ServiceLogger) Get(ctx context.Context, id string) (output types.Credential, err error) {
	defer func(begin time.Time) {
		_ = svc.l.Log(
			"service", "Credentials",
			"method", "Get",
			"input", id,
			"output", output,
			"err", err,
			"took", time.Since(begin),
		)
	}(time.Now())

	output, err = svc.next.Get(ctx, id)
	return output, err
}

func (svc ServiceLogger) Create(ctx context.Context, ci types.CredentialInput) (output types.Credential, err error) {
	defer func(begin time.Time) {
		_ = svc.l.Log(
			"service", "Credentials",
			"method", "Create",
			"input", ci,
			"output", output,
			"err", err,
			"took", time.Since(begin),
		)
	}(time.Now())

	output, err = svc.next.Create(ctx, ci)
	return output, err
}

func (svc ServiceLogger) Update(ctx context.Context, id string, ci types.CredentialInput) (output types.Credential, err error) {
	defer func(begin time.Time) {
		_ = svc.l.Log(
			"service", "Credentials",
			"method", "Update",
			"input", []interface{}{id, ci},
			"output", output,
			"err", err,
			"took", time.Since(begin),
		)
	}(time.Now())

	output, err = svc.next.Update(ctx, id, ci)
	return output, err
}

func (svc ServiceLogger) Delete(ctx context.Context, id string) (err error) {
	defer func(begin time.Time) {
		_ = svc.l.Log(
			"service", "Credentials",
			"method", "Delete",
			"input", id,
			"err", err,
			"took", time.Since(begin),
		)
	}(time.Now())

	err = svc.next.Delete(ctx, id)
	return err
}

func (svc ServiceLogger) DumpDB(ctx context.Context) (output []byte, err error) {
	defer func(begin time.Time) {
		_ = svc.l.Log(
			"service", "Credentials",
			"method", "Dump",
			"output", output,
			"err", err,
			"took", time.Since(begin),
		)
	}(time.Now())

	output, err = svc.next.DumpDB(ctx)
	return output, err
}
