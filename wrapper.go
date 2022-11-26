package recovery // import "go.unistack.org/micro-wrapper-recovery/v3"

import (
	"context"
	"fmt"

	"go.unistack.org/micro/v3/server"
)

type wrapper struct {
	serverHandlerFunc    func(context.Context, server.Request, interface{}, error) error
	serverSubscriberFunc func(context.Context, server.Message, error) error
	/*
		clientCallFunc       func(context.Context, string, client.Request, interface{}, client.CallOptions, error) error
		clientClient         func(client.Client, error) error
	*/
}

func NewServerHandlerWrapper(fn func(context.Context, server.Request, interface{}, error) error) server.HandlerWrapper {
	handler := &wrapper{
		serverHandlerFunc: fn,
	}
	return handler.HandlerFunc
}

func (w *wrapper) HandlerFunc(fn server.HandlerFunc) server.HandlerFunc {
	return func(ctx context.Context, req server.Request, rsp interface{}) (err error) {
		defer func() {
			r := recover()
			switch verr := r.(type) {
			case nil:
				return
			case error:
				err = w.serverHandlerFunc(ctx, req, rsp, verr)
			default:
				err = w.serverHandlerFunc(ctx, req, rsp, fmt.Errorf("%v", r))
			}
		}()
		err = fn(ctx, req, rsp)
		return err
	}
}

func NewServerSubscriberWrapper(fn func(context.Context, server.Message, error) error) server.SubscriberWrapper {
	handler := &wrapper{
		serverSubscriberFunc: fn,
	}
	return handler.SubscriberFunc
}

func (w *wrapper) SubscriberFunc(fn server.SubscriberFunc) server.SubscriberFunc {
	return func(ctx context.Context, msg server.Message) (err error) {
		defer func() {
			r := recover()
			switch verr := r.(type) {
			case nil:
				return
			case error:
				err = w.serverSubscriberFunc(ctx, msg, verr)
			default:
				err = w.serverSubscriberFunc(ctx, msg, fmt.Errorf("%v", r))
			}
		}()
		err = fn(ctx, msg)
		return err
	}
}

/*
func NewClientWrapper() client.Wrapper {
	return func(c client.Client) client.Client {
		handler := &wrapper{
			clientClient: c,
		}
		return handler
	}
}

func NewCallWrapper() client.CallWrapper {
	return func(fn client.CallFunc) client.CallFunc {
		handler := &wrapper{
			clientCallFunc: fn,
		}
		return handler.CallFunc
	}
}

func (w *wrapper) CallFunc(ctx context.Context, addr string, req client.Request, rsp interface{}, opts client.CallOptions) (err error) {
	defer func() {
		r := recover()
		switch verr := r.(type) {
		case nil:
			return
		case error:
			err = w.clientCallFunc(ctx, addr, req, rsp, opts, verr)
		default:
			err = w.clientCallFunc(ctx, addr, req, rsp, opts, fmt.Errorf("%v", r))
		}
	}()
	err = w.CallFunc(ctx, addr, req, rsp, opts)
	return err
}

func (w *wrapper) Call(ctx context.Context, req client.Request, rsp interface{}, opts ...client.CallOption) error {
	defer func() {
		r := recover()
		switch verr := r.(type) {
		case nil:
			return
		case error:
			err = w.clientClient.Call(ctx, addr, req, rsp, opts, verr)
		default:
			err = w.clientClient.Call(ctx, addr, req, rsp, opts, fmt.Errorf("%v", r))
		}
	}()
	err = w.clientClient.Call(ctx, req, rsp, opts...)
	return err
}

func (w *wrapper) Stream(ctx context.Context, req client.Request, opts ...client.CallOption) (client.Stream, error) {

	stream, err := w.Client.Stream(ctx, req, opts...)

	return stream, err
}

func (w *wrapper) Publish(ctx context.Context, p client.Message, opts ...client.PublishOption) error {

	err := w.Client.Publish(ctx, p, opts...)


	return err
}
*/
