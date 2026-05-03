package variables

import (
	"fmt"
	"regexp"

	reqerrors "reqium/internal/errors"
	"reqium/internal/models"
)

type TemplateVariableResolver struct {
	pattern *regexp.Regexp
}

func NewTemplateVariableResolver() *TemplateVariableResolver {
	return &TemplateVariableResolver{pattern: regexp.MustCompile(`\{\{\s*([a-zA-Z0-9_.-]+)\s*\}\}`)}
}

func (r *TemplateVariableResolver) Resolve(input string, variables map[string]string) (string, error) {
	var missing string
	result := r.pattern.ReplaceAllStringFunc(input, func(match string) string {
		parts := r.pattern.FindStringSubmatch(match)
		if len(parts) != 2 {
			return match
		}
		value, ok := variables[parts[1]]
		if !ok {
			missing = parts[1]
			return match
		}
		return value
	})
	if missing != "" {
		return "", fmt.Errorf("%w: %s", reqerrors.ErrVariableNotFound, missing)
	}
	return result, nil
}

func (r *TemplateVariableResolver) ResolveRequest(req models.Request, variables map[string]string) (models.Request, error) {
	var err error
	req.URL, err = r.Resolve(req.URL, variables)
	if err != nil {
		return models.Request{}, err
	}

	headers := make(map[string]string, len(req.Headers))
	for key, value := range req.Headers {
		resolved, err := r.Resolve(value, variables)
		if err != nil {
			return models.Request{}, err
		}
		headers[key] = resolved
	}
	req.Headers = headers

	if len(req.Body) > 0 {
		body, err := r.Resolve(string(req.Body), variables)
		if err != nil {
			return models.Request{}, err
		}
		req.Body = []byte(body)
	}

	return req, nil
}
