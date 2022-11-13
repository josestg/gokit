package httpx

import (
	"net/http"
	"net/http/httptest"
	"reflect"
	"testing"
)

func TestChain(t *testing.T) {
	executionTrace := make([]int, 0)

	factory := func(index int) Middleware {
		return func(handler Handler) HandlerFunc {
			return func(w http.ResponseWriter, r *http.Request) error {
				executionTrace = append(executionTrace, index)

				defer func() {
					executionTrace = append(executionTrace, -1*index)
				}()

				return handler.ServeHTTP(w, r)
			}
		}
	}

	rootHandler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		executionTrace = append(executionTrace, 0)
	})

	chain := NewChain(factory(1), factory(2), factory(3))
	handler1 := chain.Extend(factory(4), factory(5)).ToHandler(rootHandler)
	handler2 := chain.ToHandler(rootHandler)

	req := httptest.NewRequest(http.MethodGet, "/", nil)
	rec := httptest.NewRecorder()
	handler1.ServeHTTP(rec, req)

	expected := []int{1, 2, 3, 4, 5, 0, -5, -4, -3, -2, -1}
	if len(executionTrace) != len(expected) {
		t.Fatalf("expected length %d; got %d", len(expected), len(executionTrace))
	}

	if !reflect.DeepEqual(executionTrace, expected) {
		t.Fatalf("expected %v; got %v", expected, executionTrace)
	}

	// reset execution trace for next test
	executionTrace = make([]int, 0)
	expected = []int{1, 2, 3, 0, -3, -2, -1}
	handler2.ServeHTTP(rec, req)

	if len(executionTrace) != len(expected) {
		t.Fatalf("expected length %d; got %d", len(expected), len(executionTrace))
	}

	if !reflect.DeepEqual(executionTrace, expected) {
		t.Fatalf("expected %v; got %v", expected, executionTrace)
	}
}
