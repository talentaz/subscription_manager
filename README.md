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
object Transactions {
  id: uuid
  userId: uuid
  priceId
  userPlanId
  status: varchar(5)
}

object UserPlans {
  id: uuid
  userId: uuid
  planId: int4
}


Transactions --> UserPlans
```

## Pricing Component

**TO BE COMPLETED AFTER THE PLANS COMPONENT IS INTEGRATED WITH STRIPE.COM**

## Plans Component - Frontend
**What needs to be done**
1. In officekube portal repo modify the code of the pricing.tsx component to handle the logic of a user switching from one plan to another as follows:
   - The component should disable the button Sign Up button for the plan that the user is currently on. To determine what the current plan is the component should call an API endpoint of the subscription manager service GET /plans/current at the time of its loading.
   - Replace button Notify Me for the plan Solo with the button SIGN UP.
   - When the user clicks on Sign Up button of any other plan a Switch Plan dialog should pop up (refer to the design here https://www.figma.com/file/vDqU6NBTspvomTQGN4nel0/3D-Workspace-(Community)?node-id=437%3A490). The dialog should be created using Syncfusion React library - https://ej2.syncfusion.com/react/documentation/dialog/getting-started/).
   - In the dialog when the user clicks button Switch make a GET call to the endpoint /payments/checkout of the subscription manager backend (see below) and pass a query parameter named price_id with the value that depends on which button SIGN UP has been clicked by the user before they got to the dialog. For the plan Enthusiast that value should be set "free", for the plan Solo it should be "price_1LECZyKUSkDFrC1EroX3h7NW".
   - If the user clicked Cancel simply return them back to the Plans page.

## Plans Component - Backend
**What needs to be done**
The service is based on the following stack:
- Language: go-lang
- Web Framework: Gin Web Framework
- Configuration Manager: Viper
- ORM: GORM
- DB Backend: PostgreSQL
- Code Generator: Open API Code Generator

### 1. **Endpoint GET /plans**
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


### 2. **Endpoint GET /payments/checkout**
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


### 3. **Endpoint POST /payments/stripewebhook**
Use the [stripe guide](https://stripe.com/docs/payments/checkout/fulfill-orders) for the Stripe Checkout integration to add an endpoint /payments/stripewebhook and its implementation.
When doing so:
  - **Make sure to add the endpoint to the module's openapi.yml first.** 
  - **Create your event handler**. Implement the mentioned in the guide function for the /webhook in a separate file **go/api_payments_stripewebhook.go** and using the **GIN web framework, not the net/http module**. 
    - **Verify events came from Stripe**. The value for the variable endpointSecret should be loaded from the service's configuration file.
    - **Define a product to sell**. For the parameter Price use the valie of the query parameter price_id.
    - **Fullfill the order**. The function FulfillOrder will be empty for now.


### 3. **Function FulfillOrder**
**TO BE COMPLETED**. Need to think about the link to the account service and db design.

**TO BE COMPLETED**: add logic of switching between plans (i.e. handling a case if the previous subscription/plan needs to be cancelled first and the case when the user switches to/from the free plan).
