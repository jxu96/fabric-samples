/*
Copyright © 2022 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"fmt"

	"github.com/hyperledger/fabric-gateway/pkg/client"
	"github.com/jxu96/fabric-samples/chaincode/test-dataset-metadata-ledger/gateway"
	"github.com/spf13/cobra"
)

// queryCmd represents the query command
var queryCmd = &cobra.Command{
	Use:   "query",
	Short: "A brief description of your command",
	Long: `A longer description that spans multiple lines and likely contains examples
and usage of using your command. For example:

Cobra is a CLI library for Go that empowers applications.
This application is a tool to generate the needed files
to quickly create a Cobra application.`,
	Run: func(cmd *cobra.Command, args []string) {
		privateMode, err := cmd.Flags().GetBool("private")

		coll, err := cmd.Flags().GetString("collection")
		cobra.CheckErr(err)

		gatewayConfig := getGatewayConfig()
		gw := gateway.NewFabricGateway()
		cobra.CheckErr(gw.WithConfiguration(gatewayConfig))
		defer gw.Client.Close()
		defer gw.Gateway.Close()

		for _, key := range args {
			if privateMode {
				result, err := gw.GetContract("basic", "mychannel").Evaluate(
					"QueryPrivate",
					client.WithArguments(key),
				)
				cobra.CheckErr(err)
				fmt.Printf("Result: %s\n", string(result))
			} else {
				result, err := gw.GetContract("basic", "mychannel").Evaluate(
					"Query",
					client.WithArguments(coll, key),
				)
				cobra.CheckErr(err)
				fmt.Printf("Result: %s\n", string(result))
			}
		}
	},
}

func init() {
	rootCmd.AddCommand(queryCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// queryCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	queryCmd.Flags().StringP("collection", "c", "", "collection in which to query data")
	queryCmd.Flags().BoolP("private", "p", false, "query private collection")
}
