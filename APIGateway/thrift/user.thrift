namespace go user

typedef i64 ID

struct UserRequest {
    1: ID userId (api.query = 'userId')
    2: string username (api.body = 'username')
    3: string token (api.header = 'token')
}

struct UserResponse {
    1: string text (api.body = 'text')
    2: i32 httpCode (api.http_code = 'true')
    3: string setCookie (api.cookie = 'setCookie')
}

service UserService {
    UserResponse GetUser(1: UserRequest req) (
        api.get = '/users/:userId',
        api.baseurl = 'example.com',
        api.param = 'true',
        api.serializer = 'json'
    )
    
    UserResponse UpdateUser(1: UserRequest req) (
        api.post = '/users/:userId',
        api.baseurl = 'example.com',
        api.param = 'true',
        api.serializer = 'json'
    )
}