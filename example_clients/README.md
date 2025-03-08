# Client Implementation Details

### Disclaimer
The implementation you see in this folder is not complete, and assumes that
only one client will ever send forefeit/draw requests. Additionally, it is not
complete.

Additionally, it was rushed/hacked together, and should **only** be used as
reference, and ***never in production***.

# How to Implement a CBMR Client

## 1: Initialization
When the client first starts up, it should send a message to the server that
looks like this:
```
GET /connect?name=[player 1 name] HTTP 1.1
```

~~This will *only* error out **if and only if** a player already exists with that
name on the server.~~

## 2: When the round ends/forefeit/draw

### 2.1: Round ending
TODO.

### 2.2: Forefeit/draw
On forefeit/draw, the following happens
```
player1 tells server it wants to forefeit/draw ----|
                                                   |
     server messages player2 about it <------------|
            |
            |-----> player2 either accepts/denies, sends back to server
                                            |
  match ends,                               |
 winner decided by   <----------------------|
 data["winner"] of
  JSON object
```

## 3: Security
**PLAYER 1 AND PLAYER 2 SHOULD NEVER DIRECTLY EXCHANGE MESSAGES!!!!!**

## 4: TODO
TODO