## Golang http handler

## Installation
```
go get -u github.com/bugfan/srv
```

## Usage
```
// Refer to [this test file](https://github.com/bugfan/srv/blob/master/srv_test.go "test file")


addr := ":8080"
s := srv.New("addr")

// set head middleware
s.AddHeadHandler(&zlog{}, &auth{})

// set tail middleware
s.AddTailHandler(&zlog{})

// set your handler
s.Handle("/", &yourHandler{})
s.Handle("/ws/", &yourWebsocketHandler{})
s.Handle("/static/", &yourStaticHandler{})

/*
* if hava tls certificate data
*/
// keyData := []byte("xxxx")
// certData := []byte("xxxx")
// s.SetTLSConfigFromBytes(certData, keyData)

/*
* if hava tls config
*/
// tlsconfig := &tls.Config{}
// s.SetTLSConfig(tlsconfig)

// listen s
/*
* method 1
* run directly
*/

s.Run()

/*
* method 2
* if you have own listener
*/
// http.ListenAndServe(addr, s)

```
