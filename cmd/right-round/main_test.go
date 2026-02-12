package main

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestRootCmd_InvalidTypeFlag(t *testing.T) {
	cmd := rootCmd()
	cmd.SetArgs([]string{"--type", "invalid"})

	err := cmd.Execute()
	assert.Error(t, err)
}

func TestRootCmd_Version(t *testing.T) {
	cmd := rootCmd()
	assert.Equal(t, "dev", cmd.Version)
}

func TestRootCmd_Flags(t *testing.T) {
	cmd := rootCmd()

	typeFlag := cmd.Flags().Lookup("type")
	assert.NotNil(t, typeFlag)
	assert.Equal(t, "", typeFlag.DefValue)

	groupFlag := cmd.Flags().Lookup("group")
	assert.NotNil(t, groupFlag)
	assert.Equal(t, "", groupFlag.DefValue)
}
