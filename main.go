package main

import (
	"context"
	"flag"
	"fmt"
	"os"
	"strings"

	dialogflow "cloud.google.com/go/dialogflow/apiv2"
	"cloud.google.com/go/dialogflow/apiv2/dialogflowpb"
	"google.golang.org/api/iterator"
	"google.golang.org/api/option"
)

func main() {
	fmt.Println(`
	8""8""8      8""""8                 8   8  8                   
	8  8  8 eeee 8    8 e   e eeee eeee 8   8  8 eeeee eeeee  eeee 
	8e 8  8 8  8 8eeee8 8   8 8    8    8e  8  8 8   8 8   8  8    
	88 8  8 8e   88     8eee8 8eee 8eee 88  8  8 8eee8 8eee8e 8eee 
	88 8  8 88   88     88  8 88   88   88  8  8 88  8 88   8 88   
	88 8  8 88e8 88     88  8 88ee 88ee 88ee8ee8 88  8 88   8 88ee 
																																 `)

	if len(os.Args) == 1 {
		helpScreen()
		return
	}

	version := "3.1"
	textPtr := flag.String("projectId", "", "Google Cloud Project ID")
	langPtr := flag.String("lang", "", "Dialogflow Language Code. eg: en-AU")
	namePtr := flag.String("name", "", "Genesys Cloud Bot Flow Name")
	keyPathPtr := flag.String("keyPath", "", "Google Cloud Key Path eg: /path/to/key.json")
	flag.Bool("version", false, "Version of McPheeWare CLI")
	flag.Bool("help", false, "Help Screen")

	flag.Parse()

	switch os.Args[1] {
	case "-projectId":
		if namePtr == nil || langPtr == nil || textPtr == nil || keyPathPtr == nil {
			fmt.Println("Missing required parameters type -help for more information")
			return
		}
		if *namePtr == "botFlow" || *namePtr == "digitalBotFlow" {
			fmt.Println("Invalid Bot Flow Name")
			return
		}
		buildDigitalBot(*textPtr, *langPtr, *namePtr, *keyPathPtr)
	case "-version":
		fmt.Println("Version: ", version)
	default:
		helpScreen()
		return
	}
}

func helpScreen() {
	fmt.Println(`
The below parameters are required to build a Genesys Cloud Bot Flow:

-projectId: Google Cloud Project ID
-lang: Dialogflow Language Code. eg: en-AU
-name: Genesys Cloud Bot Flow Name
-keyPath: Google Cloud Key Path eg: /path/to/key.json

To check the version of the CLI, use the -version flag
																																 `)
}

func buildDigitalBot(projectId, lang, flowName, keyPath string) {
	intents, err := ListIntents(projectId, lang, keyPath)
	if err != nil {
		fmt.Println("Error: ", err)
		return
	}

	fmt.Println("Received intents: ", len(intents))

	var allVariables = ""
	var allTasks = ""
	var allIntents = ""
	var entityNameReferences = ""
	var allUtterances = ""
	var allEntities = ""
	var allEntityTypes = ""
	var allSlots = []string{}

	for _, intent := range intents {
		var displayName = intent.DisplayName

		// clean up display name with invalid characters
		if strings.Contains(displayName, " ") || strings.Contains(displayName, "-") || strings.Contains(displayName, ".") || strings.Contains(displayName, "@") {
			displayName = strings.ReplaceAll(displayName, " ", "_")
			displayName = strings.ReplaceAll(displayName, "-", "_")
			displayName = strings.ReplaceAll(displayName, ".", "_")
			displayName = strings.ReplaceAll(displayName, "@", "_")
		}
		if !strings.Contains(displayName, "Knowledge.KnowledgeBase") {
			createTask := createTask(displayName)
			allTasks += fmt.Sprintf("\n%s", createTask)
			createIntent := createIntent(displayName)
			allIntents += fmt.Sprintf("\n%s", createIntent)
			var allSegments = ""
			if len(intent.TrainingPhrases) > 0 {
				entityNameReferences = ""
				for _, trainingPhrase := range intent.TrainingPhrases {
					var segment = "            - segments:\n"
					fmt.Println("Training Phrase: ", trainingPhrase.Parts)
					for _, part := range trainingPhrase.Parts {
						// escape quotes if they are in the text
						if strings.Contains(part.Text, "\"") {
							part.Text = strings.ReplaceAll(part.Text, "\"", "\\\"")
						}
						// add text to segment
						segment += fmt.Sprintf("                - text: \"%s\"\n", part.Text)
						// if text contains entity, add entity to segment
						if part.EntityType != "" {
							// clean up entity name with invalid characters
							if strings.Contains(part.EntityType, "@") || strings.Contains(part.EntityType, ".") || strings.Contains(part.EntityType, " ") || strings.Contains(part.EntityType, "-") {
								part.EntityType = strings.ReplaceAll(part.EntityType, "@", "")
								part.EntityType = strings.ReplaceAll(part.EntityType, ".", "_")
								part.EntityType = strings.ReplaceAll(part.EntityType, " ", "_")
								part.EntityType = strings.ReplaceAll(part.EntityType, "-", "_")
							}
							segment += fmt.Sprintf(`                  entity:
                    name: %s`+"\n", part.EntityType)
							// add entity to entities if not already there as well as variables
							if !contains(allSlots, part.EntityType) {
								allSlots = append(allSlots, part.EntityType)
								allVariables += fmt.Sprintf("\n%s", createVariable(part.EntityType))
								allEntities += fmt.Sprintf("\n%s", createEntity(part.EntityType))
								allEntityTypes += fmt.Sprintf("\n%s", createEntityType(part.EntityType))
								entityNameReferences += fmt.Sprintf("\n            - %s", part.EntityType)
							}
						}
					}
					allSegments += fmt.Sprintf("\n%s", segment)
				}

			} else {
				createSegment := createSegment("No utterance")
				allSegments += fmt.Sprintf("\n%s", createSegment)
			}
			if entityNameReferences == "" {
				entityNameReferences = " []"
			}
			createUtterances := createUtterances(allSegments, entityNameReferences, displayName)
			allUtterances += fmt.Sprintf("\n%s", createUtterances)
			fmt.Println("Completed Intent: ", displayName)
		}
	}

	for _, slot := range allSlots {
		fmt.Println("Slot created: ", slot)
	}

	// Get all entities and blank out if none exist with []
	if allVariables == "" {
		allVariables = " []"
	}
	if allTasks == "" {
		allTasks = " []"
	}
	if allIntents == "" {
		allIntents = " []"
	}
	if allUtterances == "" {
		allUtterances = " []"
	}
	if allEntities == "" {
		allEntities = " []"
	}
	if allEntityTypes == "" {
		allEntityTypes = " []"
	}

	createYaml := createYaml(flowName, allVariables, allTasks, allIntents, allUtterances, allEntities, allEntityTypes)
	os.WriteFile(fmt.Sprintf("%s.yaml", flowName), []byte(createYaml), 0777)
	fmt.Println("Flow created: ", flowName)
}

func contains(slice []string, item string) bool {
	for _, s := range slice {
		if s == item {
			return true
		}
	}
	return false
}

func ListIntents(projectID, lang, keyPath string) ([]*dialogflowpb.Intent, error) {
	ctx := context.Background()

	intentsClient, clientErr := dialogflow.NewIntentsClient(ctx, option.WithCredentialsFile(keyPath))
	if clientErr != nil {
		return nil, clientErr
	}
	defer intentsClient.Close()

	if projectID == "" {
		return nil, fmt.Errorf("received empty project (%s)", projectID)
	}

	parent := fmt.Sprintf("projects/%s/agent", projectID)

	request := dialogflowpb.ListIntentsRequest{Parent: parent, IntentView: dialogflowpb.IntentView_INTENT_VIEW_FULL, LanguageCode: lang}

	intentIterator := intentsClient.ListIntents(ctx, &request)
	var intents []*dialogflowpb.Intent

	for intent, status := intentIterator.Next(); status != iterator.Done; {
		if len(intents) > 1000 {
			fmt.Println("Error: Does your api key have API admin access??")
			os.Exit(1)
		}
		intents = append(intents, intent)
		intent, status = intentIterator.Next()
	}

	return intents, nil
}

func createVariable(name string) string {
	// Create a new variable
	var variable = fmt.Sprintf(`    - stringVariable:
        name: Slot.%s
        initialValue:
          noValue: true
        isInput: true
        isOutput: true`, name)
	return variable
}

func createTask(name string) string {
	// Create a new task
	var task = fmt.Sprintf(`    - task:
        name: %s
        actions:
          - exitBotFlow:
              name: Exit Bot Flow`, name)
	return task
}

func createIntent(name string) string {
	// Create a new intent
	var intent = fmt.Sprintf(`      - intent:
          confirmation:
            exp: "MakeCommunication(\n  ToCommunication(ToCommunication(\"I think you want to %s, is that correct?\")))"
          name: %s
          task:
            name: %s`, name, name, name)
	return intent
}

func createSegment(text string) string {
	var segment = fmt.Sprintf(`            - segments:
                - text: %s`, text)
	return segment
}

func createEntity(name string) string {
	var entity = fmt.Sprintf(`        - name: %s
          type: %sType`, name, name)
	return entity
}

func createEntityType(name string) string {
	var entityType = fmt.Sprintf(`        - name: %sType
          description: The description of the Entity Type "%sType"
          mechanism:
            type: Regex
            restricted: true
            items: []`, name, name)
	return entityType
}

func createUtterances(segments, entityNameReferences, name string) string {
	// Create a new utterance
	var utterance = fmt.Sprintf(`        - utterances:%s
              source: User
          entityNameReferences:%s
          name: %s`, segments, entityNameReferences, name)
	return utterance
}

func createYaml(flowName, variables, task, intent, utterance, entities, entityTypes string) string {
	// Create the yaml file
	var yaml = fmt.Sprintf(`digitalBot:
  name: %s
  defaultLanguage: en-au
  startUpRef: "/digitalBot/bots/bot[Initial Greeting_10]"
  bots:
    - bot:
        name: Initial Greeting
        refId: Initial Greeting_10
        actions:
          - waitForInput:
              name: Wait for Input
              question:
                exp: "MakeCommunication(\n  ToCommunication(ToCommunication(\"What would you like to do?\")))"
              knowledgeSearchResult:
                noValue: true
              noMatch:
                exp: "MakeCommunication(\n  ToCommunication(ToCommunication(\"Tell me again what you would like to do.\")))"
  variables:%s
  tasks:%s
  settingsBotFlow:
    intentSettings:%s
  settingsNaturalLanguageUnderstanding:
    nluDomainVersion:
      intents:%s
      entities:%s
      entityTypes:%s
      language: en-au
      languageVersions: {}
    mutedUtterances: []
`, flowName, variables, task, intent, utterance, entities, entityTypes)
	return yaml
}
