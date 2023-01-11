# Groupie tracker

<ins>By : **{{ .ProjectName }}**</ins>

Group :
- BOIZOT Maxime
- SAROUILLE Alexandre
- BONNELL Nathan

Trello : https://trello.com/b/rIyCGasZ/groupie-tracker

<ins>Subject:</ins>  
Groupie Trackers consists on receiving a given API and manipulate the data contained in it, in order to create a site, displaying the information.

It will be given an API, that consists in four parts:
* The first one, artists, containing information about some bands and artists like their name(s), image, in which year they began their activity, the date of their first album and the members.
* The second one, locations, consists in their last and/or upcoming concert locations.
* The third one, dates, consists in their last and/or upcoming concert dates.
* And the last one, relation, does the link between all the other parts, artists, dates and locations.

Given all this you should build a user friendly website where you can display the bands info through several data visualizations (examples : blocks, cards, tables, list, pages, graphics, etc). It is up to you to decide how you will display it.

This project also focuses on the creation of events/actions and on their visualization.

The event/action we want you to do is known as a client call to the server (client-server). We can say it is a feature of your choice that needs to trigger an action. This action must communicate with the server in order to recieve information, ([request-response])(https://en.wikipedia.org/wiki/Request%E2%80%93response)
An event consists in a system that responds to some kind of action triggered by the client, time, or any other factor.

## Launch of the project

Do the command : 

```bash	
go run webserver.go
```

Then go to the address : http://localhost:8080/
to see the result.
