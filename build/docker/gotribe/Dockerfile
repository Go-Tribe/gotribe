# Copyright 2024 Innkeeper GoTribe <https://www.gotribe.cn>. All rights reserved.
# Use of this source code is governed by a MIT style
# license that can be found in the LICENSE file. The original repo for
# this file is https://www.gotribe.cn


FROM BASE_IMAGE
LABEL maintainer="<my@dengmengmian.com>"

WORKDIR /opt/gotribe

RUN ln -sf /usr/share/zoneinfo/Asia/Shanghai /etc/localtime && \
      echo "Asia/Shanghai" > /etc/timezone && \
      mkdir -p /var/log/app

COPY gotribe /opt/gotribe/bin/

ENTRYPOINT ["/opt/gotribe/bin/gotribe"]
CMD ["-c", "/etc/gotribe/gotribe.yaml"]
