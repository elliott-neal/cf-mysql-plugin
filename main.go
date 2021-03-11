package main

import (
	"code.cloudfoundry.org/cli/plugin"
	"fmt"
	"github.com/elliott-neal/cf-mysql-plugin/cfmysql"
	"os"
)

func main() {
	if len(os.Args) == 1 {
		fmt.Fprintf(os.Stderr, "This executable is a cf plugin. "+
			"Run `cf install-plugin %s` to install it\nand `cf mysql service-name` "+
			"to use it.\n",
			os.Args[0])
		os.Exit(1)
	}

	mysqlPlugin := newPlugin()
	plugin.Start(mysqlPlugin)

	os.Exit(mysqlPlugin.GetExitCode())
}

func newPlugin() *cfmysql.MysqlPlugin {
	httpClientFactory := cfmysql.NewHttpClientFactory()
	osWrapper := cfmysql.NewOsWrapper()
	requestDumper := cfmysql.NewRequestDumper(osWrapper, os.Stderr)
	http := cfmysql.NewHttpWrapper(httpClientFactory, requestDumper)
	apiClient := cfmysql.NewApiClient(http)

	sshRunner := cfmysql.NewSshRunner()
	netWrapper := cfmysql.NewNetWrapper()
	waiter := cfmysql.NewPortWaiter(netWrapper)
	randWrapper := cfmysql.NewRandWrapper()
	cfService := cfmysql.NewCfService(apiClient, sshRunner, waiter, http, randWrapper, os.Stderr)

	execWrapper := cfmysql.NewExecWrapper()
	ioUtilWrapper := cfmysql.NewIoUtilWrapper()
	runner := cfmysql.NewMysqlRunner(execWrapper, ioUtilWrapper, osWrapper)

	portFinder := cfmysql.NewPortFinder()

	return cfmysql.NewMysqlPlugin(cfmysql.PluginConf{
		In:          os.Stdin,
		Out:         os.Stdout,
		Err:         os.Stderr,
		CfService:   cfService,
		PortFinder:  portFinder,
		MysqlRunner: runner,
	})
}
