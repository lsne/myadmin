/*
 * @Author: Liu Sainan
 * @Date: 2023-12-10 17:02:32
 */

package db

import "time"

type S3Bucket struct {
	Name         string
	CreationDate time.Time
}

type S3Object struct {
	Key          string
	LastModified time.Time
	Size         int64
	StorageClass string
}
