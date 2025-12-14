package v1

import (
	"context"
	"fmt"
	"net/http"
	"net/url"

	corev3 "github.com/envoyproxy/go-control-plane/envoy/config/core/v3"
	authv3 "github.com/envoyproxy/go-control-plane/envoy/service/auth/v3"
	typev3 "github.com/envoyproxy/go-control-plane/envoy/type/v3"
	statusv3 "google.golang.org/genproto/googleapis/rpc/status"
	"google.golang.org/grpc/codes"

	authv1 "github.com/nkolesnikov999/micro2-OK/shared/pkg/proto/auth/v1"
)

const (
	SessionCookieName = "X-Session-Uuid"

	HeaderUserUUID    = "X-User-Uuid"
	HeaderUserLogin   = "X-User-Login"
	HeaderContentType = "content-type"
	HeaderAuthStatus  = "X-Auth-Status"

	HeaderCookie        = "cookie"
	HeaderAuthorization = "authorization"

	ContentTypeJSON = "application/json"

	AuthStatusDenied = "denied"
)

func (a *api) Check(ctx context.Context, req *authv3.CheckRequest) (*authv3.CheckResponse, error) {
	if req == nil {
		return a.denyRequest("request cannot be nil", typev3.StatusCode_Unauthorized), nil
	}

	sessionUUID, err := a.extractSessionUUID(req)
	if err != nil {
		return a.denyRequest("missing or invalid session", typev3.StatusCode_Unauthorized), nil
	}

	whoamiResp, err := a.Whoami(ctx, &authv1.WhoamiRequest{
		SessionUuid: sessionUUID,
	})
	if err != nil {
		return a.denyRequest("invalid session", typev3.StatusCode_Unauthorized), nil
	}

	return a.allowRequest(whoamiResp), nil
}

func (a *api) extractSessionUUID(req *authv3.CheckRequest) (string, error) {
	if req.Attributes == nil || req.Attributes.Request == nil {
		return "", fmt.Errorf("no HTTP request found")
	}

	headers := req.Attributes.Request.Http.Headers

	if cookieHeader, ok := headers[HeaderCookie]; ok && cookieHeader != "" {
		sessionUUID := a.extractSessionFromCookies(cookieHeader)
		if sessionUUID != "" {
			return sessionUUID, nil
		}
	}

	return "", fmt.Errorf("session uuid not found in cookies")
}

func (a *api) extractSessionFromCookies(cookieHeader string) string {
	req := &http.Request{Header: make(http.Header)}
	req.Header.Add(HeaderCookie, cookieHeader)

	if cookie, err := req.Cookie(SessionCookieName); err == nil {
		var sessionUUID string
		sessionUUID, err = url.QueryUnescape(cookie.Value)
		if err != nil {
			return cookie.Value
		}

		return sessionUUID
	}

	return ""
}

func (a *api) allowRequest(whoamiResp *authv1.WhoamiResponse) *authv3.CheckResponse {
	headers := []*corev3.HeaderValueOption{
		{
			Header: &corev3.HeaderValue{
				Key:   HeaderUserUUID,
				Value: whoamiResp.User.Uuid,
			},
		},
		{
			Header: &corev3.HeaderValue{
				Key:   HeaderUserLogin,
				Value: whoamiResp.User.Info.Login,
			},
		},
	}

	return &authv3.CheckResponse{
		Status: &statusv3.Status{Code: 0},
		HttpResponse: &authv3.CheckResponse_OkResponse{
			OkResponse: &authv3.OkHttpResponse{
				Headers:         headers,
				HeadersToRemove: []string{HeaderCookie, HeaderAuthorization},
			},
		},
	}
}

func (a *api) denyRequest(message string, statusCode typev3.StatusCode) *authv3.CheckResponse {
	return &authv3.CheckResponse{
		Status: &statusv3.Status{Code: int32(codes.Unauthenticated)},
		HttpResponse: &authv3.CheckResponse_DeniedResponse{
			DeniedResponse: &authv3.DeniedHttpResponse{
				Status: &typev3.HttpStatus{
					Code: statusCode,
				},
				Body: fmt.Sprintf(`{"error": "%s"}`, message),
				Headers: []*corev3.HeaderValueOption{
					{
						Header: &corev3.HeaderValue{
							Key:   HeaderContentType,
							Value: ContentTypeJSON,
						},
					},
					{
						Header: &corev3.HeaderValue{
							Key:   HeaderAuthStatus,
							Value: AuthStatusDenied,
						},
					},
				},
			},
		},
	}
}
