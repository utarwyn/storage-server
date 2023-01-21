<h1 align="center">Storage Server</h1>

<h4 align="center">
Lightweight and production-ready service to serve your static files.
<br>
Inspired from <a href="https://github.com/PierreZ/goStatic">goStatic</a>.
</h4>

<p align="center">
   <a href="https://hub.docker.com/r/utarwyn/storage-server">
      <img src="https://img.shields.io/docker/pulls/utarwyn/storage-server.svg" alt="Docker Pulls">
   </a>
   <a href="https://github.com/utarwyn/storage-server/actions/workflows/publish.yml">
      <img src="https://img.shields.io/github/actions/workflow/status/utarwyn/storage-server/publish.yml?label=docker%20build" alt="Docker Build status">
   </a>
   <a href="https://github.com/utarwyn/storage-server/blob/main/LICENSE">
      <img src="https://img.shields.io/github/license/utarwyn/storage-server" alt="License">
   </a>
</p>

"Storage Server" is a **self-hosted service** to serve your files for everyone. This package is intended to be
as light as possible, without dependency and very simple. It includes **key features** we can expect from a web
service, like caching, compression, logging and security. All configurable with ease!

You have a problem with the service or want to have a new feature? Feel free to open an issue.


Key features
------------

- Serve your files with an embedded static web server
- Upload/Delete files through protected HTTP routes
- Enable caching/CORS policy on specific files
- Expose files list of specific directories as a json
- Automatic GZIP compression of sent files
- Light and efficient as possible
- No dependency :tada:

Installation
------------

This service is designed to be used as a Docker image. All major platforms are supported. \
So, you have to use the image `utarwyn/storage-server`.

You can also build the executable from source code, just like a normal Golang program.


Usage
-----

Example on how to run it:

```
docker run -p 80:8080 -e PORT=8080 -v path/to/files:/srv/http utarwyn/storage-server
```

Available options:

| Command parameter   | Env variable        | Default   | Description                                  |
|---------------------|---------------------|-----------|----------------------------------------------|
| port                | PORT                | 8043      | Listening port                               |
| base-path           | BASE_PATH           | /srv/http | Directory where files are stored             |
| client-secret       | CLIENT_SECRET       | -         | Secret key used to access privileged routes  |
| enable-logging      | ENABLE_LOGGING      | false     | Enable log request                           |
| caching-directories | CACHING_DIRECTORIES | -         | List of directories to cache                 |
| expose-directories  | EXPOSE_DIRECTORIES  | -         | List of directories to expose as a json file |
| allow-origins       | ALLOW_ORIGINS       | -         | List of origins to allow using CORS policy   |

License
-------

"Storage Server" is open-sourced software licensed under the [MIT license][1].

---
> GitHub [@utarwyn][2] &nbsp;&middot;&nbsp; Twitter [@Utarwyn][3]


[1]: https://github.com/utarwyn/storage-server/blob/main/LICENSE

[2]: https://github.com/utarwyn

[3]: https://twitter.com/Utarwyn
