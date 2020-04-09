package cmd

import (
	"fmt"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"log"
	"os"
)

var cfgFile string

func init() {
	cobra.OnInitialize(initConfig)
	rootCmd.PersistentFlags().StringVarP(&cfgFile, "config", "C", "", "configuration file (defaults are \"/etc/dns-tools/dns-tools-config.json\" and \"./dns-tools-config.json\")")
	rootCmd.AddCommand(signCmd)
	rootCmd.AddCommand(verifyCmd)
	rootCmd.AddCommand(digestCmd)
	rootCmd.AddCommand(resetPKCS11KeysCmd)
	Log = log.New(os.Stderr, "[dns-tools]", log.Ldate|log.Ltime)
}

var Log *log.Logger

var rootCmd = &cobra.Command{
	Use:   "dns-tools",
	Short: "Signs a DNS zone using a PKCS11 Device",
	Long: `Allows to sign a DNS zone using a PKCS#11 device.
	
	For more information, visit "https://github.com/niclabs/dns-tools".`,
}

func Execute() {

	if err := rootCmd.Execute(); err != nil {
		Log.Printf("Error: %s", err)
		os.Exit(1)
	}
}

func initConfig() {
	if cfgFile != "" {
		viper.SetConfigFile(cfgFile)
	} else {
		viper.AddConfigPath("/etc/dns-tools/")
		viper.AddConfigPath("./")
		viper.SetConfigName("dns-tools-config")
	}

	viper.AutomaticEnv()
	if err := viper.ReadInConfig(); err == nil {
		Log.Println("Using config file:", viper.ConfigFileUsed())
	}
}

// filesExist returns an error if any of the paths received as args does not point to a readable file.
func filesExist(filepaths ...string) error {
	for _, path := range filepaths {
		_, err := os.Stat(path)
		if err != nil || os.IsNotExist(err) {
			return fmt.Errorf("File %s doesn't exist or it has not reading permissions\n", path)
		}
	}
	return nil
}
