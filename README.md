# Real-time Map

_Real-time Map_ displays real-time positions of public transport vehicles in Helsinki. It's a showcase for Proto.Actor - an ultra-fast distributed actors solution for Go, C#, and Java/Kotlin.

_This repository containes the Go version of the sample_

The app features:
* Real-time positions of vehicles.
* Vehicle trails.
* Geofencing notifications (vehicle entering and exiting the area).
* Vehicles in geofencing areas per public transport company.
* Horizontal scaling.

The goals of this app are:
1. Showing what Proto.Actor can do.
1. Presenting a semi-real-world use case of the distributed actor model.
1. Aiding people in learning how to use Proto.Actor.

**[Find more about Proto.Actor here.](https://proto.actor/)**

![image](https://user-images.githubusercontent.com/1219044/132653003-58733735-f49a-4615-adb5-36552b1415c1.png)

## Running the app


Configure Mapbox:
1. Create an account on [Mapbox](https://www.mapbox.com/).
1. Copy a token from: main dashbaord / access tokens / default public token.
1. Paste the token in `frontend\src\mapboxConfig.ts`.

Start Backend:
```
cd backend
go run main.go
```

Start frontend:
```
cd frontend
npm install
npm run serve
```

The app is available on [localhost:8080](http://localhost:8080/).

## How does it work?

Please refer to the [.NET version README](https://github.com/asynkron/realtimemap-dotnet#how-does-it-work) for a detailed description of the architecture.