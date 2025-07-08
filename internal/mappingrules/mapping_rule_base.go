package mappingrules

import "slices"

type MappingRuleBase struct {
	Name  string
	Modes []string
}

func NewMappingRuleBase(
	name string,
	modes []string,
) MappingRuleBase {
	if len(modes) == 0 {
		modes = []string{"*"}
	}

	return MappingRuleBase{
		Name:  name,
		Modes: modes,
	}
}

func (rule *MappingRuleBase) modeCheck(mode *string) bool {
	if rule.Modes[0] == "*" {
		return true
	}
	return slices.Contains(rule.Modes, *mode)
}
