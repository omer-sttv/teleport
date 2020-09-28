/*
Copyright 2020 Gravitational, Inc.

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package local

import (
	"context"

	"github.com/gravitational/teleport/lib/backend"
	"github.com/gravitational/teleport/lib/services"

	"github.com/gravitational/trace"
)

func (s *IdentityService) GetAppWebSession(ctx context.Context, req services.GetAppWebSessionRequest) (services.WebSession, error) {
	if err := req.Check(); err != nil {
		return nil, trace.Wrap(err)
	}

	item, err := s.Get(ctx, backend.Key(webPrefix, sessionsPrefix, appsPrefix, req.Username, req.ParentHash, req.SessionID))
	if err != nil {
		return nil, trace.Wrap(err)
	}

	session, err := services.GetWebSessionMarshaler().UnmarshalWebSession(item.Value)
	if err != nil {
		return nil, trace.Wrap(err)
	}
	return session, nil
}

func (s *IdentityService) GetAppWebSessions(ctx context.Context) ([]services.WebSession, error) {
	startKey := backend.Key(webPrefix, sessionsPrefix, appsPrefix)
	result, err := s.GetRange(ctx, startKey, backend.RangeEnd(startKey), backend.NoLimit)
	if err != nil {
		return nil, trace.Wrap(err)
	}

	out := make([]services.WebSession, len(result.Items))
	for i, item := range result.Items {
		session, err := services.GetWebSessionMarshaler().UnmarshalWebSession(item.Value)
		if err != nil {
			return nil, trace.Wrap(err)
		}
		out[i] = session
	}
	return out, nil
}

func (s *IdentityService) UpsertAppWebSession(ctx context.Context, session services.WebSession) error {
	value, err := services.GetWebSessionMarshaler().MarshalWebSession(session)
	if err != nil {
		return trace.Wrap(err)
	}
	item := backend.Item{
		Key:     backend.Key(webPrefix, sessionsPrefix, appsPrefix, session.GetUser(), session.GetParentHash(), session.GetName()),
		Value:   value,
		Expires: session.GetExpiryTime(),
	}

	if _, err = s.Put(ctx, item); err != nil {
		return trace.Wrap(err)
	}
	return nil
}

func (s *IdentityService) DeleteAppWebSession(ctx context.Context, req services.DeleteAppWebSessionRequest) error {
	if err := s.Delete(ctx, backend.Key(webPrefix, sessionsPrefix, appsPrefix, req.Username, req.ParentHash, req.SessionID)); err != nil {
		return trace.Wrap(err)
	}
	return nil
}

func (s *IdentityService) DeleteAllAppWebSessions(ctx context.Context) error {
	startKey := backend.Key(webPrefix, sessionsPrefix, appsPrefix)
	if err := s.DeleteRange(ctx, startKey, backend.RangeEnd(startKey)); err != nil {
		return trace.Wrap(err)
	}
	return nil
}

func (s *IdentityService) GetAppSession(ctx context.Context, sessionID string) (services.AppSession, error) {
	item, err := s.Get(ctx, backend.Key(sessionsPrefix, appsPrefix, sessionID))
	if err != nil {
		return nil, trace.Wrap(err)
	}

	session, err := services.GetAppSessionMarshaler().UnmarshalAppSession(item.Value)
	if err != nil {
		return nil, trace.Wrap(err)
	}
	return session, nil
}

func (s *IdentityService) GetAppSessions(ctx context.Context) ([]services.AppSession, error) {
	startKey := backend.Key(sessionsPrefix, appsPrefix)
	result, err := s.GetRange(ctx, startKey, backend.RangeEnd(startKey), backend.NoLimit)
	if err != nil {
		return nil, trace.Wrap(err)
	}

	out := make([]services.AppSession, len(result.Items))
	for i, item := range result.Items {
		session, err := services.GetAppSessionMarshaler().UnmarshalAppSession(item.Value)
		if err != nil {
			return nil, trace.Wrap(err)
		}
		out[i] = session
	}
	return out, nil
}

func (s *IdentityService) UpsertAppSession(ctx context.Context, session services.AppSession) error {
	value, err := services.GetAppSessionMarshaler().MarshalAppSession(session)
	if err != nil {
		return trace.Wrap(err)
	}
	item := backend.Item{
		Key:     backend.Key(sessionsPrefix, appsPrefix, session.GetName()),
		Value:   value,
		Expires: session.Expiry(),
	}

	if _, err = s.Put(ctx, item); err != nil {
		return trace.Wrap(err)
	}
	return nil
}

func (s *IdentityService) DeleteAppSession(ctx context.Context, sessionID string) error {
	if err := s.Delete(ctx, backend.Key(sessionsPrefix, appsPrefix, sessionID)); err != nil {
		return trace.Wrap(err)
	}
	return nil
}

func (s *IdentityService) DeleteAllAppSessions(ctx context.Context) error {
	startKey := backend.Key(sessionsPrefix, appsPrefix)
	if err := s.DeleteRange(ctx, startKey, backend.RangeEnd(startKey)); err != nil {
		return trace.Wrap(err)
	}
	return nil
}
