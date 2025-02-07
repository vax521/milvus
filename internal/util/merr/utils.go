// Licensed to the LF AI & Data foundation under one
// or more contributor license agreements. See the NOTICE file
// distributed with this work for additional information
// regarding copyright ownership. The ASF licenses this file
// to you under the Apache License, Version 2.0 (the
// "License"); you may not use this file except in compliance
// with the License. You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

// errutil package provides utility for errors handling.
package merr

import (
	"context"
	"fmt"
	"strings"

	"github.com/cockroachdb/errors"
	"github.com/milvus-io/milvus-proto/go-api/commonpb"
)

// Declare a success status, avoid create it every time
var successStatus = &commonpb.Status{}

// Code returns the error code of the given error,
// WARN: DO NOT use this for now
func Code(err error) int32 {
	if err == nil {
		return 0
	}

	cause := errors.Cause(err)
	switch cause := cause.(type) {
	case milvusError:
		return cause.code()

	default:
		if errors.Is(cause, context.Canceled) {
			return CanceledCode
		} else if errors.Is(cause, context.DeadlineExceeded) {
			return TimeoutCode
		} else {
			return errUnexpected.code()
		}
	}
}

func IsRetriable(err error) bool {
	return Code(err)&retriableFlag != 0
}

// Status returns a status according to the given err,
// returns Success status if err is nil
func Status(err error) *commonpb.Status {
	if err == nil {
		return successStatus
	}

	return &commonpb.Status{
		Code:   Code(err),
		Reason: err.Error(),
		// Deprecated, for compatibility, set it to UnexpectedError
		ErrorCode: commonpb.ErrorCode_UnexpectedError,
	}
}

// Error returns a error according to the given status,
// returns nil if the status is a success status
func Error(status *commonpb.Status) error {
	if status.GetCode() == 0 && status.GetErrorCode() == commonpb.ErrorCode_Success {
		return nil
	}

	// use code first
	code := status.GetCode()
	if code == 0 {
		return newMilvusError(fmt.Sprintf("legacy error code:%d, reason: %s", status.GetErrorCode(), status.GetReason()), errUnexpected.errCode, false)
	}

	return newMilvusError(status.GetReason(), code, code&retriableFlag != 0)
}

// Service related
func WrapErrServiceNotReady(stage string, msg ...string) error {
	err := errors.Wrapf(ErrServiceNotReady, "stage=%s", stage)
	if len(msg) > 0 {
		err = errors.Wrap(err, strings.Join(msg, "; "))
	}
	return err
}

func WrapErrServiceUnavailable(reason string, msg ...string) error {
	err := errors.Wrap(ErrServiceUnavailable, reason)
	if len(msg) > 0 {
		err = errors.Wrap(err, strings.Join(msg, "; "))
	}
	return err
}

func WrapErrServiceMemoryLimitExceeded(predict, limit float32, msg ...string) error {
	err := errors.Wrapf(ErrServiceMemoryLimitExceeded, "predict=%v, limit=%v", predict, limit)
	if len(msg) > 0 {
		err = errors.Wrap(err, strings.Join(msg, "; "))
	}
	return err
}

func WrapErrServiceRequestLimitExceeded(limit int32, msg ...string) error {
	err := errors.Wrapf(ErrServiceRequestLimitExceeded, "limit=%v", limit)
	if len(msg) > 0 {
		err = errors.Wrap(err, strings.Join(msg, "; "))
	}
	return err
}

// Collection related
func WrapErrCollectionNotFound(collection any, msg ...string) error {
	err := wrapWithField(ErrCollectionNotFound, "collection", collection)
	if len(msg) > 0 {
		err = errors.Wrap(err, strings.Join(msg, "; "))
	}
	return err
}

func WrapErrCollectionNotLoaded(collection any, msg ...string) error {
	err := wrapWithField(ErrCollectionNotLoaded, "collection", collection)
	if len(msg) > 0 {
		err = errors.Wrap(err, strings.Join(msg, "; "))
	}
	return err
}

// Partition related
func WrapErrPartitionNotFound(partition any, msg ...string) error {
	err := wrapWithField(ErrPartitionNotFound, "partition", partition)
	if len(msg) > 0 {
		err = errors.Wrap(err, strings.Join(msg, "; "))
	}
	return err
}

func WrapErrPartitionNotLoaded(partition any, msg ...string) error {
	err := wrapWithField(ErrPartitionNotLoaded, "partition", partition)
	if len(msg) > 0 {
		err = errors.Wrap(err, strings.Join(msg, "; "))
	}
	return err
}

// ResourceGroup related
func WrapErrResourceGroupNotFound(rg any, msg ...string) error {
	err := wrapWithField(ErrResourceGroupNotFound, "rg", rg)
	if len(msg) > 0 {
		err = errors.Wrap(err, strings.Join(msg, "; "))
	}
	return err
}

// Replica related
func WrapErrReplicaNotFound(id int64, msg ...string) error {
	err := wrapWithField(ErrReplicaNotFound, "replica", id)
	if len(msg) > 0 {
		err = errors.Wrap(err, strings.Join(msg, "; "))
	}
	return err
}

// Channel related
func WrapErrChannelNotFound(name string, msg ...string) error {
	err := wrapWithField(ErrChannelNotFound, "channel", name)
	if len(msg) > 0 {
		err = errors.Wrap(err, strings.Join(msg, "; "))
	}
	return err
}

// Segment related
func WrapErrSegmentNotFound(id int64, msg ...string) error {
	err := wrapWithField(ErrSegmentNotFound, "segment", id)
	if len(msg) > 0 {
		err = errors.Wrap(err, strings.Join(msg, "; "))
	}
	return err
}

func WrapErrSegmentNotLoaded(id int64, msg ...string) error {
	err := wrapWithField(ErrSegmentNotLoaded, "segment", id)
	if len(msg) > 0 {
		err = errors.Wrap(err, strings.Join(msg, "; "))
	}
	return err
}

func WrapErrSegmentLack(id int64, msg ...string) error {
	err := wrapWithField(ErrSegmentLack, "segment", id)
	if len(msg) > 0 {
		err = errors.Wrap(err, strings.Join(msg, "; "))
	}
	return err
}

// Index related
func WrapErrIndexNotFound(msg ...string) error {
	err := error(ErrIndexNotFound)
	if len(msg) > 0 {
		err = errors.Wrap(err, strings.Join(msg, "; "))
	}
	return err
}

// Node related
func WrapErrNodeNotFound(id int64, msg ...string) error {
	err := wrapWithField(ErrNodeNotFound, "node", id)
	if len(msg) > 0 {
		err = errors.Wrap(err, strings.Join(msg, "; "))
	}
	return err
}

func WrapErrNodeOffline(id int64, msg ...string) error {
	err := wrapWithField(ErrNodeOffline, "node", id)
	if len(msg) > 0 {
		err = errors.Wrap(err, strings.Join(msg, "; "))
	}
	return err
}

func WrapErrNodeLack(expectedNum, actualNum int64, msg ...string) error {
	err := errors.Wrapf(ErrNodeLack, "expectedNum=%d, actualNum=%d", expectedNum, actualNum)
	if len(msg) > 0 {
		err = errors.Wrap(err, strings.Join(msg, "; "))
	}
	return err
}

// IO related
func WrapErrIoKeyNotFound(key string, msg ...string) error {
	err := errors.Wrapf(ErrIoKeyNotFound, "key=%s", key)
	if len(msg) > 0 {
		err = errors.Wrap(err, strings.Join(msg, "; "))
	}
	return err
}

func WrapErrIoFailed(key string, msg ...string) error {
	err := errors.Wrapf(ErrIoFailed, "key=%s", key)
	if len(msg) > 0 {
		err = errors.Wrap(err, strings.Join(msg, "; "))
	}
	return err
}

// Parameter related
func WrapErrParameterInvalid[T any](expected, actual T, msg ...string) error {
	err := errors.Wrapf(ErrParameterInvalid, "expected=%v, actual=%v", expected, actual)
	if len(msg) > 0 {
		err = errors.Wrap(err, strings.Join(msg, "; "))
	}
	return err
}

func WrapErrParameterInvalidRange[T any](lower, upper, actual T, msg ...string) error {
	err := errors.Wrapf(ErrParameterInvalid, "expected in (%v, %v), actual=%v", lower, upper, actual)
	if len(msg) > 0 {
		err = errors.Wrap(err, strings.Join(msg, "; "))
	}
	return err
}

func wrapWithField(err error, name string, value any) error {
	return errors.Wrapf(err, "%s=%v", name, value)
}
