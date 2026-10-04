package executor

import (
	"net/http"
	"testing"

	"github.com/router-for-me/CLIProxyAPI/v8/internal/clienterror"
)

// Real upstream bodies captured from production relays. Each must rotate to the
// next credential instead of being returned to the client as a request fault.
func TestNewOpenAICompatStatusError_CredentialFailuresRotate(t *testing.T) {
	cases := []struct {
		name   string
		status int
		body   string
		want   int
	}{
		{"b.ai balance=0 as 400", http.StatusBadRequest,
			`{"error":{"type":"api_error","code":"insufficient_user_quota","message":"credit insufficient balance: balance=0 required=3084","param":""}}`,
			http.StatusPaymentRequired},
		{"relay out of quota as 403", http.StatusForbidden,
			`{"error":{"message":"Out of quota (balance=0, required=702)"}}`,
			http.StatusPaymentRequired},
		{"billing_error as 400", http.StatusBadRequest,
			`{"error":{"message":"Request exceeds your free allowance of 114,661 weighted tokens.","type":"billing_error","param":null,"code":null}}`,
			http.StatusPaymentRequired},
		{"bailian overdue payment as 400", http.StatusBadRequest,
			`{"error":"Access denied, please make sure your account is in good standing. For details, see: https://help.aliyun.com/zh/model-studio/error-code#overdue-payment"}`,
			http.StatusPaymentRequired},
		{"groq TPM cap as 413", http.StatusRequestEntityTooLarge,
			`{"error":{"message":"Request too large for model openai/gpt-oss-120b on tokens per minute (TPM): Limit 8000, Requested 8887, please reduce your message size and try again.","type":"tokens","code":"rate_limit_exceeded"}}`,
			http.StatusTooManyRequests},
		{"orcarouter free prompt cap as 400", http.StatusBadRequest,
			`{"error":{"message":"This prompt is longer than the free tier allows for a single request.","type":"invalid_request_error","code":"free_rate_limited","metadata":{"reason":"err_free_prompt_cap","retryable":false}}}`,
			http.StatusTooManyRequests},
		{"cloudflare rejects tool-result turns as 400", http.StatusBadRequest,
			`{"errors":[{"message":"AiError: Bad input: Error: oneOf at '/' not met, 0 matches: required properties at '/' are 'prompt'","code":5006}],"success":false}`,
			http.StatusTooManyRequests},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			err := newOpenAICompatStatusError(tc.status, http.Header{}, []byte(tc.body))
			if err.StatusCode() != tc.want {
				t.Fatalf("status = %d, want %d", err.StatusCode(), tc.want)
			}
			if clienterror.IsRequestFault(err.StatusCode(), err) {
				t.Fatalf("classified as request fault; conductor would not rotate to the next credential")
			}
		})
	}
}

func TestNewOpenAICompatStatusError_RealRequestFaultUnchanged(t *testing.T) {
	err := newOpenAICompatStatusError(http.StatusBadRequest, http.Header{},
		[]byte(`{"error":{"message":"field messages is required","type":"invalid_request_error"}}`))
	if err.StatusCode() != http.StatusBadRequest {
		t.Fatalf("status = %d, want 400", err.StatusCode())
	}
	if !clienterror.IsRequestFault(err.StatusCode(), err) {
		t.Fatalf("a malformed request must stay a request fault and not sweep the credential pool")
	}
}
