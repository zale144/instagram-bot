from werkzeug.wrappers import Request, Response
from werkzeug.serving import run_simple
from jsonrpc import JSONRPCResponseManager, dispatcher

import face_detect
import os

@Request.application
def application(request):
    dispatcher["Num.Faces"] = lambda s: face_detect.number_of_faces(s["url"])
    response = JSONRPCResponseManager.handle(request.data, dispatcher)
    return Response(response.json, mimetype='application/json')

if __name__ == '__main__':
    host = os.environ['RPC_URI']
    run_simple(host, 4000, application)
