# Instagram Bot

Instagram Bot is a concept app that demonstrates how you can promote your website/platform on Instagram.

Initially I have written it as a monolithic web application, but I've rewritten it to follow the Microservice architecture.


## Features

- [Send invites to your website to Instagram users](#send-invites-to-your-website-to-instagram-users)
- [Batch process users based on a hashtag](#batch-process-users-based-on-a-hashtag)
- [Follow Instagram users](#follow-instagram-users)
- [No account creation, login with your Instagram credentials](#login-with-your-instagram-credentials)
- [Uses face detection to select a suitable background image](#background-image-selected-by-face-detection)

## Technologies/libraries used

- [Golang](https://golang.org)
- [Micro](https://micro.mu)
- [gRPC](https://grpc.io)
- [Kubernetes](https://kubernetes.io)
- [Docker](https://www.docker.com)
- [My fork of Goinsta](https://github.com/zale144/goinsta)
- [Open Graph](http://ogp.me)
- [Labstack Echo](https://echo.labstack.com)
- [Python](https://www.python.org)
- [OpenCV](https://opencv.org)
- [wkhtmltopdf](https://wkhtmltopdf.org)
- [PostgreSQL](https://www.postgresql.org)
- [Consul (before moving to Kubernetes)](https://www.consul.io)
- [React.js](https://reactjs.org)
- [Gorm](http://gorm.io)
- [dep](https://golang.github.io/dep)
- [Bootstrap](https://getbootstrap.com)


## Send invites to your website to Instagram users

Once you login you will see a list of all Instagram users that you follow. You can use the search bar to filter
them. After clicking on a specific user, they will receive a Direct Message on Instagram with a link to your (TODO)
website registration page, and in the link an OG tag will be embedded, containing an image of what their profile
on your website would potentially look like. It will show their pictures and basic info as it is on their Instagram
profile.

To achieve this,
- A request is sent to the **api** microservice to process the user by username
- The **api** microservice sends a request to **htmltoimage** microservice to create an image of your
    HTML template with that user's pictures and info applied to it.
- **htmltoimage** does this by executing the **wkhtmltoimage** binary and passing
    the image parameters and the URL to the fake HTML profile as input
- **wkhtmltoimage** consumes the URL, runs a detached browser to render the fake HTML profile, and 
    generates an image out of it
- The fake HTML profile will be served from the **web** microservice.
- When requested, **web** composes that page by applying user info it fetches from **session**
    to an HTML template (in the future this template should be customizable)
- **session** gets basic info from the user's Instagram profile, like the profile image, biography,
    username, full name, media feed, from Instagram API, via **goinsta**.
- **session** selects the background image according to criteria: it's landscape, it's a picture of the user,
    or at least it's a picture of a person
- To detect if it's a picture of a person, it calls the **facedetect** microservice with a link to the image.
- The image created by **htmltoimage** is returned to **api** and saved as a file.
- Then, **api** continues by calling **session** to send an Instagram Direct message to the user with the link
    to a page that should (TODO) redirect to your registration page, with that generated image of the
    fake profile as an *OG tag*, which will show up nicely in the message.
- This page will again be served from the **web** service.

## Batch process users based on a hashtag

By navigating to the **Jobs** tab, you will be able to create and view batch processing jobs.
Issuing a job requires you to provide an Instagram **#hashtag**, limit and a message.
All users associated with the **hashtag**, within limit provided, will then receive a Direct message
on Instagram with the invitation to your website, with a custom message. Meaning, they will be processed
just like [this](#send-invites-to-your-website-to-instagram-users).

To achieve this,
- A request is sent to **api**, to create a Job.
- The job will be saved to the database and run in the background, so you don't have to wait around.
- The job will: 
    * Request **session** to get all users by hashtag within provided limit.
    * [Process](#send-invites-to-your-website-to-instagram-users) them each by each if they fulfill two conditions: they're not processed yet, the user is not you.
    * Save the processed users the the database


## Follow Instagram users

The search bar can be used to find Instagram users that are not in your followed list, and if found, you can
follow them by clicking on the *Follow* button.

## Login with your Instagram credentials

In order to login, simply provide your Instagram account credentials (they're not secretly saved to a database).
The login request happens like this,

- A request is sent to **web** service to create a cookie for you, and you will gain access to the app's main page.
- **web** sends a call to **session** to create a session for you
- **session** calls the Instagram API to log you in and cache the session
- **web** calls **api** to create a JWT token for you
- The JWT token is then saved to your browser's local storage. This is how you will have authorization
    to access the **api** as a web service
- If successful, you will be redirected to the app's main page

## Background image selected by face detection

Besides the username, full name, profile picture, biography, **session** gets also user's media from the Instagram API.
A background image for the fake profile will be chosen from those media. To choose a suitable background image, 
criteria will be respected, such as: 
- It's a landscape image, 
- It's a picture of the user (if user is tagged in the picture)
- It's a picture of a person. 

To detect faces in the image,

- **session** calls the **facedetect** microservice by RPC with a URL to the image as a parameter.
- **facedetect** runs a simple RPC server written in Python and uses machine learning with OpenCV 
    for detecting faces in images.
- **facedetect** returns the number of faces it found in the image from the provided URL.
    
