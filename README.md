# Notification Service

### Main idea and use cases
- notification system that can handle multiple types of notifications on different channels
- the ui will be a screen in the settings panel of an existing app where the notifications will be split by
notification type(like different transactions, password changes, security code confirmations...)
- for each notification type the user can configure to which notification channels he wants to receive the notification on
- if he configures multiple -> he'll receive a notification for each channel
- I assume that we need to handle priority notifications as well because some notifications might not need to 
be delivered immediately but some like getting a security code might or a misused payment card

### Infrastructure
- Everything is containerized and I have a `docker-compose.yml` that contains all the configurations needed for starting everything 
- Main components:
  - API service -> that will handle the creation of the models in the database and publishing of a notification message
  - Kafka -> using it as a message broker for producing and consuming notification messages
  - Postgres -> storing all the entities and user information
  - Notification Consumer -> consumes and processes notifications that are published in the kafka

### Configuration
- all the configuration sits in one place in a `.env` file and is being parsed in one struct `common/config.go`
- this can be changed in the future so each component has its own configuration

### API service
- simple service for handling user authorization and creation of the entities
- not all endpoints are done, only some of them in order to showcase how this might work
  - creating notifications, templates and the configured user notification channels(not restful, needs to change for prod)
  - endpoints have to be done for updating and deleting all of these entities 
  - sending notifications(not restfull but easy to test, this will be removed in a prod environment)
- validation for the requests may be better as well, not all check if the given entities actually exist in the database(should be changed for prod)
- currently all endpoints are accessible to all authenticated people
  - this has to be changed in the future so that users with admin or prod-ops roles can access the create template and notification endpoints
- authentication happens when you call the authenticate endpoint with an available user in the database
  - you get a jwt token that you can then later user for calling the other endpoints
  - `auth` folder contains all the code needed for this and `api/middleware.go` contains the middleware handler for validating the jwt
- `make generate-swagger-docs` -> will generate the swagger docs for the handlers (stored in the `docs` folder)
  - when you run the api the swagger docs will be available at `localhost:8080/swagger/*` where you can a look on the documented endpoints
- `make generate-mocks` -> generates all the mocks for all the packages
- there is some small amount of code for creating example data for easier local development that has to be removed in the future

### Kafka 
- I'm using serama as a client library for kafka in go
- everything regarding kafka is located in the folder `kafka` 
- I currently have one kafka topic that will be used for publishing notification messages and each message consist of:
  - notification_id -> so we can get everything needed for the notification from the database
    - keep in mind all the notifications in the database are intended to be configured by an admin and to not change much
  - username -> the username of the user that has to receive the notification
  - priority -> the priority of the notification(a predefined set that can be changed -> low, medium, high)
    - this could be changed to numbers in the future so you can infinite priorities
- `At least once` message delivery -> I handle by manually committing the offsets when processing a message
  - if the consumer fails to process the message, panics or returns an error -> the message will be redelivered
- `Priority` of the notification -> the kafka topic is partitioned by priority 
  - I've currently set it up so the kafka topic I'm using has 10 partitions but this can be changed in the `docker-compose.yml`
  - I'm using a custom partitioning strategy that sends the messages to specific partitions based on their priority
    - `Bucket Priority Pattern` -> I divide my partitions into buckets, each bucket handles a different priority
      - the bucket size depends on the priority(higher priority messages have more partitions to go to)
    - for each message, I take its priority and determine to what bucket of partitions does it need to go based on a predefined distribution
      - I do that by first getting from kafka how much partitions are there for the given topic and then dividing them based on the distribution
      - The current distribution is -> {high: 50%, medium: 30%, low: 20%}
      - When I know to what bucket of partitions the message needs to go -> then I determine the specific partition randomly
- Each consumer has to part of one `consumer group`
  - because each message will be delivered to one consumer in each consumer group 
  - because you can set a rebalancing strategy in kafka if a given consumer is removed or added a new one
- `Sticky rebalancing strategy` -> kafka tries to assign only one partition for each consumer and tries to keep it that way
  - this is needed because of the bucket priority pattern
- The key here is to have the (number of partitions) >= (number of consumers)
  - if the number if consumers is smaller than the number of partitions then multiple partitions can be assigned to one consumer
  - in this way kafka will rebalance the consumers appropriately based on the partition count 
  - and publishing the messages in specific partitions achieves the handling of the messages based on their priority
- In this way -> the partitions and consumers can be scaled independently and the rebalancing logic is offloaded to kafka
  - scaling partitions -> in the `docker-compose.yml`, increase the partition number in the `kafka-topics-init` service
  - scaling the consumers -> in the `docker-compose.yml`, increase the number of replicas in the `notification-consumer` service
  - currently the partitions are 10 and the consumers are 5 -> based on the above algorithm
    - 5 of them will receive high priority, 3 of them medium, 2 of them low
    - the perfect distribution would be to scale the consumers to 10 if this scale doesn't work for your system, then each consumer will be assigned only one partition
- There are other ways to achieve the same distribution based on priority but they are a bit more inefficient and harder to scale
  - Have multiple consumer groups on the same topic and each consumer group can filter out messages they don't consume
    - really inefficient because all the consumers need to handle every message -> network load, unessecary resource utilization
  - Have multiple kafka topics for each notification priority
    - hard to add new add new priorities, as you need to add more topics and some consumers can be left idle and most likely the resources will not be equally utilized
    - perfect if this is crucial and you want to have different behaviours based on different priorities
  - Have one main topic that handles all of the notifications and router service that can then route on other kafka topics based on priority
    - unnecessary infrastructure and processing
    - again the problem is adding new priorities as you need to add new kafka topics 

### Postgres
- `repositories` - all the code for handling the entities stored in postgres
- I'm using `GORM` as a database client in order to easily map the entities in the database into golang struct and work with them  
- `repositories/client.go` contains the code for the creation and migration of tables
- each entity has its separate interface and object for handling the data so it can be changed in the future if needed
- tests are conducted by spinning up an sqlLite instance in memory and running the GORM library against it
- all the models are stored in the `repositories/models.go`, it can be split in the future for easier readability
- for local development I spin a small postgres container in the `docker-compose.yml`'
- some of the entities can be simplified regarding their representation in the database because currently it's not the most optimal way
  - i did it for ease of implementation but this has to be changed in the future

### Entities
- user: (username, password_hash, notification_settings)
  - notification settings contains the information needed when a user has a notification channel configured
  - I imagine that in order to be able to select a given notification channel you would to set it up 
  - this should be handled from separate endpoints
  - currently the whole object is stored in one column -> this could be split in the future so it exists in a separate table
- template: (id, channel, template)
  - will be configurable by business people, prod ops, admins and so on 
  - channel -> the notification channel this template is for
  - template -> holds the actual template -> for now it's a string field -> it can be configured to be a link to a template stored in a separate place in the future
- notification: (id, type, priority, channel_to_template_map)
  - will be configurable by business people, prod ops, admins and so on
  - type -> the type of the notification -> user did something, amount received, transaction failed...
  - priority -> what's the priority of the notification, how urgent is it (high, medium, low)
  - channel_to_template_map -> a structure for holding a mapping between the notification channel and it's corresponding template to this specific notification type
    - I imagine that for each notification type the user can select on which of his configured notification channels he wants to receive it
    - I imagine that he can select multiple notification channels for the same notification type
    - I imagine that the available channels for each notification types can be stored in the same map and quickly determine what channels are available for what notification 
- user_notification_channels: (username, notification_id, channels)
  - channels -> holds information for this specific notification what are the notification channels that the user has configured 
  - you can quickly determine to what channel you should send the given notification

### Notification Consumer
- it has one main component: notification factory and its split into several components
  - metadata generator, content generation, notification sender
- metadata generator, content generation, notification sender operate using the strategy design pattern
  - each component has different concrete implementations based on the notification
  - for each one of them I have a separate entrypoint/dummy implementation that switches between the implementations
- metadata generator -> generates metadata based on the notification type(everything needed to build a notification)
  - currently the concrete implementations are empty -> this will depend on the different notifications the system will support
- content generator -> generates the content based on the notification channel
  - I imagine that most of the time the way you build the content will stay the same depending on the notification channel
  - it's a separate component on purpose as it could be changed or merged in the metadata generator if the use cases for the notifications differ more on notification type rather then channel
- notification sender -> sends the notification based on the notification channel
  - currently each sender return success even if it's not successful just for ease of local development -> this has to be changed in the future
  - I've setup dummy clients for the email and sms senders that are located in `docker-compose.yml` for easier development
- notification factory -> uses all the components above to stitch the business logic together
  - gets the notification and user entities from the database
  - generate the notification metadata using the metadata generator
  - extracts the configured user notification channels for this notification type
  - extracts the user's notification settings for each channel
  - for each available channel in the notification and that the user has set up
    - get the template from the database(this can be extracted and done in bulk beforehand so we minimize db calls)
    - generate the notification content using the content generator based on the metadata and template
    - get the user notification settings based on the notification channel
    - send the notification based on the notification channel, user notification settings and content
