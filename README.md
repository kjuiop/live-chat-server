# live-chat-server

---

live-chat-server 는 채팅 메시지를 주고 받을 수 있는 서버입니다.

live-chat-server 는 단독 서버 형태로 운영되며, 데이터 흐름은 아래와 같습니다.

![image](https://github.com/user-attachments/assets/b64af16b-e320-49c9-8187-6cadc1b12c3c)

- 사용자는 http api 를 통해 chat-room CRUD 기능을 수행합니다.
- 채팅방 생성 이후 사용자가 해당 채팅방에 입장하면 채팅을 진행할 수 있습니다.
    - join, leave, chat

## Tech Stack

- Golang(1.22), Gin, Gorilla, Redis


## TEST

```shell
make test
```


## BUILD

```shell
make build
```


## RUN
```shell
./live-chat-server
```

<br />

