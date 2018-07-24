import cv2
import numpy as np
import urllib.request

def url_to_image(path):
    req = urllib.request.Request(path, headers={'User-Agent': 'Mozilla/5.0'})
    r = urllib.request.urlopen(req)

    image = np.asarray(bytearray(r.read()), dtype="uint8")
    image = cv2.imdecode(image, cv2.COLOR_BGR2GRAY)
    # return the image
    return image

def number_of_faces(url):
    cascPath = "haarcascade_frontalface_default.xml"
    # create the haar cascade
    faceCascade = cv2.CascadeClassifier(cascPath)
    img = url_to_image(url)

    # detect faces in the image
    faces = faceCascade.detectMultiScale(
        img,
        scaleFactor=1.1,
        minNeighbors=5,
        minSize=(30, 30)
        )
    num_of_faces = len(faces)
    print(num_of_faces)
    return num_of_faces 
