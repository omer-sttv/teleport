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

package jwt

import (
	"fmt"
	"testing"
	"time"

	"github.com/gravitational/teleport/lib/utils"
	"github.com/jonboulle/clockwork"

	"gopkg.in/check.v1"
)

type Suite struct{}

var _ = fmt.Printf
var _ = check.Suite(&Suite{})

func TestJWT(t *testing.T) { check.TestingT(t) }

func (s *Suite) SetUpSuite(c *check.C) {
	utils.InitLoggerForTests()
}
func (s *Suite) TearDownSuite(c *check.C) {}
func (s *Suite) SetUpTest(c *check.C)     {}
func (s *Suite) TearDownTest(c *check.C)  {}

func (s *Suite) TestSignAndVerify(c *check.C) {
	publicKey, privateKey, err := GenerateKeypair()
	c.Assert(err, check.IsNil)

	// Create a new key that can sign and verify tokens.
	key, err := New(&Config{
		PublicKey:  publicKey,
		PrivateKey: privateKey,
	})
	c.Assert(err, check.IsNil)

	// Sign a token with the new key.
	token, err := key.Sign(&SignParams{
		Username: "foo@example.com",
		Roles:    []string{"foo", "bar"},
		Expiry:   1 * time.Minute,
	})
	c.Assert(err, check.IsNil)

	// Verify that the token can be validated and values match expected values.
	claims, err := key.Verify(token)
	c.Assert(err, check.IsNil)
	c.Assert(claims.GetUsername(), check.Equals, "foo@example.com")
	c.Assert(claims.GetRoles(), check.DeepEquals, []string{"foo", "bar"})
}

// TestPublicOnlyVerify checks that a non-signing key used to validate a JWT
// can be created.
func (s *Suite) TestPublicOnlyVerify(c *check.C) {
	publicKey, privateKey, err := GenerateKeypair()
	c.Assert(err, check.IsNil)

	// Create a new key that can sign and verify tokens.
	key, err := New(&Config{
		PublicKey:  publicKey,
		PrivateKey: privateKey,
	})
	c.Assert(err, check.IsNil)

	// Sign a token with the new key.
	token, err := key.Sign(&SignParams{
		Username: "foo@example.com",
		Roles:    []string{"foo", "bar"},
		Expiry:   1 * time.Minute,
	})
	c.Assert(err, check.IsNil)

	// Create a new key that can only verify tokens and make sure the token
	// values match the expected values.
	key, err = NewPublic(&Config{
		PublicKey: publicKey,
	})
	c.Assert(err, check.IsNil)
	claims, err := key.Verify(token)
	c.Assert(err, check.IsNil)
	c.Assert(claims.GetUsername(), check.Equals, "foo@example.com")
	c.Assert(claims.GetRoles(), check.DeepEquals, []string{"foo", "bar"})

	// Make sure this key returns an error when trying to sign.
	_, err = key.Sign(&SignParams{})
	c.Assert(err, check.NotNil)
}

// TestExpiry checks that tokens expiration works.
func (s *Suite) TestExpiry(c *check.C) {
	publicKey, privateKey, err := GenerateKeypair()
	c.Assert(err, check.IsNil)

	clock := clockwork.NewFakeClockAt(time.Now())

	// Create a new key that can be used to sign and verify tokens.
	key, err := New(&Config{
		Clock:      clock,
		PublicKey:  publicKey,
		PrivateKey: privateKey,
	})
	c.Assert(err, check.IsNil)

	// Sign a token with a 1 minute expiration.
	token, err := key.Sign(&SignParams{
		Username: "foo@example.com",
		Roles:    []string{"foo", "bar"},
		Expiry:   1 * time.Minute,
	})
	c.Assert(err, check.IsNil)

	// Verify that the token is still valid.
	claims, err := key.Verify(token)
	c.Assert(err, check.IsNil)
	c.Assert(claims.GetUsername(), check.Equals, "foo@example.com")
	c.Assert(claims.GetRoles(), check.DeepEquals, []string{"foo", "bar"})

	// Advance time by two minutes and verify the token is no longer valid.
	clock.Advance(2 * time.Minute)
	_, err = key.Verify(token)
	c.Assert(err, check.NotNil)
}
