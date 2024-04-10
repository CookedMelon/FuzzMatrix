#!/usr/bin/env python3
import aiohttp
from aiohttp import web
import os
import random
import sys
import json

MIN_PATH = '/out/min_queue'

async def handle_get(request):
    path = request.match_info.get('path', "")
    if path == "/out/queue/seeds":
        dic = {}
        counter = 0
        files = os.listdir(os.curdir + MIN_PATH)
        for file in files:
            if counter == 3:
                break
            if not os.path.isdir(file):
                with open(os.curdir + MIN_PATH + "/" + file, 'rb') as f:
                    try:
                        file_context = f.read()
                        dic[counter] = str(file_context)
                        counter += 1
                    except Exception as e:
                        print(str(e))
        print(dic)
        return web.Response(text=json.dumps(dic), content_type='application/json')
    # 处理特定路径 'block_freq.json?_=...'
    else:
        # 如果path为static.json或以out/fuzzer_stats开头或以out/block_freq.json开头
        # if path == 'static.json' or path.startswith('out/fuzzer_stats') or path.startswith('out/block_freq.json'):
        #     return web.FileResponse(os.curdir + '/'+path)
        # # 其余的404
        # else:
        #     return web.Response(status=404)

        # 为了方便测试，直接返回路径下全部文件
        return web.FileResponse(os.curdir + '/'+path)

async def handle_post(request):
    data = await request.read()
    path = os.curdir + '/out/queue/' + ''.join([str(random.randint(0, 9)) for _ in range(10)])
    with open(path, 'wb') as f:
        f.write(data)
    return web.Response(status=200)

async def cors_middleware(app, handler):
    async def cors_handler(request):
        response = await handler(request)
        response.headers['Access-Control-Allow-Origin'] = '*'
        return response
    return cors_handler

app = web.Application()
app.router.add_get('/{path:.*}', handle_get)
app.router.add_post('/queue', handle_post)
app.middlewares.append(cors_middleware)

if __name__ == '__main__':
    port = int(sys.argv[1]) if len(sys.argv) > 1 else 8000
    web.run_app(app, port=port)
