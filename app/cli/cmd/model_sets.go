package cmd

import (
	"fmt"
	"plandex/api"
	"plandex/term"
	"github.com/spf13/cobra"
)

var modelSetsCmd = &cobra.Command{
	Use:   "model-sets",
	Short: "Manage model sets",
}

var listModelSetsCmd = &cobra.Command{
	Use:   "list",
	Short: "List all model sets",
	Run:   listModelSets,
}

var createModelSetCmd = &cobra.Command{
	Use:   "create",
	Short: "Create a model set",
	Run:   createModelSet,
}

func init() {
	RootCmd.AddCommand(modelSetsCmd)
	modelSetsCmd.AddCommand(listModelSetsCmd)
	modelSetsCmd.AddCommand(createModelSetCmd)
	modelSetsCmd.AddCommand(deleteModelSetCmd)
}

var deleteModelSetCmd = &cobra.Command{
	Use:   "delete [name-or-index]",
	Short: "Delete a model set by name or index",
	Args:  cobra.MaximumNArgs(1),
	Run:   deleteModelSet,
}

func deleteModelSet(cmd *cobra.Command, args []string) {
	var modelSets []shared.ModelSet
	var err error

	term.StartSpinner("Fetching model sets...")
	modelSets, err = api.Client.ListModelSets()
	term.StopSpinner()

	if err != nil {
		term.OutputErrorAndExit("Error fetching model sets: %v", err)
		return
	}

	if len(modelSets) == 0 {
		fmt.Println("No model sets available to delete.")
		return
	}

	var setToDelete *shared.ModelSet

	if len(args) == 1 {
		input := args[0]
		// Try to parse input as index
		index, err := strconv.Atoi(input)
		if err == nil && index > 0 && index <= len(modelSets) {
			setToDelete = &modelSets[index-1]
		} else {
			// Search by name
			for _, s := range modelSets {
				if s.Name == input {
					setToDelete = &s
					break
				}
			}
		}
	}

	if setToDelete == nil {
		fmt.Println("Select a model set to delete:")
		for i, set := range modelSets {
			fmt.Printf("%d: %s\n", i+1, set.Name)
		}
		var selectedIndex int
		fmt.Scanln(&selectedIndex)
		if selectedIndex < 1 || selectedIndex > len(modelSets) {
			fmt.Println("Invalid selection.")
			return
		}
		setToDelete = &modelSets[selectedIndex-1]
	}

	term.StartSpinner(fmt.Sprintf("Deleting model set '%s'...", setToDelete.Name))
	err = api.Client.DeleteModelSet(setToDelete.Id)
	term.StopSpinner()

	if err != nil {
		term.OutputErrorAndExit("Error deleting model set: %v", err)
		return
	}

	fmt.Printf("Model set '%s' deleted successfully.\n", setToDelete.Name)
}

func listModelSets(cmd *cobra.Command, args []string) {
	term.StartSpinner("Fetching model sets...")
	modelSets, err := api.Client.ListModelSets()
	term.StopSpinner()

	if err != nil {
		term.OutputErrorAndExit("Error fetching model sets: %v", err)
		return
	}

	fmt.Println("Available Model Sets:")
	for _, set := range modelSets {
		fmt.Printf("- %s: %s\n", set.Name, set.Description)
	}
}

func createModelSet(cmd *cobra.Command, args []string) {
	set := &shared.ModelSet{}

	name, err := term.GetUserStringInput("Enter model set name:")
	if err != nil {
		term.OutputErrorAndExit("Error reading model set name: %v", err)
		return
	}
	set.Name = name

	description, err := term.GetUserStringInput("Enter description:")
	if err != nil {
		term.OutputErrorAndExit("Error reading description: %v", err)
		return
	}
	set.Description = description

	// Additional details would be prompted here similarly

	term.StartSpinner("Creating model set...")
	err = api.Client.CreateModelSet(set)
	term.StopSpinner()

	if err != nil {
		term.OutputErrorAndExit("Error creating model set: %v", err)
		return
	}

	fmt.Println("Model set created successfully.")
}


