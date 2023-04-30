# Groupie Tracker
This is a web application written in golang tracking the [heroku api](https://groupietrackers.herokuapp.com/api/artists) that provides a lot of informations about musical bands.

## Installation

Use ``git`` to install it from our repository.

```bash
$ git clone https://github.com/XinJann/GroupieTracker.git
$ cd GroupieTracker
$ go build
```

## Usage
(Make sure to be in the ``groupie-tracker`` directory)
```bash
$ go run main.go
```

## Objectives 

- [x] Extract and display all bands data from the api.
- [x] Multiple filters such as sort by member size or concert location.
- [x] A searchbar with autocompletion of band name or member name.
- [x] Use of caches so filters will be quicker overtime.
- [x] Refresh the api and all the caches every 24 hours.
- [x] Log crucial informations such as filling/emptying caches.
- [x] Display the location of all concerts of a certain band on an earth map.
- [ ] Make some coffee.
