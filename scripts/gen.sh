#!/bin/bash

THRIFT_PACKAGE_PREFIX=kylexu.com/thrift_client_pool/api
THRIFT_IDL=../idl/add.thrift
thrift -r --gen go:ignore_initialisms,package_prefix=${THRIFT_PACKAGE_PREFIX} -out ../ ${THRIFT_IDL}
