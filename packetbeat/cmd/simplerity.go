package cmd

import (
	"github.com/spf13/cobra"

	"github.com/aliksend/beats/packetbeat/simplerity"
)

func simplerityCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "simplerity",
		Short: "Simplerity integration",
	}
	cmd.AddCommand(loginCommand(), selectCommand(), loadCommand())
	return cmd
}

func loginCommand() *cobra.Command {
	var login, password string
	cmd := &cobra.Command{
		Use:   "login",
		Short: "Login to simplerity",
		Run: func(cmd *cobra.Command, args []string) {
			simplerity.Login(login, password)
		},
	}
	cmd.Flags().StringVarP(&login, "login", "l", "", "login (required)")
	cmd.MarkFlagRequired("login")
	cmd.Flags().StringVarP(&password, "password", "p", "", "password (required)")
	cmd.MarkFlagRequired("password")
	return cmd
}

func selectCommand() *cobra.Command {
	var agentID string
	cmd := &cobra.Command{
		Use:   "select",
		Short: "Select agent to use",
		Run: func(cmd *cobra.Command, args []string) {
			simplerity.Select(agentID)
		},
	}
	cmd.Flags().StringVarP(&agentID, "agent", "a", "", "agent id (required)")
	cmd.MarkFlagRequired("agent-id")
	return cmd
}

func loadCommand() *cobra.Command {
	var configFileName string
	cmd := &cobra.Command{
		Use:   "load",
		Short: "Load config",
		Run: func(cmd *cobra.Command, args []string) {
			simplerity.LoadConfig(configFileName)
		},
	}
	cmd.Flags().StringVarP(&configFileName, "save-to-file", "", "packetbeat.yml", "name of file to save config to")
	return cmd
}
