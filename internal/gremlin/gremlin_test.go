package gremlin

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"git.annabunches.net/annabunches/joyful/internal/configparser"
	"github.com/stretchr/testify/require"
)

func TestParseAndConvertSimpleButton(t *testing.T) {
	p := parseFixture(t, fixture(
		`<mode>Base</mode>`,
		inputXML("Base", "root-base"),
		vjoyAction("map-base", 1),
		rootAction("root-base", "map-base"),
	))
	result, err := Convert(p, Options{})
	require.NoError(t, err)
	require.Equal(t, map[string]int{"button": 1}, result.Counts)
	require.Contains(t, string(result.Files["buttons.yml"]), "button: 0")
}

func TestUnknownActionFails(t *testing.T) {
	_, err := Parse(strings.NewReader(fixture(
		`<mode>Base</mode>`,
		inputXML("Base", "root-base"),
		`<action id="mystery" type="execute-command"><property type="string"><name>command</name><value>oops</value></property></action>`,
		rootAction("root-base", "mystery"),
	)))
	require.ErrorContains(t, err, `unknown executable action type "execute-command"`)
}

func TestActionReferenceCycleFails(t *testing.T) {
	_, err := Parse(strings.NewReader(fixture(
		`<mode>Base</mode>`,
		inputXML("Base", "root-base"),
		rootAction("nested", "root-base"),
		rootAction("root-base", "nested"),
	)))
	require.ErrorContains(t, err, "action reference cycle")
}

func TestInheritanceOverride(t *testing.T) {
	p := parseFixture(t, fixture(
		`<mode>Base</mode><mode parent="Base">Child</mode>`,
		inputXML("Base", "root-base")+inputXML("Child", "root-child"),
		vjoyAction("map-base", 1)+vjoyAction("map-child", 2),
		rootAction("root-base", "map-base")+rootAction("root-child", "map-child"),
	))
	result, err := Convert(p, Options{})
	require.NoError(t, err)
	yaml := string(result.Files["buttons.yml"])
	require.Contains(t, yaml, "button: 0")
	require.Contains(t, yaml, "button: 1")
	require.Equal(t, 2, result.Counts["button"])
}

func TestHighButtonRemap(t *testing.T) {
	p := parseFixture(t, fixture(
		`<mode>Base</mode>`, inputXML("Base", "root-base"),
		vjoyAction("map-base", 124), rootAction("root-base", "map-base"),
	))
	result, err := Convert(p, Options{})
	require.NoError(t, err)
	require.Contains(t, string(result.Files["buttons.yml"]), "button: 59")
	require.Contains(t, strings.Join(result.Warnings, "\n"), "remapped vJoy button 124")
}

func TestMacroIsOmittedWithWarning(t *testing.T) {
	p := parseFixture(t, fixture(
		`<mode>Base</mode>`, inputXML("Base", "root-base"),
		`<action id="macro" type="macro"><macro-action type="key"/></action>`,
		rootAction("root-base", "macro"),
	))
	result, err := Convert(p, Options{})
	require.NoError(t, err)
	require.Empty(t, result.Counts)
	require.Contains(t, strings.Join(result.Warnings, "\n"), "omitted macro action macro")
}

func TestFullProfileGoldenTreeAndDeterminism(t *testing.T) {
	repoRoot := filepath.Clean(filepath.Join("..", ".."))
	input := filepath.Join(repoRoot, "joystick_gremlin_profile.xml")
	first := filepath.Join(t.TempDir(), "first")
	second := filepath.Join(t.TempDir(), "second")
	result, err := Generate(input, first, Options{})
	require.NoError(t, err)
	_, err = Generate(input, second, Options{})
	require.NoError(t, err)
	require.Equal(t, 139, totalCount(result.Counts))
	require.Equal(t, 12, result.Counts["axis"])
	require.Zero(t, result.Counts["axis-to-relaxis"])
	require.Equal(t, 1, result.Counts["mode-shift"])
	require.Len(t, result.Warnings, 15)
	axes := string(result.Files["axes.yml"])
	require.Contains(t, axes, "name: \"L-X-Axis\"\n    modes:")
	require.NotContains(t, axes, "device: \"left\"\n      axis: X\n      inverted: true")
	require.Contains(t, axes, "name: \"L-Y-Axis\"")
	require.NotContains(t, axes, "device: \"left\"\n      axis: Y\n      inverted: true")
	require.NotContains(t, axes, "device: \"left\"\n      axis: Z\n      inverted: true")
	require.NotContains(t, axes, "device: \"right\"\n      axis: Z\n      inverted: true")
	require.Contains(t, axes, "name: \"Right mini-stick Y\"")
	require.Contains(t, axes, "device: \"right\"\n      axis: RY\n      inverted: true\n    output:\n      device: \"right-vjoy\"\n      axis: RX")
	require.Contains(t, axes, "name: \"Right mini-stick X\"")
	require.Contains(t, axes, "device: \"right\"\n      axis: RX\n    output:\n      device: \"right-vjoy\"\n      axis: Z")
	require.Contains(t, axes, "name: \"Left twist\"")
	require.Contains(t, axes, "device: \"left\"\n      axis: RZ\n      inverted: true\n    output:\n      device: \"left-vjoy\"\n      axis: RY")
	require.Contains(t, axes, "name: \"Right twist\"")
	require.Contains(t, axes, "device: \"right\"\n      axis: RZ\n      inverted: true\n    output:\n      device: \"right-vjoy\"\n      axis: RY")
	require.NotContains(t, axes, "axis-to-relaxis")
	require.NotContains(t, axes, "device: \"mouse\"")

	for _, name := range outputNames {
		golden, err := os.ReadFile(filepath.Join(repoRoot, "joystick_gremlin_profile", name))
		require.NoError(t, err, name)
		generated, err := os.ReadFile(filepath.Join(first, name))
		require.NoError(t, err, name)
		repeated, err := os.ReadFile(filepath.Join(second, name))
		require.NoError(t, err, name)
		require.Equal(t, golden, generated, name)
		require.Equal(t, generated, repeated, name)
	}

	parsed, err := configparser.ParseConfig(first)
	require.NoError(t, err)
	require.Len(t, parsed.Devices, 4)
	require.Len(t, parsed.Modes, 4)
	require.Len(t, parsed.Rules, 139)
}

func parseFixture(t *testing.T, xml string) *Profile {
	t.Helper()
	p, err := Parse(strings.NewReader(xml))
	require.NoError(t, err)
	return p
}

func fixture(modes, inputs, leafActions, roots string) string {
	return fmt.Sprintf(`<profile version="14">
<inputs>%s</inputs><settings/><logical-device/><library>%s%s</library>
<modes>%s</modes><scripts/><devices><device><device-id>left-id</device-id><device-name>VKB OT L</device-name></device></devices>
</profile>`, inputs, leafActions, roots, modes)
}

func inputXML(mode, root string) string {
	return fmt.Sprintf(`<input><device-id>left-id</device-id><input-type>button</input-type><mode>%s</mode><input-id>1</input-id><action-configuration><root-action>%s</root-action><behavior>button</behavior></action-configuration></input>`, mode, root)
}

func rootAction(id, child string) string {
	return fmt.Sprintf(`<action id="%s" type="root"><actions><action-id>%s</action-id></actions></action>`, id, child)
}

func vjoyAction(id string, button int) string {
	return fmt.Sprintf(`<action id="%s" type="map-to-vjoy">
<property type="int"><name>vjoy-device-id</name><value>1</value></property>
<property type="int"><name>vjoy-input-id</name><value>%d</value></property>
<property type="input_type"><name>vjoy-input-type</name><value>button</value></property>
<property type="string"><name>action-label</name><value>Test</value></property>
</action>`, id, button)
}
