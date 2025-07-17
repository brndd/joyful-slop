package mappingrules

import (
	"testing"

	"github.com/stretchr/testify/suite"
)

type MappingRuleBaseTests struct {
	suite.Suite
}

func TestRunnerMappingRuleBaseTests(t *testing.T) {
	suite.Run(t, new(MappingRuleBaseTests))
}

func (t *MappingRuleBaseTests) TestNewMappingRuleBase() {
	t.Run("No Modes", func() {
		base := NewMappingRuleBase("foo", []string{})
		t.Equal("foo", base.Name)
		t.EqualValues([]string{"*"}, base.Modes)
	})

	t.Run("Has Modes", func() {
		base := NewMappingRuleBase("foo", []string{"bar", "baz"})
		t.Equal("foo", base.Name)
		t.Contains(base.Modes, "bar")
		t.Contains(base.Modes, "baz")
		t.NotContains(base.Modes, "*")
	})
}

func (t *MappingRuleBaseTests) TestModeCheck() {
	t.Run("* works on all modes", func() {
		base := NewMappingRuleBase("", []string{})
		mode := "bar"
		t.True(base.modeCheck(&mode))
		mode = "baz"
		t.True(base.modeCheck(&mode))
	})

	t.Run("single mode only works in that mode", func() {
		base := NewMappingRuleBase("", []string{"bar"})
		mode := "bar"
		t.True(base.modeCheck(&mode))
		mode = "baz"
		t.False(base.modeCheck(&mode))
	})

	t.Run("multiple modes work in each mode", func() {

	})
}
