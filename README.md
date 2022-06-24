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

## Pricing Component

** TO BE COMPLETED AFTER THE PLANS COMPONENT IS INTEGRATED WITH STRIPE.COM **

## Plans Component - Frontend
**What needs to be done**
1. In officekube portal repo modify the code of the pricing.tsx component to handle the logic of a user switching from one plan to another as follows:
   - The component should disable the button Sign Up button for the plan that the user is currently on. To determine what the current plan is the component should call an API endpoint of the billing service GET /plans/current at the time of its loading.
   - When the user clicks on Sign Up button of any other plan a Switch Plan dialog should pop up (refer to the design here https://www.figma.com/file/vDqU6NBTspvomTQGN4nel0/3D-Workspace-(Community)?node-id=437%3A490). The dialog should be created using Syncfusion React library - https://ej2.syncfusion.com/react/documentation/dialog/getting-started/).

** Resume after reviewing the stripe docs **


## Plans Component - Backend
**What needs to be done**
The service is based on the following stack:
- Language: go-lang
- Web Framework: Gin Web Framework
- Configuration Manager: Viper
- ORM: GORM
- DB Backend: PostgreSQL
- Code Generator: Open API Code Generator

1. **Endpoint GET /plans**, as a part of this assignment implement /plans  (refer [openapi.yml](https://gitlab.dev.workspacenow.cloud/platform/subscription-manager/-/blob/main/api/openapi.yml)). Assume that all available plans are stored in a table AvailablePlans persisted in the PostgreSQL db to which the service has read/write access.
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

2. ** Endpoint /create-checkout-session **. Use the [quickstart guide](https://stripe.com/docs/checkout/quickstart) for the Stripe Checkout integration to add an endpoint /payments/create-checkout-session and its implementation.
When doing so:
- ** Make sure to add the endpoint to the module's openapi.yml first. **
- Implement the mentioned in the guide function createCheckoutSession in a separate file go/api_payments_checkout.go and using the GIN web framework (not the net/http module).
** Resume using the quickstart **
