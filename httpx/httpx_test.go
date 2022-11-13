package httpx

import (
	"errors"
	"net/http"
	"net/http/httptest"
	"reflect"
	"testing"
)

func TestServeMux(t *testing.T) {
	mux := NewServeMux()

	t.Run("with no params", func(t *testing.T) {
		h := &mockHandler{t: t}
		mux.HandleFunc(http.MethodGet, "/a/b/c", h.Handler(nil))

		rec := httptest.NewRecorder()
		req := httptest.NewRequest(http.MethodGet, "/a/b/c", nil)
		mux.ServeHTTP(rec, req)

		h.verifyMethod(http.MethodGet)
		h.verifyPath("/a/b/c")
		h.verifyNumCalls(1)
		h.verifyParams()
	})

	t.Run("with params", func(t *testing.T) {
		h := &mockHandler{t: t}
		mux.HandleFunc(http.MethodGet, "/first/:var1/third/:var2", h.Handler(nil))

		rec := httptest.NewRecorder()
		req := httptest.NewRequest(http.MethodGet, "/first/123/third/456", nil)
		mux.ServeHTTP(rec, req)

		h.verifyMethod(http.MethodGet)
		h.verifyPath("/first/123/third/456")
		h.verifyNumCalls(1)
		h.verifyParams(param{key: "var1", val: "123"}, param{key: "var2", val: "456"})
	})

	t.Run("with non-nil error return", func(t *testing.T) {
		defer func() {
			if r := recover(); r == nil {
				t.Errorf("expected a panic")
			}
		}()

		h := &mockHandler{t: t}
		mux.HandleFunc(http.MethodGet, "/", h.Handler(errors.New("any error")))

		rec := httptest.NewRecorder()
		req := httptest.NewRequest(http.MethodGet, "/", nil)
		mux.ServeHTTP(rec, req)

		h.verifyMethod(http.MethodGet)
		h.verifyPath("/")
		h.verifyNumCalls(1)
		h.verifyParams()
	})
}

func TestServeMux_WithGlobalChain(t *testing.T) {

	tracer := newTracer(t)
	chain := NewChain(tracer.factory(1), tracer.factory(2))

	mux := NewServeMuxWithChain(chain)

	h := &mockHandler{t: t}
	mux.HandleFunc(http.MethodGet, "/a/b/c", h.Handler(nil))

	rec := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodGet, "/a/b/c", nil)
	mux.ServeHTTP(rec, req)

	h.verifyMethod(http.MethodGet)
	h.verifyPath("/a/b/c")
	h.verifyNumCalls(1)
	h.verifyParams()

	expectedExecutionTrace := []int{1, 2, 2, 1}
	tracer.verifyExecutionTrace(expectedExecutionTrace)
}

func TestServeMux_WithLocalChain(t *testing.T) {

	tracer := newTracer(t)

	mux := NewServeMux()

	h := &mockHandler{t: t}
	mux.HandleFunc(http.MethodGet, "/a/b/c", h.Handler(nil), tracer.factory(1), tracer.factory(2))

	rec := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodGet, "/a/b/c", nil)
	mux.ServeHTTP(rec, req)

	h.verifyMethod(http.MethodGet)
	h.verifyPath("/a/b/c")
	h.verifyNumCalls(1)
	h.verifyParams()

	expectedExecutionTrace := []int{1, 2, 2, 1}
	tracer.verifyExecutionTrace(expectedExecutionTrace)
}

func TestServeMux_WithGlobalAndLocalChain(t *testing.T) {

	tracer := newTracer(t)

	mux := NewServeMuxWithChain(NewChain(tracer.factory(1), tracer.factory(2)))

	h := &mockHandler{t: t}
	mux.HandleFunc(http.MethodGet, "/a/b/c", h.Handler(nil), tracer.factory(3), tracer.factory(4))

	rec := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodGet, "/a/b/c", nil)
	mux.ServeHTTP(rec, req)

	h.verifyMethod(http.MethodGet)
	h.verifyPath("/a/b/c")
	h.verifyNumCalls(1)
	h.verifyParams()

	expectedExecutionTrace := []int{1, 2, 3, 4, 4, 3, 2, 1}
	tracer.verifyExecutionTrace(expectedExecutionTrace)
}

type param struct {
	key, val string
}

type mockHandler struct {
	path   string
	method string
	params Params
	calls  int
	t      *testing.T
}

func (m *mockHandler) verifyMethod(method string) {
	if m.method != method {
		m.t.Errorf("expected method %v; got method %v", method, m.method)
	}
}

func (m *mockHandler) verifyPath(path string) {
	if m.path != path {
		m.t.Errorf("expected path %v; got path %v", path, m.path)
	}
}

func (m *mockHandler) verifyNumCalls(n int) {
	if m.calls != n {
		m.t.Errorf("expected called %v times; got %v times", n, m.calls)
	}
}
func (m *mockHandler) verifyParams(params ...param) {
	for _, p := range params {
		if m.params.Get(p.key) != p.val {
			m.t.Errorf("expected has params name %s with value %s", p.key, p.val)
		}
	}
}

func (m *mockHandler) Handler(returnsError error) HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) error {
		m.calls++
		m.method = r.Method
		m.params = ParamsFromContext(r.Context())
		m.path = r.URL.Path
		return returnsError
	}
}

type tracer struct {
	t              *testing.T
	executionTrace []int
}

func newTracer(t *testing.T) *tracer {
	return &tracer{
		t:              t,
		executionTrace: make([]int, 0),
	}
}

func (t *tracer) factory(index int) Middleware {
	return func(handler Handler) HandlerFunc {
		return func(w http.ResponseWriter, r *http.Request) error {
			t.executionTrace = append(t.executionTrace, index)

			defer func() {
				t.executionTrace = append(t.executionTrace, index)
			}()

			return handler.ServeHTTP(w, r)
		}
	}
}

func (t *tracer) verifyExecutionTrace(expectedExecutionTrace []int) {
	t.t.Helper()
	if !reflect.DeepEqual(t.executionTrace, expectedExecutionTrace) {
		t.t.Errorf("expected execution trace is %v; got %v", expectedExecutionTrace, t.executionTrace)
	}
}
