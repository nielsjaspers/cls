package arguments

import (
	"fmt"
	"path/filepath"
	"strings"

	"github.com/nielsjaspers/cls/pkg"
	"github.com/spf13/cobra"
)

type FileData struct {
    Content []byte
    Extension [15]byte  // Maximum of 15 characters for file extension
    Filename [255]byte  // Maximum of 255 characters for filename
}

// ExecuteCommand runs the root command and returns any file content processed.
func ExecuteCommand() (FileData, error) {
    var fileData FileData

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
			// Read the file content and store it in fileData.Content
			content, err := filehandler.FileUpload(args[0])
			if err != nil {
				return fmt.Errorf("failed to read file: %w", err)
			}
            fileData.Content = content
			fmt.Printf("File content read: %d bytes\n", len(fileData.Content))
            
            // Read file name and store it in fileData.Filename
            fName := filepath.Base(args[0])
            var fNameBytes [255]byte
            copy(fNameBytes[:], []byte(fName))
            fileData.Filename = fNameBytes

            // Read file extension and store it in fileData.Extension
            fExt := filepath.Ext(fName)
            var fExtBytes [15]byte
            copy(fExtBytes[:], []byte(fExt))
            fileData.Extension = fExtBytes

			return nil
		},
	}

	rootCmd.AddCommand(shareCmd)

	if err := rootCmd.Execute(); err != nil {
        panic(err)
		// return nil, err
	}

	return fileData, nil
}

