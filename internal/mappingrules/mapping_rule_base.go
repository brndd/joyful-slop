package mappingrules

import "slices"

type MappingRuleBase struct {
	Name  string
	Modes []string
}

func (rule *MappingRuleBase) modeCheck(mode *string) bool {
	if rule.Modes[0] == "*" {
		return true
	}
	return slices.Contains(rule.Modes, *mode)
}
