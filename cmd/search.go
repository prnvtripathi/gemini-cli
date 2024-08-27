package cmd

import (
	"bufio"   // for reading input
	"context" // for context.Background()
	"fmt"     // for Println()
	"log"     // for Fatal()
	"os"      // for Getenv()
	"strings" // for Join()

	"github.com/google/generative-ai-go/genai" // for NewClient()
	"github.com/spf13/cobra"                   // for Command(), MinimumNArgs(), AddCommand()
	"google.golang.org/api/option"             // for WithAPIKey()
)

var searchCmd = &cobra.Command{ // searchCmd is a pointer to a Command struct
	Use:   "search",                                                 // Use is a string that represents the command name
	Short: "A command to search for a query using the Gemini model", // Short is a string that represents the command description
	Args:  cobra.MinimumNArgs(0),                                    // Args is a function that sets the minimum number of arguments required for the command
	Run: func(cmd *cobra.Command, args []string) {
		checkAndSetAPIKey() // Check if the API key is set
		startChat()         // Start the chat loop
	},
	// Run is a function that is executed when the command is called
}

func init() {
	rootCmd.AddCommand(searchCmd) // AddCommand adds a child command to the parent command
}

// checkAndSetAPIKey checks if GEMINI_API_KEY is set, and executes the appropriate script if it's not
func checkAndSetAPIKey() {
	apiKey := os.Getenv("GEMINI_API_KEY")
	if apiKey == "" {
		log.Fatal("GEMINI_API_KEY is not set. Please set the API key using `export GEMINI_API_KEY=<your-api-key>`") // If the API key is not set, log the error
	}
}

// startChat handles the chat interaction with the user
func startChat() {
	ctx := context.Background()                                                         // Creating a context
	client, err := genai.NewClient(ctx, option.WithAPIKey(os.Getenv("GEMINI_API_KEY"))) // Creating a new client
	if err != nil {
		log.Fatal(err) // If there is an error, log the error
	}
	defer client.Close() // Close the client after the function ends

	model := client.GenerativeModel("gemini-1.5-flash") // Creating a generative model

	scanner := bufio.NewScanner(os.Stdin)
	fmt.Println("Start chatting! Type 'exit' or 'quit' to end the chat.")

	for {
		fmt.Print("you: ")
		scanner.Scan()
		userInput := scanner.Text()

		if strings.ToLower(userInput) == "exit" || strings.ToLower(userInput) == "quit" { // If the user types 'exit' or 'quit', break the loop
			fmt.Println("gemini: Goodbye!")
			break
		}

		resp, err := model.GenerateContent(ctx, genai.Text(userInput)) // Generating content
		if err != nil {
			log.Fatal(err)
		}

		fmt.Printf("gemini: %s\n", (resp.Candidates[0].Content.Parts[0])) // Printing the generated content

	}
}
