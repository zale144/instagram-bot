import cv2
import numpy as np
import urllib.request
import sys

def url_to_image(path):
    req = urllib.request.Request(path, headers={'User-Agent': 'Mozilla/5.0'})
    r = urllib.request.urlopen(req)

    image = np.asarray(bytearray(r.read()), dtype="uint8")
    image = cv2.imdecode(image, cv2.COLOR_BGR2GRAY)
 
	# return the image
    return image
# Get user supplied values
imagePath = sys.argv[1]
cascPath = "haarcascade_frontalface_default.xml"

# Create the haar cascade
faceCascade = cv2.CascadeClassifier(cascPath)

img = url_to_image(imagePath)

# Detect faces in the image
faces = faceCascade.detectMultiScale(
    img,
    scaleFactor=1.1,
    minNeighbors=5,
    minSize=(30, 30)
    #flags = cv2.CV_HAAR_SCALE_IMAGE
)

print(len(faces))

# Draw a rectangle around the faces
#for (x, y, w, h) in faces:
#    cv2.rectangle(img, (x, y), (x+w, y+h), (0, 255, 0), 2)

#cv2.imshow("Faces found", img)
#cv2.waitKey(0)