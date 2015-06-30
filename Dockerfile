FROM progrium/busybox
MAINTANER amir@scalock.com
ADD bin/batten /
CMD ["/batten", "check"]