# AutoUpdate

本工具作用于服务器或者需要自动更新相应工具的地方，只需要将更新地址放入AutoUpdateLink.txt文件并指定下载目录即可
选项:

```shell
    ./autoupdate -link=链接文件的地址 -dir=下载目录
```

例如

```shell
    ./autoupdate -link=/tmp/link -dir=/home/downloads
```

此外，如果需要更新的文件链接会变动，可以手动加入规则，并在main函数中调用下载即可
