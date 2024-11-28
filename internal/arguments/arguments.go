package arguments

import (
	"fmt"
	"github.com/nielsjaspers/cls/pkg"
	"github.com/spf13/cobra"
	"path/filepath"
)

type FileData struct {
	Content   []byte
	Extension [15]byte  // Maximum of 15 characters for file extension
	Filename  [255]byte // Maximum of 255 characters for filename
}

// ExecuteCommand runs the root command and returns any file content processed.
func ExecuteCommand() (FileData, []string, error) {
	var fileData FileData

	// Contains the data used by the commands
	// list "" ""
	// get <remote> <local>
	var commandData [3]string

	rootCmd := &cobra.Command{
		Use:   "cls",
		Short: "Command-Line file sharing",
	}

	shareCmd := &cobra.Command{
		Use:     "share <file-path>",
		Aliases: []string{"-s", "s"},
		Short:   "Share a file to a remote location",
		Long:    "Provide the file path to share the file with the server.",
		Example: "cls share <path/to/file>",
		Args:    cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			commandData[0] = "SHARE_FILE_SHARE_FILE" // Share file marker for server

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

	listAllCmd := &cobra.Command{
		Use:     "list",
		Aliases: []string{"-l", "-ls", "l", "ls"},
		Short:   "List all files on the remote location",
		Long:    "Request a list from the server with all files currently on storage",
		Example: "cls list",
		Args:    cobra.ExactArgs(0),
		Run: func(cmd *cobra.Command, args []string) {
			commandData[0] = "LIST_ALL_LIST_ALL" // List all files marker for server
		},
	}

	getFileCmd := &cobra.Command{
		Use:     "get <remote/file> <local/path>",
		Aliases: []string{"-g", "g"},
		Short:   "Request a single file from the remote location",
		Long:    "Request a single file from the server and store it to a given local path",
		Example: "cls get <remote/file> <local/path>",
		Args:    cobra.ExactArgs(2),
		Run: func(cmd *cobra.Command, args []string) {
			commandData[0] = "GET_FILE_GET_FILE" // Get file marker for server
			commandData[1] = args[0]             // Remote path to file
			commandData[2] = args[1]             // Path to local destination
		},
	}

	rootCmd.AddCommand(shareCmd)
	rootCmd.AddCommand(listAllCmd)
	rootCmd.AddCommand(getFileCmd)

	if err := rootCmd.Execute(); err != nil {
		panic(err)
	}

	return fileData, commandData[:], nil
}
