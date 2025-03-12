# googlebot-migrate-genesysbot

This is a CLI tool that helps in migrating Google Bots to Genesys Cloud Bots

Incase you don't have GO installed and don't want to build this from source I have done a `build` of this CLI tool for each operating system Linux, MacOS & Windows. These files can be found in each of the folders in this repo.

To enable these to be ran as a system level command you will need to put them in your system path in my case as im on Linux I put mine in

```
/usr/local/bin
```

Depending on your operating system you may have a different directory. If you don't want to install it into the system directory you can also just run it from teh directory you save it to. But in this case you will need to be in that directory to run the program.

If you run it without any parameters you will be taken to the `help screen`

```

	8""8""8      8""""8                 8   8  8
	8  8  8 eeee 8    8 e   e eeee eeee 8   8  8 eeeee eeeee  eeee
	8e 8  8 8  8 8eeee8 8   8 8    8    8e  8  8 8   8 8   8  8
	88 8  8 8e   88     8eee8 8eee 8eee 88  8  8 8eee8 8eee8e 8eee
	88 8  8 88   88     88  8 88   88   88  8  8 88  8 88   8 88
	88 8  8 88e8 88     88  8 88ee 88ee 88ee8ee8 88  8 88   8 88ee


The below parameters are required to build a Genesys Cloud Bot Flow or Knowledge Base CSV:

-type: digitalBot or knowledgeBase
-projectId: Google Cloud Project ID
-lang: Dialogflow Language Code. eg: en-AU
-name: Genesys Cloud Bot Flow Name or CSV file name
-keyPath: Google Cloud Key Path eg: /path/to/key.json

To check the version of the CLI, use the -version flag

```

## key.json

You will need to create a service account `json` key to the Google project that holds the Dialogflow ES agent. This is what is then used to authenticate with the apis to download the information. In my example i have it downloaded into my Downloads folder and its called "key.json" yours may well have different name and location.

As the for the permissions required it needs to have `Dialogflow API Admin.`

![](/docs/images/role.png?raw=true)

Ensure you then create a json key for the service account

![](/docs/images/json-key.png?raw=true)

## using the CLI

To run the program you will need to supply the parameters like the below:

```
mcpheeware-migrate -type digitalBot -projectId my-projectId -lang en-AU -name test -keyPath ~/Downloads/key.json
```

This will then output a `yaml` file that can then be used by [Archy](https://developer.genesys.cloud/devapps/archy/) to either create a `architect` digitalBotFlow input file type or directly PUSH the flow to Genesys Cloud

## Importing the flow

If you want to create a file that you can use the GUI to import then you can run:

```
archy createImportFile --file fileName.yaml
```

and it will create an import file you can use in architect.

If you already have `Archy` connected to your own ORG then you can simply `create` the flow directly. NOTE: I recommend to create not publish as there will be validation errors that will need to be addressed most likely in the flow.

```
archy create --file fileName.yaml
```

Once you have the flow imported into Genesys Cloud you can then edit it from there. While the Intents, Utterances, Slots, Data & Reusable Tasks get created you will still need to build out the actually use case of what needs to be "done" in each Reusable Task.

## Result

So this is an example of a covid-19 BOT that I had in my google dialogueflow ES

![](/docs/images/google.png?raw=true)

Once I run the program to convert it to a Genesys Digital Bot Flow it passes in all the intents

![](/docs/images/intents.png?raw=true)

And even creates and puts the `Slots` or as google calls them `entities`

![](/docs/images/utterances.png?raw=true)

To keep it simple I have made each slot type be `RegEx` so you will need to then apply the slot type you want as well as the RegEx expression if you want it to stay as that type. You will also need to build out the outcome in the reusable task on what actually needs to happen in the intent.

## Knowledge CSV Import feature

If you have used the new CSV file feature in version 4+ of this CLI you can then use it to import directly into the Genesys Cloud KnowledgeBase as articles using the Web GUI import option. `NOTE: right now I have made this support 100 training phrases` in the article but this can be expanded if required in the CSV file creation columns.

![](/docs/images/upload_csv.png?raw=true)

Then once uploaded you will be able to vew them from within the articles as well as see the training phrases setup as Phrasings in teh article itself

![](/docs/images/knowledgeBase.png?raw=true)
