package main

import (
	"context"
	"fmt"
	"net/http"
	"net/url"
	"os"
	"path"
	"strings"
	"time"

	"github.com/tencentyun/cos-go-sdk-v5"
)

func main() {
	secretID := mustGetEnv("COS_SECRET_ID")
	secretKey := mustGetEnv("COS_SECRET_KEY")
	sourceBucket := getEnv("COS_SOURCE_BUCKET", "7072-prod-5g94mx7a3d07e78c-1353115175")
	sourceRegion := getEnv("COS_SOURCE_REGION", "ap-shanghai")
	destBucket := getEnv("COS_DEST_BUCKET", "anyuyinian-1353115175")
	destRegion := getEnv("COS_DEST_REGION", "ap-shanghai")
	prefix := os.Getenv("COS_PREFIX")
	dryRun := strings.EqualFold(os.Getenv("COS_DRY_RUN"), "true")

	sourceClient := newCOSClient(secretID, secretKey, sourceBucket, sourceRegion)
	destClient := newCOSClient(secretID, secretKey, destBucket, destRegion)
	ctx := context.Background()

	if err := ensureBucketAccessible(ctx, sourceClient, "source"); err != nil {
		fatalf("source bucket check failed: %v", err)
	}
	if err := ensureBucketAccessible(ctx, destClient, "destination"); err != nil {
		fatalf("destination bucket check failed: %v", err)
	}

	fmt.Printf("start migration source=%s(%s) dest=%s(%s) prefix=%q dryRun=%v\n", sourceBucket, sourceRegion, destBucket, destRegion, prefix, dryRun)

	marker := ""
	total := 0
	success := 0

	for {
		result, _, err := sourceClient.Bucket.Get(ctx, &cos.BucketGetOptions{
			Prefix:  prefix,
			Marker:  marker,
			MaxKeys: 1000,
		})
		if err != nil {
			fatalf("list source bucket failed: %v", err)
		}

		for _, item := range result.Contents {
			total++
			key := item.Key
			if key == "" {
				continue
			}

			fmt.Printf("[%d] %s\n", total, key)
			if dryRun {
				success++
				continue
			}

			sourceURL := fmt.Sprintf("%s/%s", sourceClient.BaseURL.BucketURL.Host, encodeCOSKey(key))
			if _, _, err := destClient.Object.Copy(ctx, key, sourceURL, nil); err != nil {
				fatalf("copy %q failed: %v", key, err)
			}
			success++
		}

		if !result.IsTruncated {
			break
		}
		marker = result.NextMarker
		if marker == "" {
			fatalf("source bucket returned truncated result without next marker")
		}
	}

	fmt.Printf("migration completed total=%d success=%d dryRun=%v\n", total, success, dryRun)
}

func newCOSClient(secretID, secretKey, bucket, region string) *cos.Client {
	bucketURL, err := url.Parse(fmt.Sprintf("https://%s.cos.%s.myqcloud.com", bucket, region))
	if err != nil {
		fatalf("invalid bucket URL for %s: %v", bucket, err)
	}

	return cos.NewClient(&cos.BaseURL{BucketURL: bucketURL}, &http.Client{
		Timeout: 120 * time.Second,
		Transport: &cos.AuthorizationTransport{
			SecretID:  secretID,
			SecretKey: secretKey,
		},
	})
}

func ensureBucketAccessible(ctx context.Context, client *cos.Client, name string) error {
	_, err := client.Bucket.Head(ctx)
	if err != nil {
		return fmt.Errorf("%s bucket is not accessible: %w", name, err)
	}
	return nil
}

func encodeCOSKey(key string) string {
	parts := strings.Split(key, "/")
	for index, part := range parts {
		if part == "" {
			continue
		}
		parts[index] = url.PathEscape(part)
	}
	encoded := path.Clean(strings.Join(parts, "/"))
	if strings.HasPrefix(key, "/") {
		encoded = "/" + encoded
	}
	if encoded == "." {
		return ""
	}
	return encoded
}

func mustGetEnv(key string) string {
	value := os.Getenv(key)
	if value == "" {
		fatalf("missing required environment variable %s", key)
	}
	return value
}

func getEnv(key, fallback string) string {
	value := os.Getenv(key)
	if value == "" {
		return fallback
	}
	return value
}

func fatalf(format string, args ...interface{}) {
	fmt.Fprintf(os.Stderr, format+"\n", args...)
	os.Exit(1)
}
