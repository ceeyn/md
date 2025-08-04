#!/bin/bash
# app name
ROOT_DIR=`pwd`
RELEASE_NAME=`git remote -v | grep fetch | awk '{print $2}' | tr '/' '\n' | tail -n 1 | sed 's/.git//'`
if [ -n "${BUILD_IMAGE}" ]; then
    RELEASE_NAME=${BUILD_IMAGE}
fi
if [ ! -n "${RELEASE_NAME}" ]; then
    RELEASE_NAME=${ROOT_DIR##*/}
    has=`echo "${RELEASE_NAME}" | grep qci`
    if [ -n "${has}" ]; then
        echo "[ERROR] please set the environment variable: BUILD_IMAGE "
        exit 1;
    fi
fi
echo "[INFO] app name: ${RELEASE_NAME} "

# mac 不执行 写profile，流水线 必须用的，BUILD_IMAGE 环境变量
mac=`uname | grep Darwin`
if [ ! -n "${mac}" ]; then
    echo "BUILD_IMAGE=${RELEASE_NAME}" >> $QCI_ENV_FILE
fi

# app config env
CONF_ENV="${ENV}"
if [ ! -n "${CONF_ENV}" ]; then
  echo "[ERROR] please set the environment variable: ENV"
  exit 1
fi
if [ ! -e ${ROOT_DIR}/conf/${CONF_ENV} ]; then
    echo "[ERROR] env conf dir ${ROOT_DIR}/conf/${CONF_ENV} not found ..."
    exit 1
fi
echo "[INFO] app env: ${CONF_ENV} "

# build dir
BUILD_DIR=${ROOT_DIR}
echo "[INFO] build dir:  ${BUILD_DIR} "

# app
APP=${BUILD_DIR}/${RELEASE_NAME}
if [ -e $APP ]; then
    rm -rf $APP
fi

# build flag
flags=""
flags=${flags}" -X '${RELEASE_NAME}/config.gitversion=$(git rev-parse --short HEAD)'"
flags=${flags}" -X '${RELEASE_NAME}/config.buildstamp=$(date +%s)'"
flags=${flags}" -X '${RELEASE_NAME}/config.goversion=$(go version)'"
flags=${flags}" -X '${RELEASE_NAME}/config.sysuname=$(uname -mns)'"
flags=${flags}" -X '${RELEASE_NAME}/config.gitbranch=$(git branch | grep "\*" | grep " .*" -o | sed "s/\ //g")'"
flags=${flags}" -X google.golang.org/protobuf/reflect/protoregistry.conflictPolicy=warn"
#flags=${flags}" -X github.com/gogo/protobuf/reflect/protoregistry.conflictPolicy=warn"
#flags=${flags}" -X github.com/golang/protobuf/reflect/protoregistry.conflictPolicy=warn"

# go build
GO111MODULE=on go build -ldflags "$flags" -o $RELEASE_NAME
if [ ! -e $APP ]; then
    echo "[ERROR] go build ${APP} failed ..."
    exit 1
else
    echo "[INFO] go build ${APP} success ..."
fi

# build release path
BUILD_PRO_DIR=${BUILD_DIR}/release/${RELEASE_NAME}

# 目录结构
mkdir -p $BUILD_PRO_DIR
mkdir -p $BUILD_PRO_DIR/bin
mkdir -p $BUILD_PRO_DIR/sbin
mkdir -p $BUILD_PRO_DIR/conf
mkdir -p $BUILD_PRO_DIR/log

# copy conf
cp -r ${ROOT_DIR}/conf/${CONF_ENV}/* $BUILD_PRO_DIR/conf/

# copy app
cp ${APP} $BUILD_PRO_DIR/bin/

# copy sbin
cp -r ${ROOT_DIR}/sbin/* $BUILD_PRO_DIR/sbin/
chmod 777 $BUILD_PRO_DIR/sbin/*

# 更新 织云包
if [ "${ENABLE_EVN_ZHIYUN}"x = "True"x ]; then
    OPERATION_NAME='IM_Lconn'
    cd ${BUILD_DIR}/release/
    tar -zcvf 'new_package.tar.gz' *
    echo "[INFO] zhiyun release name: ${RELEASE_NAME}_${ENV} , product: ${OPERATION_NAME} "
    echo "[INFO] zhiyun URL: https://yun.isd.com/index.php/package/versions/?product=${OPERATION_NAME}&package=${RELEASE_NAME}_${ENV}"
    qci-plugin zhiyun_simple_submit --clean "true" --version "DEFAULT" --product "${OPERATION_NAME}" --name "${RELEASE_NAME}_${ENV}" --description "[QCI] branch:${QCI_REPO_BRANCH}, tag: ${QCI_REPO_TAG}, git_commit_message:${QCI_COMMIT_MESSAGE}" --tarball "${BUILD_DIR}/release/new_package.tar.gz"
    if [[ $? != 0 ]]; then
        echo "[ERROR] zhiyun submit failed ..."
        exit 1
    fi
fi

echo "[INFO] build shell finish ..."
