package cmd

import (
	"os"
	"strings"

	"github.com/civo/cli/config"
	"github.com/civo/cli/utility"
	"github.com/spf13/cobra"
)

var loadBalancerShowCmd = &cobra.Command{
	Use:     "show",
	Aliases: []string{"get", "inspect"},
	Example: `civo loadbalancer show ID/NAME -o custom -f "ID: Name"`,
	Args:    cobra.MinimumNArgs(1),
	Short:   "Show load balancer",
	Long: `Show a specified load balancer.
If you wish to use a custom format, the available fields are:

	* id
	* name
	* algorithm
	* pulic_ip
	* state
	* private_ip
	* firewall_id
	* cluter_id
	* external_traffic_policy
	* session_affinity
	* session_affinity_config_timeout`,
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

		lb, err := client.GetLoadBalancer(args[0])
		if err != nil {
			utility.Error("%s", err)
			os.Exit(1)
		}

		ow := utility.NewOutputWriter()
		ow.StartLine()

		ow.AppendDataWithLabel("id", lb.ID, "ID")
		ow.AppendDataWithLabel("name", lb.Name, "Name")
		ow.AppendDataWithLabel("algorithm", lb.Algorithm, "Algorithm")
		ow.AppendDataWithLabel("public_ip", lb.PublicIP, "Public IP")
		ow.AppendDataWithLabel("state", lb.State, "State")

		if outputFormat == "json" || outputFormat == "custom" {
			ow.AppendDataWithLabel("private_ip", lb.PrivateIP, "Private IP")
			ow.AppendDataWithLabel("firewall_id", lb.FirewallID, "Firewall ID")
			ow.AppendDataWithLabel("cluster_id", lb.ClusterID, "Cluster ID")
			ow.AppendDataWithLabel("external_traffic_policy", lb.ExternalTrafficPolicy, "External Traffic Policy")
			ow.AppendDataWithLabel("session_affinity", lb.SessionAffinity, "Session Affinity")
			ow.AppendDataWithLabel("session_affinity_config_timeout", string(lb.SessionAffinityConfigTimeout), "Session Affinity ConfigT imeout")
		}

		var backendList []string
		for _, backend := range lb.Backends {
			backendList = append(backendList, backend.IP)
		}

		ow.AppendData("Backends", strings.Join(backendList, ", "))

		switch outputFormat {
		case "json":
			ow.ToJSON(lb, prettySet)
		case "custom":
			ow.WriteCustomOutput(outputFields)
		default:
			ow.WriteTable()
		}
	},
}
