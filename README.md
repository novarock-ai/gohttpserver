# gohttpserver

> NOTE: 本项目基于 `https://github.com/codeskyblue/gohttpserver` 进行了大量修改，不保证和原仓库功能完全兼容


- 基于 Vue2.x 以及 Golang 实现的高性能文件服务器
- 使用 streamlit 或对文件在线预览有需求的可以参考[该项目(streamlit-file-browser)](https://github.com/pragmatic-streamlit/streamlit-file-browser)


## Requirements
Tested with go-1.16

## Screenshots
![screen](./demo.gif)

## Features

1. [x] 文件基本的增删改查以及上传 
3. [x] zip 归档
4. [x] 虚拟滚动
5. [x] 增量搜索
6. [x] 支持精细到用户以及增,删,改,查等单独功能粒度的权限控制, 需要前置网关配合
7. [x] 自定义静态文件 CDN
8. [x] 安全 path 正则, 既可通过正则指定开发者希望或者不希望用户可以访问的 path
9. [x] PinRoot, 既自由指定 path 前缀


## Installation
当前请自行使用 golang 编译目标平台的二进制程序

## Usage
```bash
gohttpserver -r ./var -p 8000 --pin-root
```


## Docker Usage

```bash
# 该命令会自动编译 arm64 版本的镜像, 可以自己根据需求进行修改或添加
$ make package
```

## Advanced Feature
### PinRoot
```bash
gohttpserver -r ./your_path -p 8000 --pin-root
```
开启 `--pin-root` 后, 当访问某个 path 的文件系统时, 会以第一次访问的 path 作为 root path。
如当访问 `localhost:8000/nested_dir` 时会直接访问到 `root/nested_dir`，当未开启 `--pin-root` 时，点击回退按钮时会回退到 root 目录，而当开了 `--pin-root` 时，点击回退按钮则无法回退。


### Custom CDN
```bash
gohttpserver -r ./your_path -p 8000 --custom-cdn "https://xxxx"
```
指定此选项后，所有涉及到的静态资源都会走指定的 CDN，但必须保证改 CDN 上有全部依赖。

### SafeSymlinkPattern
```bash
gohttpserver -r ./your_path -p 8000 --safe-symlink-pattern "^/a/[^/]+/" --safe-symlink-pattern "^/b/[^/]+/"
```
支持指定匹配 path，可以传 List

### Authorization
为支持在项目运行时由业务侧动态控制不同用户或不同场景下的权限，因此该功能通过 query 传参的方式进行权限控制。如业务代码中可以通过如下方式创建服务器的访问链接，以在 streamlit 中使用为例，其他场景下原理类似，其中 `http://localhost:8000?access={access_all}` 就是要交给用户的服务器链接，其中对于权限通过 access 参数控制:
```python
    import streamlit as st
    from streamlit_file_browser import st_file_browser
    st.set_page_config(layout='wide')
    
    access_all = 0
    access_upload = 0b10000000
    access_delete = 0b01000000
    access_folder = 0b00100000
    access_download = 0b00010000
    access_archive = 0b00001000
    access_preview = 0b00000100
    
    access_all |= access_upload
    access_all |= access_delete
    access_all |= access_folder
    access_all |= access_download
    access_all |= access_archive
    access_all |= access_preview
    
    print(1111, access_all)

    st_file_browser('./',
        key="deep1",
        use_static_file_server=True,
        static_file_server_path=f'http://localhost:8000?access={access_all}',
    )
```

不过 query 是一种不安全的方式，因此建议在访问链接之前加一层鉴权网关，可以参考如下代码：
```python
import jwt
from fastapi import FastAPI, Request, Response
import httpx

app = FastAPI()

SECRET_KEY = "my_secret_key"

def decode_jwt(token: str):
    try:
        decoded_payload = jwt.decode(token, SECRET_KEY, algorithms=["HS256"])
        return decoded_payload
    except jwt.ExpiredSignatureError:
        return "Token has expired"
    except jwt.InvalidTokenError:
        return "Invalid token"


TARGET_URL = "http://localhost:8000" # gohttpserver path

@app.api_route("/{path:path}", methods=["GET", "POST", "PUT", "DELETE", "PATCH", "OPTIONS", "HEAD"])
async def proxy(request: Request, path: str):
    async with httpx.AsyncClient() as client:
        url = f"{TARGET_URL}/{path}"
        
        if not path.startswith("-/assets") and not path.startswith("-/sysinfo") and path != "favicon.ico":
            token = request.query_params.get("token")
            decode_result = decode_jwt(token)
            request_access = int(decode_result.get("access"))
            if not request_access:
                return Response(content="Access Denied", status_code=403)
            current_access = int(request.query_params.get("access"))
            if request_access != current_access:
                return Response(content="Access Denied", status_code=403)

        method = request.method
        data = await request.body()

        headers = dict(request.headers)
        headers.pop("host", None)

        response = await client.request(
            method=method,
            url=url,
            headers=headers,
            content=data,
            params=request.query_params
        )

        return Response(
            content=response.content,
            status_code=response.status_code,
            headers=dict(response.headers)
        )

if __name__ == "__main__":
    import uvicorn
    uvicorn.run(app, host="0.0.0.0", port=9999)
```
然后访问链接则需要加上 `token` 作为 query 参数：
```python
    import streamlit as st
    from streamlit_file_browser import st_file_browser
    st.set_page_config(layout='wide')
    
    access_all = 0
    access_upload = 0b10000000
    access_delete = 0b01000000
    access_folder = 0b00100000
    access_download = 0b00010000
    access_archive = 0b00001000
    access_preview = 0b00000100
    
    access_all |= access_upload
    access_all |= access_delete
    access_all |= access_folder
    access_all |= access_download
    access_all |= access_archive
    access_all |= access_preview
    
    print(1111, access_all)
    
    import jwt
    import datetime
    
    SECRET_KEY = "my_secret_key"
    
    def encode_jwt(payload: dict):
        token = jwt.encode(payload, SECRET_KEY, algorithm="HS256")
        return token


    payload = {
        "access": access_all, 
        "exp": datetime.datetime.now() + datetime.timedelta(hours=1)
    }
    file_server_token = encode_jwt(payload)

    st_file_browser('./',
        key="deep1",
        use_static_file_server=True,
        static_file_server_path=f'http://localhost:9999?access={access_all}&token={file_server_token}',
    )
```

更具体的使用方法可参考[该链接](https://github.com/pragmatic-streamlit/streamlit-file-browser/issues/38)

## LICENSE
This project is licensed under [MIT](LICENSE).
