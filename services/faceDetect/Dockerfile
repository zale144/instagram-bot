FROM kavolorn/opencv

COPY . .

RUN apt-get update \
    && apt-get install -y curl \
    && curl -O https://bootstrap.pypa.io/get-pip.py \
    && python3 get-pip.py \
    && pip install werkzeug \
    && pip install json-rpc \
    && pip install requests
