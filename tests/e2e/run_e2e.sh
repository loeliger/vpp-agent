#!/bin/bash
set -euo pipefail

echo "preparing E2E test"
set -x

# compile vpp-agent
go build -v -o ./tests/e2e/vpp-agent.test \
  -ldflags "-X github.com/ligato/vpp-agent/vendor/github.com/ligato/cn-infra/agent.BuildVersion=TEST_E2E" \
  ./cmd/vpp-agent

# compile e2e test suite
go test -c -o ./tests/e2e/e2e.test ./tests/e2e

# start image
# TODO: do not run docker image with pid=host,
#  because any other vpp running now breaks tests
cid=$(docker run -d -it \
	-v $PWD/tests/e2e/e2e.test:/e2e.test:ro \
	-v $PWD/tests/e2e/vpp-agent.test:/vpp-agent:ro \
	-v $PWD/tests/e2e/grpc.conf:/etc/grpc.conf:ro \
	-v /var/run/docker.sock:/var/run/docker.sock \
	--label e2e.test="$*" \
	--pid="host" \
	--privileged \
	--env KVSCHEDULER_GRAPHDUMP=true \
	--env VPP_IMG="$VPP_IMG" \
	--env GRPC_CONFIG=/etc/grpc.conf \
	${DOCKER_ARGS-} \
	"$VPP_IMG" bash)

set +x

cleanup() {
	echo "stopping test container"
	set -x
	docker stop -t 2 "$cid" >/dev/null
	docker rm "$cid" >/dev/null
}

vppver=$(docker exec -i "$cid" dpkg-query -f '${Version}' -W vpp)

trap 'cleanup' EXIT

echo "============================================================="
echo -e " E2E TEST - VPP \e[1;33m${vppver}\e[0m"
echo "============================================================="

# run e2e test
if docker exec -i "$cid" /e2e.test -test.v $*; then
	echo >&2 "-------------------------------------------------------------"
	echo >&2 -e " \e[32mPASSED\e[0m (took: ${SECONDS}s)"
	echo >&2 "-------------------------------------------------------------"
	exit 0
else
	res=$?
	echo >&2 "-------------------------------------------------------------"
	echo >&2 -e " \e[31mFAILED!\e[0m (exit code: $res)"
	echo >&2 "-------------------------------------------------------------"

	# dump container logs
	logs=$(docker logs --tail 10 "$cid")
	if [[ -n "$logs" ]]; then
		echo >&2 -e "\e[1;30m${logs}\e[0m"
	fi

	exit $res
fi
