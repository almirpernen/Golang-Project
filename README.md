Posts can have many likes from users, and users can like many posts. This is represented by the PostLike struct, which creates a many-to-many relationship between posts and users through a composite unique index on UserID and PostID.

### Comment - CommentLike (Many-to-Many)

Comments can have many likes, and users can like many comments, similarly managed by the CommentLike struct, indicating a many-to-many relationship between comments and users with a unique composite index on UserID and CommentID.

### User - User (Followers/Followings) (Many-to-Many)

This is a self-referencing many-to-many relationship where users can follow and be followed by many other users. The Followers and Followings slices in the User struct represent this relationship through a join table (`user_followers`). The many2many tag specifies the name of the join table, and joinForeignKey`/`JoinReferences tags specify the columns in the join table representing the following and follower users, respectively.

## API Structure

### Authentication

#### Signup

- Method: POST
- Endpoint: /signup
- Body:

   {
    "username": "your_username",
    "password": "your_password"
  }
  
#### Signin

- Method: POST
- Endpoint: /signin
- Body:
{
  "username": "your_username",
  "password": "your_password"
}

#### Users

Get Users List (protected)

- Method: GET
- Endpoint: /users

Get User by ID (protected)

- Method: GET
- Endpoint: /users/:id

Delete User by ID (protected)

- Method: DELETE
- Endpoint: /users/:id

Follow User (protected)

- Method: POST
- Endpoint: /users/:id/follow

Unfollow User (protected)

- Method: POST
- Endpoint: /users/:id/unfollow

#### Posts

Create post(protected)

- Method: POST
- Endpoint: /post
- Body:
{
  "content": "your_content"
}

List posts

- Method: GET
- Endpoint: /post

GET post by ID

- Method: GET
- Endpoint: /post/:id

Update post(protected)

- Method: PUT
- Endpoint: /post/:id
{
  "content": "your_updated_content"
}

Delete post(protected)

- Method: DELETE
- Endpoint: /post/:id

Like post(protected)

- Method: POST
- Endpoint: /post/:id/like

Unlike post(protected)

- Method: POST
- Endpoint: /post/:id/unlike

#### Comments

Create comment(protected)

- Method: POST
- Endpoint: /comment/:id
- Body:
{
  "content": "your_content"
}

List comments

- Method: GET
- Endpoint: /comment

Comment comment by ID

- Method: GET
- Endpoint: /comment/:id

Update comment(protected)

- Method: PUT
- Endpoint: /comment/:id
{
  "content": "your_updated_content"
}

Delete comment(protected)

- Method: DELETE
- Endpoint: /comment/:id

Like comment(protected)

- Method: POST
- Endpoint: /comment/:id/like

Unlike comment(protected)

- Method: POST
- Endpoint: /comment/:id/unlike
