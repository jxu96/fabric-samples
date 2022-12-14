/*
Copyright © 2022 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"encoding/json"
	"os"

	"github.com/hyperledger/fabric-gateway/pkg/client"
	"github.com/jxu96/fabric-samples/chaincode/test-dataset-metadata-ledger/gateway"
	"github.com/spf13/cobra"
)

// registerCmd represents the register command
var registerCmd = &cobra.Command{
	Use:   "register",
	Short: "Register metadata to the ledger",
	Long: `A longer description that spans multiple lines and likely contains examples
and usage of using your command. For example:

Cobra is a CLI library for Go that empowers applications.
This application is a tool to generate the needed files
to quickly create a Cobra application.`,
	Run: func(cmd *cobra.Command, args []string) {
		// read metadata and encode to base64
		mdPath, err := cmd.Flags().GetString("metadata")
		cobra.CheckErr(err)
		md, err := os.ReadFile(mdPath)
		cobra.CheckErr(err)

		// get collections and encode to base64
		collections, err := cmd.Flags().GetStringArray("collection")
		public, err := cmd.Flags().GetBool("public")
		if public {
			collections = append(collections, "")
		}
		cobra.CheckErr(err)
		bs, err := json.Marshal(collections)
		cobra.CheckErr(err)

		transientData := map[string][]byte{
			"metadata":    md,
			"collections": bs,
		}

		gatewayConfig := getGatewayConfig()
		gw := gateway.NewFabricGateway()
		cobra.CheckErr(gw.WithConfiguration(gatewayConfig))
		defer gw.Client.Close()
		defer gw.Gateway.Close()
		_, err = gw.GetContract("basic", "mychannel").Submit(
			"Register",
			client.WithTransient(transientData),
		)
		cobra.CheckErr(err)
	},
}

func init() {
	rootCmd.AddCommand(registerCmd)

	registerCmd.Flags().String("metadata", "", "path to metadata file")
	registerCmd.Flags().StringArray("collection", []string{}, "collections in which to register the metadata")
	registerCmd.Flags().BoolP("public", "p", true, "register in public ledger")
}
