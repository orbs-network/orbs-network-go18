package test

import (
	"github.com/stretchr/testify/require"
	"strings"
	"testing"
)

func TestHelp(t *testing.T) {
	out, err := GammaCli().Run("help")
	t.Log(out)
	require.Error(t, err, "help should exit nonzero")
	require.NotEmpty(t, out, "help output should not be empty")
	require.True(t, strings.Contains(out, "start-local"))
	require.True(t, strings.Contains(out, "stop-local"))

	out2, err := GammaCli().Run()
	require.Error(t, err, "run without arguments should exit nonzero")
	require.Equal(t, out, out2, "help output should be equal")
}

func TestVersion(t *testing.T) {
	out, err := GammaCli().Run("version")
	require.NoError(t, err, "version should succeed")
	require.True(t, strings.Contains(out, "version"))
	t.Log(out)
}
