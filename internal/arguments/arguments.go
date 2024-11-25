package arguments

import (
	"fmt"

	"github.com/spf13/cobra"
)

type Args struct {
	Command   *cobra.Command
	FileBytes []byte
	FilePath  string
}

func InitArgs() *Args {
    args := &Args{
        Command: rootCmd(),
    }
    return args
}

func rootCmd() *cobra.Command {
	rootCmd := &cobra.Command{
		Use:   "cls",
		Short: "Command-Line file sharing",
	}

	rootCmd.AddCommand(testCmd())

	return rootCmd
}

func testCmd() *cobra.Command {
    var testCmd = &cobra.Command{
        Use: "test",
        Short: "Short test message",
        Args: cobra.ExactArgs(1),
        Run: func(cmd *cobra.Command, args []string) {
            var str = args[0]
            fmt.Printf("Your argument: %v", str)
        },
    }
    return testCmd
    
}
