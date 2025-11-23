package rbac

import (
	"github.com/casbin/casbin/v2"
)

type Service struct {
	enforcer *casbin.Enforcer
}

func NewService(modelPath, policyPath string) (*Service, error) {
	enforcer, err := casbin.NewEnforcer(modelPath, policyPath)
	if err != nil {
		return nil, err
	}

	// Load policies from file
	if err := enforcer.LoadPolicy(); err != nil {
		return nil, err
	}

	return &Service{
		enforcer: enforcer,
	}, nil
}

func (s *Service) Enforce(sub, obj, act string) (bool, error) {
	return s.enforcer.Enforce(sub, obj, act)
}

func (s *Service) AddPolicy(sub, obj, act string) (bool, error) {
	return s.enforcer.AddPolicy(sub, obj, act)
}

func (s *Service) AddGroupingPolicy(sub, role string) (bool, error) {
	return s.enforcer.AddGroupingPolicy(sub, role)
}
