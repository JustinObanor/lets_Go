API for a students dormitory.
This README explains for only a student path. But its the same for provisions, room, and workers.
This will be an application programming interface that allows users to Create, Read, Update, and Delete specific resources.

## Idea
You can make a table or anything you prefer, that will show all available students/provisions/room/workers in the database. Each row of a particular resource, should have options for viewing, editting, or deleting it

In short, a table of resources(students, provisions, room, workers). In each row, i can read, edit, or delete it

Somewhere on the screen, maybe below should be a place for creating a new resource, or for signup

## Types of users
```

-admin          can perform any operation
-user           only with correct credentials which are in http Authorization header

```


Admins have all rights.
Only users with correct credential in the http Authorization header are allowed to Create, Update or Delete a specific resource. 


## Post Signup
User need to signup first
There should be a place where user can input -
```
                {
                     "username": "string",
                     "password": "string"
                }
``` 
And then submit

```
-url            "https://localhost:8080/signup" 
```

In this page, the user is required to pass in a username and a password in the request body

Feel free to add any other design you like

Also make a page for error when username already exists

After signup, the user will see all the available students


## Post student 
This is for making a new student
We use Basic HTTP Authorization header(Basic + base64)

This is passed in the request header as - 
```
Basic oWVzflDrjRHBVuH0I=
```

Also make a page for incorrect credentials(i think)

If credentials are correct, there should be a place where user can input -
```
            {
            "ID": 0,
            "firstName": "string",
            "lastName": "string",
            "UUID": 0,
            "studRoom": {
                "id": 0,
                "room": 0
            },
            "studFloor": {
                "id": 0,
                "floor": 0
            }
            }
```

```
-url            "https://localhost:8080/student" 
```

## Post provision
```
{
  "id": 0,
  "bedhseet": 0,
  "pillow": 0,
  "towel": 0,
  "blanket": 0,
  "curtain": 0
}
```
```
-url            "https://localhost:8080/provision" 
```

## Post room
```
{
  "id": 0,
  "room": 0,
  "chairs": 0,
  "tables": 0,
  "shelves": 0
}
```

```
-url            "https://localhost:8080/room" 
```

## Post worker
```
{
  "id": 0,
  "firstName": "string",
  "lastName": "string",
  "workFloor": {
    "id": 0,
    "Floor": {
      "floor": 0,
      "code": 0
    }
  },
  "workDays": "string"
}
```

```
-url            "https://localhost:8080/worker" 
```


User should pass in the id, firstname, lastname.... and then submit it

Feel free to add any other design you like

## Get all student
Ths should show all the students(the table of students)
No Authorization required
This url list all the students available in the database


```
-url            "https://localhost:8080/student" 
```

```
[
            {
            "ID": 1,
            "firstName": "string",
            "lastName": "string",
            "date": "string",
            "UUID": 1,
            "studRoom": {
                "id": 1,
                "room": 0
            },
            "studFloor": {
                "id": 1,
                "floor": 0
            }
            },

            {
            "ID": 2,
            "firstName": "string",
            "lastName": "string",
            "date": "string",
            "UUID": 2,
            "studRoom": {
                "id": 2,
                "room": 0
            },
            "studFloor": {
                "id": 2,
                "floor": 0
            }
            }
]            
```


Feel free to add any other design you like



## Get student
This should show one student
No Authorization required

If student exists, then we see page of student

```
            {
            "ID": 1,
            "firstName": "string",
            "lastName": "string",
            "date": "string",
            "UUID": 1,
            "studRoom": {
                "id": 1,
                "room": 0
            },
            "studFloor": {
                "id": 1,
                "floor": 1
            }
            }
```

```
-url            "https://localhost:8080/student/1" 
```

## Get provision
```
{
  "id": 0,
  "bedhseet": "string",
  "pillow": "string",
  "towel": "string",
  "blanket": "string",
  "curtain": "string"
}
```

## Get room
```
{
  "id": 0,
  "room": "string",
  "chairs": "string",
  "tables": "string",
  "shelves": "string"
}
```

## Get worker
```
{
  "id": 1,
  "firstName": "Micheal",
  "lastName": "Angelo",
  "workFloor": {
    "id": 1,
    "floor": {
      "floor": 13,
      "code": 13
    }
  },
  "workDays": "15-10-1999"
}
```


User should pass in the id

Feel free to add any other design you like


## Put student
Authorization required

We us Basic HTTP Authorization header(Basic base64)

This is passed in the request header as - 
```
Basic oWVzflDrjRHBVuH0I=
```

```
            {
            "ID": 1,
            "firstName": "string",
            "lastName": "string",
            "studRoom": {
                "id": 1,
                "room": 0
            },
            "studFloor": {
                "id": 1,
                "floor": 1
            }
            }
```

```
-url            "https://localhost:8080/student/1" 
```

Feel free to add any other design you like

## Put provision
```
{
  "id": 0,
  "bedhseet": 0,
  "pillow": 0,
  "towel": 0,
  "blanket": 0,
  "curtain": 0
}
```

## Put room
```
{
  "id": 0,
  "room": 0,
  "chairs": 0,
  "tables": 0,
  "shelves": 0
}
```

## put worker
```
{
  "id": 0,
  "firstName": "string",
  "lastName": "string",
  "workFloor": {
    "id": 0,
    "floor": {
      "floor": 0,
      "code": 0
    }
  },
  "workDays": "string"
}
```


## Delete student
Authorization required

We us Basic HTTP Authorization header(Basic base64)

This is passed in the request header as - 
```
Basic oWVzflDrjRHBVuH0I=
```


```
-url            "https://localhost:8080/student/1" 
```

Feel free to add any other design you like




