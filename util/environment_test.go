//   Copyright 2016 Wercker Holding BV
//
//   Licensed under the Apache License, Version 2.0 (the "License");
//   you may not use this file except in compliance with the License.
//   You may obtain a copy of the License at
//
//       http://www.apache.org/licenses/LICENSE-2.0
//
//   Unless required by applicable law or agreed to in writing, software
//   distributed under the License is distributed on an "AS IS" BASIS,
//   WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
//   See the License for the specific language governing permissions and
//   limitations under the License.

package util

import (
	"testing"

	"github.com/stretchr/testify/suite"
)

type EnvironmentSuite struct {
	TestSuite
}

func (s *EnvironmentSuite) SetupTest() {
	s.TestSuite.SetupTest()
}

func TestEnvironmentSuite(t *testing.T) {
	suiteTester := new(EnvironmentSuite)
	suite.Run(t, suiteTester)
}

func (s *EnvironmentSuite) TestPassthru() {
	env := NewEnvironment("X_PUBLIC=foo", "XXX_PRIVATE=bar", "NOT=included")
	s.Equal(1, len(env.GetPassthru().Ordered()))
	s.Equal(1, len(env.GetHiddenPassthru().Ordered()))
}

func (s *EnvironmentSuite) TestInterpolate() {
	env := NewEnvironment("PUBLIC=foo", "X_PRIVATE=zed", "XXX_OTHER=otter")
	env.Update(env.GetPassthru().Ordered())
	env.Hidden.Update(env.GetHiddenPassthru().Ordered())

	// this is impossible to set because the order the variables are applied is
	// defined by the caller
	//env.Update([][]string{[]string{"X_PUBLIC", "bar"}})
	//tt.Equal(env.Interpolate("$PUBLIC"), "foo", "Non-prefixed should alias any X_ prefixed vars.")
	s.Equal(env.Interpolate("${PUBLIC}"), "foo", "Alternate shell style vars should work.")

	// NB: stipping only works because we cann Update with the passthru
	// function above
	s.Equal(env.Interpolate("$PRIVATE"), "zed", "Xs should be stripped.")
	s.Equal(env.Interpolate("$OTHER"), "otter", "XXXs should be stripped.")
	s.Equal(env.Interpolate("one two $PUBLIC bar"), "one two foo bar", "interpolation should work in middle of string.")
}

func (s *EnvironmentSuite) TestOrdered() {
	env := NewEnvironment("PUBLIC=foo", "X_PRIVATE=zed")
	expected := [][]string{[]string{"PUBLIC", "foo"}, []string{"X_PRIVATE", "zed"}}
	s.Equal(env.Ordered(), expected)
}

func (s *EnvironmentSuite) TestExport() {
	env := NewEnvironment("PUBLIC=foo", "X_PRIVATE=zed")
	expected := []string{`export PUBLIC="foo"`, `export X_PRIVATE="zed"`}
	s.Equal(env.Export(), expected)
}

func (s *EnvironmentSuite) TestLoadFile() {
	env := NewEnvironment("PUBLIC=foo")
	expected := [][]string{
		[]string{"PUBLIC", "foo"},
		[]string{"A", "1"},
		[]string{"B", "2"},
		[]string{"C", "3"},
		[]string{"D", "4"},
		[]string{"E", "5"},
		[]string{"F", "6"},
		[]string{"G", "7"},
	}
	env.LoadFile("../tests/environment-test-load-file.env")
	s.Equal(8, len(env.Ordered()), "Should only load 8 valid lines.")
	s.Equal("foo", env.Get("PUBLIC"), "LoadFile should ignore keys already set in env.")
	s.Equal(expected, env.Ordered(), "LoadFile should maintain order.")
}
