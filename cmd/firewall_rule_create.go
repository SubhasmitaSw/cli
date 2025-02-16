package cmd

import (
	"fmt"
	"strings"

	"github.com/civo/civogo"
	"github.com/civo/cli/config"
	"github.com/civo/cli/utility"

	"os"

	"github.com/spf13/cobra"
)

var protocol, startPort, endPort, direction, label, cidr string

var firewallRuleCreateCmd = &cobra.Command{
	Use:     "create",
	Aliases: []string{"new", "add"},
	Short:   "Create a new firewall rule",
	Args:    cobra.MinimumNArgs(1),
	Example: "civo firewall rule create FIREWALL_NAME/FIREWALL_ID [flags]",
	Run: func(cmd *cobra.Command, args []string) {
		utility.EnsureCurrentRegion()

		client, err := config.CivoAPIClient()
		if regionSet != "" {
			client.Region = regionSet
		}
		if err != nil {
			utility.Error("Creating the connection to Civo's API failed with %s", err)
			os.Exit(1)
		}

		firewall, err := client.FindFirewall(args[0])
		if err != nil {
			utility.Error("%s", err)
			os.Exit(1)
		}

		newRuleConfig := &civogo.FirewallRuleConfig{
			FirewallID: firewall.ID,
			Protocol:   protocol,
			StartPort:  startPort,
			Cidr:       strings.Split(cidr, ","),
			Label:      label,
		}

		// Check the rule address, if the input is different
		// from (inbound or outbound) then we will generate an error
		if direction == "ingress" {
			newRuleConfig.Direction = direction
		} else if direction == "" {
			utility.Error("'--direction' flag can't be empty")
			os.Exit(1)
		} else {
			utility.Error("'--direction' flag only support 'ingress' as of now, not '%s'", direction)
			os.Exit(1)
		}

		if endPort == "" {
			newRuleConfig.EndPort = startPort
		} else {
			newRuleConfig.EndPort = endPort
		}

		rule, err := client.NewFirewallRule(newRuleConfig)
		if err != nil {
			utility.Error("%s", err)
			os.Exit(1)
		}

		ow := utility.NewOutputWriterWithMap(map[string]string{"id": rule.ID, "name": rule.Label})

		switch outputFormat {
		case "json":
			ow.WriteSingleObjectJSON(prettySet)
		case "custom":
			ow.WriteCustomOutput(outputFields)
		default:
			if rule.Label == "" {
				if newRuleConfig.EndPort == newRuleConfig.StartPort {
					fmt.Printf("Created a firewall rule allowing access to port %s from %s with ID %s\n", utility.Green(newRuleConfig.StartPort), utility.Green(strings.Join(newRuleConfig.Cidr, ", ")), rule.ID)
				} else {
					fmt.Printf("Created a firewall rule allowing access to ports %s-%s from %s with ID %s\n", utility.Green(newRuleConfig.StartPort), utility.Green(newRuleConfig.EndPort), utility.Green(strings.Join(newRuleConfig.Cidr, ", ")), rule.ID)
				}
			} else {
				if newRuleConfig.EndPort == newRuleConfig.StartPort {
					fmt.Printf("Created a firewall rule called %s allowing access to port %s from %s with ID %s\n", utility.Green(rule.Label), utility.Green(newRuleConfig.StartPort), utility.Green(strings.Join(newRuleConfig.Cidr, ", ")), rule.ID)
				} else {
					fmt.Printf("Created a firewall rule called %s allowing access to ports %s-%s from %s with ID %s\n", utility.Green(rule.Label), utility.Green(newRuleConfig.StartPort), utility.Green(newRuleConfig.EndPort), utility.Green(strings.Join(newRuleConfig.Cidr, ", ")), rule.ID)
				}
			}
		}
	},
}
