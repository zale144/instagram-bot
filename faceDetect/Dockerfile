 
# Dockerfile
FROM kavolorn/opencv

COPY . .

RUN apt-get update
RUN apt-get install -y curl
RUN curl -O https://bootstrap.pypa.io/get-pip.py
RUN python3 get-pip.py

ENTRYPOINT ["python3","-u","face_detect.py"]
