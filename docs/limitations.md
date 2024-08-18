<h1 align="center"> hermes-crypto-core</h1>
<p align="center"><img alt="hermes-crypto-core" src="../assets/hermes-crypto-logo.svg" width="200"></p>

# Logic & Limitations
Here we take a look at some of the logic behind the app, the limitations thereof and overall improvement flow.

## Current Solution Logic
Currently we assume:
 - After 60 seconds, the BTC price WILL change.
 - That a small threshold; give or take 5 seconds, is OK between vote placed & result checked.

With that, it means that when a user places their bet and a new vote is created - that vote is either in the following states:
- Result check pending
- Ready to check result
- Expired/more than +-60 seconds has passed (give or take 5 seconds)

Users are not allowed to place another vote (enforced by the API) until their last bet has been resolved. If they close their browser and come back to the last bet OR an error occurred on the last bet, we will attempt to resolve it for them regardless if it is "ready to check" OR "expired". This is a known limitation right now.


### Improvements to Solution
Considering the assumptions and limitations above, there are ways we can improve on the solution as it stands, but with various Pros and Cons.

#### Expired Votes
We handle expired votes the same way as ready to check votes; but this can be improved by adding an additional user flow where we inform the user that their vote has expired & that no change will be made to the end total of their score. Feedback to the user is critical in this case so that the vote does not just "disappear", but for now the compromise is to resolve the last vote when it is possible.

#### Result Calculation
The result for the last vote is calculated via the API itself. Blocking the user from placing a new vote is also done via API logic making sure that there is no way to "hack" the system. That being said, a better solution for more "fairly" calculating votes would be implementing a background job. 

##### Pros & Cons of Improvements
Both of these improvements are interlinked - depending on which solution is implemented. If we have a background job, the chances of an expired vote is lower as we do not rely on the client to call the API to calculate the results. However we may still want to implement an "expired" flow for the event where the async flow fails; but that brings up larger questions as to how to we compensate a user when the failure is our fault?

Pros of Current Implementation:
- As it stands the current solution is easy to follow (logic wise) and the reduction in complexity means we can add more features and enrich it later on.
- More complex "edge case" logic in the F/E can become technical debt to remove once the underlying technical issue has been resolved. Sometimes less work upfront is also less work later to refine it.
- Since the logic already resides in the B/E layer of this app, it makes it easy to add in an async function to more fairly handle vote resolution.
- Simple, easy and working solutions means we can iterate faster, add more features with less "bugginess" and "edge-cases" but also help us figure out with usage what those features should be. It is a fine balance between assumption, engineering for the end-user but also not preempting scenarios that do not exist.

Cons of Current Implementation:
- This discrepancy in vote resolution reduces fairness overall towards a users total score.
- The lack of feedback to the user about the "expired" vote also contributes to that lack of fairness; albeit the user is not able to do anything about it/action anything off of that error.
- I would have liked the separation of the business logic on when to calculate the users' result to reside on the B/E > reducing the reliance on the user needing their browser to be open or to have a stable internet connection.



## Feature Improvements
For general features of the application, I would like to add the following (not in order of importance):
- Make the flip clock look even more like a flip clock (and start from 1 minute for countdown).
- Add a user page that each signed in user can view to see their stats.
- Add a delete my profile functionality of the F/E.
- Use OAuth for signup.
- Include a graph on page load to show historical data of BTC as well as the users vote data in relation to that.



## Technical Improvements
- Use something like tailwind or styled components to get a better grip on the styling side of the F/E.
- Even more useful tests; focus on splitting up between end-to-end tests, integration tests & unit tests.
- Redesign of the DB > DynamoDB with one table was used for ease of implementation/getting this up and running quickly _but_ I would like to properly split out the votes from the users and add more data to capture.
- Add swagger docs for the API itself.
- Implement a "vote processing" service that picks up and calculates votes for items on a queue as soon as they hit the 60 second age mark.
- Make it easier to start up app locally & look at infrastructure as code for deploying it again in future.
- Use AWS Key Management Service for keeping some of the keys that we need, and configure the environment as the variable to fetch them by.
- Overkill for this scenario - but better logging. Using something like DataDog or just step 1; improving the logs overall
- Added a cmd folder for the B/E to help split up some of the setup and then just initializing everything in the main.go. Not an issue for the app at its current scale - but will be if it grows.
