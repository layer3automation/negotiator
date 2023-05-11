# Layer3 Negotiator

The main usage of the negotiator is to connect two different networks that are connected by a pure L2 link. You can learn more about the use case scenario below.

The negotiator is meant to be a fixed interface to start and handle a negotiation, it doesn't implement any kind of code to configure network devices.

## Interfaces

The project defines an API interface with the following endpoints:

* `POST /api/v1/start_negotiation` is used to start a negotiation. It is called when you want to start a negotiation with another configurator, internally it calls the remote negotiator on `POST /api/v1/handle_negotiation`, which is explained below.
* `POST /api/v1/handle_negotiation` is used to handle a negotiation request from another negotiator.

Additionally, the negotiator has a gRPC connection to a configuration agent that is in charge to configure the network devices. You can learn more [here](https://github.com/layer3automation/configuration_agent_template)

The big picture is, in the case of a single router to be configured:
![architecture](/images/Architecture.png)

## Use case scenario

The negotiator is developed to be part of a bigger picture, which is the Internet Exchange federation network fabric. The use case is described by the following picture which shows how the communication happens and how the different networks are connected.

![full diagram](/images/FullDiagram.png)

## Next Steps

The configurator is far from being completed. Particularly the following aspects need significant improvements:

* Parameters passed to the configuration agent. I've included in those parameters only the ones useful for the initial development but of course, they are not all the possible useful parameters that will cover every possible configuration scenario.
* Parameters passed between negotiators (body of `POST /api/v1/handle_negotiation`). Now only the essential values are passed but they will not cover all the possible use cases.
* Parameters needed to start the negotiation. Now I've included only the essential ones.
