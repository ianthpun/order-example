# Order Orchestrator
The order orchestrator is a service that manages the workflow orchestration of an Order,
which consists of order validation/confirmation/cancellation/expiry, payment processing, and asset delivery within Dapper.

## Design Approach

The service is structured to be in a Clean archrecture design while using an embedded workflow engine.


[insert graph]

Essentially, GRPC will be used to signal different functions of the workflow, which then would be run against the workflow engine.
The workflow engine will then orchestrate the application layer in order to process the order workflow to something similar to this:

[insert graph]

## notes
- credit transfers are currently handled through the [main monorepo](https://github.com/dapperlabs/dapper-flow-api/), where we have copied over the protobuf files from. We are currently not able to import them from the monorepo as the import is too large.

    