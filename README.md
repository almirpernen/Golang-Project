# Final Project Golang 2024 Spring: Go Blog with Comments and Authentication

This project is based on Golang, and target audience of this app is drivers who loves their cars, and wanto lead their social media.

## Authors


 - Izbassarov Arman 22B030577
 - Pernen Almir 22B030577


## Database Models Overview
<img width="712" alt="Снимок экрана 2024-05-16 в 23 32 46" src="https://github.com/almirpernen/eshop/assets/123065546/5793b9a6-8921-48b7-b9b5-947e917576d8">


### User Model

- gorm.Model: Inherits fields ID, CreatedAt, UpdatedAt, DeletedAt from GORM's base model.
- Username: string, stores the username of the user.
- Password: string, stores the encrypted password (not exported in JSON).
- Posts: Slice of Post, represents a one-to-many relationship with Post (A user can have many posts). Uses UserID as the foreign key.
- Comments: Slice of Comment, represents a one-to-many relationship with Comment (A user can author many comments). Uses UserID as the foreign key.
- Followers: Slice of User, represents a many-to-many relationship with other User entities, indicating users who follow this user. Utilizes a join table user_followers, with FollowingID as the join foreign key and FollowerID as the join reference.
- Followings: Slice of User, represents a many-to-many relationship with other User entities, indicating users this user follows. Utilizes the same join table user_followers, with FollowerID as the join foreign key and FollowingID as the join reference.

### Post Model

- gorm.Model: Inherits fields ID, CreatedAt, UpdatedAt, DeletedAt.
- Content: string, stores the content of the post.
- UserID: uint, foreign key linking back to the User who authored the post.
- User: User, represents the many-to-one relationship with User.
- Comments: Slice of Comment, represents a one-to-many relationship with Comment (A post can have many comments). Uses PostID as the foreign key.
- LikesCount: int, not a database field (`gorm:"-"`) but used to store the count of likes a post has received.

### Comment Model

- gorm.Model: Inherits fields ID, CreatedAt, UpdatedAt, DeletedAt.
- Content: string, stores the content of the comment.
- UserID: uint, foreign key linking back to the User who authored the comment.
- PostID: uint, foreign key linking to the Post the comment belongs to.
- LikesCount: int, not a database field (`gorm:"-"`) but used to store the count of likes a comment has received.

### PostLike Model

- UserID: uint, part of a composite unique index with PostID, represents the user who liked the post.
- PostID: uint, part of a composite unique index with UserID, represents the post that was liked.

### CommentLike Model

- UserID: uint, part of a composite unique index with CommentID, represents the user who liked the comment.
- CommentID: uint, part of a composite unique index with UserID, represents the comment that was liked.

## Database Relationships

### User - Post (One-to-Many)

A user can have many posts. This is represented by the Posts slice in the User struct, with a foreignKey tag pointing to UserID in the Post struct, indicating that multiple posts can belong to a single user.

### Post - Comment (One-to-Many)

A post can have many comments. This relationship is shown by the Comments slice in the Post struct, with a foreignKey tag pointing to PostID in the Comment struct, indicating that multiple comments can be associated with a single post.

### User - Comment (One-to-Many)

A user can have many comments. This is represented similarly to posts, where the Comments slice in the User struct has a foreignKey pointing to UserID in the Comment struct, indicating a user can author multiple comments.

### Post - User (Many-to-One)

Many posts belong to one user. This inverse relationship of the first point is represented by the User field in the Post struct, which points back to the owning user. The foreignKey:UserID indicates the association's direction.

### Comment - User (Many-to-One)

Many comments belong to one user. Similar to posts, this is the inverse relationship of the third point, where each Comment struct has a User field pointing back to the commenter.

### Post - PostLike (Many-to-Many)

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
``` json

   {
    "username": "your_username",
    "password": "your_password"
  }
```
#### Signin

- Method: POST
- Endpoint: /signin
- Body:
``` json

{
  "username": "your_username",
  "password": "your_password"
}
```

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
```json
{
  "content": "your_content"
}
```

List posts

- Method: GET
- Endpoint: /post

GET post by ID

- Method: GET
- Endpoint: /post/:id

Update post(protected)

- Method: PUT
- Endpoint: /post/:id

``` json
{
  "content": "your_updated_content"
}
```

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
``` json

{
  "content": "your_content"
}
``` 

List comments

- Method: GET
- Endpoint: /comment

Comment comment by ID

- Method: GET
- Endpoint: /comment/:id

Update comment(protected)

- Method: PUT
- Endpoint: /comment/:id
``` json

{
  "content": "your_updated_content"
}
```
Delete comment(protected)

- Method: DELETE
- Endpoint: /comment/:id

Like comment(protected)

- Method: POST
- Endpoint: /comment/:id/like

Unlike comment(protected)

- Method: POST
- Endpoint: /comment/:id/unlike
