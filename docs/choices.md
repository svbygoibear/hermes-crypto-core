<h1 align="center"> hermes-crypto-core</h1>
<p align="center"><img alt="hermes-crypto-core" src="../assets/hermes-crypto-logo.svg" width="200"></p>

# Technical Decisions
Wondering why some technologies were used instead of others? Well here we answer some of those questions.

## AWS Lambda
I've decided to go serverless with this project and use AWS Lambda to achieve this. Why, and why do we technically only have 1 function that exposes this whole project?
-  It simplifies deployment and management, as we only have one Lambda function to configure, monitor, and update, reducing operational overhead. 
- This approach also leverages AWS Lambda's automatic scaling, allowing the backend to handle varying levels of traffic without needing manual intervention. 
- Additionally, by bundling everything into one function, we can take advantage of Go's fast startup times and performance within Lambda's environment, ensuring that your API remains responsive. 
- Finally, this setup can be cost-effective, as we only pay for the compute time when this function is actually invoked, eliminating the need for maintaining always-on servers.

For a bigger API design, I would perhaps create a Lambda function per API domain, but that was out of scope for the scale of this project currently. On top of that, for the size of this project it is also easier to keep all our functionality in one repo.

## Using gin
I am not always in favour of frameworks, bug `go-gin` is incredibly powerful. I use it because;
- It is a lightweight and high-performance web framework, designed to deliver fast response times and handle a large number of requests efficiently, making it ideal for building APIs and web services. 
- Its simplicity and minimalism allow us to write clean and concise code, reducing the complexity often associated with web frameworks. 
- It also provides a comprehensive set of features out-of-the-box, such as routing, middleware support, and JSON validation, which streamline the development process and reduce the need for third-party libraries. For a demo project like this, it means we are able to hit the ground running without sacrificing code quality.

## DynamoDB
Most of the data we need to store would make almost more sense to store in a relational database and leverage some of the built in features around referential integrity; but for the cost and time purposes of this project DynamoDB was chosen. On top of that, it is also fully managed which reduces worries around maintenance, scaling or backups. This simplifies the infrastructure and makes sense as this is a non-critical application. On top of that, its flexible data model allows us to store and retrieve data without needing to define a strict schema upfront, which is ideal for small projects where requirements may evolve over time (such as this one). 

## CI/CD
To keep things simple and easy, as part of this pipeline there are a few main actions:
- `Drafting Releases`: Release drafter makes our lives easier by automatically drafting releases based on the tags as described in the contribution guide. What it means is that when we are ready to officially create a release version, we don't have to worry about what has changed.
- `Tests on PR`: Writing & maintaining tests are incredibly important. Writing a test for the same of writing a test is not useful, but it is equally bad to have no tests at all. This action makes it easier for reviewers to check on the quality of code, see if tests were added/changed and help devs give them confidence in the changes they make.
- `Deploy after merge`: Again, the ability to deploy when actions fire off are powerful. As part of this workflow, when changes are pushed into `main`, we start kicking off a build and deploy of this application to AWS Lambda. In future this script can be changed so that changes to `main` goes to our test environment and changes to something like a `prod` branch goes into the final production environment.
