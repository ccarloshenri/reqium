package app

import (
	"context"
	"time"

	"reqium/internal/enums"
	"reqium/internal/interfaces"
	"reqium/internal/models"
)

type RunnerService struct {
	collections interfaces.CollectionRepository
	requests    *RequestService
	resolver    interfaces.VariableResolver
	envs        *EnvironmentService
}

func NewRunnerService(collections interfaces.CollectionRepository, requests *RequestService, resolver interfaces.VariableResolver, envs *EnvironmentService) *RunnerService {
	return &RunnerService{collections: collections, requests: requests, resolver: resolver, envs: envs}
}

func (s *RunnerService) Run(ctx context.Context, collectionName string, envName string, timeout time.Duration) ([]models.RunResult, error) {
	collection, err := s.collections.Get(ctx, collectionName)
	if err != nil {
		return nil, err
	}

	variables, err := s.envs.Variables(ctx, envName)
	if err != nil {
		return nil, err
	}

	results := make([]models.RunResult, 0, len(collection.Requests))
	for _, saved := range collection.Requests {
		req := saved.ToRequest(timeout)
		if s.resolver != nil {
			req, err = s.resolver.ResolveRequest(req, variables)
			if err != nil {
				results = append(results, failedRunResult(saved, req, err))
				continue
			}
		}

		response, err := s.requests.Send(ctx, req)
		result := models.RunResult{
			RequestName: saved.Name,
			Method:      req.Method,
			URL:         req.URL,
			StatusCode:  response.StatusCode,
			Duration:    response.Duration,
			Status:      enums.RunnerPassed,
		}
		if err != nil || response.StatusCode >= 400 {
			result.Status = enums.RunnerFailed
		}
		if err != nil {
			result.Error = err.Error()
		}
		results = append(results, result)
	}

	return results, nil
}

func failedRunResult(saved models.SavedRequest, req models.Request, err error) models.RunResult {
	return models.RunResult{
		RequestName: saved.Name,
		Method:      req.Method,
		URL:         req.URL,
		Status:      enums.RunnerFailed,
		Error:       err.Error(),
	}
}
