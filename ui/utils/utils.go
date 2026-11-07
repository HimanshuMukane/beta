package utils

import (
	"context"
	"crypto/rand"
	"fmt"
	"regexp"
	"slices"
	"strings"
	"time"

	"github.com/a-h/templ"

	twmerge "github.com/Oudwins/tailwind-merge-go"
)

var (
	twV4Var = regexp.MustCompile(`-\((--[A-Za-z0-9-]+)\)`)
	twV3Var = regexp.MustCompile(`-\[var\((--[A-Za-z0-9-]+)\)\]`)
)

func lastTopLevelColon(s string) int {
	depth, last := 0, -1
	for i, r := range s {
		switch r {
		case '[':
			depth++
		case ']':
			depth--
		case ':':
			if depth == 0 {
				last = i
			}
		}
	}
	return last
}

func twImportantToV3(token string) string {
	if !strings.HasSuffix(token, "!") {
		return token
	}
	base := strings.TrimSuffix(token, "!")
	i := lastTopLevelColon(base)
	return base[:i+1] + "!" + base[i+1:]
}

func twImportantToV4(token string) string {
	i := lastTopLevelColon(token)
	if i+1 >= len(token) || token[i+1] != '!' {
		return token
	}
	return token[:i+1] + token[i+2:] + "!"
}

func twMapTokens(s string, f func(string) string) string {
	if !strings.Contains(s, "!") {
		return s
	}
	tokens := strings.Fields(s)
	for i, t := range tokens {
		tokens[i] = f(t)
	}
	return strings.Join(tokens, " ")
}

func twReorder(input, merged string) string {
	survivors := strings.Fields(merged)
	if len(survivors) < 2 {
		return merged
	}
	remaining := make(map[string]int, len(survivors))
	for _, c := range survivors {
		remaining[c]++
	}
	ordered := make([]string, 0, len(survivors))
	tokens := strings.Fields(input)
	for i := len(tokens) - 1; i >= 0; i-- {
		if remaining[tokens[i]] == 0 {
			continue
		}
		remaining[tokens[i]]--
		ordered = append(ordered, tokens[i])
	}

	if len(ordered) != len(survivors) {
		return merged
	}
	slices.Reverse(ordered)
	return strings.Join(ordered, " ")
}

func CN(classes ...any) string {
	flat := flattenClasses(classes)
	merged := make([]string, len(flat))
	for i, c := range flat {
		merged[i] = twMapTokens(twV4Var.ReplaceAllString(c, "-[var($1)]"), twImportantToV3)
	}
	input := strings.Join(merged, " ")
	out := twMapTokens(twReorder(input, twmerge.Merge(merged...)), twImportantToV4)
	return twV3Var.ReplaceAllString(out, "-($1)")
}

func flattenClasses(values []any) []string {
	var out []string
	for _, v := range values {
		switch t := v.(type) {
		case nil:
		case string:
			if t != "" {
				out = append(out, t)
			}
		case []string:
			for _, c := range t {
				if c != "" {
					out = append(out, c)
				}
			}
		case []any:
			out = append(out, flattenClasses(t)...)
		case map[string]bool:
			keys := make([]string, 0, len(t))
			for c, on := range t {
				if on && c != "" {
					keys = append(keys, c)
				}
			}
			slices.Sort(keys)
			out = append(out, keys...)
		}
	}
	return out
}

func Ptr[T any](v T) *T {
	return &v
}

func IfElse[T any](condition bool, trueValue T, falseValue T) T {
	if condition {
		return trueValue
	}
	return falseValue
}

func Map[T, U any](s []T, f func(T) U) []U {
	out := make([]U, len(s))
	for i, v := range s {
		out[i] = f(v)
	}
	return out
}

func MergeAttributes(attrs ...templ.Attributes) templ.Attributes {
	merged := templ.Attributes{}
	for _, attr := range attrs {
		for k, v := range attr {
			merged[k] = v
		}
	}
	return merged
}

type directionKey struct{}

func WithDirection(ctx context.Context, direction string) context.Context {
	return context.WithValue(ctx, directionKey{}, direction)
}

func Direction(ctx context.Context) string {
	if direction, ok := ctx.Value(directionKey{}).(string); ok && direction != "" {
		return direction
	}
	return "ltr"
}

func DirectionAttributes(ctx context.Context) templ.Attributes {
	if direction, ok := ctx.Value(directionKey{}).(string); ok && direction != "" {
		return templ.Attributes{"data-templ-direction": direction}
	}
	return templ.Attributes{}
}

func RandomID() string {
	return fmt.Sprintf("id-%s", rand.Text())
}

var ScriptVersion = fmt.Sprintf("%d", time.Now().Unix())

var ScriptURL = func(path string) string {
	return path + "?v=" + ScriptVersion
}
