package arguments

import (
	"fmt"

	"github.com/nielsjaspers/cls/pkg"
	"github.com/spf13/cobra"
)

// ExecuteCommand runs the root command and returns any file content processed.
func ExecuteCommand() ([]byte, error) {
	var fileContent []byte

	rootCmd := &cobra.Command{
		Use:   "cls",
		Short: "Command-Line file sharing",
	}

	shareCmd := &cobra.Command{
		Use:     "share <file-path>",
		Short:   "Share a file to a remote location",
		Long:    "Provide the file path to share the file with the server.",
		Example: "cls share <path/to/file>",
		Args:    cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			// Read the file content and store it in fileContent
			data, err := filehandler.FileUpload(args[0])
			if err != nil {
				return fmt.Errorf("failed to read file: %w", err)
			}
			fileContent = data
			fmt.Printf("File content read: %d bytes\n", len(fileContent))
			return nil
		},
	}

	rootCmd.AddCommand(shareCmd)

	if err := rootCmd.Execute(); err != nil {
		return nil, err
	}

	return fileContent, nil
}

