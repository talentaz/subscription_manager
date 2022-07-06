# Introduction 
The purpose of this document is to describe the subscrintion manager module, its functions, and design.

# Overview
The pricing manager is a backend service that implements billing, payment history, and subscription features. It implements them through integration with the payment gateway stripe.com

# Design
From the user perspective, the module will have the following pages:
- **Pricing**. The page is available to both registered and unauthorized users and describes all available plans.
- **Plans**. The page is available to registered users only (as a part of user account?) and allows the user to upgrade/downgrade/cancel their current plan.
- **Billing**. The page is available to registered users only (as a part of user account?) and allows the user to view upcoming/current billing statement as well as the history of payments.

To accept payments the module will use [Stripe Checkout](https://stripe.com/docs/payments/checkout) integration method (from the [3 available](https://stripe.com/docs/payments/online-payments)).

## Database Structure
The service has the following table structure:
```plantuml
object transactions {
  id: uuid
  userId: uuid
  priceId: varchar(100)
  userPlanId: uuid
  sessionId: varchar(100)
  customerId: varchar(100)
  status: varchar(50)
  created_ts: timestamp
  last_modified_ts: timestamp
}

object user_plans {
  id: uuid
  userId: uuid
  planId: int4
  customerId: varchar(100)
  priceId: varchar(100)  
  subscriptionId: varchar(100)
  status: varchar(50)
  created_ts: timestamp
  last_modified_ts: timestamp
}

object available_plans {
  id: int4
  name: varchar(50)
  description: varchar(250)
  price: numeric
  recurrence: int4
  priceId: varchar(100)
}


transactions "N " --> "1 " user_plans
user_plans "N " --> "1 " available_plans
```

## Pricing Component
The pricing page on the web-site as well the portal component pricing.tsx should be implemented as per this [design](https://www.figma.com/file/vDqU6NBTspvomTQGN4nel0/3D-Workspace-(Community)?node-id=437%3A490).

## Plans Component - Frontend
**What needs to be done**
1. In officekube portal repo modify the code of the pricing.tsx component to handle the logic of a user switching from one plan to another as follows:
   - The component should disable the button Sign Up button for the plan that the user is currently on. To determine what the current plan is the component should call an API endpoint of the subscription manager service GET /plans/current at the time of its loading.
   - Replace button Notify Me for the plan Solo with the button SIGN UP.
   - When the user clicks on Sign Up button of any other plan a Switch Plan dialog should pop up (refer to this [design](https://www.figma.com/file/vDqU6NBTspvomTQGN4nel0/3D-Workspace-(Community)?node-id=437%3A490)). The dialog should be created using [Syncfusion React library](https://ej2.syncfusion.com/react/documentation/dialog/getting-started/).
   - In the dialog when the user clicks button Switch make a GET call to the endpoint /payments/checkout of the subscription manager backend (see below) and pass a query parameter named price_id with the value that depends on which button SIGN UP has been clicked by the user before they got to the dialog. For the plan Enthusiast that value should be set "free", for the plan Solo it should be "price_1LECZyKUSkDFrC1EroX3h7NW".
   - If the user clicked Cancel simply return them back to the Plans page.

### Success and Failure Pages
After a user has been redirected to the Stripe checkout page, Stripe will redirect the user back to either a success or a failure page indicating whether the user has successfully signed up for our subscription.
The pricing.tsx will be responsible for showing either success or failure. For that to work the page might receive a URL parameter named checkoutResult. Modify the page as follows:
1. Use the react-router-dom library to retrieve the value of the URL parameter checkoutResult immediately after the page has been loaded into a web-browser.
2. If the parameter is equal to "success" then show a popup message (using Dialog component from Syncfusion library) with a button OK and a message "Thank you for your subscription!". When the user clicks OK, the dialog should be closed.
3. If the parameter is equal to "failure" then show a popup message (using Dialog component from Syncfusion library) with a button OK and a message "Sorry, something went wrong. Please try again later or contact us!". When the user clicks OK, the dialog should be closed.
4. If the parameter is not set to any value then no action should be taken.

## Plans Component - Backend
**What needs to be done**
The service is based on the following stack:
- Language: go-lang
- Web Framework: Gin Web Framework
- Configuration Manager: Viper
- ORM: GORM
- DB Backend: PostgreSQL
- Code Generator: Open API Code Generator

### 1. Endpoint GET /plans
As a part of this assignment implement /plans  (refer [openapi.yml](https://gitlab.dev.workspacenow.cloud/platform/subscription-manager/-/blob/main/api/openapi.yml)). Assume that all available plans are stored in a table AvailablePlans persisted in the PostgreSQL db to which the service has read/write access.
The service should make a call into the DB and retrieve all records from the mentioned table where field active is equal to "true". The table AbailablePlans has the following structure:

- id int primary key auto incremental
- name char 50 required
- description char 250
- price real required
- recurrence int required default 30

For testing purpose the table can be populated with the following records:

|id|name|description|price|recurrence|
|-|-|-|-|-|
|1|Enthusiast||0|30|
|2|Solo||10|30|
|3|Expert||30|30|
|4|Team||100|30|

**Development Approach**
1. Generate a code skeleton for the application using the following command:
openapi-generator-cli generate --package-name workspaceEngine -g go-gin-server -i openapi.yml 
1. Implement the endpoint GET /plans using Gin Web Framework and GORM ORM (for access to DB).
1. Avoid hard-coding service configuration (e.g. db connection parameters). The service configuration should be persisted in a YAML file subscription_manager.yml.


### 2. Endpoint GET /payments/checkout
Use the [quickstart guide](https://stripe.com/docs/checkout/quickstart) for the Stripe Checkout integration to add an endpoint /payments/checkout and its implementation.
When doing so:
  - **Make sure to add the endpoint to the module's openapi.yml first.** The endpoint should receive a string query parameter named price_id.
  - Implement the mentioned in the guide function createCheckoutSession in a separate file **go/api_payments_checkout.go** and using the **GIN web framework, not the net/http module**. 
    - **Set up the server**. Use this for stripe.Key
    `sk_test_51LECU4KUSkDFrC1EovIjSi4jNsHRwz3eT8CggtBRBfPtLORVIMd7Md1sDxDe71lGvO0AR1bMJXO6uNbxDFFru4Yx00dkzMY092`
    - **Create a Checkout Session**. The function should extract from the gin context a value of the price_id parameter and use it as per the guide. In code of the function do not use the variable domain and instead populate parameters SuccessURL and CancelURL with values of the config parameters with the same names (SuccessURL and CancelURL)
    - **Define a product to sell**. For the parameter Price use the valie of the query parameter price_id.
    - **Choose the mode**. Use the mode subscription.
    - **Supply success and cancel URLs**. See the note for **Create a Checkout Session** above.


### 3. Endpoint POST /payments/stripewebhook
Use the [stripe guide](https://stripe.com/docs/payments/checkout/fulfill-orders) for the Stripe Checkout integration to add an endpoint /payments/stripewebhook and its implementation.
When doing so:
  - **Make sure to add the endpoint to the module's openapi.yml first.** 
  - **Create your event handler**. Implement the mentioned in the guide function for the /webhook in a separate file **go/api_payments_stripewebhook.go** and using the **GIN web framework, not the net/http module**. 
    - **Verify events came from Stripe**. The value for the variable endpointSecret should be loaded from the service's configuration file.
    - **Define a product to sell**. For the parameter Price use the valie of the query parameter price_id.
    - **Fullfill the order**. The function FulfillOrder will be empty for now.


### 4. Function FulfillOrder
1. Create tables transactions and user_plans and update the table available_plans as per the [db design](https://gitlab.dev.workspacenow.cloud/platform/subscription-manager/-/edit/main/README.md#database-structure) and their models in the code.
2. In security.go add a function GetUserId that should return a string containing a user id. The implementation of the function replicates almost 100% the code of a function IsApiAuthenticated and additionally extracts the user id from idToken.Subject and returns it. Make sure that both function re-use the same code (rather than copying and pasting it).
3. In handler for the endpoint /payments/checkout right after retrieving the price_id add the code that will:
   - Pull a record from the table user_plans using user id (use function GetUserId), price_id. If such a record exists and its status is equal to "CURRENT" return http code 208.
4. In the same handler for the endpoint /payments/checkout right after creating a stripe session but before redirecting a user to s.URL insert code that would perform the following:
   - Pull a record from the table available_plans where priceId = price_id.
   - If a record from the table user_plans (pulled in step 3) exists update its status to "CHECKOUT", planId, priceId, and the field last_modified_ts to current timestamp. Othwerise, create a new record with properly populated fields id (newly generated uuid), userId, planId (plan_id), priceId, created_ts (set to current timestamp), and status (set to CHECKOUT).
   - Create a new record in the table transactions and populate its fields userId, priceId, sessionId (s.ID), userPlandId (id of the record from the table user_plans), status (CHECKOUT), and created_ts (set to current timestamp).
5. Modify function FulfillOrder as follows:
   - Using sample code [here](https://stripe.com/docs/payments/checkout/custom-success-page), extract customer ID from the stripe session.
   - Pull a record from the table transactions where sessionId = session.ID and update its fields customerId (session.customer.ID), status = CURRENT, last_modified_ts = current timestamp.
   - Pull a record from the table user_plans where id = userPlanId (pulled from transactions) and update its fields customerId (session.customer.ID), subscriptionId (session.subscription.ID), status = CURRENT, last_modified_ts = current timestamp.

### 5. Endpoint GET /plans/current
Implement the endpoint /payments/current as follows:
- Create the endpoint handler in a separate file go/api_payments_current.go and using the GIN web framework. 
- Secure the endpoint with a call to IsApiAuthenticated().
- Pull a user id using the function GetUserId
- Retrieve a record from user_plans where userId == user id and status == 'CURRENT'.
- If no record is found then return an http code 404.
- If a record has been found using its field planId retrieve a matching record from the table available_plans.
- Create an instance of the model APlan, populate its properties with proper values from the user_plans and available_plans records and return the model along with http code 200.

### 6. Upgrading/Downgrading Subscription
If a user is already subscribed and requests to upgrade/downgrade their subscription the following needs to be done:

1. Update the table user_plans (add fields subscriptionId and priceId) as per the [db design](https://gitlab.dev.workspacenow.cloud/platform/subscription-manager/-/edit/main/README.md#database-structure) and its model in the code.
2. In code of the handler for the endpoint /payments/checkout where a record from the table user_plans is retrieved and before updating its status to "CHECKOUT" and the field last_modified_ts to current timestamp store the value of its field priceId into a local variable. If the record has not been found, modify the code creating a new record with properly populated fields id (newly generated uuid), userId, planId (plan_id), created_ts (set to current timestamp), status (set to CHECKOUT), **and the field priceId set to the passed in value**.

**FIXME: The logic of checking the current plan (see table transactions) is messed up. Need to fix it first**.

### 7. Cancelling Subscription


**TO BE COMPLETED**: add logic of switching between plans (i.e. handling a case if the previous subscription/plan needs to be cancelled first and the case when the user switches to/from the free plan).
