package v1

import (
	"context"
	"fmt"
	"net/http"
	"net/url"
	"strings"

	corev3 "github.com/envoyproxy/go-control-plane/envoy/config/core/v3"
	envoyauth "github.com/envoyproxy/go-control-plane/envoy/service/auth/v3"
	typev3 "github.com/envoyproxy/go-control-plane/envoy/type/v3"
	"go.uber.org/zap"
	statusv3 "google.golang.org/genproto/googleapis/rpc/status"
	"google.golang.org/grpc/codes"

	"github.com/nkolesnikov999/micro2-OK/platform/pkg/logger"
	iamauth "github.com/nkolesnikov999/micro2-OK/shared/pkg/proto/auth/v1"
)

const (
	SessionCookieName = "session_uuid"

	HeaderUserUUID    = "X-User-Uuid"
	HeaderUserLogin   = "X-User-Login"
	HeaderContentType = "content-type"
	HeaderAuthStatus  = "X-Auth-Status"
	HeaderSessionUUID = "X-Session-Uuid"

	HeaderCookie        = "cookie"
	HeaderAuthorization = "authorization"

	ContentTypeJSON = "application/json"

	AuthStatusDenied = "denied"
)

func (a *api) Check(ctx context.Context, req *envoyauth.CheckRequest) (*envoyauth.CheckResponse, error) {
	logger.Info(ctx, "Check method called",
		zap.Bool("req_is_nil", req == nil),
	)

	if req == nil {
		return a.denyRequest("request cannot be nil", typev3.StatusCode_Unauthorized), nil
	}

	if req.Attributes != nil && req.Attributes.Request != nil {
		logger.Info(ctx, "Check request details",
			zap.String("path", req.Attributes.Request.Http.Path),
			zap.String("method", req.Attributes.Request.Http.Method),
		)
	}

	sessionUUID, err := a.extractSessionUUID(ctx, req)
	if err != nil {
		logger.Info(ctx, "failed to extract session UUID",
			zap.Error(err),
		)
		return a.denyRequest("missing or invalid session", typev3.StatusCode_Unauthorized), nil
	}

	logger.Info(ctx, "calling Whoami with session UUID",
		zap.String("session_uuid", sessionUUID),
	)

	whoamiResp, err := a.Whoami(ctx, &iamauth.WhoamiRequest{
		SessionUuid: sessionUUID,
	})
	if err != nil {
		logger.Info(ctx, "Whoami failed",
			zap.String("session_uuid", sessionUUID),
			zap.Error(err),
		)
		return a.denyRequest("invalid session", typev3.StatusCode_Unauthorized), nil
	}

	logger.Info(ctx, "Whoami successful, allowing request",
		zap.String("user_uuid", whoamiResp.User.Uuid),
	)

	return a.allowRequest(whoamiResp, sessionUUID), nil
}

func (a *api) extractSessionUUID(ctx context.Context, req *envoyauth.CheckRequest) (string, error) {
	if req.Attributes == nil || req.Attributes.Request == nil {
		return "", fmt.Errorf("no HTTP request found")
	}

	headers := req.Attributes.Request.Http.Headers

	logger.Info(ctx, "extractSessionUUID: checking headers",
		zap.String("path", req.Attributes.Request.Http.Path),
	)

	cookieHeader, ok := headers[HeaderCookie]
	if !ok || cookieHeader == "" {
		logger.Info(ctx, "cookie header not found",
			zap.Bool("has_cookie", false),
		)
		return "", fmt.Errorf("cookie header not found")
	}

	logger.Info(ctx, "found cookie header",
		zap.Bool("has_cookie", true),
	)

	sessionUUID := a.extractSessionFromCookies(cookieHeader)
	if sessionUUID != "" {
		logger.Info(ctx, "extracted session UUID",
			zap.String("session_uuid", sessionUUID),
		)
		return sessionUUID, nil
	}

	logger.Info(ctx, "failed to extract session UUID from cookie",
		zap.Bool("has_cookie", cookieHeader != ""),
	)
	return "", fmt.Errorf("session uuid not found in cookies")
}

func (a *api) extractSessionFromCookies(cookieHeader string) string {
	if cookieHeader == "" {
		return ""
	}

	// Парсим cookie заголовок вручную для большей надежности
	// Формат: "session_uuid=value" или "session_uuid=value; other=value"
	cookies := parseCookies(cookieHeader)
	if sessionUUID, ok := cookies[SessionCookieName]; ok {
		// Декодируем URL-encoded значение, если нужно
		if decoded, err := url.QueryUnescape(sessionUUID); err == nil {
			return decoded
		}
		return sessionUUID
	}

	// Fallback: используем стандартный метод Go
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

// parseCookies парсит cookie заголовок в map
func parseCookies(cookieHeader string) map[string]string {
	cookies := make(map[string]string)

	// Разделяем по точке с запятой
	parts := strings.Split(cookieHeader, ";")
	for _, part := range parts {
		part = strings.TrimSpace(part)
		// Разделяем на имя и значение
		kv := strings.SplitN(part, "=", 2)
		if len(kv) == 2 {
			name := strings.TrimSpace(kv[0])
			value := strings.TrimSpace(kv[1])
			cookies[name] = value
		}
	}

	return cookies
}

func (a *api) allowRequest(whoamiResp *iamauth.WhoamiResponse, sessionUUID string) *envoyauth.CheckResponse {
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

	if sessionUUID != "" {
		headers = append(headers, &corev3.HeaderValueOption{
			Header: &corev3.HeaderValue{
				Key:   HeaderSessionUUID,
				Value: sessionUUID,
			},
		})
	}

	return &envoyauth.CheckResponse{
		Status: &statusv3.Status{Code: 0},
		HttpResponse: &envoyauth.CheckResponse_OkResponse{
			OkResponse: &envoyauth.OkHttpResponse{
				Headers:         headers,
				HeadersToRemove: []string{HeaderCookie, HeaderAuthorization},
			},
		},
	}
}

func (a *api) denyRequest(message string, statusCode typev3.StatusCode) *envoyauth.CheckResponse {
	return &envoyauth.CheckResponse{
		Status: &statusv3.Status{Code: int32(codes.Unauthenticated)},
		HttpResponse: &envoyauth.CheckResponse_DeniedResponse{
			DeniedResponse: &envoyauth.DeniedHttpResponse{
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
