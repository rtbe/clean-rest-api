# Description

This is overview of project`s architecture and some of architectural decisions that has been made during constructing it.

## Structure of a project

```/delivery``` - Contains means of communications for project with outer world (REST api).

```/domain``` - Contains project business entities and their use cases.

```/internal``` - Contains libraries (util packages) that helps to address different infrastructure project needs. Such libraries provides:

- Logging.
- Validating data.
- Managing connections to different databases.
- Testing.
- Documenting.

```/repository``` - Defines the set of interfaces for each of project services that handles persistent storage independent of particular technology used for that need. These interfaces serve as Ports (in terms of hexagonal architecture) or Driven adapters.
Using of that pattern helps us to invert dependencies between our application needs in persistent storage of data for different services with concrete implementation of that behavior. That helps us to provide modularity for our persistent storage needs, so we can:

1. Make clear and clean boundaries between data storage and busyness logic (Single responsibility, separate of concerns - S in SOLID), so decisions and their implications do not propagate through whole project. In other words provide decoupling through dependency inversion (D in SOLID).
1. Hide implementation details into interface implementation so we can use different approaches to storing data for each of project services.
1. When building project you can delay making a decision about use of concrete technology for storage project's data, and just use plain built-in map. That can help in speeding up of building prototypes.

To implement persistent data storage i decided to use relational data model storage provided by PostgreSQL, but for storing users authentication data i decided to use non-relational model storage provided by MongoDB.
