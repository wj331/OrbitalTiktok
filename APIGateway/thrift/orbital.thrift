namespace go orbital

struct Person {
    1: string name,
    2: i32 age,
}

struct Request {
	1: string message
}

struct Response {
	1: string message
}

service PeopleService {
    Person editPerson(1: Person person)
    Response echo(1: Request req)
}
