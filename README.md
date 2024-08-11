<h1 align="center"> hermes-crypto</h1>
<p align="center"><img alt="hermes-crypto" src="./assets/hermes-crypto-logo.svg" width="200"></p>

Hermes (called Mercury in Roman mythology) was considered the messenger of the Olympic gods. He possesses the ability to influence outcomes and tip the scales in favor of those who seek his favor. As the god of luck, he brings both fortune and misfortune to those who dare to test their luck.

`hermes-crypto` is a fun page where you can ponder if the price of your coin will go up or down; place your bets, and see if the gods will be in your favour!

## This REPO
This repo contains all the services and core code (and none of the client stuff!). That means business logic, APIs and other goodies that you would normally associate with a B/E.

```bash
.
├── README.md                   <-- This instructions file
├── internal                    <-- All internal services, routing, middleware, dbs etc
│   ├── db                      <-- All logic relating to interacting with the underlying database
│   └── handlers                <-- These are our API handlers - they are the glue that keeps things together
│   └── middleware              <-- Middleware for our API > in this case error handling
│   └── models                  <-- All models used throughout this app
│   └── services                <-- External services code > custom code to interact with external APIs
└── main.go                     <-- Lambda function code, our entrypoint
```

### API
Ultimately, this is an API with all the functionality necessary to run. It follows some `REST`-like principles, and the API itself is split into domains:

#### Users
The `users` API focusses on all functions relating to users and their votes. Since user and vote entities are tied together, they are both represented by this API together.

## Status

[![Deploy Status](https://github.com/svbygoibear/hermes-crypto-core/actions/workflows/deploy-to-lambda.yaml/badge.svg?branch=main)]()


## What makes me tick?

Under the hood, I am powered by;

-   [Gin](https://gin-gonic.com/): Gin is a fast, easy to use web-framework perfect for crafting APIs at scale!
-   [Golang](https://go.dev/): Go is an open-source programming language that is easy to learn, has tons of libraries and is well used and loved.
-   [DynamoDB](https://aws.amazon.com/pm/dynamodb/?gclid=Cj0KCQjwwuG1BhCnARIsAFWBUC2rKB9_2LwJA7ornNVDCMoK519wOKVKusZJtwk-hEHraZBI_hYHYHMaAs-NEALw_wcB&trk=bf64c969-685f-4fc4-b36b-4bcbda56cee7&sc_channel=ps&ef_id=Cj0KCQjwwuG1BhCnARIsAFWBUC2rKB9_2LwJA7ornNVDCMoK519wOKVKusZJtwk-hEHraZBI_hYHYHMaAs-NEALw_wcB:G:s&s_kwcid=AL!4422!3!536324221335!e!!g!!dynamodb!12195830303!119606857560): Following the theme of easy and lightweight, this project makes use of  DynamoDB to keep track of any info.

# Installation

### Software

To properly run this project, assuming you already have git installed, you will also need to ensure that you have the following installed on your machine:

-  * AWS CLI already configured with Administrator permission
-   [`Go`](https://go.dev/doc/install): First things first, click on the link to install the latest version of Go. This project has been tested with `v1.22.6` and offers support for this version. Follow the requirements based on your OS.
-   [`node.js`](https://nodejs.org/en): Lowest possible version compatible with this project is `v18.14.0`. The current LTS is however recommended.
-   [`AWS CLI`](https://docs.aws.amazon.com/cli/latest/userguide/getting-started-install.html): Interacting with any CLI commands will require the installation of AWS CLI.
-   [`AWS SAM CLI`](https://docs.aws.amazon.com/serverless-application-model/latest/developerguide/install-sam-cli.html): We will be using the AWS SAM CLI for deploying code. This is also to run our lambda API locally for testing. Currently this has been tested on version `1.121.0`.
-   [`Git`](https://git-scm.com/book/en/v2/Getting-Started-Installing-Git): You will also need to install Git (hopefully you have already - but if not now is your chance!).
-   [`Docker`](https://www.docker.com/): To make things easier, this project has been docker-rized. No more manual db setup, just compose and GO. 
-   [`NoSQL Workbench`](https://docs.aws.amazon.com/amazondynamodb/latest/developerguide/workbench.html): If you want to connect to your local instance of DynamoDB - use AWS' NoSQL Workbench. For more info on connecting to a local db; check [this](https://medium.com/@bthiban/running-dynamodb-locally-using-docker-68c8bbed29fa) out.
-   [`AWS Lambda runtime interface emulator`](https://docs.aws.amazon.com/lambda/latest/dg/go-image.html#go-image-provided): If you want to attempt to run and debug the underlying serverless Lambda architecture, you will need to install the emulator. More info [here](https://github.com/aws/aws-lambda-runtime-interface-emulator?tab=readme-ov-file#installing) as well. Take note of your operating system.
-   [`Make`](https://makefiletutorial.com/): Another way to interact with this repo is to make use of... Make! It is still useful, even if this is not a large program.


#### Extra info
- For multiple versions of Go, have a look here: [Managing Go Installations](https://go.dev/doc/manage-install).


## License

[MIT](https://choosealicense.com/licenses/mit/)

