namespace go http

struct BizRequest {
    2: string text(api.body = 'text')
    3: i32 token(api.header = 'token')
    6: optional list<string> items(api.query = 'items')
}

struct BizResponse {
    1: i32 token(api.header = 'token')
    2: string text(api.body = 'text')
    5: i32 http_code(api.http_code = '') 
}

service BizService {
    BizResponse BizMethod1(1: BizRequest req)(api.get = '/bizservice/method1', api.baseurl = 'example.com', api.param = 'true')
    BizResponse BizMethod2(1: BizRequest req)(api.post = '/bizservice/method2', api.baseurl = 'example.com', api.param = 'true', api.serializer = 'form')
    BizResponse BizMethod3(1: BizRequest req)(api.post = '/bizservice/method3', api.baseurl = 'example.com', api.param = 'true', api.serializer = 'json')
}