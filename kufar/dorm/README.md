API for a students dormitory.
This README explains for only a student. But its the same for provisions, room, and workers
This will be an application programming interface that allows users to Create, Read, Update, and Delete specific resources.

## Idea
You can make a table or anything you prefer, that will show all available students in the database. Each row of a particular student, should have options for viewing the particular student, editting the particular student, or deleting it

Somewhere on the screen, maybe below should be a place for creating a new student, or for signup

## Types of users
```

-admin          can perform any operation
-user           only with correct credentials which are in http Authorization header

```

```

Admins have all rights.
Only users with correct credential in the http Authorization header are allowed to Create, Update or Delete a specific resource. 

```

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
We use Basic HTTP Authorization header(Basic base64)

This is passed in the request header as - 
```
Basic oWVzflDrjRHBVuH0I=
```

Also make a page for incorrect credentials

If credentials are correct, there should be a place where user can input -
```
            {
            "ID": 0,
            "Firstname": "string",
            "Lastname": "string",
            "Date": "string",
            "UUID": 0,
            "StudRoom": {
                "id": 0,
                "room": 0
            },
            "StudFloor": {
                "id": 0,
                "floor": 0
            }
            }
```

```
-url            "https://localhost:8080/student" 
```

User should pass in the id, firstname, lastname.... and then submit it

Feel free to add any other design you like

## Get all student
Ths should show all the students
No Authorization required
This url list all the students available in the database


```
-url            "https://localhost:8080/student" 
```

```
[
            {
            "ID": 1,
            "Firstname": "string",
            "Lastname": "string",
            "Date": "string",
            "UUID": 1,
            "StudRoom": {
                "id": 1,
                "room": 0
            },
            "StudFloor": {
                "id": 1,
                "floor": 1
            }
            },

            {
            "ID": 1,
            "Firstname": "string",
            "Lastname": "string",
            "Date": "string",
            "UUID": 1,
            "StudRoom": {
                "id": 1,
                "room": 0
            },
            "StudFloor": {
                "id": 1,
                "floor": 1
            }
            }
]            
```


Feel free to add any other design you like



## Get student
This should show one student
No Authorization required

The user should specify the id of the student required. example id - 1
If student exists, then we see page of student

```
            {
            "ID": 1,
            "Firstname": "string",
            "Lastname": "string",
            "Date": "string",
            "UUID": 1,
            "StudRoom": {
                "id": 1,
                "room": 0
            },
            "StudFloor": {
                "id": 1,
                "floor": 1
            }
            }
```

```
-url            "https://localhost:8080/student/1" 
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
            "Firstname": "string",
            "Lastname": "string",
            "Date": "string",
            "UUID": 1,
            "StudRoom": {
                "id": 1,
                "room": 0
            },
            "StudFloor": {
                "id": 1,
                "floor": 1
            }
            }
```

```
-url            "https://localhost:8080/student/1" 
```

Feel free to add any other design you like



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




