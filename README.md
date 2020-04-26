# Image Converter

## What for
Convert image to target type. ex: heic to jpg

## Usage
Easy to use as micro service in a system

## Building requirment
* docker

## Skills
* docker
* golang
* imagemagick

## Step by Step
The following is apply to local development or test
1. Build the image
```
cd <path-of-repo>
docker build -t image-converter ./
```

2. Start the image
```
docker run -p 80:80 -d -t image-converter
```

3. Open [Postman](https://www.postman.com/) or other tools can send `POST` method with `form-data`

4. Pass the params to endpoint `/image`
```
POST /image

form-data:
image: <the file you pick>
type: <the target type you want. ex: jpg>
```

5. Done! Receive the image return from the container
> If using Postman, the image will directly display in the response body

## Wish list
- [ ] Compress the image if too large
- [ ] Response image's url by keeping the image in the container
- [ ] Allow pdf to jpg
