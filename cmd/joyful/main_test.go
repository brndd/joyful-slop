package main

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestShouldAnnounceModeChange(t *testing.T) {
	require.True(t, shouldAnnounceModeChange("SCM Mode", "Nav Mode", false))
	require.False(t, shouldAnnounceModeChange("SCM Mode", "SCM Mode", false))
	require.False(t, shouldAnnounceModeChange("SCM Mode", "Modifier", true))
	require.False(t, shouldAnnounceModeChange("Modifier", "SCM Mode", true))
}
