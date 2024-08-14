<h1 align="center"> hermes-crypto-core</h1>
<p align="center"><img alt="hermes-crypto-core" src="../assets/hermes-crypto-logo.svg" width="200"></p>

# Critical Improvements
These improvements are rated quite high in priority as they could affect the functionality of the application to end users.

### Switch from Gecko > Binance
~~Currently we make use of [CoinGecko](https://www.coingecko.com/) to get up to date information on crypto prices. CoinGecko unfortunately have very poor limits on their free plan, thus switching to a 3rd party API such as [Binance](https://www.binance.com/) would be the smart move longer term. Currently when we run out of our total monthly API calls, the app will cease to function as this would cause an error.~~
This has been implemented.

# Overall Use Improvements
These are improvements and changes that can be applied that will expose more data to increase the feature set on the F/E or improve overall consistency.

### Add a historical crypto endpoint
This involves adding an endpoint that exposes some historical BTC data as well, so we can enable the F/E to create an interactive graph for users to see the change in value of BTC while they play the game.\

### Have a server-side trigger for score updates
Currently scores are updated using an endpoint from the F/E > but if there is a delay it means that we have a delay in processing the score. The solution here would be to implement something like a `webjob` or a queue server-side that checks for any votes that are older than 60 seconds to update them.

### Implementing Auth Flow
This is a big feature but the goal would be to create a proper auth service with OAuth integration, and use this to create an access token for users to identify themselves when they interact with the API. For this version of the project however, this is out of scope.

# Development Improvements
Development improvements may not be seen as important as critical, user-facing improvements but adding to any development improvements means better collaboration, better early issue detection as well as maintenance of the codebase which feeds into keeping the list of critical improvements _low_.

### Adding SwaggerDocs
End-users will typically not care if we have `swagger` documentation or not, but for the purpose of collaboration and treating the API like feature it is incredibly important. Thus, as part of future improvements, I would like to add swagger documentation as it:
- Provides a clear, interactive interface for developers to explore and understand the API, making it easier for them to integrate with your services. This means easier integration between this body of work and its corresponding F/E application.
- This documentation includes details about the API endpoints, request and response formats, authentication methods, and error codes, reducing the chances of miscommunication or incorrect implementation.  
- Additionally, it serves as a single source of truth, ensuring that the API documentation stays up-to-date with the actual implementation, which enhances both internal and external collaboration. 
- Finally, it can be used to test the API directly through the Swagger UI, making it easier to debug and validate endpoints.

A potential project that can be used to achieve this is called [`gin-swagger`](https://github.com/swaggo/gin-swagger) > this should be investigated for the next release. Additionally, hosting the documentation on AWS and ensuring it is updated on each release.

### Add API versioning
Speaking of making the API more robust, this projects API follows quite a few of the `REST API` standards. With that, we currently have not implemented versioning. For a future release, versioning should be included so we are able to move breaking changes to a newer versioned API, giving integration time to respond to these changes in the UI.

### Local Build Scripts
For this, I want to spend some time to use [SAM](https://aws.amazon.com/serverless/sam/) to help write some scripts so contributing developers can easily setup this project locally, develop and test on it. Currently testing has to be done via deploying this application and ensuring that we update configuration manually.

The ideal would be to keep this repo, as well as how we run it and how we deploy it, as consistent as possible.