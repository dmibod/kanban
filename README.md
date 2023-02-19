# Kanban

## [Architecture](https://www.draw.io/?lightbox=1&highlight=0000ff&edit=_blank&layers=1&nav=1&title=kanban-overview.html#Uhttps%3A%2F%2Fdrive.google.com%2Fuc%3Fid%3D1-IE9Vru5Vs3sAgcKoh9mQOlP3p4KNKEd%26export%3Ddownload)

<!--
```plantuml

@startuml diagram
Interface Abstraction {
    +operation()
}

Abstract class Implementation {
    +operationImpl()
}

Class ConcreteImplementationA {
    +operationImpl()
}

Class ConcreteImplementationB {
    +operationImpl()
}

Class RefinedAbstraction {
    -implementation: Implementation
    +operation()
}

Abstraction -> RefinedAbstraction
RefinedAbstraction *-> Implementation
ConcreteImplementationA -> Implementation
ConcreteImplementationB -> Implementation
@enduml

```
-->
![Diagram](/diagram.svg)
