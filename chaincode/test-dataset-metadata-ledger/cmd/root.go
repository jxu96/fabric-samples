/*
Copyright © 2022 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"os"
	"strings"

	"github.com/jxu96/fabric-samples/chaincode/test-dataset-metadata-ledger/gateway"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:   "test-dataset-metadata-ledger",
	Short: "A brief description of your application",
	Long: `A longer description that spans multiple lines and likely contains
examples and usage of using your application. For example:

Cobra is a CLI library for Go that empowers applications.
This application is a tool to generate the needed files
to quickly create a Cobra application.`,
	// Uncomment the following line if your bare application
	// has an action associated with it:
	// Run: func(cmd *cobra.Command, args []string) { },
}

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() {
	err := rootCmd.Execute()
	if err != nil {
		os.Exit(1)
	}
}

// mspID        = "Org1MSP"
//
//	cryptoPath   = "../../test-network/organizations/peerOrganizations/org1.example.com"
//	certPath     = cryptoPath + "/users/User1@org1.example.com/msp/signcerts/cert.pem"
//	keyPath      = cryptoPath + "/users/User1@org1.example.com/msp/keystore/"
//	tlsCertPath  = cryptoPath + "/peers/peer0.org1.example.com/tls/ca.crt"
//	peerEndpoint = "localhost:7051"
//	gatewayPeer  = "peer0.org1.example.com"
func init() {
	cobra.OnInitialize(initConfig)
	// Here you will define your flags and configuration settings.
	// Cobra supports persistent flags, which, if defined here,
	// will be global for your application.

	rootCmd.PersistentFlags().String("as", "org1.user1", "as user (default is org1.user1)")
	rootCmd.PersistentFlags().String("peer", "org1.peer0", "(default is org1.peer0)")

	// Cobra also supports local flags, which will only run
	// when this action is called directly.
	// rootCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")

	wd, err := os.Getwd()
	cobra.CheckErr(err)

	cryptoPathOrg1 := wd + "/organizations/peerOrganizations/org1.example.com"

	viper.SetDefault("org1.mspID", "Org1MSP")
	viper.SetDefault("org1.cryptoPath", cryptoPathOrg1)
	viper.SetDefault("org1.user1.certPath", cryptoPathOrg1+"/users/User1@org1.example.com/msp/signcerts/cert.pem")
	viper.SetDefault("org1.user1.keyPath", cryptoPathOrg1+"/users/User1@org1.example.com/msp/keystore/")
	viper.SetDefault("org1.peer0.tlsCertPath", cryptoPathOrg1+"/peers/peer0.org1.example.com/tls/ca.crt")
	viper.SetDefault("org1.peer0.endpoint", "localhost:7051")
	viper.SetDefault("org1.peer0.gateway", "peer0.org1.example.com")

	cryptoPathOrg2 := wd + "/organizations/peerOrganizations/org2.example.com"

	viper.SetDefault("org2.mspID", "Org2MSP")
	viper.SetDefault("org2.cryptoPath", cryptoPathOrg2)
	viper.SetDefault("org2.user1.certPath", cryptoPathOrg2+"/users/User1@org2.example.com/msp/signcerts/cert.pem")
	viper.SetDefault("org2.user1.keyPath", cryptoPathOrg2+"/users/User1@org2.example.com/msp/keystore/")
	viper.SetDefault("org2.peer0.tlsCertPath", cryptoPathOrg2+"/peers/peer0.org2.example.com/tls/ca.crt")
	viper.SetDefault("org2.peer0.endpoint", "localhost:9051")
	viper.SetDefault("org2.peer0.gateway", "peer0.org2.example.com")
}

// initConfig reads in config file and ENV variables if set.
func initConfig() {
	if cfgPath := os.Getenv("CONFIG_PATH"); cfgPath != "" {
		viper.AddConfigPath(cfgPath)
	} else {
		home, err := os.UserHomeDir()
		cobra.CheckErr(err)
		viper.AddConfigPath(home)
	}

	viper.SetConfigName(".config-fabric")
	viper.SetConfigType("yaml")

	viper.AutomaticEnv() // read in environment variables that match
	viper.ReadInConfig()
}

func getGatewayConfig() gateway.FabricGatewayConfiguration {
	user, err := rootCmd.PersistentFlags().GetString("as")
	cobra.CheckErr(err)
	peer, err := rootCmd.PersistentFlags().GetString("peer")
	cobra.CheckErr(err)
	org := strings.Split(user, ".")[0]

	return gateway.FabricGatewayConfiguration{
		MspID:        viper.GetString(org + ".mspID"),
		CertPath:     viper.GetString(user + ".certPath"),
		KeyPath:      viper.GetString(user + ".keyPath"),
		TlsCertPath:  viper.GetString(peer + ".tlsCertPath"),
		PeerEndpoint: viper.GetString(peer + ".endpoint"),
		PeerGateway:  viper.GetString(peer + ".gateway"),
	}
}
