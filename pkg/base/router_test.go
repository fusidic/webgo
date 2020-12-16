package base

import (
	"fmt"
	"reflect"
	"testing"

	"github.com/stretchr/testify/assert"
)

func newTestRouter() *router {
	r := newRouter()
	// test router.go only, there for we use addRouter() rather than
	// r.Get("/", nil)
	r.addRouter("GET", "/", nil)
	r.addRouter("GET", "/hello/:name", nil)
	r.addRouter("GET", "hello/b/c", nil)
	r.addRouter("GET", "/hi/:name", nil)
	r.addRouter("GET", "assets/*filepath", nil)
	return r
}

func TestParsePattern(t *testing.T) {
	tests := []struct {
		name    string
		pattern string
		want    []string
	}{
		{
			name:    "pattern01",
			pattern: "/p/:name",
			want:    []string{"p", ":name"},
		},
		{
			name:    "pattern02",
			pattern: "/p/*",
			want:    []string{"p", "*"},
		},
		{
			name:    "pattern03",
			pattern: "/p/*name/*",
			want:    []string{"p", "*name"},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assert.Equal(t, tt.want, parsePattern(tt.pattern))
		})
	}
	ok := reflect.DeepEqual(parsePattern("/p/:name"), []string{"p", ":name"})
	ok = ok && reflect.DeepEqual(parsePattern("/p/*"), []string{"p", "*"})
	ok = ok && reflect.DeepEqual(parsePattern("/p/*name/*"), []string{"p", "*name"})
	if !ok {
		t.Fatal("test 'func parsePattern(pattern string) []string' failed")
	}
}

func TestGetRoute(t *testing.T) {
	r := newTestRouter()
	n, ps := r.getRoute("GET", "/hello/fusidic")
	if n == nil {
		t.Fatal("nil shouln't be returned")
	}
	if n.pattern != "/hello/:name" {
		t.Fatal("should match /hello/:name")
	}
	if ps["name"] != "fusidic" {
		t.Fatal("name should be 'fusidic'")
	}
	fmt.Printf("matched path: %s, params['name']: %s \n", n.pattern, ps["name"])
}
